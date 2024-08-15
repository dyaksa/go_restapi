run:
	go run main.go

install:
	go mod download
	
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags musl -o main main.go

migrate:
	migrate -path db/migrations -database "postgresql://admin:password@localhost:5432/belajar_golang?sslmode=disable" up

rollback:
	migrate -path db/migrations -database "postgresql://admin:password@localhost:5432/belajar_golang?sslmode=disable" down

create-migration:
	migrate create -ext sql -dir db/migrations $(name)

PHONY: run migrate rollback create-migration
