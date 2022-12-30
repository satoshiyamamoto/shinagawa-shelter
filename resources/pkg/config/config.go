package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	DataSourceName string
	DatasetURLs    []string
)

func init() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	DataSourceName = fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		user,
		password,
		host,
		name,
	)
	log.Println("Configured DataSourceName:", strings.Replace(DataSourceName, fmt.Sprintf(":%s@", password), ":******@", 1))

	urls := os.Getenv("DATASET_URLS")
	DatasetURLs = strings.Split(urls, ",")
	log.Println("Configured DatasetURLs:", DatasetURLs)
}
