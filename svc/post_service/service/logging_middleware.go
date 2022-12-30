package service

import (
	"time"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/go-kit/log"
)

// implement function to return ServiceMiddleware
func newLoggingMiddleware(logger log.Logger) postManagerMiddleware {
	return func(next om.PostManager) om.PostManager {
		return loggingMiddleware{next, logger}
	}
}

type loggingMiddleware struct {
	next   om.PostManager
	logger log.Logger
}

func (m loggingMiddleware) GetPost(request om.GetPostRequest) (result om.GetPostResult, err error) {
	defer func(begin time.Time) {
		m.logger.Log(
			"method", "GetPost",
			"request", request,
			"result", result,
			"duration", time.Since(begin),
		)
	}(time.Now())
	result, err = m.next.GetPost(request)
	return
}

func (m loggingMiddleware) AddPost(request om.AddPostRequest) error {
	return m.next.AddPost(request)
}

func (m loggingMiddleware) UpdatePost(request om.UpdatePostRequest) error {
	return m.next.UpdatePost(request)
}

func (m loggingMiddleware) DeletePost(username string, url string) error {
	return m.next.DeletePost(username, url)
}
