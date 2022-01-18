package server

import (
	"context"
	"fmt"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/blogpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BlogServer struct {
	blogpb.UnimplementedBlogServiceServer
}

func (*BlogServer) FindAll(ctx context.Context, req *blogpb.FindAllRequest) (*blogpb.FindAllResponse, error) {
	var posts []model.Post
	if err := postgresdb.DB.Find(&posts).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	var PostsRes = []*blogpb.Post{}
	for _, v := range posts {
		PostsRes = append(PostsRes, &blogpb.Post{
			Id:       int32(v.Model.ID),
			ImageUrl: v.ImageUrl,
			Title:    v.Title,
			Abstract: v.Abstract,
			Body:     v.Body,
		})
	}
	return &blogpb.FindAllResponse{
		Posts: PostsRes,
	}, nil
}

func (*BlogServer) Find(ctx context.Context, req *blogpb.FindRequest) (*blogpb.FindResponse, error) {
	var post model.Post
	if err := postgresdb.DB.Where("id = ?", req.GetPostId()).Find(&post).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	if post.ID == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("پست یافت نشد"),
		)
	}
	return &blogpb.FindResponse{
		Post: &blogpb.Post{
			Id:        int32(post.Model.ID),
			ImageUrl: post.ImageUrl,
			Title:    post.Title,
			Abstract: post.Abstract,
			Body:     post.Body,
		},
	}, nil
}
