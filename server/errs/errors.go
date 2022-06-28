package errs

// Errors enum
const (
	NotSignedIn           string = "not-signed-in"
	UserDoesNotExist      string = "user-does-not-exist"
	EventNotFound         string = "event-not-found"
	FriendRequestNotFound string = "friend-request-not-found"
)

type GoogleAPIError struct {
	Code    int         `json:"code"`
	Details interface{} `json:"details"`
	Errors  interface{} `json:"errors"`
}
