build:
	protoc -I api/grpc \
               --go_out=pkg/api \
               --go-grpc_out=pkg/api \
               api/grpc/basic_service.proto