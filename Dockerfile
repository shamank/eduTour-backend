FROM golang:latest

WORKDIR /core

COPY ./app/ ./app

RUN go mod download
RUN go build -o /.bin/app ./app/cmd/main.go

CMD ["/.bin/app"]