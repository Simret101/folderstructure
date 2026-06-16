migrate-down:
	- migrate -database postgres://postgres:secret@localhost:5432/choiceloan?sslmode=disable -path internal/query/schemas -verbose down
migrate-up:
	- migrate -database postgres://postgres:secret@localhost:5432/choiceloan?sslmode=disable -path internal/query/schemas -verbose up
migrate-create:
	- migrate create -ext sql -dir internal/query/schemas -tz "UTC" $(name)
lint:
	- golangci-lint run --timeout 5m --verbose

swagger:
	-swag fmt && swag init -g cmd/main.go
sqlc:
	cd ./config && sqlc generate