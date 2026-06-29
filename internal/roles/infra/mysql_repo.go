package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
	"github.com/go-sql-driver/mysql"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository { return &MySQLRepository{db: db} }

func (r *MySQLRepository) Create(ctx context.Context, label string) (domain.Role, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO roles (label) VALUES (?)`,
		label,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return domain.Role{}, domain.ErrRoleAlreadyExists
		}
		return domain.Role{}, fmt.Errorf("insert role %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Role{}, fmt.Errorf("role last insert id %w", err)
	}

	return r.FindByID(ctx, id)
}

func (r *MySQLRepository) AddUserRole(ctx context.Context, roleID int64, userID int64) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users_roles (role_id, user_id) VALUE (?, ?);`,
		roleID, userID,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return errors.New("user role already exists")
		}
		return fmt.Errorf("insert users_roles %w", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("users_roles last insert id %w", err)
	}
	return nil
}

func (r *MySQLRepository) RemoveUserRole(ctx context.Context, roleID int64, userID int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM users_roles WHERE role_id = ? AND user_id = ?;`,
		roleID, userID,
	)

	if err != nil {
		return fmt.Errorf("delete users_roles %w", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete users_roles rows affected %w", err)
	}
	return nil
}

func (r *MySQLRepository) FetchByPage(ctx context.Context, page, perPage uint64) ([]domain.Role, error) {
	if page < 1 || perPage < 1 {
		return []domain.Role{}, nil
	}
	offset := (page - 1) * perPage
	return r.findMany(ctx,
		`SELECT id, label, created_at, updated_at FROM roles ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		perPage, offset,
	)
}

func (r *MySQLRepository) FindByID(ctx context.Context, id int64) (domain.Role, error) {
	return r.findOne(ctx,
		`SELECT id, label, created_at, updated_at FROM roles WHERE id=?`,
		id,
	)
}

func (r *MySQLRepository) SearchByLabel(ctx context.Context, label string) ([]domain.Role, error) {
	return r.findMany(ctx,
		`SELECT id, label, created_at, updated_at FROM roles WHERE label LIKE ?`,
		"%"+label+"%",
	)
}

func (r *MySQLRepository) FindByUserID(ctx context.Context, userID int64) (domain.Role, error) {
	return r.findOne(ctx,
		`SELECT roles.id, label, created_at, updated_at
				FROM roles
				    INNER JOIN users_roles ON (roles.id = users_roles.role_id)
				WHERE users_roles.user_id=?`,
		userID,
	)
}

func (r *MySQLRepository) findOne(ctx context.Context, query string, args ...any) (domain.Role, error) {
	var role domain.Role
	err := r.db.QueryRowContext(ctx, query, args).
		Scan(&role.ID, &role.Label, &role.CreatedAt, &role.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Role{}, domain.ErrRoleNotFound
	}
	if err != nil {
		return domain.Role{}, fmt.Errorf("query role %w", err)
	}
	return role, nil
}

func (r *MySQLRepository) findMany(ctx context.Context, query string, args ...any) ([]domain.Role, error) {
	rows, err := r.db.QueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("query roles %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Label, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan role %w", err)
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error %w", err)
	}

	return roles, nil
}
