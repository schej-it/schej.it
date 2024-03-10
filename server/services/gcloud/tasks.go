package gcloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	"cloud.google.com/go/cloudtasks/apiv2beta3/cloudtaskspb"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/timestamppb"
	"schej.it/server/logger"
	"schej.it/server/services/listmonk"
	"schej.it/server/utils"
)

var TasksClient *cloudtasks.Client

func InitTasks() func() {
	ctx := context.Background()

	var err error
	credsFile := os.Getenv("SERVICE_ACCOUNT_KEY_PATH")

	TasksClient, err = cloudtasks.NewClient(ctx, option.WithCredentialsFile(credsFile))
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	// Return function to close client
	return func() {
		TasksClient.Close()
	}
}

func CreateEmailTask(email string, ownerName string, eventName string, eventId string) {
	// Get listmonk url env vars
	listmonkUrl := os.Getenv("LISTMONK_URL")
	listmonkUsername := os.Getenv("LISTMONK_USERNAME")
	listmonkPassword := os.Getenv("LISTMONK_PASSWORD")
	basicAuthString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", listmonkUsername, listmonkPassword)))

	// Find if subscriber exists in listmonk
	subscriberExists, _ := listmonk.DoesUserExist(email)

	// If subscriber doesn't exist, add subscriber to listmonk
	if !subscriberExists {
		listmonk.AddUserToListmonk(email, "", "", "", nil)
	}

	// Get email template ids
	initialEmailReminderId, err := strconv.Atoi(os.Getenv("LISTMONK_INITIAL_EMAIL_REMINDER_ID"))
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	// secondEmailReminderId, err := strconv.Atoi(os.Getenv("LISTMONK_SECOND_EMAIL_REMINDER_ID"))
	// if err != nil {
	// 	logger.StdErr.Panicln(err)
	// }
	// finalEmailReminderId, err := strconv.Atoi(os.Getenv("LISTMONK_FINAL_EMAIL_REMINDER_ID"))
	// if err != nil {
	// 	logger.StdErr.Panicln(err)
	// }

	// Construct URLs
	var baseUrl string
	if utils.IsRelease() {
		baseUrl = "https://schej.it"
	} else {
		baseUrl = "http://localhost:8080"
	}
	eventUrl := fmt.Sprintf("%s/e/%s", baseUrl, eventId)
	finishedUrl := fmt.Sprintf("%s/e/%s/responded", baseUrl, eventId)

	// Create JSON object
	body, err := json.Marshal(bson.M{
		"subscriber_email": email,
		"template_id":      initialEmailReminderId,
		"data": bson.M{
			"ownerName":   ownerName,
			"eventName":   eventName,
			"eventUrl":    eventUrl,
			"finishedUrl": finishedUrl,
		},
		"content_type": "html",
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	task, err := TasksClient.CreateTask(context.Background(), &cloudtaskspb.CreateTaskRequest{
		Parent: "projects/schej-it/locations/us-central1/queues/SendReminderEmail",
		Task: &cloudtaskspb.Task{
			ScheduleTime: timestamppb.New(time.Now().Add(30 * time.Second)),
			// ScheduleTime: timestamppb.Now(),
			PayloadType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					Url:        fmt.Sprintf("%s/api/tx", listmonkUrl),
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Headers: map[string]string{
						"Authorization": fmt.Sprintf("Basic %s", basicAuthString),
						"Content-Type":  "application/json",
					},
					Body: body,
				},
			},
		},
	})

	if err != nil {
		logger.StdErr.Panicln(err)
	}

	fmt.Println("TASK NAME: ", task.Name)
}
