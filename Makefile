include .env
export $(shell sed 's/=.*//' .env)
.PHONY:
.SILENT:

build:
	go mod download && go build -o ./.bin/app ./app/cmd/main.go

run: build
	./.bin/app

migrate-up:
	migrate -path ./migrations -database 'postgres://pguser:${DB_PASSWORD}@localhost:5431/devdb?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://pguser:${DB_PASSWORD}@localhost:5431/devdb?sslmode=disable' down

