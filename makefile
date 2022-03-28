generateSwaggerDoc:
	swagger generate spec -o ./swagger.yaml

runHTTPServer:
	go run ./server/http/server.go

runGrpcServer:
	go run ./server/grpc/server.go

buildProto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./interfaceAdapters/grpc/protocol/url-service.proto

buildHTTPServer:
	go build -o ./build/http ./server/http/server.go

buildGrpcServer:
	go build -o ./build/grpc ./server/grpc/server.go