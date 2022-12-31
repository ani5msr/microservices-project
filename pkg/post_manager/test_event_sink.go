package post_manager

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

type testEventsSink struct {
	addPostEvents     map[string][]*om.Post
	updatePostEvents  map[string][]*om.Post
	deletedPostEvents map[string][]string
}

func (s *testEventsSink) OnPostAdded(username string, post *om.Post) {
	if s.addPostEvents[username] == nil {
		s.addPostEvents[username] = []*om.Post{}
	}
	s.addPostEvents[username] = append(s.addPostEvents[username], post)
}

func (s *testEventsSink) OnPostUpdated(username string, link *om.Post) {
	if s.updatePostEvents[username] == nil {
		s.updatePostEvents[username] = []*om.Post{}
	}
	s.updatePostEvents[username] = append(s.updatePostEvents[username], link)
}

func (s *testEventsSink) OnPostDeleted(username string, url string) {
	if s.deletedPostEvents[username] == nil {
		s.deletedPostEvents[username] = []string{}
	}
	s.deletedPostEvents[username] = append(s.deletedPostEvents[username], url)
}

func newPostManagerEventsSink() *testEventsSink {
	return &testEventsSink{
		map[string][]*om.Post{},
		map[string][]*om.Post{},
		map[string][]string{},
	}
}
