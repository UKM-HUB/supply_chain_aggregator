.PHONY: proto

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/auth/*.proto proto/inventory/*.proto proto/order/*.proto

run-gateway:
	go run cmd/api-gateway/main.go
