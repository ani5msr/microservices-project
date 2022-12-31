package post_checker_events

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

type Event struct {
	Username string
	Url      string
	Status   om.PostStatus
}
