FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./.bin/app ./cmd/main.go
#EXPOSE 8000
#CMD ["./.bin/app"]

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/.bin/app .bin/app
COPY --from=builder /app/configs configs/
COPY --from=builder /app/.env .

EXPOSE 8000

ENTRYPOINT [".bin/app"]

