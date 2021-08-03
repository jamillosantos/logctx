
.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: test
test: lint
	go test ./... -v -cover