GOPATH_BIN := $(shell go env GOPATH)/bin
PROTOC := PATH="$(GOPATH_BIN):$$PATH" protoc --proto_path=.

.PHONY: help proto proto-user proto-auth proto-product proto-cart proto-order proto-payment test test-user-service test-auth-service test-product-service test-cart-service test-order-service test-payment-service test-buyer-api mysql-up mysql-status mysql-down run-user-service run-buyer-api run-auth-service run-product-service run-cart-service run-order-service run-payment-service

help:
	@printf "%-24s %s\n" "proto" "Generate all Go/gRPC proto code"
	@printf "%-24s %s\n" "proto-user" "Generate user proto Go/gRPC code"
	@printf "%-24s %s\n" "proto-auth" "Generate auth proto Go/gRPC code"
	@printf "%-24s %s\n" "proto-product" "Generate product proto Go/gRPC code"
	@printf "%-24s %s\n" "proto-cart" "Generate cart proto Go/gRPC code"
	@printf "%-24s %s\n" "proto-order" "Generate order proto Go/gRPC code"
	@printf "%-24s %s\n" "proto-payment" "Generate payment proto Go/gRPC code"
	@printf "%-24s %s\n" "test" "Run auth-service, buyer-api, user-service, product-service, cart-service, order-service, payment-service and generated code tests"
	@printf "%-24s %s\n" "test-user-service" "Run user-service package tests"
	@printf "%-24s %s\n" "test-auth-service" "Run auth-service package tests"
	@printf "%-24s %s\n" "test-product-service" "Run product-service package tests"
	@printf "%-24s %s\n" "test-cart-service" "Run cart-service package tests"
	@printf "%-24s %s\n" "test-order-service" "Run order-service package tests"
	@printf "%-24s %s\n" "test-payment-service" "Run payment-service package tests"
	@printf "%-24s %s\n" "test-buyer-api" "Run buyer-api package tests"
	@printf "%-24s %s\n" "mysql-up" "Install or upgrade local MySQL Helm release"
	@printf "%-24s %s\n" "mysql-status" "Show local MySQL service status"
	@printf "%-24s %s\n" "mysql-down" "Uninstall local MySQL Helm release"
	@printf "%-24s %s\n" "run-user-service" "Run user-service locally"
	@printf "%-24s %s\n" "run-buyer-api" "Run buyer-api locally"
	@printf "%-24s %s\n" "run-auth-service" "Run auth-service locally"
	@printf "%-24s %s\n" "run-product-service" "Run product-service locally"
	@printf "%-24s %s\n" "run-cart-service" "Run cart-service locally"
	@printf "%-24s %s\n" "run-order-service" "Run order-service locally"
	@printf "%-24s %s\n" "run-payment-service" "Run payment-service locally"

proto: proto-user proto-auth proto-product proto-cart proto-order proto-payment

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

proto-product:
	@mkdir -p gen/go
	@$(PROTOC) \
		--go_out=paths=source_relative:./gen/go \
		--go-grpc_out=paths=source_relative:./gen/go \
		./proto/product/v1/product.proto

proto-cart:
	@mkdir -p gen/go
	@$(PROTOC) \
		--go_out=paths=source_relative:./gen/go \
		--go-grpc_out=paths=source_relative:./gen/go \
		./proto/cart/v1/cart.proto

proto-order:
	@mkdir -p gen/go
	@$(PROTOC) \
		--go_out=paths=source_relative:./gen/go \
		--go-grpc_out=paths=source_relative:./gen/go \
		./proto/order/v1/order.proto

proto-payment:
	@mkdir -p gen/go
	@$(PROTOC) \
		--go_out=paths=source_relative:./gen/go \
		--go-grpc_out=paths=source_relative:./gen/go \
		./proto/payment/v1/payment.proto

test:
	@go test ./backend/services/user-service/... ./backend/services/auth-service/... ./backend/services/product-service/... ./backend/services/cart-service/... ./backend/services/order-service/... ./backend/services/payment-service/... ./backend/apis/buyer-api/... ./gen/go/proto/...

test-user-service:
	@go test ./backend/services/user-service/...

test-auth-service:
	@go test ./backend/services/auth-service/...

test-product-service:
	@go test ./backend/services/product-service/...

test-cart-service:
	@go test ./backend/services/cart-service/...

test-order-service:
	@go test ./backend/services/order-service/...

test-payment-service:
	@go test ./backend/services/payment-service/...

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
	@BUYER_API_HTTP_ADDR=127.0.0.1:18080 BUYER_API_USER_SERVICE_GRPC_ADDR=127.0.0.1:19082 BUYER_API_AUTH_SERVICE_GRPC_ADDR=127.0.0.1:19081 BUYER_API_PRODUCT_SERVICE_GRPC_ADDR=127.0.0.1:19083 BUYER_API_CART_SERVICE_GRPC_ADDR=127.0.0.1:19084 BUYER_API_ORDER_SERVICE_GRPC_ADDR=127.0.0.1:19085 BUYER_API_PAYMENT_SERVICE_GRPC_ADDR=127.0.0.1:19086 go run ./backend/apis/buyer-api/cmd/buyer-api

run-auth-service:
	@AUTH_SERVICE_HTTP_ADDR=127.0.0.1:18081 AUTH_SERVICE_GRPC_ADDR=127.0.0.1:19081 go run ./backend/services/auth-service/cmd/auth-service

run-product-service:
	@PRODUCT_SERVICE_HTTP_ADDR=127.0.0.1:18083 PRODUCT_SERVICE_GRPC_ADDR=127.0.0.1:19083 go run ./backend/services/product-service/cmd/product-service

run-cart-service:
	@CART_SERVICE_HTTP_ADDR=127.0.0.1:18084 CART_SERVICE_GRPC_ADDR=127.0.0.1:19084 CART_SERVICE_PRODUCT_SERVICE_GRPC_ADDR=127.0.0.1:19083 go run ./backend/services/cart-service/cmd/cart-service

run-order-service:
	@ORDER_SERVICE_HTTP_ADDR=127.0.0.1:18085 ORDER_SERVICE_GRPC_ADDR=127.0.0.1:19085 ORDER_SERVICE_USER_SERVICE_GRPC_ADDR=127.0.0.1:19082 ORDER_SERVICE_CART_SERVICE_GRPC_ADDR=127.0.0.1:19084 go run ./backend/services/order-service/cmd/order-service

run-payment-service:
	@PAYMENT_SERVICE_HTTP_ADDR=127.0.0.1:18086 PAYMENT_SERVICE_GRPC_ADDR=127.0.0.1:19086 PAYMENT_SERVICE_ORDER_SERVICE_GRPC_ADDR=127.0.0.1:19085 go run ./backend/services/payment-service/cmd/payment-service
