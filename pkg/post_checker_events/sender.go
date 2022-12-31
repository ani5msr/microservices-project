package post_checker_events

import (
	"log"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/nats-io/nats.go"
)

type eventSender struct {
	hostname string
	nats     *nats.EncodedConn
}

func (s *eventSender) OnPostChecked(username string, url string, status om.PostStatus) {
	err := s.nats.Publish(subject, Event{username, url, status})
	if err != nil {
		log.Fatal(err)
	}
}

func NewEventSender(url string) (om.PostCheckerEvents, error) {
	ec, err := connect(url)
	if err != nil {
		return nil, err
	}
	return &eventSender{hostname: url, nats: ec}, nil
}
