package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
	"github.com/go-sql-driver/mysql"
)

type MySQLRepository struct{ db *sql.DB }

func NewMySQLRepository(db *sql.DB) *MySQLRepository { return &MySQLRepository{db: db} }

func (r *MySQLRepository) Create(ctx context.Context, name string) (domain.Cohort, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO cohorts (name) VALUES (?)`, name)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return domain.Cohort{}, domain.ErrCohortAlreadyExists
		}
		return domain.Cohort{}, fmt.Errorf("insert cohort: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Cohort{}, fmt.Errorf("cohort last insert id: %w", err)
	}
	return r.FindByID(ctx, uint64(id))
}

func (r *MySQLRepository) FetchByPage(ctx context.Context, page, perPage uint64) ([]domain.Cohort, error) {
	if page < 1 || perPage < 1 {
		return []domain.Cohort{}, nil
	}
	return r.findMany(ctx, `SELECT id, name, created_at, updated_at FROM cohorts ORDER BY created_at DESC LIMIT ? OFFSET ?`, perPage, (page-1)*perPage)
}

func (r *MySQLRepository) FindByID(ctx context.Context, id uint64) (domain.Cohort, error) {
	return r.findOne(ctx, `SELECT id, name, created_at, updated_at FROM cohorts WHERE id=?`, id)
}

func (r *MySQLRepository) SearchByName(ctx context.Context, name string) ([]domain.Cohort, error) {
	return r.findMany(ctx, `SELECT id, name, created_at, updated_at FROM cohorts WHERE name LIKE ?`, "%"+name+"%")
}

func (r *MySQLRepository) FindByUserID(ctx context.Context, userID uint64) ([]domain.Cohort, error) {
	return r.findMany(ctx, `SELECT cohorts.id, name, created_at, updated_at FROM cohorts INNER JOIN users_cohorts ON cohorts.id = users_cohorts.cohort_id WHERE users_cohorts.user_id=?`, userID)
}

func (r *MySQLRepository) AddUserCohort(ctx context.Context, cohortID, userID uint64) ([]domain.Cohort, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users_cohorts (cohort_id, user_id) VALUES (?, ?)`, cohortID, userID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return nil, domain.ErrUserCohortExists
		}
		return nil, fmt.Errorf("insert users_cohorts: %w", err)
	}
	return r.FindByUserID(ctx, userID)
}

func (r *MySQLRepository) RemoveUserCohort(ctx context.Context, cohortID, userID uint64) ([]domain.Cohort, error) {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM users_cohorts WHERE cohort_id = ? AND user_id = ?`, cohortID, userID); err != nil {
		return nil, fmt.Errorf("delete users_cohorts: %w", err)
	}
	return r.FindByUserID(ctx, userID)
}

func (r *MySQLRepository) findOne(ctx context.Context, query string, args ...any) (domain.Cohort, error) {
	var cohort domain.Cohort
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&cohort.ID, &cohort.Name, &cohort.CreatedAt, &cohort.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Cohort{}, domain.ErrCohortNotFound
	}
	if err != nil {
		return domain.Cohort{}, fmt.Errorf("query cohort: %w", err)
	}
	return cohort, nil
}

func (r *MySQLRepository) findMany(ctx context.Context, query string, args ...any) ([]domain.Cohort, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query cohorts: %w", err)
	}
	defer rows.Close()
	var cohorts []domain.Cohort
	for rows.Next() {
		var cohort domain.Cohort
		if err := rows.Scan(&cohort.ID, &cohort.Name, &cohort.CreatedAt, &cohort.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan cohort: %w", err)
		}
		cohorts = append(cohorts, cohort)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cohort rows: %w", err)
	}
	return cohorts, nil
}
