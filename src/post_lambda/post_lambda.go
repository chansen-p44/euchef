package main

import (
	"log"
	"os"
	"time"
	"github.com/chansen-p44/euchef/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func post(now time.Time, storage *euchef.DynamoStorage) error {

	// Fetch data
	data, err := euchef.FetchData(now)
	if err != nil {
		log.Println("fetchData error:", err)
		return err
	}

	// Parse Data
	items, err := euchef.ParseData(data)
	if err != nil {
		log.Println("parseData error:", err)
		return err
	}

	// Post slack
	if len(items) == 0 {
		log.Println("No menus found")
		return nil
	}

	key := getEnv("SlackAPIKey", "")
	channelID := getEnv("SlackChannelID", "")
	timestamp, err := euchef.PostSlackMessage(items, channelID, key)
	if err != nil {
		log.Println("PostSlackMessage error:", err)
		return err
	}

	// Store the timestamp in the database so we can delete it later in the day
	log.Println("Storing status in database")
	status := euchef.Status{Timestamp: timestamp, ChannelID: channelID}
	err = storage.StoreStatus(status)
	if err != nil {
		log.Println("StoreStatus error:", err)
		return nil
	}

	log.Println("All done")
	
	return nil
}

func handler(request map[string]interface{}) error {

	log.Printf("Schedule Request: %+v", request)

	endpoint := getEnv("DynamoEndPoint", "http://localhost:8000")
	tableName := getEnv("DynamoTableName", "euchef_Status")
	log.Printf("Creating dynamodb storage with endpoint: %s and tablename: %s", endpoint, tableName)
	storage := euchef.NewDynamoStorage(endpoint, tableName)

	now := time.Now()

	// For debugging
	// format := "02-01-2006"
	// now, _ := time.Parse(format, "17-01-2019")

	// Determine current date and post message if available
	post(now, storage)

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return fallback
}

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// Tell the lambda runtime where the handlers live
	lambda.Start(handler)
}
