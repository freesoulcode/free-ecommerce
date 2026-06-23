GOPATH_BIN := $(shell go env GOPATH)/bin
PROTOC := PATH="$(GOPATH_BIN):$$PATH" protoc --proto_path=.

.PHONY: help proto proto-user proto-auth test test-user-service test-buyer-api mysql-up mysql-status mysql-down run-user-service run-buyer-api run-auth-service

help:
	@printf "%-24s %s\n" "proto" "Generate all Go/gRPC proto code"
	@printf "%-24s %s\n" "proto-user" "Generate user proto Go/gRPC code"
	@printf "%-24s %s\n" "proto-auth" "Generate auth proto Go/gRPC code"
	@printf "%-24s %s\n" "test" "Run buyer-api, user-service and generated code tests"
	@printf "%-24s %s\n" "test-user-service" "Run user-service package tests"
	@printf "%-24s %s\n" "test-buyer-api" "Run buyer-api package tests"
	@printf "%-24s %s\n" "mysql-up" "Install or upgrade local MySQL Helm release"
	@printf "%-24s %s\n" "mysql-status" "Show local MySQL service status"
	@printf "%-24s %s\n" "mysql-down" "Uninstall local MySQL Helm release"
	@printf "%-24s %s\n" "run-user-service" "Run user-service locally"
	@printf "%-24s %s\n" "run-buyer-api" "Run buyer-api locally"
	@printf "%-24s %s\n" "run-auth-service" "Run auth-service locally"

proto: proto-user proto-auth

proto-user:
	@mkdir -p gen/go
	@$(PROTOC) \
		--go_out=paths=source_relative:./gen/go \
		--go-grpc_out=paths=source_relative:./gen/go \
		./proto/user/v1/user.proto

proto-auth:
	@mkdir -p gen/go
	@$(PROTOC) \
		--go_out=paths=source_relative:./gen/go \
		--go-grpc_out=paths=source_relative:./gen/go \
		./proto/auth/v1/auth.proto

test:
	@go test ./backend/services/user-service/... ./backend/apis/buyer-api/... ./gen/go/proto/...

test-user-service:
	@go test ./backend/services/user-service/...

test-buyer-api:
	@go test ./backend/apis/buyer-api/...

mysql-up:
	@helm upgrade --install local-mysql ./deploy/helm/local-mysql --namespace free-ecommerce-local --create-namespace

mysql-status:
	@kubectl -n free-ecommerce-local get svc,pods,pvc

mysql-down:
	@helm uninstall local-mysql --namespace free-ecommerce-local

run-user-service:
	@USER_SERVICE_HTTP_ADDR=127.0.0.1:18082 USER_SERVICE_GRPC_ADDR=127.0.0.1:19082 go run ./backend/services/user-service/cmd/user-service

run-buyer-api:
	@BUYER_API_HTTP_ADDR=127.0.0.1:18080 BUYER_API_USER_SERVICE_GRPC_ADDR=127.0.0.1:19082 go run ./backend/apis/buyer-api/cmd/buyer-api

run-auth-service:
	@AUTH_SERVICE_HTTP_ADDR=127.0.0.1:18081 go run ./backend/services/auth-service/cmd/auth-service
