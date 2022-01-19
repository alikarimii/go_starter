GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)
GO_MODULE := $(shell go mod edit -json | grep Path | head -n 1 | cut -d ":" -f 2 | cut -d '"' -f 2)
PROTO_DIR := src/customeraccounts/infrastructure/adapter/grpc/proto
REST_SWAGGER_TARGET_DIR := src/customeraccounts/infrastructure/adapter/rest
gen:
	@protoc \
		-I $(PROTO_DIR) \
		-I /usr/local/include \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=paths=source_relative:$(PROTO_DIR) \
		--go-grpc_out=require_unimplemented_servers=false,paths=source_relative:$(PROTO_DIR) \
		--grpc-gateway_out=paths=source_relative,generate_unbound_methods=true,logtostderr=true:$(PROTO_DIR) \
		--swagger_out=logtostderr=true:$(REST_SWAGGER_TARGET_DIR) \
		$(PROTO_DIR)/customer.proto
run:
	go run src/services/grpc/cmd/main.go

test: