DB_URL=postgresql://postgres:12345@localhost:5678/event_scheduler?sslmode=disable
MIGRATION_NAME ?= "init_schema"

postgres:
	docker run --name docker_postgres --network  event-network -e POSTGRES_PASSWORD=12345 -p 5678:5432 -d postgres
createdb:
	docker exec -it docker_postgres createdb --username=postgres --owner=postgres event_scheduler
dropdb:
	docker exec -it docker_postgres dropdb --username=postgres event_scheduler
migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up
migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(MIGRATION_NAME)

mock:
	mockery --dir=repo --name=Repository --output=mocks --case=underscore 
	mockery --dir=cache --name=Cache --output=mocks --case=underscore
	mockery --dir=tasks --name=TaskManager --output=mocks --case=underscore

run:
	go run main.go
test:
	go test -v ./handlers \
		./middleware \
		./cache \
		-coverprofile event_scheduler.out \
        && go tool cover -html=event_scheduler.out -o event_scheduler.html
build:
	go build .
db-docs:
	dbdocs build ./docs/database.dbml --project event_scheduler
doc:
	swag init

.PHONY: postgres createdb dropdb migrate-up migrate-down run mock test build db-docs doc migrate-create