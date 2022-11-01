FROM golang:alpine3.16 as builder

WORKDIR /app
RUN apk add --no-cache gcc musl-dev linux-headers git make protoc

COPY . /app/

# generate proto contracts
RUN cd /app && go get -u google.golang.org/protobuf/cmd/protoc-gen-go && go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN cd /app && go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN export PATH="$PATH:$(go env GOPATH)/bin" && cd /app && make proto-gen

RUN cd /app && go mod tidy
RUN cd /app && go build -o ./build/bin/main -gcflags="all=-N -l" ./cmd/main.go


FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/build/bin /app

ENTRYPOINT [ "/app/main" ]