package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/go-kit/kit/endpoint"
)

type post struct {
	Url         string
	Title       string
	Description string
	Status      string
	Tags        map[string]bool
	CreatedAt   string
	UpdatedAt   string
}

func newPost(source om.Post) post {
	return post{
		Url:         source.Url,
		Title:       source.Title,
		Description: source.Description,
		Status:      source.Status,
		Tags:        source.Tags,
		CreatedAt:   source.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   source.UpdatedAt.Format(time.RFC3339),
	}
}

type getPostResponse struct {
	Posts []post `json:"posts"`
	Err   string `json:"err"`
}

type deletePostRequest struct {
	Username string
	Url      string
}

type SimpleResponse struct {
	Err string `json:"err"`
}

func decodeGetPostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request om.GetPostRequest
	q := r.URL.Query()
	request.UrlRegex = q.Get("url")
	request.TitleRegex = q.Get("title")
	request.DescriptionRegex = q.Get("description")
	request.Username = q.Get("username")
	request.Tag = q.Get("tag")
	request.StartToken = q.Get("start")
	return request, nil
}

func decodeAddPostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request om.AddPostRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdatePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request om.UpdatePostRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeDeletePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request deletePostRequest
	q := r.URL.Query()
	request.Username = q.Get("username")
	request.Url = q.Get("url")
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func makeGetPostEndpoint(svc om.PostManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.GetPostRequest)
		result, err := svc.GetPost(req)
		res := getPostResponse{}
		for _, post := range result.Posts {
			res.Posts = append(res.Posts, newPost(post))
		}
		if err != nil {
			res.Err = err.Error()
			return res, err
		}
		return res, nil
	}
}

func makeAddPostEndpoint(svc om.PostManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.AddPostRequest)
		err := svc.AddPost(req)
		res := SimpleResponse{}
		if err != nil {
			res.Err = err.Error()
			return res, err
		}
		return res, nil
	}
}

func makeUpdatePostEndpoint(svc om.PostManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.UpdatePostRequest)
		err := svc.UpdatePost(req)
		res := SimpleResponse{}
		if err != nil {
			res.Err = err.Error()
			return res, err
		}
		return res, nil
	}
}

func makeDeletePostEndpoint(svc om.PostManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deletePostRequest)
		err := svc.DeletePost(req.Username, req.Url)
		res := SimpleResponse{}
		if err != nil {
			res.Err = err.Error()
			return res, err
		}
		return res, nil
	}
}
