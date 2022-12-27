package object_model

type ThoughtManager interface {
	GetThought(request GetLinksRequest) (GetLinksResult, error)
	AddThought(request AddLinkRequest) error
	UpdateThought(request UpdateLinkRequest) error
	DeleteThought(username string, url string) error
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

type NewsManager interface {
	GetNews(request GetNewsRequest) (GetNewsResult, error)
}

type ThoughtManagerEvents interface {
	OnThoughtAdded(username string, link *Link)
	OnThoughtUpdated(username string, link *Link)
	OnThoughtDeleted(username string, url string)
}

type ThoughtCheckerEvents interface {
	OnThoughtChecked(username string, url string, status LinkStatus)
}
