package http

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"shinagawa-shelter/pkg/config"
	"shinagawa-shelter/pkg/database"
	"shinagawa-shelter/pkg/model"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")

	// DefaultCategory Default category
	DefaultCategory string = ""

	// DefaultPage Default Page number
	DefaultPage = 1

	// DefaultPageSize Default Page size
	DefaultPageSize = 10

	// MaxPageSize Maximum Page size
	MaxPageSize = math.MaxInt
)

func SyncHandler(ctx context.Context) error {
	runningAt := time.Now()

	// Update the shelters
	for _, url := range config.DatasetURLs {
		log.Println("get dataset from", url)
		body, err := GetDataset(url)
		if err != nil {
			log.Println("failed to get dataset", err)
			return err
		}

		r := csv.NewReader(body)
		r.Read() // Skip header

		records, err := r.ReadAll()
		if err == io.EOF {
			log.Println("dataset is empty")
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

	// Delete the closed shelters
	shelters, err := database.FindShelters(&DefaultCategory, &DefaultPage, &MaxPageSize)
	if err != nil {
		log.Println("failed to delete closed shelters")
		return err
	}

	i := 0
	for _, s := range shelters {
		if s.UpdatedAt != nil && s.UpdatedAt.Before(runningAt) {
			database.DeleteShelter(s)
			i++
		}
	}

	if i > 0 {
		log.Printf("%d shelters deleted", i)
	}

	return nil
}

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
		page = DefaultPage
	}
	size, err := strconv.Atoi(request.QueryStringParameters["size"])
	if size < 1 || err != nil {
		size = DefaultPageSize
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
