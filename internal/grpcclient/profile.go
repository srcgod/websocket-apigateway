package grpcclient

import (
	profilev1 "github.com/srcgod/profileproto/gen/profile"
	"google.golang.org/grpc"
)

func NewProfileClient(address string) (profilev1.ProfileServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := profilev1.NewProfileServiceClient(conn)
	return client, nil
}
