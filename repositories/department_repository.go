package repositories

import (
	"context"
	"employee-manager/lib"
	"employee-manager/models"
	"employee-manager/types"
	"fmt"
	"net/http"

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
	WHERE 
		LOWER(name) LIKE '%%%s%%' 
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
	INSERT INTO departments (name, manager_id) VALUES (@name, @manager_id) RETURNING *
	`
	args := pgx.NamedArgs{
		"name": department.Name,
		"manager_id": department.ManagerId,
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	newDepartment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Department])
	if err != nil {
		return nil, err
	}

	return &newDepartment, nil
}

func (r *DepartmentRepository) PartialUpdate(id string, payload types.UpdateDepartmentProfilePayload) (*models.Department, error) {
	query, args, err := lib.BuildPartialUpdateQuery("departments", "id", id, &payload)
	if err != nil {
		return nil, err
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}

	department, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Department])
	if err != nil {
		return nil, err
	}

	return &department, nil
}

func (r *DepartmentRepository) Delete(id string) error {
	query := `DELETE FROM departments WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	commandTag, err := r.pgConn.Exec(r.ctx, query, args); 
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.NewError(http.StatusNotFound, "")
	}

	return nil
}