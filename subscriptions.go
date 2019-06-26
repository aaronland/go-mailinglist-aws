package dynamodb

import (
	"context"
	"errors"
	"github.com/aaronland/go-aws-session"
	"github.com/aaronland/go-mailinglist/database"
	"github.com/aaronland/go-mailinglist/subscription"
	aws "github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

const SUBSCRIPTIONS_DEFAULT_TABLENAME string = "subscriptions"

type DynamoDBSubscriptionsDatabaseOptions struct {
	TableName   string
	BillingMode string
	CreateTable bool
}

func DefaultDynamoDBSubscriptionsDatabaseOptions() *DynamoDBSubscriptionsDatabaseOptions {

	opts := DynamoDBSubscriptionsDatabaseOptions{
		TableName:   SUBSCRIPTIONS_DEFAULT_TABLENAME,
		BillingMode: "PAY_PER_REQUEST",
		CreateTable: false,
	}

	return &opts
}

type DynamoDBSubscriptionsDatabase struct {
	database.SubscriptionsDatabase
	client  *aws_dynamodb.DynamoDB
	options *DynamoDBSubscriptionsDatabaseOptions
}

func NewDynamoDBSubscriptionsDatabaseWithDSN(dsn string, opts *DynamoDBSubscriptionsDatabaseOptions) (database.SubscriptionsDatabase, error) {

	sess, err := session.NewSessionWithDSN(dsn)

	if err != nil {
		return nil, err
	}

	return NewDynamoDBSubscriptionsDatabaseWithSession(sess, opts)
}

func NewDynamoDBSubscriptionsDatabaseWithSession(sess *aws_session.Session, opts *DynamoDBSubscriptionsDatabaseOptions) (database.SubscriptionsDatabase, error) {

	client := aws_dynamodb.New(sess)

	if opts.CreateTable {
		_, err := CreateSubscriptionsTable(client, opts)

		if err != nil {
			return nil, err
		}
	}

	db := DynamoDBSubscriptionsDatabase{
		client:  client,
		options: opts,
	}

	return &db, nil
}

func (db *DynamoDBSubscriptionsDatabase) GetSubscriptionWithAddress(addr string) (*subscription.Subscription, error) {

	req := &aws_dynamodb.GetItemInput{
		TableName: aws.String(db.options.TableName),
		Key: map[string]*aws_dynamodb.AttributeValue{
			"address": {
				S: aws.String(addr),
			},
		},
	}

	log.Println(req)
	rsp, err := db.client.GetItem(req)

	if err != nil {
		return nil, err
	}

	var sub *subscription.Subscription

	err = aws_dynamodbattribute.UnmarshalMap(rsp.Item, &sub)

	if err != nil {
	   return nil, err
	}

	return sub, nil
}

func (db *DynamoDBSubscriptionsDatabase) AddSubscription(sub *subscription.Subscription) error {

     return PutItem(db.client, db.options, sub)
}

func (db *DynamoDBSubscriptionsDatabase) RemoveSubscription(sub *subscription.Subscription) error {
	return errors.New("Please write me")
}

func (db *DynamoDBSubscriptionsDatabase) UpdateSubscription(sub *subscription.Subscription) error {

     return PutItem(db.client, db.options, sub)
}

func (db *DynamoDBSubscriptionsDatabase) ListSubscriptionsConfirmed(ctx context.Context, callback database.ListSubscriptionsFunc) error {
	return errors.New("Please write me")
}

func (db *DynamoDBSubscriptionsDatabase) ListSubscriptionsUnconfirmed(ctx context.Context, callback database.ListSubscriptionsFunc) error {
	return errors.New("Please write me")
}
