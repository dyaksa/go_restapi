run:
	go run main.go

migrate:
	migrate -path db/migrations -database "postgresql://admin:password@localhost:5432/belajar_golang?sslmode=disable" up

rollback:
	migrate -path db/migrations -database "postgresql://admin:password@localhost:5432/belajar_golang?sslmode=disable" down

create-migration:
	migrate create -ext sql -dir db/migrations $(name)

PHONY: run migrate rollback create-migration