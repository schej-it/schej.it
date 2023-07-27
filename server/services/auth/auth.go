package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"schej.it/server/logger"
	"schej.it/server/utils"
)

type GoogleApiTokenResponse struct {
	AccessToken      string `json:"access_token"`
	IdToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Returns access, refresh, and id tokens from the auth code
func GetTokensFromAuthCode(code string) GoogleApiTokenResponse {
	// Call Google oauth token endpoint
	var redirectUri string
	if utils.IsRelease() {
		redirectUri = "https://schej.it/auth"
	} else {
		redirectUri = "http://localhost:8080/auth"
	}
	values := url.Values{
		"client_id":     {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectUri},
	}
	resp, err := http.PostForm(
		"https://oauth2.googleapis.com/token",
		values,
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	defer resp.Body.Close()

	var res GoogleApiTokenResponse

	json.NewDecoder(resp.Body).Decode(&res)
	if len(res.Error) > 0 {
		data, _ := json.MarshalIndent(res, "", "  ")
		logger.StdErr.Panicln(string(data))
	}

	return res
}
