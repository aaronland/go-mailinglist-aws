package dynamodb

import (
	"context"
	"errors"

	"github.com/aaronland/go-aws-session"
	"github.com/aaronland/go-mailinglist/confirmation"
	"github.com/aaronland/go-mailinglist/database"
	aws "github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const CONFIRMATIONS_DEFAULT_TABLENAME string = "confirmations"
const CONFIRMATIONS_DEFAULT_BILLINGMODE string = "PAY_PER_REQUEST"

type DynamoDBConfirmationsDatabaseOptions struct {
	TableName string
}

func DefaultDynamoDBConfirmationsDatabaseOptions() *DynamoDBConfirmationsDatabaseOptions {

	opts := DynamoDBConfirmationsDatabaseOptions{
		TableName: CONFIRMATIONS_DEFAULT_TABLENAME,
	}

	return &opts
}

type DynamoDBConfirmationsDatabase struct {
	database.ConfirmationsDatabase
	client  *aws_dynamodb.DynamoDB
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

	client := aws_dynamodb.New(sess)

	db := DynamoDBConfirmationsDatabase{
		client:  client,
		options: opts,
	}

	return &db, nil
}

func (db *DynamoDBConfirmationsDatabase) AddConfirmation(conf *confirmation.Confirmation) error {

	item, err := aws_dynamodbattribute.MarshalMap(conf)

	if err != nil {
		return err
	}

	req := &aws_dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(db.options.TableName),
	}

	_, err = db.client.PutItem(req)

	if err != nil {
		return err
	}

	return nil
}

func (db *DynamoDBConfirmationsDatabase) RemoveConfirmation(conf *confirmation.Confirmation) error {

	req := &aws_dynamodb.DeleteItemInput{
		TableName: aws.String(db.options.TableName),
		Key: map[string]*aws_dynamodb.AttributeValue{
			"code": {
				S: aws.String(conf.Code),
			},
		},
	}

	_, err := db.client.DeleteItem(req)

	if err != nil {
		return err
	}

	return nil
}

func (db *DynamoDBConfirmationsDatabase) GetConfirmationWithCode(code string) (*confirmation.Confirmation, error) {

	req := &aws_dynamodb.GetItemInput{
		TableName: aws.String(db.options.TableName),
		Key: map[string]*aws_dynamodb.AttributeValue{
			"code": {
				S: aws.String(code),
			},
		},
	}

	rsp, err := db.client.GetItem(req)

	if err != nil {
		return nil, err
	}

	var conf *confirmation.Confirmation

	err = aws_dynamodbattribute.UnmarshalMap(rsp.Item, &conf)

	if err != nil {
		return nil, err
	}

	if conf.Code == "" {
		return nil, new(database.NoRecordError)
	}

	return conf, nil
}

func (db *DynamoDBConfirmationsDatabase) ListConfirmations(ctx context.Context, callback database.ListConfirmationsFunc) error {
	return errors.New("Please write me")
}
