DATABASE_URL ?= postgres://gravehard:gravehard@localhost:5432/gravehard?sslmode=disable

.PHONY: run test migrate-up migrate-down seed lint sqlc swagger

run:
	cd backend && go run ./cmd/api

test:
	cd backend && go test -v -race ./...

migrate-up:
	docker run --rm --network host -v "$(PWD)/backend/migrations:/migrations" migrate/migrate:v4.18.2 -path=/migrations -database "$(DATABASE_URL)" up

migrate-down:
	docker run --rm --network host -v "$(PWD)/backend/migrations:/migrations" migrate/migrate:v4.18.2 -path=/migrations -database "$(DATABASE_URL)" down 1

seed:
	cd backend && go run ./cmd/seed

lint:
	cd backend && golangci-lint run ./...

sqlc:
	cd backend && go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0 generate

swagger:
	cd backend && go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g cmd/api/main.go -o docs
