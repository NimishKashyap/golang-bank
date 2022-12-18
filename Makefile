postgres:
	docker run --name postgresdb -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
startdb:
	docker start postgresdb
createdb:
	docker exec -it postgresdb createdb --username=root --owner=root simple_bank 
dropdb:
	docker exec -it postgresdb dropdb simple_bank 
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc: 
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: createdb dropdb postgres startdb migrateup migratedown sqlc server