package main

import (
	"testing"
	"time"

	"github.com/chansen-p44/euchef/src"
)

func TestLambdaPost(t *testing.T) {

	format := "02-01-2006"
	now, _ := time.Parse(format, "15-01-2019")

	endpoint := getEnv("DynamoEndpoint", "http://localhost:8000")
	tableName := getEnv("DynamoTableName", "euchef_Status")

	storage := euchef.NewDynamoStorage(endpoint, tableName)

	post(now, storage)
}
