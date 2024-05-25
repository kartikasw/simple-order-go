postgres:
    docker run --name postgres12 -p 2024\:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSSWORD=secret -d postgres\:12-alpine

createdb:
    docker exec -it postgres12 createdb --username=admin --owner=root simple_order_go

dropdb:
    docker exec -it postgres12 dropdb simple_order_go

migrateup:
    migrate -path migration -database "postgresql://admin:secret@localhost:2024/simple-order-go?sslmode=disable" -verbose up

migratedown:
    migrate -path migration -database "postgresql://admin:secret@localhost:2024/simple-order-go?sslmode=disable" -verbose down

test:
    go test ./... -v -cover

server:
    go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown test server
