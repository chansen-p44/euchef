package euchef

import (
	"log"
	"testing"
)

func TestStoreStatus(t *testing.T) {

	//endpoint := "https://dynamodb.eu-west-1.amazonaws.com"
	endpoint := "http://localhost:8000"
	tableName := "euchef_Status"

	dynamoStorage := NewDynamoStorage(endpoint, tableName)

	status := Status{Timestamp: "123", ChannelID: ""}

	err := dynamoStorage.StoreStatus(status)
	if err != nil {
		t.Fatal(err)
	}

}
func TestGetStatus(t *testing.T) {

	//endpoint := "https://dynamodb.eu-west-1.amazonaws.com"
	endpoint := "http://localhost:8000"
	tableName := "euchef_Status"

	dynamoStorage := NewDynamoStorage(endpoint, tableName)

	status, err := dynamoStorage.GetStatus()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Staus: %+v", status)
}
