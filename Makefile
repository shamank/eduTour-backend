.PHONY:
.SILENT:

build:
	go mod download && go build -o ./.bin/app ./app/cmd/main.go

run: build
	./.bin/app