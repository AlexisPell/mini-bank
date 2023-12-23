# create container
psql_run:
	docker run --name simplebank-postgres -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=123456 -d postgres:15-alpine

# create db in the container
psql_create_db:
	docker exec -it simplebank-postgres createdb --username=user simple_bank

# start container
psql_up:
	docker start simplebank-postgres

# stop container
psql_down:
	docker stop simplebank-postgres

# run migrations
migrateup:
	migrate -path db/migration -database "postgresql://user:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

# remove migrations
migratedown:
	migrate -path db/migration -database "postgresql://user:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test ./... -v -cover

.PHONY: psql_run psql_up psql_down psql_create_db migrateup migratedown sqlc test