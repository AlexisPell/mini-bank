# create container
psql_run:
	docker run --name minibank-postgres -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=123456 -d postgres:15-alpine

# create db in the container
psql_create_db:
	docker exec -it minibank-postgres createdb --username=user minibank

# start container
psql_up:
	docker start minibank-postgres

# stop container
psql_down:
	docker stop minibank-postgres

# run migrations
migrate_up:
	migrate -path db/migration -database "postgresql://user:123456@localhost:5432/minibank?sslmode=disable" -verbose up

# remove migrations
migrate_down:
	migrate -path db/migration -database "postgresql://user:123456@localhost:5432/minibank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test ./... -v -cover

server:
	go run cmd/minibank/main.go

.PHONY: psql_run psql_up psql_down psql_create_db migrate_up migrate_down sqlc test server