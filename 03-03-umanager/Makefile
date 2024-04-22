APP=umanager
SHELL := /bin/bash
DOCKER_COMPOSE_FILE=docker-compose.yaml

.PHONY: build
build:
	go build -o bin/links-srv cmd/links-srv/main.go
	go build -o bin/users-srv cmd/users-srv/main.go
	go build -o bin/api-srv cmd/api-gw/main.go

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: lint
lint:
	golangci-lint run --timeout 5m -v ./...

.PHONY: genid
genid:
	go run cmd/genid/main.go

.PHONY: generate
generate:
	protoc --go_out=pkg/pb --go_opt=paths=source_relative --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	--proto_path=./pkg/pb ./pkg/pb/common.proto

	protoc --go_out=pkg/pb --go_opt=paths=source_relative --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	--proto_path=./pkg/pb ./pkg/pb/users.proto

	protoc --go_out=pkg/pb --go_opt=paths=source_relative --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	--proto_path=./pkg/pb ./pkg/pb/links.proto

	go generate ./...

.PHONY: install
install:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@1.56.2

.PHONY: test
test:
	go test -short ./...

.PHONY: integration
integration:
	go test -race ./...

start: ## Start all deployed services
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

stop: ## Stop all deployed services
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop

.PHONY: migrate-up
migrate-up:
	 migrate -source "file://./migrations" -database "postgres://localhost:5434/users?sslmode=disable&user=postgres&password=postgres" up

.PHONY: migrate-down
migrate-down:
	 migrate -source "file://./migrations" -database \
	 "postgres://localhost:5434/users?sslmode=disable&user=postgres&password=postgres" down