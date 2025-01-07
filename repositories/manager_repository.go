package repositories

import (
	"context"
	"employee-manager/models"
	"employee-manager/types"
	"fmt"
	"net/http"
	"reflect"
	"strings"

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

func (r *ManagerRepository) PartialUpdate(id string, payload types.UpdateManagerProfilePayload) (*models.Manager, error) {
	query, args, err := buildUpdateQuery("managers", "id", id, &payload)
	if err != nil {
		return nil, err
	}
	query += " RETURNING *"
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

func GetJSONTagName(field reflect.StructField) string {
	tag := field.Tag.Get("db")

	parts := strings.Split(tag, ",")

	return parts[0]
}

func buildUpdateQuery(tableName, idField, idValue string, data interface{}) (string, pgx.NamedArgs, error) {
	val := reflect.ValueOf(data).Elem()
	typ := reflect.TypeOf(data).Elem()

	query := fmt.Sprintf("UPDATE %s SET ", tableName)
	args := pgx.NamedArgs{}
	var setClauses []string
	index := 1

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := GetJSONTagName(fieldType)

		if fieldName == "" || fieldName == "-" {
			continue
		}

		if !fieldValue.IsNil() {
			setClauses = append(setClauses, fmt.Sprintf("%s = @%s", fieldName, fieldName))
			args[fieldName] = fieldValue.Elem().Interface() // Dereference the pointer
			index++
		}
	}

	if len(setClauses) == 0 {
		return "", nil, models.NewError(http.StatusBadRequest, "no fields to update")
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE %s = @%s", idField, idField)
	args[idField] = idValue

	return query, args, nil
}