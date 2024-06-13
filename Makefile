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


.PHONY: ent
ent:
	go generate ./infra/persistence/dal/db/ent && \
	go run entgo.io/ent/cmd/ent generate \
		--feature sql/execquery \
		--template ./infra/persistence/dal/db/extension/data_upload.tmpl ./infra/persistence/dal/db/ent/schema;