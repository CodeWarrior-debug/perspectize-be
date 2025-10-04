.PHONY: build run test docker-up docker-down clean

build:
	go build -o bin/perspectize-server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v -cover ./...

docker-up:
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f postgres

clean:
	rm -rf bin/
	docker-compose down -v
