run:
	go run cmd/main.go

air:
	air

build:
	go build -o bin/api-template ./cmd/main.go

swagger:
	swag init -g cmd/main.go --output docs

migrate:
	go run cmd/main.go migrate

migrate-refresh:
	go run cmd/main.go migrate:refresh

fix-dirty:
	go run cmd/main.go fix:dirty

seed:
	go run cmd/main.go db:seed

test:
	go test ./...

.PHONY: run air build swagger migrate migrate-refresh fix-dirty seed test

