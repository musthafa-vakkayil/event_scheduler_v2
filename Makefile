postgres:
	docker run --name docker_postgres -e POSTGRES_PASSWORD=12345 -p 5678:5432 -d postgres
createdb:
	docker exec -it docker_postgres createdb --username=postgres --owner=postgres event_scheduler
dropdb:
	docker exec -it docker_postgres dropdb --username=postgres event_scheduler
migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:12345@localhost:5678/event_scheduler?sslmode=disable" -verbose up
migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:12345@localhost:5678/event_scheduler?sslmode=disable" -verbose down
sqlc:
	sqlc generate
run:
	go run main.go
.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc run