package post_manager_events

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

type Event struct {
	EventType om.PostManagerEventTypeEnum
	Username  string
	Link      *om.Link
}
