run:
	go run ./cmd/

build:
	go build -o ./bin/main ./cmd/

start: build
	./bin/main