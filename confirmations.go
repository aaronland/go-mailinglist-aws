package dynamodb

import (
	"context"
	"errors"
	"github.com/aaronland/go-aws-session"
	"github.com/aaronland/go-mailinglist/confirmation"
	"github.com/aaronland/go-mailinglist/database"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBConfirmationsDatabaseOptions struct {
	TableName   string
	BillingMode string
}

func DefaultDynamoDBConfirmationsDatabaseOptions() *DynamoDBConfirmationsDatabaseOptions {

	opts := DynamoDBConfirmationsDatabaseOptions{
		TableName:   "confirmations",
		BillingMode: "PAY_PER_REQUEST",
	}

	return &opts
}

type DynamoDBConfirmationsDatabase struct {
	database.ConfirmationsDatabase
	client  *dynamodb.DynamoDB
	options *DynamoDBConfirmationsDatabaseOptions
}

func NewDynamoDBConfirmationsDatabaseWithDSN(dsn string, opts *DynamoDBConfirmationsDatabaseOptions) (database.ConfirmationsDatabase, error) {

	sess, err := session.NewSessionWithDSN(dsn)

	if err != nil {
		return nil, err
	}

	return NewDynamoDBConfirmationsDatabaseWithSession(sess, opts)
}

func NewDynamoDBConfirmationsDatabaseWithSession(sess *aws_session.Session, opts *DynamoDBConfirmationsDatabaseOptions) (database.ConfirmationsDatabase, error) {

	client := dynamodb.New(sess)

	_, err := CreateConfirmationsTable(client, opts)

	if err != nil {
		return nil, err
	}

	db := DynamoDBConfirmationsDatabase{
		client:  client,
		options: opts,
	}

	return &db, nil
}

func (db *DynamoDBConfirmationsDatabase) AddConfirmation(conf *confirmation.Confirmation) error {
	return errors.New("Please write me")
}

func (db *DynamoDBConfirmationsDatabase) RemoveConfirmation(conf *confirmation.Confirmation) error {
	return errors.New("Please write me")
}

func (db *DynamoDBConfirmationsDatabase) GetConfirmationWithCode(code string) (*confirmation.Confirmation, error) {
	return nil, errors.New("Please write me")
}

func (db *DynamoDBConfirmationsDatabase) ListConfirmations(ctx context.Context, callback database.ListConfirmationsFunc) error {
	return errors.New("Please write me")
}
