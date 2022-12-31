package post_manager_client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

type deletePostRequest struct {
	Username string
	Url      string
}

type SimpleResponse struct {
	Err string
}

func decodeSimpleResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp SimpleResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

type EndpointSet struct {
	GetPostEndpoint    endpoint.Endpoint
	AddPostEndpoint    endpoint.Endpoint
	UpdatePostEndpoint endpoint.Endpoint
	DeletePostEndpoint endpoint.Endpoint
}

func decodeGetPostResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var res om.GetPostResult
	err := json.NewDecoder(r.Body).Decode(&res)
	return res, err
}

func (s EndpointSet) GetPost(req om.GetPostRequest) (result om.GetPostResult, err error) {
	res, err := s.GetPostEndpoint(context.Background(), req)
	if err != nil {
		return
	}
	result = res.(om.GetPostResult)

	return
}

func (s EndpointSet) AddPost(req om.AddPostRequest) (err error) {
	resp, err := s.AddPostEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)

	if response.Err != "" {
		err = errors.New(response.Err)
	}
	return
}

func (s EndpointSet) UpdatePost(req om.UpdatePostRequest) (err error) {
	resp, err := s.UpdatePostEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)

	if response.Err != "" {
		err = errors.New(response.Err)
	}
	return
}

func (s EndpointSet) DeletePost(username string, url string) (err error) {
	resp, err := s.DeletePostEndpoint(context.Background(), &deletePostRequest{Username: username, Url: url})
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)

	if response.Err != "" {
		err = errors.New(response.Err)
	}
	return
}
