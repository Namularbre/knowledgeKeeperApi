package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Namularbre/knowledgeKeeperApi/docs"
	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/app"
	authinfra "github.com/Namularbre/knowledgeKeeperApi/internal/auth/infra"
	authhttp "github.com/Namularbre/knowledgeKeeperApi/internal/auth/infra/http"
	authsql "github.com/Namularbre/knowledgeKeeperApi/internal/auth/infra/sql"
	cohortapp "github.com/Namularbre/knowledgeKeeperApi/internal/cohort/app"
	cohortinfra "github.com/Namularbre/knowledgeKeeperApi/internal/cohort/infra"
	cohorthttp "github.com/Namularbre/knowledgeKeeperApi/internal/cohort/infra/http"
	cohortsql "github.com/Namularbre/knowledgeKeeperApi/internal/cohort/infra/sql"
	"github.com/Namularbre/knowledgeKeeperApi/internal/config"
	"github.com/Namularbre/knowledgeKeeperApi/internal/infra/db"
	httpserver "github.com/Namularbre/knowledgeKeeperApi/internal/infra/http"
	rolesapp "github.com/Namularbre/knowledgeKeeperApi/internal/roles/app"
	rolesdomain "github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
	rolesinfra "github.com/Namularbre/knowledgeKeeperApi/internal/roles/infra"
	roleshttp "github.com/Namularbre/knowledgeKeeperApi/internal/roles/infra/http"
	rolesql "github.com/Namularbre/knowledgeKeeperApi/internal/roles/infra/sql"
	subjectsapp "github.com/Namularbre/knowledgeKeeperApi/internal/subjects/app"
	subjectsinfra "github.com/Namularbre/knowledgeKeeperApi/internal/subjects/infra"
	subjectshttp "github.com/Namularbre/knowledgeKeeperApi/internal/subjects/infra/http"
	subjectssql "github.com/Namularbre/knowledgeKeeperApi/internal/subjects/infra/sql"
)

// @title           knowledgeKeeperApi
// @version         1.0.0
// @description     Personal knowledge keeper API. Auth (register/login/refresh), and metadata endpoints.
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in   header
// @name Authorization
// @description Type "Bearer {token}" where {token} is the access_token returned by /auth/login.
func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	maria, err := db.NewMariaDB(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.User,
		cfg.DB.Password,
	)
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer func() {
		_ = maria.Close()
	}()

	bootCtx, bootCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer bootCancel()

	if err := maria.Ping(bootCtx); err != nil {
		log.Fatalf("db ping error: %v", err)
	}
	log.Println("DB connection OK")

	if err := maria.ApplySchema(bootCtx, authsql.Schema); err != nil {
		log.Fatalf("schema apply error: %v", err)
	}
	log.Println("Auth schema applied")

	if err := maria.ApplySchema(bootCtx, rolesql.Schema); err != nil {
		log.Fatalf("schema apply error: %v", err)
	}
	log.Println("Roles schema applied")

	if err := maria.ApplySchema(bootCtx, subjectssql.Schema); err != nil {
		log.Fatalf("schema apply error: %v", err)
	}
	log.Println("Subjects schema applied")
	if err := maria.ApplySchema(bootCtx, cohortsql.Schema); err != nil {
		log.Fatalf("schema apply error: %v", err)
	}
	log.Println("Cohorts schema applied")

	users := authinfra.NewMySQLUserRepository(maria.DB())
	refreshes := authinfra.NewMySQLRefreshTokenRepository(maria.DB())
	hasher := authinfra.NewBcryptHasher(0)
	issuer := authinfra.NewJWTIssuer(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	roles := rolesinfra.NewMySQLRepository(maria.DB())
	subjects := subjectsinfra.NewMySQLRepository(maria.DB())
	cohorts := cohortinfra.NewMySQLRepository(maria.DB())

	authHandlers := authhttp.Handlers{
		Register: authhttp.RegisterHandler{UC: app.RegisterUser{Users: users, Hasher: hasher}},
		Login: authhttp.LoginHandler{UC: app.LoginUser{
			Users:         users,
			RefreshTokens: refreshes,
			Hasher:        hasher,
			Tokens:        issuer,
			RefreshTTL:    cfg.JWT.RefreshTTL,
		}},
		Refresh: authhttp.RefreshHandler{UC: app.RefreshSession{
			Users:         users,
			RefreshTokens: refreshes,
			Tokens:        issuer,
			RefreshTTL:    cfg.JWT.RefreshTTL,
		}},
		Me: authhttp.MeHandler{UC: app.Me{
			Users: users,
		}},
	}

	rolesHandlers := roleshttp.Handlers{
		CreateRole: roleshttp.CreateRoleHandler{UC: rolesapp.CreateRole{
			Roles: roles,
		}},
		FindByID: roleshttp.FindByIDHandler{UC: rolesapp.FindByID{
			Roles: roles,
		}},
		FindByUserID: roleshttp.FindByUserIDHandler{UC: rolesapp.FindByUserID{
			Roles: roles,
		}},
		AddUserRole: roleshttp.AddUserRoleHandler{UC: rolesapp.AddUserRole{
			Roles: roles,
		}},
		RemoveUserRole: roleshttp.RemoveUserRoleHandler{UC: rolesapp.RemoveUserRole{
			Roles: roles,
		}},
		SearchByLabel: roleshttp.SearchByLabelHandler{UC: rolesapp.SearchByLabel{
			Roles: roles,
		}},
	}

	subjectsHandlers := subjectshttp.Handlers{
		CreateSubject:     subjectshttp.CreateSubjectHandler{UC: subjectsapp.CreateSubject{Subjects: subjects}},
		FindByID:          subjectshttp.FindByIDHandler{UC: subjectsapp.FindByID{Subjects: subjects}},
		FindByUserID:      subjectshttp.FindByUserIDHandler{UC: subjectsapp.FindByUserID{Subjects: subjects}},
		AddUserSubject:    subjectshttp.AddUserSubjectHandler{UC: subjectsapp.AddUserSubject{Subjects: subjects}},
		RemoveUserSubject: subjectshttp.RemoveUserSubjectHandler{UC: subjectsapp.RemoveUserSubject{Subjects: subjects}},
		SearchByName:      subjectshttp.SearchByNameHandler{UC: subjectsapp.SearchByName{Subjects: subjects}},
	}

	cohortHandlers := cohorthttp.Handlers{
		CreateCohort:     cohorthttp.CreateCohortHandler{UC: cohortapp.CreateCohort{Cohorts: cohorts}},
		FindByID:         cohorthttp.FindByIDHandler{UC: cohortapp.FindByID{Cohorts: cohorts}},
		FindByUserID:     cohorthttp.FindByUserIDHandler{UC: cohortapp.FindByUserID{Cohorts: cohorts}},
		AddUserCohort:    cohorthttp.AddUserCohortHandler{UC: cohortapp.AddUserCohort{Cohort: cohorts}},
		RemoveUserCohort: cohorthttp.RemoveUserCohortHandler{UC: cohortapp.RemoveUserCohort{Cohorts: cohorts}},
		SearchByName:     cohorthttp.SearchByNameHandler{UC: cohortapp.SearchByName{Cohorts: cohorts}},
	}

	server := httpserver.NewServer(cfg.Port)
	server.RegisterRoutes(func(mux *http.ServeMux) {
		mux.Handle("/auth/register", httpserver.LogMiddleware(authHandlers.Register))
		mux.Handle("/auth/login", httpserver.LogMiddleware(authHandlers.Login))
		mux.Handle("/auth/refresh", httpserver.LogMiddleware(authHandlers.Refresh))
		mux.Handle("/auth/me", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, authHandlers.Me)))

		mux.Handle("/roles/create", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, authhttp.RequireAnyRole(roles, rolesdomain.RoleAdmin)(rolesHandlers.CreateRole))))
		mux.Handle("/roles/findbyid", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, rolesHandlers.FindByID)))
		mux.Handle("/roles/finduserroles", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, rolesHandlers.FindByUserID)))
		mux.Handle("/roles/adduserrole", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, authhttp.RequireAnyRole(roles, rolesdomain.RoleAdmin)(rolesHandlers.AddUserRole))))
		mux.Handle("/roles/removeuserrole", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, authhttp.RequireAnyRole(roles, rolesdomain.RoleAdmin)(rolesHandlers.RemoveUserRole))))
		mux.Handle("/roles/searchbylabel", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, rolesHandlers.SearchByLabel)))

		mux.Handle("/subjects/create", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, subjectsHandlers.CreateSubject)))
		mux.Handle("/subjects/findbyid", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, subjectsHandlers.FindByID)))
		mux.Handle("/subjects/findusersubjects", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, subjectsHandlers.FindByUserID)))
		mux.Handle("/subjects/addusersubject", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, subjectsHandlers.AddUserSubject)))
		mux.Handle("/subjects/removeusersubject", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, subjectsHandlers.RemoveUserSubject)))
		mux.Handle("/subjects/searchbyname", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, subjectsHandlers.SearchByName)))

		mux.Handle("/cohorts/create", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, cohortHandlers.CreateCohort)))
		mux.Handle("/cohorts/findbyid", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, cohortHandlers.FindByID)))
		mux.Handle("/cohorts/findusercohorts", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, cohortHandlers.FindByUserID)))
		mux.Handle("/cohorts/addusercohort", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, cohortHandlers.AddUserCohort)))
		mux.Handle("/cohorts/removeusercohort", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, cohortHandlers.RemoveUserCohort)))
		mux.Handle("/cohorts/searchbyname", httpserver.LogMiddleware(authhttp.RequireBearer(issuer, cohortHandlers.SearchByName)))
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
