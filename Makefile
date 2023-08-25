generate_grpc:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	book/book.proto

all: generate_grpc build run
	
build:
	go build .

run:
	./grpc-curd