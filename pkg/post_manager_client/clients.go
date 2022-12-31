package post_manager_client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	om "github.com/ani5msr/microservices-project/pkg/object_model"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewClient(baseURL string) (om.PostManager, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	getPostEndpoint := httptransport.NewClient(
		"GET",
		copyURL(u, "/posts"),
		encodeGetPostRequest,
		decodeGetPostResponse).Endpoint()

	addPostEndpoint := httptransport.NewClient(
		"POST",
		copyURL(u, "/posts"),
		encodeHTTPGenericRequest,
		decodeSimpleResponse).Endpoint()

	updatePostEndpoint := httptransport.NewClient(
		"PUT",
		copyURL(u, "/posts"),
		encodeHTTPGenericRequest,
		decodeSimpleResponse).Endpoint()

	deletePostEndpoint := httptransport.NewClient(
		"DELETE",
		copyURL(u, "/posts"),
		encodeDeletePostRequest,
		decodeSimpleResponse).Endpoint()

	// Returning the EndpointSet as an interface relies on the
	// EndpointSet implementing the Service methods. That's just a simple bit
	// of glue code.
	return EndpointSet{
		GetPostEndpoint:    getPostEndpoint,
		AddPostEndpoint:    addPostEndpoint,
		UpdatePostEndpoint: updatePostEndpoint,
		DeletePostEndpoint: deletePostEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// Extract the request details from the incoming request and add them as query arguments
func encodeGetPostRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(om.GetPostRequest)
	urlRegex := url.QueryEscape(r.UrlRegex)
	titleRegex := url.QueryEscape(r.TitleRegex)
	descriptionRegex := url.QueryEscape(r.DescriptionRegex)
	username := url.QueryEscape(r.Username)
	tag := url.QueryEscape(r.Tag)
	startToken := url.QueryEscape(r.StartToken)

	q := req.URL.Query()
	q.Add("url", urlRegex)
	q.Add("title", titleRegex)
	q.Add("description", descriptionRegex)
	q.Add("username", username)
	q.Add("tag", tag)
	q.Add("start", startToken)
	req.URL.RawQuery = q.Encode()
	return encodeHTTPGenericRequest(ctx, req, request)
}

func encodeDeletePostRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(*deletePostRequest)
	q := req.URL.Query()
	q.Add("username", r.Username)
	q.Add("url", r.Url)
	req.URL.RawQuery = q.Encode()
	return encodeHTTPGenericRequest(ctx, req, request)
}
