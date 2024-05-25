postgres:
	docker run --name postgres -p 2024\:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=secret -d postgres:latest

createdb:
	docker exec -it postgres createdb --username=admin --owner=admin simple-order-go

dropdb:
	docker exec -it postgres dropdb simple-order-go

migrateup:
	migrate -path migration -database "postgresql\://admin\:secret@localhost\:2024/simple-order-go?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql\://admin\:secret@localhost\:2024/simple-order-go?sslmode=disable" -verbose down

test:
	go test ./... -v -cover

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown test server
