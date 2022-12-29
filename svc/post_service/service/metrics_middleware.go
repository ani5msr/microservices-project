package service

import (
	"strings"
	"time"

	"github.com/ani5msr/microservices-project/pkg/metrics"
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/prometheus/client_golang/prometheus"
)

// implement function to return ServiceMiddleware
func newMetricsMiddleware() postManagerMiddleware {
	return func(next om.PostManager) om.PostManager {
		m := metricsMiddleware{next,
			map[string]prometheus.Counter{},
			map[string]prometheus.Summary{}}
		methodNames := []string{"GetLinks", "AddLink", "UpdateLink", "DeleteLink"}
		for _, name := range methodNames {
			m.requestCounter[name] = metrics.NewCounter("link", strings.ToLower(name)+"_count", "count # of requests")
			m.requestLatency[name] = metrics.NewSummary("link", strings.ToLower(name)+"_summary", "request summary in milliseconds")

		}
		return m
	}
}

type metricsMiddleware struct {
	next           om.PostManager
	requestCounter map[string]prometheus.Counter
	requestLatency map[string]prometheus.Summary
}

func (m metricsMiddleware) recordMetrics(name string, begin time.Time) {
	m.requestCounter[name].Inc()
	durationMilliseconds := float64(time.Since(begin).Nanoseconds() * 1000000)
	m.requestLatency[name].Observe(durationMilliseconds)
}

func (m metricsMiddleware) GetPost(request om.GetPostsRequest) (result om.GetPostsResult, err error) {
	defer func(begin time.Time) {
		m.recordMetrics("GetPost", begin)
	}(time.Now())
	result, err = m.next.GetPost(request)
	return
}

func (m metricsMiddleware) AddPost(request om.AddPostRequest) error {
	defer func(begin time.Time) {
		m.recordMetrics("AddPost", begin)
	}(time.Now())
	return m.next.AddPost(request)
}

func (m metricsMiddleware) UpdatePost(request om.UpdatePostRequest) error {
	defer func(begin time.Time) {
		m.recordMetrics("UpdatePost", begin)
	}(time.Now())
	return m.next.UpdatePost(request)
}

func (m metricsMiddleware) DeletePost(username string, url string) error {
	defer func(begin time.Time) {
		m.recordMetrics("DeletePost", begin)
	}(time.Now())
	return m.next.DeletePost(username, url)
}
