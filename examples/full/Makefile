bad:
	protoc -I api/grpc \
               -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
               -I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
               --go_out=internal/pkg/api \
               --go-grpc_out=internal/pkg/api \
               --grpc-gateway_out=internal/pkg/api \
               --grpc-gateway-ts_out=internal/pkg/web/@vervstack/matreshka \
               api/grpc/*.proto

get-clients-good:
	protoc  -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc \
           	-I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
           	-I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
           	--go_out=paths=source_relative:internal/pkg/api/clients/ \
           	--go-grpc_out=paths=source_relative:internal/pkg/api/clients/ \
           	proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc/*.proto

modules:
	protoc -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc \
		   -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
		   -I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
		   --go_out=module=go.vervstack.ru:internal/api/clients \
		   --go-grpc_out=module=go.vervstack.ru:internal/api/clients \
		   proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc/*.proto


new-qm:
	protoc -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc \
		 	-I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
		   	-I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
		   	--go_out=module=go.vervstack.ru:internal/api/clients \
		   	--go-grpc_out=module=go.vervstack.ru:internal/api/clients \
		   	proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc/*.proto

old:
	protoc -I . \
		   -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
		   -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95 \
		   -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc \
		   -I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
		   --go_out=paths=source_relative,module=go.vervstack.ru:internal/api/clients \
		   --go-grpc_out=module=go.vervstack.ru:internal/api/clients \
		   proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc/*.proto
# 1. Удалён -I .
# 2. Удалён экстра импорт без пути

q:
	protoc -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
		   -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc \
		   -I proto_modules/mod/go.redsock.ru/protoc-gen-npm/v0.0.17 \
		   --go_out=module=go.vervstack.ru:internal/api/clients \
		   --go-grpc_out=module=go.vervstack.ru:internal/api/clients \
		   proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95/api/grpc/*.proto