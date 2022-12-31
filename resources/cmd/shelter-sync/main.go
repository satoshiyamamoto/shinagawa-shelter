package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"shinagawa-shelter/pkg/config"
	"shinagawa-shelter/pkg/database"
	"shinagawa-shelter/pkg/http"
	"shinagawa-shelter/pkg/model"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) error {

	for _, url := range config.DatasetURLs {
		log.Println("get dataset from", url)
		body, err := http.GetDataset(url)
		if err != nil {
			log.Println("failed to get dataset", err)
			return err
		}

		r := csv.NewReader(body)
		r.Read() // truncate header

		records, err := r.ReadAll()
		if err == io.EOF {
			log.Println("dataaset is empty")
			return nil
		}
		if err != nil {
			log.Println("failed to parse dataaset", err)
		}

		for _, record := range records {
			shelter := model.NewShelter(record)
			database.MergeShelter(shelter)
		}

		log.Printf("%d shelters synced", len(records))
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
