package main

import (
	"testing"

	"github.com/chansen-p44/euchef/src"
)

func TestLambdaDelete(t *testing.T) {

	endpoint := getEnv("DynamoEndpoint", "http://localhost:8000")
	tableName := getEnv("DynamoTableName", "euchef_Status")

	storage := euchef.NewDynamoStorage(endpoint, tableName)

	delete(storage)
}
