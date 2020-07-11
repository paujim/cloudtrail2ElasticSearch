package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

var lambdaHanlder = LambdaHandler{&http.Client{}}

func main() {
	lambda.Start(lambdaHanlder.Handle)
}
