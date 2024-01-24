package dao

import (
	"errors"
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
	ItemID     string `json:"itemID"`
	BucketType string `json:"bucketType"`
	ItemType   string `json:"itemType"`
	ItemReason string `json:"itemReason"`
	CreatedOn  time.Time
}

func PutItem(item *Item) error {
	if item.ItemID == "" {
		item.ItemID = uuid.New().String()
	}
	item.CreatedOn = time.Now()

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

func GetItem(itemID string) (Item, error) {
	if itemID == "" {
		return Item{}, errors.New("Invalid ItemID")
	}

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("TLPDS"),
		Key: map[string]*dynamodb.AttributeValue{
			"itemID": {
				S: aws.String(itemID),
			},
		},
	})
	if err != nil {
		return Item{}, err
	}
	if result.Item == nil {
		return Item{}, nil
	}

	item := new(Item)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return Item{}, err
	}
	return *item, nil
}
