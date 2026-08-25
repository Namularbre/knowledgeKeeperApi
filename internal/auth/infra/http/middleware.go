package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
	rolesdomain "github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type ctxKey struct{}

var userIDKey = ctxKey{}

// RequireBearer wraps an http.Handler so it only executes when a valid
// access token is present in the Authorization header. The authenticated
// user ID is injected into the request context.
func RequireBearer(tokens domain.TokenIssuer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing_bearer_token")
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		userID, err := tokens.ParseAccessToken(raw)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid_access_token")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFrom extracts the authenticated user ID injected by RequireBearer.
// The second return value indicates whether the context carried an ID.
func UserIDFrom(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDKey).(int64)
	return v, ok
}

// RequireAnyRole wraps a handler so it only executes when the authenticated
// user has at least one of the allowed roles. It must be wrapped by
// RequireBearer so the authenticated user ID is available in the context.
func RequireAnyRole(roles rolesdomain.Repository, allowed ...string) func(http.Handler) http.Handler {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, label := range allowed {
		allowedSet[strings.ToLower(strings.TrimSpace(label))] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFrom(r.Context())
			if !ok || userID < 0 {
				writeError(w, http.StatusUnauthorized, "unauthenticated")
				return
			}

			assigned, err := roles.FindByUserID(r.Context(), uint64(userID))
			if err != nil {
				writeError(w, http.StatusInternalServerError, "role_lookup_failed")
				return
			}

			for _, role := range assigned {
				if _, ok := allowedSet[strings.ToLower(strings.TrimSpace(role.Label))]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeError(w, http.StatusForbidden, "insufficient_role")
		})
	}
}
