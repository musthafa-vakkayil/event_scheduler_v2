DB_URL=postgresql://postgres:12345@localhost:5678/event_scheduler?sslmode=disable

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
mock:
	mockery --dir=repo --name=Repository --output=mocks --case=underscore
run:
	go run main.go
test:
	go test -v ./handlers \
		./middleware \
		-coverprofile event_scheduler.out \
        && go tool cover -html=event_scheduler.out -o event_scheduler.html
build:
	go build .
db-docs:
	dbdocs build ./docs/database.dbml --project event_scheduler

.PHONY: postgres createdb dropdb migrate-up migrate-down run mock test build db-docs