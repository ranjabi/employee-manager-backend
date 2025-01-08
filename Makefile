POSTGRESQL_URL=postgres://dev:123456@localhost:5432/employee_manager?sslmode=disable

test:
	echo $(1)
	
mig.create:
	migrate create -ext sql -dir db/migrations -seq $(n)

mig.up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

mig.down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down

mig.fix:
	migrate -path db/migrations -database ${POSTGRESQL_URL} force $(v)

fmt:
	gofmt -w .