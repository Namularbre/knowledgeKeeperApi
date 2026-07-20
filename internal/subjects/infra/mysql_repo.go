package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
	"github.com/go-sql-driver/mysql"
)

type MySQLRepository struct{ db *sql.DB }

func NewMySQLRepository(db *sql.DB) *MySQLRepository { return &MySQLRepository{db: db} }

func (r *MySQLRepository) Create(ctx context.Context, name string) (domain.Subject, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO subjects (name) VALUES (?)`, name)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return domain.Subject{}, domain.ErrSubjectAlreadyExists
		}
		return domain.Subject{}, fmt.Errorf("insert subject: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Subject{}, fmt.Errorf("subject last insert id: %w", err)
	}
	return r.FindByID(ctx, uint64(id))
}

func (r *MySQLRepository) FetchByPage(ctx context.Context, page, perPage uint64) ([]domain.Subject, error) {
	if page < 1 || perPage < 1 {
		return []domain.Subject{}, nil
	}
	return r.findMany(ctx, `SELECT id, name, created_at, updated_at FROM subjects ORDER BY created_at DESC LIMIT ? OFFSET ?`, perPage, (page-1)*perPage)
}

func (r *MySQLRepository) FindByID(ctx context.Context, id uint64) (domain.Subject, error) {
	return r.findOne(ctx, `SELECT id, name, created_at, updated_at FROM subjects WHERE id=?`, id)
}

func (r *MySQLRepository) SearchByName(ctx context.Context, name string) ([]domain.Subject, error) {
	return r.findMany(ctx, `SELECT id, name, created_at, updated_at FROM subjects WHERE name LIKE ?`, "%"+name+"%")
}

func (r *MySQLRepository) FindByUserID(ctx context.Context, userID uint64) ([]domain.Subject, error) {
	return r.findMany(ctx, `SELECT subjects.id, name, created_at, updated_at FROM subjects INNER JOIN users_subjects ON subjects.id = users_subjects.subject_id WHERE users_subjects.user_id=?`, userID)
}

func (r *MySQLRepository) AddUserSubject(ctx context.Context, subjectID, userID uint64) ([]domain.Subject, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users_subjects (subject_id, user_id) VALUES (?, ?)`, subjectID, userID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return nil, domain.ErrUserSubjectExists
		}
		return nil, fmt.Errorf("insert users_subjects: %w", err)
	}
	return r.FindByUserID(ctx, userID)
}

func (r *MySQLRepository) RemoveUserSubject(ctx context.Context, subjectID, userID uint64) ([]domain.Subject, error) {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM users_subjects WHERE subject_id = ? AND user_id = ?`, subjectID, userID); err != nil {
		return nil, fmt.Errorf("delete users_subjects: %w", err)
	}
	return r.FindByUserID(ctx, userID)
}

func (r *MySQLRepository) findOne(ctx context.Context, query string, args ...any) (domain.Subject, error) {
	var subject domain.Subject
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&subject.ID, &subject.Name, &subject.CreatedAt, &subject.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Subject{}, domain.ErrSubjectNotFound
	}
	if err != nil {
		return domain.Subject{}, fmt.Errorf("query subject: %w", err)
	}
	return subject, nil
}

func (r *MySQLRepository) findMany(ctx context.Context, query string, args ...any) ([]domain.Subject, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query subjects: %w", err)
	}
	defer rows.Close()
	var subjects []domain.Subject
	for rows.Next() {
		var subject domain.Subject
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.CreatedAt, &subject.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan subject: %w", err)
		}
		subjects = append(subjects, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("subject rows: %w", err)
	}
	return subjects, nil
}
