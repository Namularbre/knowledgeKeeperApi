package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
	"github.com/go-sql-driver/mysql"
)

type MySqlRoleRepository struct {
	db *sql.DB
}

func NewMySqlRoleRepository(db *sql.DB) *MySqlRoleRepository { return &MySqlRoleRepository{db: db} }

func (r *MySqlRoleRepository) Create(ctx context.Context, label string) (domain.Role, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO roles (label) VALUES (?)`,
		label,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return domain.Role{}, errors.New("role already exists")
		}
		return domain.Role{}, fmt.Errorf("insert user %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Role{}, fmt.Errorf("role last insert id %w", err)
	}

	return r.FindById(ctx, id)
}

func (r *MySqlRoleRepository) AddUserRole(ctx context.Context, roleID int64, userID int64) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users_roles (role_id, user_id) VALUE (?, ?);`,
		roleID, userID,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return errors.New("role already exists")
		}
		return fmt.Errorf("insert users_roles %w", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("users_roles last insert id %w", err)
	}
	return nil
}

func (r *MySqlRoleRepository) RemoveUserRole(ctx context.Context, roleID int64, userID int64) error {
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

func (r *MySqlRoleRepository) FindById(ctx context.Context, id int64) (domain.Role, error) {
	return r.findOne(ctx,
		`SELECT id, label, created_at, updated_at FROM roles WHERE id=?`,
		id,
	)
}

func (r *MySqlRoleRepository) FindByLabel(ctx context.Context, label string) (domain.Role, error) {
	return r.findOne(ctx,
		`SELECT id, label, created_at, updated_at FROM roles WHERE label=?`,
		label,
	)
}

func (r *MySqlRoleRepository) FindByUserID(ctx context.Context, userID int64) (domain.Role, error) {
	return r.findOne(ctx,
		`SELECT roles.id, label, created_at, updated_at 
				FROM roles 
				    INNER JOIN users_roles ON (roles.id = users_roles.role_id)
				WHERE users_roles.user_id=?`,
		userID,
	)
}

func (r *MySqlRoleRepository) findOne(ctx context.Context, query string, args any) (domain.Role, error) {
	var role domain.Role
	err := r.db.QueryRowContext(ctx, query, args).
		Scan(&role.ID, &role.Label, &role.CreatedAt, &role.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Role{}, domain.ErrRoleNotFound
	}
	if err != nil {
		return domain.Role{}, fmt.Errorf("query user %w", err)
	}
	return role, nil
}
