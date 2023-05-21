FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app/
COPY ./configs ./configs/
COPY ./migrations ./migrations/
COPY ./docs ./docs/
COPY .env .


#RUN go mod download
RUN GOOS=linux go build -o ./.bin/app ./app/cmd/main.go
#
#RUN chmod 777  ./.bin/app
#EXPOSE 8080
CMD ["./.bin/app"]

#FROM alpine:latest
#
#WORKDIR /root/
#
#
#
#COPY --from=0 /root/service/.bin/app .
#COPY --from=0 /root/service/configs configs/
#
#RUN ls
#
#CMD ["./app"]