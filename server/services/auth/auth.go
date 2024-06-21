package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/logger"
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
func GetTokensFromAuthCode(code string, origin string) GoogleApiTokenResponse {
	// Call Google oauth token endpoint
	redirectUri := fmt.Sprintf("%s/auth", origin)
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

type GoogleApiAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	Error       bson.M `json:"error"`
}

func RefreshAccessToken(refreshToken string) GoogleApiAccessTokenResponse {
	values := url.Values{
		"client_id":     {os.Getenv("CLIENT_ID")},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
	}

	resp, err := http.PostForm(
		"https://oauth2.googleapis.com/token",
		values,
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	defer resp.Body.Close()

	var res GoogleApiAccessTokenResponse
	json.NewDecoder(resp.Body).Decode(&res)

	return res
}

type RefreshAccessTokenData struct {
	TokenResponse GoogleApiAccessTokenResponse
	Email         string
	Error         *interface{}
}

func RefreshAccessTokenAsync(refreshToken string, email string, c chan RefreshAccessTokenData) {
	// Recover from panics
	defer func() {
		if err := recover(); err != nil {
			c <- RefreshAccessTokenData{Error: &err}
		}
	}()

	tokenResponse := RefreshAccessToken(refreshToken)

	c <- RefreshAccessTokenData{tokenResponse, email, nil}
}
