package repositories

import (
	"context"
	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/types"
	"fmt"

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

func (r *ManagerRepository) Save(manager models.Manager) (*models.Manager, error) {
	query := `
	INSERT INTO managers (
		email,
		password
	) 
	VALUES (
		LOWER(@email),
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

func (r *ManagerRepository) FindById(id string) (*models.Manager, error) {
	query := `SELECT * FROM managers WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	manager, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Manager])
	if err != nil {
		return nil, err
	}

	return &manager, nil
}

func (r *ManagerRepository) FindByEmail(email string) (*models.Manager, error) {
	query := `SELECT * FROM managers WHERE email = @email`
	args := pgx.NamedArgs{
		"email": email,
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	manager, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Manager])
	if err != nil {
		return nil, err
	}

	return &manager, nil
}

func (r *ManagerRepository) PartialUpdate(id string, payload types.UpdateManagerProfilePayload) (*models.Manager, error) {
	query, args, err := lib.BuildPartialUpdateQuery("managers", "id", id, &payload)
	if err != nil {
		return nil, err
	}
	
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}

	manager, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Manager])
	if err != nil {
		return nil, err
	}

	return &manager, nil
}
