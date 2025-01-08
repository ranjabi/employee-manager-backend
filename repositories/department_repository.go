package repositories

import (
	"context"
	"employee-manager/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DepartmentRepository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
}

func NewDepartmentRepository(ctx context.Context, pgConn *pgxpool.Pool) DepartmentRepository {
	return DepartmentRepository{ctx, pgConn}
}

func (r *DepartmentRepository) Save(department models.Department) (*models.Department, error) {
	query := `
	INSERT INTO departments (name) VALUES (@name) RETURNING *
	`
	args := pgx.NamedArgs{
		"name": department.Name,
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	newDepartment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Department])
	if err != nil {
		return nil, err
	}

	return &newDepartment, nil
}