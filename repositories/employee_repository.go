package repositories

import (
	"context"
	"employee-manager/models"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeRepository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
}

func NewEmployeeRepository(ctx context.Context, pgConn *pgxpool.Pool) EmployeeRepository {
	return EmployeeRepository{ctx, pgConn}
}

func (r *EmployeeRepository) Save(employee models.Employee) (*models.Employee, error) {
	query := `
	INSERT INTO employees (
		identity_number
		,name
		,employee_image_uri
		,gender
		,department_id
	) 
	VALUES (
		@identity_number
		,@name
		,@employee_image_uri
		,@gender
		,@department_id
	)
	RETURNING *
	`
	args := pgx.NamedArgs{
		"identity_number":    employee.IdentityNumber,
		"name":               employee.Name,
		"employee_image_uri": employee.EmployeeImageUri,
		"gender":             employee.Gender,
		"department_id":      employee.DepartmentId,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}

	newEmployee, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Employee])
	if err != nil {
		return nil, err
	}

	return &newEmployee, nil
}
