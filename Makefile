run:
	go run ./cmd/

build:
	go build -o ./bin/main ./cmd/

start: build
	./bin/main

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

test:
	go test -v ./...