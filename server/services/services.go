package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/db"
	"schej.it/server/logger"
	"schej.it/server/models"
)

// Calls the given url with the given method using the user's google api access token
func CallGoogleApi(user *models.User, method string, url string, body *bson.M) *http.Response {
	db.RefreshUserTokenIfNecessary(user, nil)

	// Format body as a buffer if not nil
	var bodyBuffer *bytes.Buffer
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyBuffer = bytes.NewBuffer(bodyBytes)
	} else {
		bodyBuffer = nil
	}

	// Construct request
	var req *http.Request
	if bodyBuffer != nil {
		req, _ = http.NewRequest(method, url, bodyBuffer)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.CalendarAccounts[user.Email].AccessToken))

	// Execute request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	return response
}
