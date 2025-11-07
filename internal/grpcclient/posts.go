package grpcclient

import (
	postsv1 "github.com/srcgod/postproto/gen/posts"
	"google.golang.org/grpc"
)

func NewPostsClient(address string) (postsv1.PostServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := postsv1.NewPostServiceClient(conn)
	return client, nil
}
