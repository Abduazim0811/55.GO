package handler

import (
	"context"
	"database/sql"
	"encoding/json"

	pb "55.GO/genproto/tutorial"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	DB *sql.DB
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	insertQuery := `INSERT INTO users (name, age, email, address, phone_numbers, occupation, company, is_active)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	address, err := json.Marshal(req.Address)
	if err != nil {
		return nil, err
	}
	phoneNumbers, err := json.Marshal(req.PhoneNumbers)
	if err != nil {
		return nil, err
	}
	res, err := s.DB.Exec(insertQuery, req.Name, req.Age, req.Email, address, phoneNumbers, req.Occupation, req.Company, req.IsActive)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user := &pb.User{
		Id:           int32(id),
		Name:         req.Name,
		Age:          req.Age,
		Email:        req.Email,
		Address:      req.Address,
		PhoneNumbers: req.PhoneNumbers,
		Occupation:   req.Occupation,
		Company:      req.Company,
		IsActive:     req.IsActive,
	}
	return &pb.GetUserResponse{User: user}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	selectQuery := `SELECT id, name, age, email, address, phone_numbers, occupation, company, is_active FROM users WHERE id = ?`
	row := s.DB.QueryRow(selectQuery, req.Id)
	var user pb.User
	var address, phoneNumbers []byte
	err := row.Scan(&user.Id, &user.Name, &user.Age, &user.Email, &address, &phoneNumbers, &user.Occupation, &user.Company, &user.IsActive)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(address, &user.Address)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(phoneNumbers, &user.PhoneNumbers)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{User: &user}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *empty.Empty) (*pb.ListUsersResponse, error) {
	selectQuery := `SELECT id, name, age, email, address, phone_numbers, occupation, company, is_active FROM users`
	rows, err := s.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		var address, phoneNumbers []byte
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Email, &address, &phoneNumbers, &user.Occupation, &user.Company, &user.IsActive)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(address, &user.Address)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(phoneNumbers, &user.PhoneNumbers)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return &pb.ListUsersResponse{Users: users}, nil
}
