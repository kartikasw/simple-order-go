migrateup:
	migrate -path migration -database "postgresql://admin:secret@localhost:2024/order_assignment?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://admin:secret@localhost:2024/order_assignment?sslmode=disable" -verbose down

test:
	go test ./... -v -cover

server:
	go run main.go

mock:
	mockgen -package mockService -destination internal/service/mock/order_service.go simple-order-go/internal/service IOrderService

.PHONY: migrateup migratedown test server
