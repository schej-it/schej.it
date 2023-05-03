package contacts

import (
	"encoding/json"
	"fmt"

	"schej.it/server/errs"
	"schej.it/server/logger"
	"schej.it/server/models"
	"schej.it/server/services"
)

func SearchContacts(user *models.User, query string) ([]models.UserProfile, *errs.GoogleAPIError) {
	type Person struct {
		Names []struct {
			FamilyName string `json:"familyName"`
			GivenName  string `json:"givenName"`
		} `json:"names"`
		Photos []struct {
			Url string `json:"url"`
		} `json:"photos"`
		EmailAddresses []struct {
			Value string `json:"value"`
		} `json:"emailAddresses"`
	}

	// Search contacts
	response := services.CallGoogleApi(
		user,
		"GET",
		fmt.Sprintf("https://people.googleapis.com/v1/people:searchContacts?query=%s&pageSize=5&readMask=names,emailAddresses,photos", query),
		nil,
	)
	defer response.Body.Close()

	// Parse response
	contactsData := struct {
		Results []struct {
			Person Person `json:"person"`
		} `json:"results"`
		Error *errs.GoogleAPIError `json:"error"`
	}{}
	if err := json.NewDecoder(response.Body).Decode(&contactsData); err != nil {
		logger.StdErr.Panicln(err)
	}

	directoryData := struct {
		People []Person             `json:"people"`
		Error  *errs.GoogleAPIError `json:"error"`
	}{}
	if len(query) > 0 {
		// Search Directory
		response = services.CallGoogleApi(
			user,
			"GET",
			fmt.Sprintf("https://people.googleapis.com/v1/people:searchDirectoryPeople?query=%s&pageSize=5&readMask=names,emailAddresses,photos&sources=DIRECTORY_SOURCE_TYPE_DOMAIN_PROFILE", query),
			nil,
		)
		defer response.Body.Close()

		// Parse response
		if err := json.NewDecoder(response.Body).Decode(&directoryData); err != nil {
			logger.StdErr.Panicln(err)
		}
	}

	// Throw error if contacts access has not been granted
	if contactsData.Error != nil {
		return nil, contactsData.Error
	} else if directoryData.Error != nil && directoryData.Error.Code == 403 {
		// Need to check if code is 403 because error 400 occurs when user is not a GSuite user (which is okay and should not throw an error)
		return nil, directoryData.Error
	}

	// Format list of contacts search results
	contacts := make([]models.UserProfile, 0)
	for _, result := range contactsData.Results {
		var userProfile models.UserProfile
		userProfile.FirstName = result.Person.Names[0].GivenName
		userProfile.LastName = result.Person.Names[0].FamilyName
		if len(result.Person.Photos) > 0 {
			userProfile.Picture = result.Person.Photos[0].Url
		}

		for _, email := range result.Person.EmailAddresses {
			userProfile.Email = email.Value
			contacts = append(contacts, userProfile)
		}
	}
	for _, person := range directoryData.People {
		var userProfile models.UserProfile
		userProfile.FirstName = person.Names[0].GivenName
		userProfile.LastName = person.Names[0].FamilyName
		if len(person.Photos) > 0 {
			userProfile.Picture = person.Photos[0].Url
		}

		for _, email := range person.EmailAddresses {
			userProfile.Email = email.Value
			contacts = append(contacts, userProfile)
		}
	}

	return contacts, nil
}
