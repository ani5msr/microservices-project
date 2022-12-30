package post_manager_events

import (
	"log"

	"github.com/nats-io/nats.go"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

type eventSender struct {
	hostname string
	nats     *nats.EncodedConn
}

func (s *eventSender) OnPostAdded(username string, link *om.Post) {
	event := Event{om.PostAdded, username, link}
	log.Printf("[post manager events]OnPostAdded(), sending to subject: %s, event: %v\n", event)
	err := s.nats.Publish(subject, event)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *eventSender) OnPostUpdated(username string, link *om.Post) {
	err := s.nats.Publish(subject, Event{om.PostUpdated, username, link})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *eventSender) OnPostDeleted(username string, url string) {
	// Ignore link delete events
}

func NewEventSender(url string) (om.PostManagerEvents, error) {
	ec, err := connect(url)
	if err != nil {
		return nil, err
	}
	return &eventSender{hostname: url, nats: ec}, nil
}
