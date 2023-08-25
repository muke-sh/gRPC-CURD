package main

import (
	"context"

	"log"

	"github.com/muke-sh/grpc-curd/book"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	book.UnimplementedBookServiceServer
}

func (s *Server) CreateBook(ctx context.Context, cbr *book.CreateBookRequest) (*book.Reponse, error) {

	bookCollection, err := GetCollection(MongoURI, "grpc", "books")

	if err != nil {
		log.Printf("error while getting a collection handle: %s", err)
		return nil, status.Error(codes.Internal, "internal server error occured")
	}

	res, err := bookCollection.InsertOne(ctx, cbr.Book)

	if err != nil {
		log.Printf("error while inserting a book: %s", err)
		return nil, status.Error(codes.Internal, "internal server error occured")
	}

	log.Printf("Book created, Id: %s", res.InsertedID)

	return &book.Reponse{
		Msg:    "Book created",
		Status: "OK",
	}, status.Error(codes.OK, "internal error occured")
}

func (s *Server) FindBook(ctx context.Context, fbr *book.FindBookRequest) (*book.Book, error) {

	oid, err := primitive.ObjectIDFromHex(fbr.Id)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid book id")
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: oid},
	}

	bookCollection, err := GetCollection(MongoURI, "grpc", "books")

	if err != nil {
		log.Printf("error while getting a collection handle: %s", err)
		return nil, status.Error(codes.Internal, "internal server error occured")
	}
	book := Book{}

	err = bookCollection.FindOne(ctx, filter).Decode(&book)

	if err != nil {
		log.Printf("error while reading database: %s", err)
		return nil, status.Error(codes.Internal, "internal server error occured")
	}

	rpcResponse, err := book.toRPCResponse()

	if err != nil {
		log.Printf("error while converting bson to rpcResponse: %s", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return rpcResponse, nil
}

func (s *Server) DeleteBook(ctx context.Context, fbr *book.FindBookRequest) (*book.Reponse, error) {

	oid, err := primitive.ObjectIDFromHex(fbr.Id)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid book id")
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: oid},
	}

	bookCollection, err := GetCollection(MongoURI, "grpc", "books")

	if err != nil {
		log.Printf("error while getting a collection handle: %s", err)
		return nil, status.Error(codes.Internal, "internal server error occured")
	}

	res, err := bookCollection.DeleteOne(ctx, filter)

	if err != nil {
		log.Printf("error while deleting book: %s", err)
		return nil, status.Error(codes.Internal, "internal server error occured")
	}

	if res.DeletedCount == 0 {
		return nil, status.Error(codes.NotFound, "book does not exists or already deleted")
	}

	return &book.Reponse{
		Msg:    "Book deleted successfully",
		Status: "OK",
	}, nil
}

func (s *Server) UpdateBook(ctx context.Context, ubr *book.UpdateBookRequest) (*book.Book, error) {

	oid, err := primitive.ObjectIDFromHex(ubr.Id)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid book id")
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: oid},
	}

	update := bson.D{
		primitive.E{Key: "$set", Value: ubr.Book},
	}

	bookCollection, err := GetCollection(MongoURI, "grpc", "books")

	if err != nil {
		log.Printf("error while getting a collection handle: %s", err)
		return nil, status.Error(codes.Internal, "internal server error occured")
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	res := bookCollection.FindOneAndUpdate(ctx, filter, update, opts)

	if err != nil {
		log.Printf("error while deleting book: %s", err)
		return nil, status.Error(codes.Internal, "internal server error occured")
	}

	if res.Err() == mongo.ErrNoDocuments {
		log.Printf("no document found with id: %s, error: %s", ubr.Id, err)
		return nil, status.Error(codes.NotFound, "no documents found with provided id")
	}

	updatedBook := book.Book{}

	err = res.Decode(&updatedBook)

	if err != nil {
		log.Printf("error while decoding update response: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error occured")
	}

	return &updatedBook, nil
}
