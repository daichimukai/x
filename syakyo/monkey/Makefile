_all: run

.PHONY: test
test:
	@go test -v -race ./...

.PHONY: generate
generate:
	@go generate ./...

.PHONY: run
run:
	@go run main.go
