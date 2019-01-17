package main

import (
	"log"
	"os"
	"github.com/chansen-p44/euchef/src"
	"github.com/aws/aws-lambda-go/lambda"
)

func delete(storage *euchef.DynamoStorage) error {

	// Store the timestamp in the database so we can delete it later in the day
	log.Println("Reading status from database")
	status, err := storage.GetStatus()
	if err != nil {
		log.Println("GetStatus error:", err)
		return err
	}

	key := getEnv("SlackAPIKey", "xoxp-522221885958-522221886902-521024060229-e270a8fe90bf16917390cb01d6ebc0df")
	
	log.Println("Deleting message with timestamp:", status.Timestamp)
	err = euchef.DeleteSlackMessage(status.Timestamp, status.ChannelID, key)
	if err != nil {
		log.Println("PostSlackMessage error:", err)
		return err
	}
	
	log.Println("All done")

	// Store slack post id
	return nil
}

func handler(request map[string]interface{}) error {

	log.Printf("Schedule Request: %+v", request)

	endpoint := getEnv("DynamoEndPoint", "http://localhost:8000")
	tableName := getEnv("DynamoTableName", "euchef_Status")
	log.Printf("Creating dynamodb storage with endpoint: %s and tablename: %s", endpoint, tableName)
	storage := euchef.NewDynamoStorage(endpoint, tableName)

	// Delete any current messages
	delete(storage)

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
