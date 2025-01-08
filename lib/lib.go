package lib

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

func SetJsonResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

func GetJSONTagName(field reflect.StructField) string {
	tag := field.Tag.Get("db")

	parts := strings.Split(tag, ",")

	return parts[0]
}

func BuildPartialUpdateQuery(tableName, idField, idValue string, data interface{}) (string, pgx.NamedArgs, error) {
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
		query = fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE %s = @%s
		`, tableName, idField, idField)
		args = pgx.NamedArgs{
			idField: idValue,
		}
		return query, args, nil
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE %s = @%s", idField, idField)
	args[idField] = idValue
	query += " RETURNING *"

	return query, args, nil
}