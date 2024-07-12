API_PROTO_FILES=$(shell find idl -name *.proto)

.PHONY: pb
# generate grpc code
pb:
	protoc --proto_path=. \
           --proto_path=./third_party \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           --grpc-gateway_out=paths=source_relative:. \
           $(API_PROTO_FILES)

# sql/execquery 执行raw sql, sql/lock 行锁
.PHONY: ent
ent:
	go generate ./infra/persistence/dal/db/ent && \
	go run -mod=mod entgo.io/ent/cmd/ent generate \
		--feature sql/execquery,sql/lock,sql/upsert \
		--template ./infra/persistence/dal/db/extension/data_upload.tmpl \
		./infra/persistence/dal/db/ent/schema;

lint:
	golangci-lint run