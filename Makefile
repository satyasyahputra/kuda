include .env
export

tidy:
	@go mod tidy

worker:
	@go run ./cmd/worker/main.go