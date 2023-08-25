package main

import (
	"encoding/json"

	"github.com/muke-sh/grpc-curd/book"
)

type Book struct {
	Id     string `bson:"_id" json:"id"`
	Title  string `bson:"title" json:"title"`
	Author string `bson:"author" json:"author"`
	ISBN   string `bson:"isbn" json:"isbn"`
	Price  int64  `bson:"price" json:"price"`
}

func (b *Book) toRPCResponse() (*book.Book, error) {

	jsonBytes, err := json.Marshal(b)

	if err != nil {
		return nil, err
	}

	rpcResponse := book.Book{}

	err = json.Unmarshal(jsonBytes, &rpcResponse)

	if err != nil {
		return nil, err
	}

	return &rpcResponse, nil
}
