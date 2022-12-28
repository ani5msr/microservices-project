package post_manager

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

type PostStore interface {
	GetPost(request om.GetPostsRequest) (om.GetPostsResult, error)
	AddPost(request om.AddPostRequest) (*om.Post, error)
	UpdatePost(request om.UpdatePostRequest) (*om.Post, error)
	DeletePost(username string, url string) error
	SetPostStatus(username, url string, status om.PostStatus) error
}
