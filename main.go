package main

import (
	"fmt"
	"log"
	"net"

	"github.com/muke-sh/grpc-curd/book"
	"google.golang.org/grpc"
)

func main() {
	port := 9000
	address := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("cannot create listener on port: %d", port)
	}

	_, err = connectDB(MongoURI)

	if err != nil {
		log.Fatalf("unable to connect to mongodb instance: %s", MongoURI)
	}

	grpcServer := grpc.NewServer()
	book.RegisterBookServiceServer(grpcServer, &Server{})

	log.Printf("Starting a gRPC server on port: %d\n", port)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("cannot start a gRPC server on port: %d", port)
	}

}
