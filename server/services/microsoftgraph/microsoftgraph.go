package microsoftgraph

import (
	"encoding/json"

	"schej.it/server/logger"
	"schej.it/server/models"
	"schej.it/server/services"
)

type UserInfo struct {
	FirstName string `json:"givenName"`
	LastName  string `json:"surname"`
	Email     string `json:"mail"`
}

func GetUserInfo(user *models.User, calendarAuth *models.OAuth2CalendarAuth) UserInfo {
	response := services.CallApi(
		user,
		calendarAuth,
		"GET",
		"https://graph.microsoft.com/v1.0/me?$select=givenName,surname,mail",
		nil,
	)
	defer response.Body.Close()

	userResponse := struct {
		GivenName string `json:"givenName"`
		Surname   string `json:"surname"`
		Mail      string `json:"mail"`
	}{}

	if err := json.NewDecoder(response.Body).Decode(&userResponse); err != nil {
		logger.StdErr.Panicln(err)
	}

	return UserInfo{
		FirstName: userResponse.GivenName,
		LastName:  userResponse.Surname,
		Email:     userResponse.Mail,
	}
}
