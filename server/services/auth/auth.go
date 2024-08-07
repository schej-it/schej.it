package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/logger"
	"schej.it/server/models"
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

// If access token has expired, get a new token for the primary account as well as all other calendar accounts, update the user object, and save it to the database
// `accounts` specifies for which accounts to refresh access tokens. If `accounts` is nil or empty, then update tokens for all accounts
func RefreshUserTokenIfNecessary(u *models.User, accounts models.Set[string]) {
	refreshTokenChan := make(chan RefreshAccessTokenData)
	numAccountsToUpdate := 0

	// If `accounts` is nil, then update tokens for all accounts
	updateAllAccounts := len(accounts) == 0

	// Refresh calendar account access tokens if necessary
	for _, account := range u.CalendarAccounts {
		if account.CalendarType == models.GoogleCalendarType { // Only refresh access tokens for Google calendar accounts
			accountAuth := account.GoogleCalendarAuth

			if _, ok := accounts[account.Email]; ok || updateAllAccounts {
				if time.Now().After(accountAuth.AccessTokenExpireDate.Time()) && len(accountAuth.RefreshToken) > 0 {
					go RefreshAccessTokenAsync(accountAuth.RefreshToken, account.Email, refreshTokenChan)
					numAccountsToUpdate++
				}
			}
		}
	}

	// Update access tokens as responses are received
	for i := 0; i < numAccountsToUpdate; i++ {
		res := <-refreshTokenChan

		if res.Error != nil {
			continue
		}

		accessTokenExpireDate := utils.GetAccessTokenExpireDate(res.TokenResponse.ExpiresIn)

		calendarAccountKey := utils.GetCalendarAccountKey(res.Email, models.GoogleCalendarType)
		if calendarAccount, ok := u.CalendarAccounts[calendarAccountKey]; ok {
			calendarAccount.GoogleCalendarAuth.AccessToken = res.TokenResponse.AccessToken
			calendarAccount.GoogleCalendarAuth.AccessTokenExpireDate = primitive.NewDateTimeFromTime(accessTokenExpireDate)
			u.CalendarAccounts[calendarAccountKey] = calendarAccount
		}
	}

	// Update user object if accounts were updated
	if numAccountsToUpdate > 0 {
		db.UsersCollection.FindOneAndUpdate(
			context.Background(),
			bson.M{"_id": u.Id},
			bson.M{"$set": u},
		)
	}
}
