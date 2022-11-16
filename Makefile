.PHONY: run lint test
# Run
run:
	go run cmd/main.go

# Lint
lint:
	golangci-lint run -c .golangci.yml


# Test
test:
	go test ./... -v -cover