package repositories

import (
	"context"
	"employee-manager/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ManagerRepository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
}

func NewManagerRepository(ctx context.Context, pgConn *pgxpool.Pool) ManagerRepository {
	return ManagerRepository{ctx, pgConn}
}

func (r *ManagerRepository) CreateManager(manager models.Manager) (*models.Manager, error) {
	query := `
		INSERT INTO managers (
			email,
			password
		) 
		VALUES (
			@email,
			@password
		)
		RETURNING id, email
	`
	args := pgx.NamedArgs{
		"email":    manager.Email,
		"password": manager.Password,
	}

	var newManager models.Manager
	if err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&newManager.Id, &newManager.Email); err != nil {
		return nil, err
	}

	return &newManager, nil
}
