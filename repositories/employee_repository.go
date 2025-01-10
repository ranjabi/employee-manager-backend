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

func (r *EmployeeRepository) GetAllEmployee(offset int, limit int, identityNumber string, name string, gender string, departmentId string) ([]models.Employee, error) {
	query := fmt.Sprintf(`
	SELECT * 
	FROM employees
	WHERE 
		LOWER(identity_number) LIKE '%%%s%%'
		AND LOWER(name) LIKE '%%%s%%'
		AND gender = '%s'
		AND department_id = '%s'
	ORDER BY created_at
	LIMIT @limit
	OFFSET @offset
	`, identityNumber, name, gender, departmentId)
	args := pgx.NamedArgs{
		"limit": limit,
		"offset": offset,
	}
	
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}

	employees, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Employee])
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *EmployeeRepository) PartialUpdate(identityNumber string, payload types.UpdateEmployeePayload) (*models.Employee, error) {
	query, args, err := lib.BuildPartialUpdateQuery("employees", "identity_number", identityNumber, &payload)
	if err != nil {
		return nil, err
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}
	fmt.Printf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n", query, args, rows)
	employee, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Employee])
	if err != nil {
		return nil, models.NewError(http.StatusNotFound, "identityNumber is not found")
	}
	

	return &employee, nil
}