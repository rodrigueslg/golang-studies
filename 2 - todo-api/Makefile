PROTO_DIR = ./api/proto

setup:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

build:
	protoc -I ${PROTO_DIR} \
	--go_out ${PROTO_DIR} \
	--go_opt paths=source_relative \
	--go-grpc_out ${PROTO_DIR} \
	--go-grpc_opt paths=source_relative \
	--grpc-gateway_out ${PROTO_DIR} \
	--grpc-gateway_opt paths=source_relative \
	${PROTO_DIR}/*.proto \
	--experimental_allow_proto3_optional

clean:
	rm -f -r ${PROTO_DIR}/*.pb.go
	rm -f -r ${PROTO_DIR}/*.pb.gw.go
