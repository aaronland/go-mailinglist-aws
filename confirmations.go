package dynamodb

import (
	"context"
	"errors"
	"fmt"

	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	"github.com/aaronland/go-mailinglist/confirmation"
	"github.com/aaronland/go-mailinglist/database"
	aws "github.com/aws/aws-sdk-go/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const CONFIRMATIONS_DEFAULT_TABLENAME string = "confirmations"
const CONFIRMATIONS_DEFAULT_BILLINGMODE string = "PAY_PER_REQUEST"

type DynamoDBConfirmationsDatabase struct {
	database.ConfirmationsDatabase
	client *aws_dynamodb.DynamoDB
	table  string
}

func init() {
	ctx := context.Background()
	database.RegisterConfirmationsDatabase(ctx, "awsdynamodb", NewDynamoDBConfirmationsDatabase)
}

func NewDynamoDBConfirmationsDatabase(ctx context.Context, uri string) (database.ConfirmationsDatabase, error) {

	client, err := aa_dynamodb.NewClientWithURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create session, %w", err)
	}

	table := CONFIRMATIONS_DEFAULT_TABLENAME

	db := DynamoDBConfirmationsDatabase{
		client: client,
		table:  table,
	}

	return &db, nil
}

func (db *DynamoDBConfirmationsDatabase) AddConfirmation(ctx context.Context, conf *confirmation.Confirmation) error {

	item, err := aws_dynamodbattribute.MarshalMap(conf)

	if err != nil {
		return fmt.Errorf("Failed to marshal confirmation, %w", err)
	}

	req := &aws_dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(db.table),
	}

	_, err = db.client.PutItem(req)

	if err != nil {
		return fmt.Errorf("Failed to put confirmation, %w", err)
	}

	return nil
}

func (db *DynamoDBConfirmationsDatabase) RemoveConfirmation(ctx context.Context, conf *confirmation.Confirmation) error {

	req := &aws_dynamodb.DeleteItemInput{
		TableName: aws.String(db.table),
		Key: map[string]*aws_dynamodb.AttributeValue{
			"code": {
				S: aws.String(conf.Code),
			},
		},
	}

	_, err := db.client.DeleteItem(req)

	if err != nil {
		return fmt.Errorf("Failed to delete confirmation, %w", err)
	}

	return nil
}

func (db *DynamoDBConfirmationsDatabase) GetConfirmationWithCode(ctx context.Context, code string) (*confirmation.Confirmation, error) {

	req := &aws_dynamodb.GetItemInput{
		TableName: aws.String(db.table),
		Key: map[string]*aws_dynamodb.AttributeValue{
			"code": {
				S: aws.String(code),
			},
		},
	}

	rsp, err := db.client.GetItem(req)

	if err != nil {
		return nil, fmt.Errorf("Failed to get confirmation, %w", err)
	}

	var conf *confirmation.Confirmation

	err = aws_dynamodbattribute.UnmarshalMap(rsp.Item, &conf)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal confirmation, %w", err)
	}

	if conf.Code == "" {
		return nil, new(database.NoRecordError)
	}

	return conf, nil
}

func (db *DynamoDBConfirmationsDatabase) ListConfirmations(ctx context.Context, callback database.ListConfirmationsFunc) error {
	return errors.New("Please write me")
}
