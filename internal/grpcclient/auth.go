package grpcclient

import (
	authv1 "github.com/srcgod/authproto/gen/user"
	"google.golang.org/grpc"
)

func NewAuthClient(address string) (authv1.AuthClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := authv1.NewAuthClient(conn)
	return client, nil
}
