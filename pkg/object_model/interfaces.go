package object_model

type PostManager interface {
	GetPost(request GetPostRequest) (GetPostResult, error)
	AddPost(request AddPostRequest) error
	UpdatePost(request UpdatePostRequest) error
	DeletePost(username string, url string) error
}

type UserManager interface {
	Register(user User) error
	Login(username string, authToken string) (session string, err error)
	Logout(username string, session string) error
}

type SocialGraphManager interface {
	Follow(followed string, follower string) error
	Unfollow(followed string, follower string) error

	GetFollowing(username string) (map[string]bool, error)
	GetFollowers(username string) (map[string]bool, error)

	//AcceptFollowRequest(followed string, follower string) error
	//RejectFollowRequest(followed string, follower string) error
	//KickFollower(followed string, follower string) error
}

type FeedManager interface {
	GetFeed(request GetFeedRequest) (GetFeedResult, error)
}

type PostManagerEvents interface {
	OnPostAdded(username string, post *Post)
	OnPostUpdated(username string, post *Post)
	OnPostDeleted(username string, url string)
}

type PostCheckerEvents interface {
	OnPostChecked(username string, url string, status PostStatus)
}
