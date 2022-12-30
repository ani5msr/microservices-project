package post_manager_events

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

func Listen(url string, sink om.PostManagerEvents) (err error) {
	conn, err := connect(url)
	if err != nil {
		return
	}

	conn.QueueSubscribe(subject, queue, func(e *Event) {
		switch e.EventType {
		case om.PostAdded:
			{
				sink.OnPostAdded(e.Username, e.Post)
			}
		case om.LinkUpdated:
			{
				sink.OnPostUpdated(e.Username, e.Post)
			}
		default:
			// Ignore other event types
		}
	})

	return
}
