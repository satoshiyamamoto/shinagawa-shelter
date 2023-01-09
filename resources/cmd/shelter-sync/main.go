package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/satoshiyamamoto/shinagawa-shelter/pkg/http"
)

func main() {
	lambda.Start(http.SyncHandler)
}
