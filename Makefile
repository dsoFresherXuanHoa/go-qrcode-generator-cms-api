dev:
	air -d

swag:
	swag init --parseDependency

grpc:
	protoc "./proto/qrcode.proto" --go-grpc_out="./proto/gRPC" --go_out="./proto/gRPC"

evans:
	evans --proto "./proto/qrcode.proto" repl --host localhost --port 5000