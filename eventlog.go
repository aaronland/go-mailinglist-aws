package dynamodb

import (
	"context"
	"fmt"

	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	"github.com/aaronland/go-mailinglist/database"
	"github.com/aaronland/go-mailinglist/eventlog"
	aws "github.com/aws/aws-sdk-go/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const EVENTLOGS_DEFAULT_TABLENAME string = "eventlogs"
const EVENTLOGS_DEFAULT_BILLINGMODE string = "PAY_PER_REQUEST"

type DynamoDBEventLogsDatabase struct {
	database.EventLogsDatabase
	client *aws_dynamodb.DynamoDB
	table  string
}

func init() {
	ctx := context.Background()
	database.RegisterEventLogsDatabase(ctx, "awsdynamodb", NewDynamoDBEventLogsDatabase)
}

func NewDynamoDBEventLogsDatabase(ctx context.Context, uri string) (database.EventLogsDatabase, error) {

	client, err := aa_dynamodb.NewClientWithURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create session, %w", err)
	}

	table := EVENTLOGS_DEFAULT_TABLENAME

	db := DynamoDBEventLogsDatabase{
		client: client,
		table:  table,
	}

	return &db, nil
}

func (db *DynamoDBEventLogsDatabase) AddEventLog(ctx context.Context, l *eventlog.EventLog) error {

	item, err := aws_dynamodbattribute.MarshalMap(l)

	if err != nil {
		return err
	}

	req := &aws_dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(db.table),
	}

	_, err = db.client.PutItem(req)

	if err != nil {
		return err
	}

	return nil
}
