package http

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"shinagawa-shelter/pkg/config"
	"shinagawa-shelter/pkg/database"
	"shinagawa-shelter/pkg/model"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

func SyncHandler(ctx context.Context) error {

	for _, url := range config.DatasetURLs {
		log.Println("get dataset from", url)
		body, err := GetDataset(url)
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

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")

	// DefaultPaze Default Page number
	DefaultPaze = 1

	// DefaultPazeSize Default Page size
	DefaultPazeSize = 10
)

func ApiHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	category := request.QueryStringParameters["category"]
	page, err := strconv.Atoi(request.QueryStringParameters["page"])
	if page < 1 || err != nil {
		page = DefaultPaze
	}
	size, err := strconv.Atoi(request.QueryStringParameters["size"])
	if size < 1 || err != nil {
		size = DefaultPazeSize
	}

	shelters, err := database.FindShelters(&category, &page, &size)
	if err != nil {
		log.Println("failed to find shelters", err)
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	b, err := json.Marshal(shelters)
	if err != nil {
		log.Println("failed to respond json", err)
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	return events.APIGatewayProxyResponse{
		Body:       string(b),
		StatusCode: 200,
	}, nil
}
