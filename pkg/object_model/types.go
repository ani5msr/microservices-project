package object_model

import "time"

type PostManagerEventTypeEnum int

const (
	PostAdded PostManagerEventTypeEnum = iota
	PostUpdated
	PostDeleted
)

const (
	PostStatusPending = "pending"
	PostStatusValid   = "valid"
	PostStatusInvalid = "invalid"
)

type PostStatus = string

type Post struct {
	Url         string
	Title       string
	Description string
	Status      PostStatus
	Tags        map[string]bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetPostRequest struct {
	UrlRegex         string
	TitleRegex       string
	DescriptionRegex string
	Username         string
	Tag              string
	StartToken       string
}

type GetPostResult struct {
	Posts         []Post
	NextPageToken string
}

type AddPostRequest struct {
	Url         string
	Title       string
	Description string
	Username    string
	Tags        map[string]bool
}

type UpdatePostRequest struct {
	Url         string
	Title       string
	Description string
	Username    string
	AddTags     map[string]bool
	RemoveTags  map[string]bool
}

type User struct {
	Email string
	Name  string
}

type PostManagerEvent struct {
	EventType PostManagerEventTypeEnum
	Username  string
	Url       string
	Timestamp time.Time
}

type GetFeedRequest struct {
	Username   string
	StartToken string
}

type GetFeedResult struct {
	Events    []*PostManagerEvent
	NextToken string
}

type CheckPostRequest struct {
	Username string
	Url      string
}
