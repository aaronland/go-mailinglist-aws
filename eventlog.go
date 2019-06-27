package dynamodb

import (
	// "context"
	// "errors"
	"github.com/aaronland/go-aws-session"
	"github.com/aaronland/go-mailinglist/database"
	// "github.com/aaronland/go-mailinglist/eventlog"
	// aws "github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	// aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	_ "log"
)

const EVENTLOGS_DEFAULT_TABLENAME string = "eventlogs"

type DynamoDBEventLogsDatabaseOptions struct {
	TableName   string
	BillingMode string
	CreateTable bool
}

func DefaultDynamoDBEventLogsDatabaseOptions() *DynamoDBEventLogsDatabaseOptions {

	opts := DynamoDBEventLogsDatabaseOptions{
		TableName:   EVENTLOGS_DEFAULT_TABLENAME,
		BillingMode: "PAY_PER_REQUEST",
		CreateTable: false,
	}

	return &opts
}

type DynamoDBEventLogsDatabase struct {
	database.EventLogsDatabase
	client  *aws_dynamodb.DynamoDB
	options *DynamoDBEventLogsDatabaseOptions
}

func NewDynamoDBEventLogsDatabaseWithDSN(dsn string, opts *DynamoDBEventLogsDatabaseOptions) (database.EventLogsDatabase, error) {

	sess, err := session.NewSessionWithDSN(dsn)

	if err != nil {
		return nil, err
	}

	return NewDynamoDBEventLogsDatabaseWithSession(sess, opts)
}

func NewDynamoDBEventLogsDatabaseWithSession(sess *aws_session.Session, opts *DynamoDBEventLogsDatabaseOptions) (database.EventLogsDatabase, error) {

	client := aws_dynamodb.New(sess)

	if opts.CreateTable {
		_, err := CreateEventLogsTable(client, opts)

		if err != nil {
			return nil, err
		}
	}

	db := DynamoDBEventLogsDatabase{
		client:  client,
		options: opts,
	}

	return &db, nil
}
