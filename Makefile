run:
	go run cmd/main.go

migrate.up:
	go run migration/main/main.go up

migration.down:
	go run migration/main/main.go rollback

test:
	go test ./src/business/usecase/...