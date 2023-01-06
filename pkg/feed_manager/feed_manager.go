package feed_manager

import (
	"errors"
	"strconv"
	"time"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/ani5msr/microservices-project/pkg/post_manager_events"
)

type FeedManager struct {
	feedStore Store
}

func (m *FeedManager) GetFeed(req om.GetFeedRequest) (resp om.GetFeedResult, err error) {
	if req.Username == "" {
		err = errors.New("user name can't be empty")
		return
	}

	startIndex := 0
	if req.StartToken != "" {
		startIndex, err := strconv.Atoi(req.StartToken)
		if err != nil || startIndex < 0 {
			err = errors.New("invalid start token: " + req.StartToken)
			return resp, err
		}
	}

	events, nextIndex, err := m.feedStore.GetFeed(req.Username, startIndex)
	if err != nil {
		return
	}

	resp.Events = events
	if nextIndex != -1 {
		resp.NextToken = strconv.Itoa(nextIndex)
	}

	return
}

func (m *FeedManager) OnPostAdded(username string, post *om.Post) {
	event := &om.PostManagerEvent{
		EventType: om.PostAdded,
		Username:  username,
		Url:       post.Url,
		Timestamp: time.Now().UTC(),
	}
	m.feedStore.AddEvent(username, event)
}

func (m *FeedManager) OnPostUpdated(username string, post *om.Post) {
	event := &om.PostManagerEvent{
		EventType: om.PostUpdated,
		Username:  username,
		Url:       post.Url,
		Timestamp: time.Now().UTC(),
	}
	m.feedStore.AddEvent(username, event)
}

func (m *FeedManager) OnPostDeleted(username string, url string) {
	event := &om.PostManagerEvent{
		EventType: om.PostDeleted,
		Username:  username,
		Url:       url,
		Timestamp: time.Now().UTC(),
	}
	m.feedStore.AddEvent(username, event)
}

func NewFeedManager(store Store, natsHostname string, natsPort string) (om.FeedManager, error) {
	nm := &FeedManager{feedStore: store}
	if natsHostname != "" {
		natsUrl := natsHostname + ":" + natsPort
		err := post_manager_events.Listen(natsUrl, nm)
		if err != nil {
			return nil, err
		}
	}

	return nm, nil
}
