package euchef

import (
	"log"
	"time"
	"errors"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const statusUUID = "{EUCHEF-STATUS}"

// Status type
type Status struct {
	UUID string
	Timestamp string
	ChannelID string
}

// DynamoStorage type
type DynamoStorage struct{
	Endpoint string
	TableName string
	DB 	*dynamodb.DynamoDB
}

// NewDynamoStorage creates a new DynamoStorage object
func NewDynamoStorage(endpoint, tableName string) *DynamoStorage {
	return &DynamoStorage{Endpoint: endpoint, TableName: tableName}
}

// CreateSession creates a new dynamodb session
func (ds *DynamoStorage) CreateSession() (error){
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String(ds.Endpoint)},
	)
	if err != nil {
		log.Println("session.NewSession error:",err)
		return err
	}

	ds.DB = dynamodb.New(sess)

	return nil
}


// ListTables lists the tables in the database
func (ds *DynamoStorage) ListTables() ([]*string, error){
	// Get the list of tables
	result, err := ds.DB.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println("dbSvc.ListTables error:",err)
		return nil, err
	}

	return result.TableNames, nil
}

// StoreStatus stores the specified status in the database
func (ds *DynamoStorage) StoreStatus(status Status) error {
	
	// Create a session we can use to communicate with dynamodb
	err := ds.CreateSession()
	if err != nil {
		log.Println("createSession error:", err)
		return err
	}
	
	// List the table in the database (we want to see if our table is already present)
	tables, err := ds.ListTables()
	if err != nil {
		log.Println("ListTables error:", err)
		return err
	}
	
	// See if our table is present
	found := false
	for _, val := range tables {
		if *val == ds.TableName {
			found = true
			break
		}
	}

	if found == false {
		log.Println("Table not found - creating it")
		// Create table
		err := ds.CreateTable()
		if err != nil {
			log.Println("CreateTable error:", err)
			return err
		}
	}

	status.UUID = statusUUID

	av, err := dynamodbattribute.MarshalMap(status)
    if err != nil {
        log.Println("Got error marshalling map:", err)
        return err
    }

    // Create item in table Movies
    input := &dynamodb.PutItemInput{
        Item: av,
        TableName: aws.String(ds.TableName),
    }

    _, err = ds.DB.PutItem(input)
    if err != nil {
        log.Println("Got error calling PutItem:", err)
		return err
	}

	return nil
}


// GetTableStatus gets the status of the specified table
func (ds *DynamoStorage) GetTableStatus() (status string, err error) {

    input := &dynamodb.DescribeTableInput{
        TableName: aws.String(ds.TableName),
    }

    // If the table is being created, we might have to wait a bit for it to show up
    loopCount := 0
    keepLooping := true
    for keepLooping {
        loopCount++

        result, err := ds.DB.DescribeTable(input)
        if err != nil {
            // Inspect the errpr
            if aerr, ok := err.(awserr.Error); ok {
                switch aerr.Code() {
                case dynamodb.ErrCodeResourceNotFoundException:
                    log.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
                    // Wait a bit before we try again
                    time.Sleep(1 * time.Second)
                    continue                    
                case dynamodb.ErrCodeInternalServerError:
                    log.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
                default:
                    log.Println(aerr.Error())
                }
            } else {
                // Print the error, cast err to awserr.Error to get the Code and
                // Message from an error.
                log.Println(err.Error())
            }

            return "", err
        }

        // Get the current status
        tableStatus := result.Table.TableStatus
        return *tableStatus, nil
    }

    return "", errors.New("Timeout")
}

// WaitForStatus waits for a table to reach a specified status
func (ds *DynamoStorage) WaitForStatus(waitstatus string) error {
    keepWaiting := true
    for keepWaiting {
        status, err := ds.GetTableStatus()
        if err != nil {
            log.Println(err)
            return err
        }
        log.Println("Status:", status)

        if status == waitstatus {
            keepWaiting = false
        } else {
            time.Sleep(1 * time.Second)
        }
    }
    return nil
}

// CreateTable creates our table
func (ds *DynamoStorage) CreateTable() (err error) {
	
    input := &dynamodb.CreateTableInput{
        AttributeDefinitions: []*dynamodb.AttributeDefinition{
            {
                AttributeName: aws.String("UUID"),
                AttributeType: aws.String("S"),
            },
        },
        KeySchema: []*dynamodb.KeySchemaElement{
            {
                AttributeName: aws.String("UUID"),
                KeyType:       aws.String("HASH"),
            },
        },
        BillingMode: aws.String("PAY_PER_REQUEST"),
        TableName: aws.String(ds.TableName),
    }

    // DynamoDBlocal doesn't support BillingMode, so just give it some provisioned capacity 
    if strings.Contains(ds.Endpoint, "localhost") {
        input.ProvisionedThroughput = &dynamodb.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(5),
            WriteCapacityUnits: aws.Int64(5),
        }
    }

    _, err = ds.DB.CreateTable(input)
    if err != nil {
        log.Println(err)
        return 
    }

    // Wait for the table to be created and active
    err = ds.WaitForStatus("ACTIVE")
    if err != nil {
        log.Println(err)
        return
    }
	return 
}


// GetStatus gets the status from the database
func (ds *DynamoStorage) GetStatus() (*Status, error){

	// Create a session we can use to communicate with dynamodb
	err := ds.CreateSession()
	if err != nil {
		log.Println("createSession error:", err)
		return nil,err
	}

	input := &dynamodb.GetItemInput{
        TableName: aws.String(ds.TableName),
        Key: map[string]*dynamodb.AttributeValue{
            "UUID": {S: aws.String(statusUUID)},
        },
    }

    result, err := ds.DB.GetItem(input)
    if err != nil {
        log.Println("GetItem error:",err)
        return nil, err
    }

    item := &Status{}
    err = dynamodbattribute.UnmarshalMap(result.Item, item)

    if err != nil {
		log.Println("Failed to unmarshal Record:", err)
		return nil,err
    }


	return item, nil
	
}