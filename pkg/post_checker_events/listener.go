package post_checker_events

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

func Listen(url string, sink om.PostCheckerEvents) (err error) {
	conn, err := connect(url)
	if err != nil {
		return
	}

	conn.QueueSubscribe(subject, queue, func(e *Event) {
		sink.OnPostChecked(e.Username, e.Url, e.Status)
	})

	return
}
