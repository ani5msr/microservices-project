package feed_manager

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

type Store interface {
	GetFeed(username string, startIndex int) (events []*om.PostManagerEvent, nextIndex int, err error)
	AddEvent(username string, event *om.PostManagerEvent) (err error)
}
