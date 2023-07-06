package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
	"schej.it/server/logger"
)

// Send email to the given email
func SendEmail(toEmail string, subject string, body string, contentType string) {
	if contentType == "" {
		contentType = "text/plain"
	}

	appPassword := os.Getenv("GMAIL_APP_PASSWORD")
	fromEmail := os.Getenv("SCHEJ_EMAIL_ADDRESS")

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, body)

	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, appPassword)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logger.StdErr.Println(err)
	}
}

func AddUserToMailchimp(email string, firstName string, lastName string) {
	// Adds the given user to the default mailchimp audience
	apiKey := os.Getenv("MAILCHIMP_API_KEY")

	body, _ := json.Marshal(bson.M{
		"email_address": email, "status": "subscribed", "merge_fields": bson.M{
			"FNAME": firstName,
			"LNAME": lastName,
		},
		"tags": bson.A{"user"},
	})
	bodyBuffer := bytes.NewBuffer(body)

	req, _ := http.NewRequest("POST", "https://us21.api.mailchimp.com/3.0/lists/b5c79106b4/members", bodyBuffer)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
	}
	defer resp.Body.Close()
}

func AddUserToMailjet(email string, firstName string, lastName string, picture string) {
	// Adds the given user to the Mailjet contact list
	apiKey := os.Getenv("MAILJET_API_KEY")
	apiSecret := os.Getenv("MAILJET_API_SECRET")
	listId := os.Getenv("MAILJET_LIST_ID")

	// Create contact
	// POST https://api.mailjet.com/v3/REST/contact {"EMAIL", email}
	// contactId = result.Data[0].ID (integer)
	body, _ := json.Marshal(bson.M{
		"Email": email,
	})
	bodyBuffer := bytes.NewBuffer(body)

	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/contact", bodyBuffer)
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}
	defer resp.Body.Close()

	result := struct {
		Data []struct {
			ID int
		}
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.StdErr.Println(err)
		return
	}

	if len(result.Data) == 0 {
		return
	}
	contactId := result.Data[0].ID

	// Update contact metadata
	// PUT https://api.mailjet.com/v3/REST/contactdata/$contactId
	// { "Data": [{"Name": "firstname", "Value":"first name!"}] }
	body, _ = json.Marshal(bson.M{
		"Data": bson.A{
			bson.M{
				"Name":  "firstname",
				"Value": firstName,
			},
			bson.M{
				"Name":  "lastname",
				"Value": lastName,
			},
			bson.M{
				"Name":  "picture",
				"Value": picture,
			},
		},
	})
	bodyBuffer = bytes.NewBuffer(body)

	req, _ = http.NewRequest("PUT", fmt.Sprintf("https://api.mailjet.com/v3/REST/contactdata/%d", contactId), bodyBuffer)
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}
	defer resp.Body.Close()

	// Add contact to "schej users" contact list
	// POST https://api.mailjet.com/v3/REST/contact/$contactId/managecontactslists
	// '{ "ContactsLists": [{"Action": "addforce", "ListID": "10219365"}] }'
	body, _ = json.Marshal(bson.M{
		"ContactsLists": bson.A{
			bson.M{
				"Action": "addforce",
				"ListID": listId,
			},
		},
	})
	bodyBuffer = bytes.NewBuffer(body)

	req, _ = http.NewRequest("POST", fmt.Sprintf("https://api.mailjet.com/v3/REST/contact/%d/managecontactslists", contactId), bodyBuffer)
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}
	defer resp.Body.Close()
}
