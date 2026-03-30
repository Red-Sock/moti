bad:
	protoc -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc \
               -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
               -I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
               --go_out=internal/api/clients \
               --go-grpc_out=internal/api/clients \
               matreshka_api.proto \
               matreshka_common.proto

# TODO THIS IS THE ONE WE WANT
# 1. Main repo
# 2. all other deps
# 3. plugins with ops
# 4. local repo's path to proto files (no wildcard)
best:
	protoc -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95 \
               -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
               -I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
               --go_out=module=go.vervstack.ru/matreshka/pkg:internal/api/clients \
               --go-grpc_out=module=go.vervstack.ru/matreshka/pkg:internal/api/clients \
               api/grpc/matreshka_api.proto \
               api/grpc/matreshka_common.proto
differnt-root-best:
	protoc -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc/ \
               -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
               -I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
               --go_out=module=go.vervstack.ru/matreshka/pkg:internal/api/clients \
               --go-grpc_out=module=go.vervstack.ru/matreshka/pkg:internal/api/clients \
               matreshka_api.proto \
               matreshka_common.proto

generated:
	protoc -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95 \
               -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
               -I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
               --go_out=module=go.vervstack.ru/matreshka/pkg:internal/api/clients \
               --go-grpc_out=module=go.vervstack.ru/matreshka/pkg:internal/api/clients \
               api/grpc/matreshka_api.proto \
               api/grpc/matreshka_common.proto \
               api/old/matreshka_api.proto \
               api/old/matreshka_common.proto
