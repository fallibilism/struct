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

proto:
	cd ./protocol && make generate_go

test:
	go test -v ./...
