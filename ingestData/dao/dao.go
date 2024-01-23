package dao

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

// Add a book record to DynamoDB.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))

type Item struct {
	ItemID     string  `json:"itemID"`
	ItemType   string `json:itemType`
	ItemReason string `json:itemReason`
	CreatedOn  time.Time
}

func PutItem(item *Item) error {
	if item.ItemID == "" {
		item.ItemID = uuid.New().String()
	}

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
