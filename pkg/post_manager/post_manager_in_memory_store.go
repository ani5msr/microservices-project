package post_manager

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
)

// User posts are a map of url:TaggedLink
type UserPosts map[string]*om.Post

// Link store is a map of username:UserLinks
type inMemoryPostStore struct {
	posts map[string]UserPosts
}

func NewInMemoryPostStore() PostStore {
	return &inMemoryPostStore{map[string]UserPosts{}}
}

func (m *inMemoryPostStore) GetPost(request om.GetPostRequest) (result om.GetPostResult, err error) {
	result.Posts = []om.Post{}
	userPosts := m.posts[request.Username]
	if userPosts == nil {
		return
	}

	// Prepare complied regexes
	var urlRegex *regexp.Regexp
	var titleRegex *regexp.Regexp
	var descriptionRegex *regexp.Regexp
	if request.UrlRegex != "" {
		urlRegex, err = regexp.Compile(request.UrlRegex)
		if err != nil {
			return
		}
	}

	if request.TitleRegex != "" {
		titleRegex, err = regexp.Compile(request.UrlRegex)
		if err != nil {
			return
		}
	}

	if request.DescriptionRegex != "" {
		descriptionRegex, err = regexp.Compile(request.UrlRegex)
		if err != nil {
			return
		}
	}

	for _, post := range userPosts {
		// Check wach link against the regular expressions
		if urlRegex != nil && !urlRegex.MatchString(post.Url) {
			continue
		}

		if titleRegex != nil && !titleRegex.MatchString(post.Title) {
			continue
		}

		if descriptionRegex != nil && !descriptionRegex.MatchString(post.Description) {
			continue
		}

		// If there no tag was requested add link immediately and continue
		if request.Tag == "" {
			result.Posts = append(result.Posts, *post)
			continue
		}

		// Add link only if it has the request tag
		if post.Tags[request.Tag] {
			result.Posts = append(result.Posts, *post)
		}
	}

	return
}

func (m *inMemoryPostStore) AddPost(request om.AddPostRequest) (link *om.Post, err error) {
	if request.Url == "" {
		err = errors.New("URL can't be empty")
		return
	}

	if request.Username == "" {
		err = errors.New("user name can't be empty")
		return
	}

	userPosts := m.posts[request.Username]
	if userPosts == nil {
		m.posts[request.Username] = UserPosts{}
		userPosts = m.posts[request.Username]
	} else {
		if userPosts[request.Url] != nil {
			msg := fmt.Sprintf("user %s already has a link for %s", request.Username, request.Url)
			err = errors.New(msg)
			return
		}
	}

	link = &om.Post{
		Url:         request.Url,
		Title:       request.Title,
		Description: request.Description,
		Status:      om.PostStatusPending,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Tags:        request.Tags,
	}
	userPosts[request.Url] = link

	return
}

func (m *inMemoryPostStore) UpdatePost(request om.UpdatePostRequest) (post *om.Post, err error) {
	userPosts := m.posts[request.Username]
	if userPosts == nil || userPosts[request.Url] == nil {
		msg := fmt.Sprintf("User %s doesn't have a link for %s", request.Username, request.Url)
		err = errors.New(msg)
		return
	}

	post = userPosts[request.Url]
	if request.Title != "" {
		post.Title = request.Title
	}

	if request.Description != "" {
		post.Description = request.Description
	}

	newTags := request.AddTags
	for t, _ := range post.Tags {
		if request.RemoveTags[t] {
			continue
		}

		newTags[t] = true
	}

	return
}

func (m *inMemoryPostStore) DeletePost(username string, url string) error {
	if url == "" {
		return errors.New("URL can't be empty")
	}

	if username == "" {
		return errors.New("User name can't be empty")
	}

	userPosts := m.posts[username]
	if userPosts == nil || userPosts[url] == nil {
		msg := fmt.Sprintf("User %s doesn't have a link for %s", username, url)
		return errors.New(msg)
	}

	delete(m.posts[username], url)
	return nil
}

func (m *inMemoryPostStore) SetPostStatus(username string, url string, status om.PostStatus) error {
	if url == "" {
		return errors.New("URL can't be empty")
	}

	if username == "" {
		return errors.New("User name can't be empty")
	}

	userPosts := m.posts[username]
	if userPosts == nil || userPosts[url] == nil {
		msg := fmt.Sprintf("User %s doesn't have a link for %s", username, url)
		return errors.New(msg)
	}

	userPosts[url].Status = status
	return nil
}
