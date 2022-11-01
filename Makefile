CURDIR := $(shell pwd)

.PHONY: proto-gen
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        ./proto/converter/*.proto

# tests section
.PHONY: test-unit
test-unit: ## Unit testing
	go test -v ./intenral/...


# linter section
.PHONY: lint
lint:
	golangci-lint --config=./golangci.yml run -v