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

.PHONY: test-integration
test-integration: docker-build docker-run
	@echo "==> Running integration tests"
	go test -v ./tests/integration/...
	make docker-stop


# linter section
.PHONY: lint
lint:
	golangci-lint --config=./golangci.yml run -v


# docker section
.PHONY: docker-build
docker-build:
	docker build -t crypto-to-fiat-converter .

.PHONY: docker-run
docker-run:
	docker run --name crypto-to-fiat-converter -d --rm -p 8080:8080 crypto-to-fiat-converter:latest

.PHONY: docker-stop
docker-stop:
	docker stop crypto-to-fiat-converter
