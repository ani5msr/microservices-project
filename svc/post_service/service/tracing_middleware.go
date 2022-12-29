package service

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/opentracing/opentracing-go"
)

// implement function to return ServiceMiddleware
func newTracingMiddleware(tracer opentracing.Tracer) linkManagerMiddleware {
	return func(next om.PostManager) om.PostManager {
		return tracingMiddleware{next, tracer}
	}
}

type tracingMiddleware struct {
	next   om.PostManager
	tracer opentracing.Tracer
}

func (m tracingMiddleware) GetPost(request om.GetPostsRequest) (result om.GetPostsResult, err error) {
	defer func(span opentracing.Span) {
		span.Finish()
	}(m.tracer.StartSpan("GetLinks"))
	result, err = m.next.GetPost(request)
	return
}

func (m tracingMiddleware) AddPost(request om.AddPostRequest) error {
	return m.next.AddPost(request)
}

func (m tracingMiddleware) UpdatePost(request om.UpdatePostRequest) error {
	return m.next.UpdatePost(request)
}

func (m tracingMiddleware) DeletePost(username string, url string) error {
	return m.next.DeletePost(username, url)
}
