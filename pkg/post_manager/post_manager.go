package post_manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/ani5msr/microservices-project/pkg/post_checker_events"
)

// Nuclio functions listen by default on port 8080 of their service IP
const Post_checker_func_url = "http://Post-checker.nuclio.svc.cluster.local:8080"

type PostManager struct {
	postStore          PostStore
	socialGraphManager om.SocialGraphManager
	eventSink          om.PostManagerEvents
	maxPostsPerUser    int64
}

func (m *PostManager) GetPost(request om.GetPostRequest) (result om.GetPostResult, err error) {
	if request.Username == "" {
		err = errors.New("user name can't be empty")
		return
	}

	result, err = m.postStore.GetPost(request)
	if result.Posts == nil {
		result.Posts = []om.Post{}
	}
	return
}

// Very wasteful way to count Posts
func (m *PostManager) getPostCount(username string) (PostCount int64, err error) {
	req := om.GetPostRequest{Username: username}
	res, err := m.GetPost(req)
	if err != nil {
		return
	}

	PostCount += int64(len(res.Posts))

	for res.NextPageToken != "" {
		req = om.GetPostRequest{Username: username, StartToken: res.NextPageToken}
		res, err = m.GetPost(req)
		if err != nil {
			return
		}

		PostCount += int64(len(res.Posts))
	}
	return
}

func triggerPostCheck(username string, url string) {
	go func() {
		checkPostRequest := &om.CheckPostRequest{Username: username, Url: url}
		data, err := json.Marshal(checkPostRequest)
		if err != nil {
			return
		}

		req, err := http.NewRequest("POST", Post_checker_func_url, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}()
}

func (m *PostManager) AddPost(request om.AddPostRequest) (err error) {
	if request.Url == "" {
		return errors.New("the URL can't be empty")
	}

	if request.Username == "" {
		return errors.New("the user name can't be empty")
	}

	PostCount, err := m.getPostCount(request.Username)
	if err != nil {
		return
	}

	if PostCount >= m.maxPostsPerUser {
		return errors.New("the user has too many Posts")
	}

	Post, err := m.postStore.AddPost(request)
	if err != nil {
		return
	}

	if m.eventSink != nil {
		followers, err := m.socialGraphManager.GetFollowers(request.Username)
		if err != nil {
			return err
		}

		for follower := range followers {
			m.eventSink.OnPostAdded(follower, Post)
		}
	}

	// Trigger Post check asynchronously (don't wait for result)
	triggerPostCheck(request.Username, request.Url)
	return
}

func (m *PostManager) UpdatePost(request om.UpdatePostRequest) (err error) {
	if request.Url == "" {
		return errors.New("the URL can't be empty")
	}

	if request.Username == "" {
		return errors.New("the user name can't be empty")
	}

	Post, err := m.postStore.UpdatePost(request)
	if err != nil {
		return
	}

	if m.eventSink != nil {
		followers, err := m.socialGraphManager.GetFollowers(request.Username)
		if err != nil {
			return err
		}

		for follower := range followers {
			m.eventSink.OnPostUpdated(follower, Post)
		}
	}

	return
}

func (m *PostManager) DeletePost(username string, url string) (err error) {
	if url == "" {
		return errors.New("the URL can't be empty")
	}

	if username == "" {
		return errors.New("the user name can't be empty")
	}

	err = m.postStore.DeletePost(username, url)
	if err != nil {
		return
	}

	if m.eventSink != nil {
		followers, err := m.socialGraphManager.GetFollowers(username)
		if err != nil {
			return err
		}

		for follower := range followers {
			m.eventSink.OnPostDeleted(follower, url)
		}
	}
	return
}

func (m *PostManager) OnPostChecked(username string, url string, status om.PostStatus) {
	m.postStore.SetPostStatus(username, url, status)
}

func NewPostManager(PostStore PostStore,
	socialGraphManager om.SocialGraphManager,
	natsUrl string,
	eventSink om.PostManagerEvents,
	maxPostsPerUser int64) (om.PostManager, error) {
	if PostStore == nil {
		return nil, errors.New("post store")
	}

	if eventSink != nil && socialGraphManager == nil {
		return nil, errors.New("social graph manager can't be nil if event sink is not nil")
	}

	Post_manager := &PostManager{
		postStore:          PostStore,
		socialGraphManager: socialGraphManager,
		eventSink:          eventSink,
		maxPostsPerUser:    maxPostsPerUser,
	}

	// Subscribe PostManager to Post checker events if nats is ocnfigured
	if natsUrl != "" {
		post_checker_events.Listen(natsUrl, Post_manager)
	}

	return Post_manager, nil
}
