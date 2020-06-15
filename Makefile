lint:
	golangci-lint run ./...

build: gen
	go build -o multiplexer main.go
