

build-race:
	go build -v -race -o ./bin/goapp ./cmd/main.go

build:
	go build -v -o ./bin/goapp ./cmd/main.go

install:
	go install -v ./...
	go test -i ./...

test-short: install
	go test -v -short ./...

test: install
	go test -v -count 1 -race -short ./...
	go test -v -count 1 -race -timeout 5s ./tests --address=127.0.0.1:3031

run: build-race
	GORACE="halt_on_error=1" ./bin/goapp
