package main

import (
	"database/sql"
	"log"
	"net"

	"55.GO/internal/handler"
	pb "55.GO/genproto/tutorial"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/mattn/go-sqlite3"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	createTableQuery := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER,
		email TEXT,
		address TEXT,
		phone_numbers TEXT,
		occupation TEXT,
		company TEXT,
		is_active BOOLEAN
	)`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &handler.Server{DB: db})

	log.Printf("Server is running on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

