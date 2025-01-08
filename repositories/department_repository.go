package repositories

import (
	"context"
	"employee-manager/models"
	"fmt"

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

func (r *DepartmentRepository) GetAllDepartment(offset int, limit int, name string) ([]models.Department, error) {
	query := fmt.Sprintf(`
	SELECT * 
	FROM departments
	WHERE LOWER(name) LIKE '%%%s%%'
	ORDER BY created_at
	LIMIT @limit
	OFFSET @offset
	`, name)
	args := pgx.NamedArgs{
		"limit": limit,
		"offset": offset,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}
	departments, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Department])
	if err != nil {
		return nil, err
	}

	return departments, nil
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