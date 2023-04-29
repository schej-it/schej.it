package errs

// Errors enum
// TODO: make these an actual type (i.e. Errors.NotSignedIn)
const (
	NotSignedIn           string = "not-signed-in"
	UserDoesNotExist      string = "user-does-not-exist"
	EventNotFound         string = "event-not-found"
	FriendRequestNotFound string = "friend-request-not-found"
	UserNotFriends        string = "user-not-friends"
	UserNotEventOwner     string = "user-not-event-owner"
)

type GoogleAPIError struct {
	Code    int         `json:"code"`
	Details interface{} `json:"details"`
	Errors  interface{} `json:"errors"`
}
