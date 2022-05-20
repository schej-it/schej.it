package errors

// Errors enum
const (
	NotSignedIn      string = "not-signed-in"
	UserDoesNotExist string = "user-does-not-exist"
)

type GoogleAPIError struct {
	Code    int         `json:"code"`
	Details interface{} `json:"details"`
	Errors  interface{} `json:"errors"`
}
