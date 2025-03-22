# Build Stage
FROM golang:1.24.1-alpine3.21 AS builder
WORKDIR /app
COPY . . 
RUN apk add curl
RUN go build -o main main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY migrations ./migrations
COPY app.env .
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]