include .env
export 

tidy:
	@go mod tidy

worker:
	@go run ./cmd/worker/main.go

endpoint:
	@go run ./cmd/endpoint/main.go

build:
	@go build -o ./build/worker ./cmd/worker/main.go
	@go build -o ./build/endpoint ./cmd/endpoint/main.go
