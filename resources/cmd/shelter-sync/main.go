package main

import (
	"shinagawa-shelter/pkg/http"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(http.SyncHandler)
}
