package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"context"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type ESClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LambdaHandler struct {
	client ESClient
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

func ToJsonString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(b)
}

func (h *LambdaHandler) Handle(ctx context.Context, event CloudTrailEvent) error {

	domain := os.Getenv("ES_HOST")
	region := os.Getenv("ES_REGION")
	esRole := os.Getenv("ES_ROLE")

	log.Printf("ES_HOST: %s\n", domain)
	log.Printf("ES_REGION: %s\n", region)
	index := "iam-user-" + event.Time.Format("2006-01-02")
	log.Printf("INDEX: %s", index)

	endpoint := domain + "/" + index + "/" + "_doc" + "/"

	service := "es"

	json := ToJsonString(event)
	log.Println("EVENT:")
	log.Print(json)
	body := strings.NewReader(json)

	credentials := stscreds.NewCredentials(session.Must(session.NewSession()), esRole)

	signer := v4.NewSigner(credentials)
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		log.Print(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	signer.Sign(req, body, service, region, time.Now())
	resp, err := h.client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(resp)
	return nil
}
