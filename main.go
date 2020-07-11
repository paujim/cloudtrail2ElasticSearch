package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

func main() {
	lambda.Start(HandleRequest)
}

type CloudTrailEvent struct {
	Version    string      `json:"version"`
	ID         string      `json:"id"`
	DetailType string      `json:"detail-type"`
	Source     string      `json:"source"`
	AccountID  string      `json:"account"`
	Time       time.Time   `json:"time"`
	Region     string      `json:"region"`
	Resources  []string    `json:"resources"`
	Detail     interface{} `json:"detail"`
}

type CloudTrailEventDetail struct {
	EventVersion      string      `json:"eventVersion"`
	UserIdentity      interface{} `json:"userIdentity"`
	EventTime         time.Time   `json:"eventTime"`
	EventSource       string      `json:"eventSource"`
	EventName         string      `json:"eventName"`
	AwsRegion         string      `json:"awsRegion"`
	SourceIPAddress   string      `json:"sourceIPAddress"`
	UserAgent         string      `json:"userAgent"`
	RequestParameters interface{} `json:"requestParameters"`
	ResponseElements  interface{} `json:"responseElements"`
	RequestID         string      `json:"requestID"`
	EventID           string      `json:"eventID"`
	EventType         string      `json:"eventType"`
}

func HandleRequest(ctx context.Context, event CloudTrailEvent) error {

	domain := os.Getenv("ES_HOST")
	region := os.Getenv("ES_REGION")

	index := "iam-user-" + event.Time.Format("2006-01-02")
	log.Println("INDEX:")
	log.Println(index)
	endpoint := domain + "/" + index + "/" + "_doc" + "/"

	service := "es"

	b, err := json.Marshal(event)
	if err != nil {
		log.Print(err)
		return err
	}

	json := string(b)
	log.Println("EVENT:")
	log.Print(json)
	body := strings.NewReader(json)

	// credentials := credentials.NewSharedCredentials("", "devops")
	credentials := credentials.NewEnvCredentials()
	signer := v4.NewSigner(credentials)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		log.Print(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	signer.Sign(req, body, service, region, time.Now())
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(resp)
	return nil
}
