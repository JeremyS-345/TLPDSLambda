package dao

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Add a book record to DynamoDB.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))

type Item struct {
	ItemID    int64     `json:"isbn"`
	CreatedOn time.Time `json:"title"`
}

func PutItem(item *Item) error {
	dynamoMappedItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}
	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("TLPDS"),
		Item:      dynamoMappedItem,
	})
	return err
}
