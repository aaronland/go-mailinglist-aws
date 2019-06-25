package dynamodb

import (
	"context"
	"errors"
	"github.com/aaronland/go-aws-session"
	"github.com/aaronland/go-mailinglist/database"
	"github.com/aaronland/go-mailinglist/subscription"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBSubscriptionsDatabaseOptions struct {
	Table string
}

func DefaultDynamoDBSubscriptionsDatabaseOptions() *DynamoDBSubscriptionsDatabaseOptions {

	opts := DynamoDBSubscriptionsDatabaseOptions{
		Table: "subscriptions",
	}

	return &opts
}

type DynamoDBSubscriptionsDatabase struct {
	database.SubscriptionsDatabase
	client  *dynamodb.DynamoDB
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

	client := dynamodb.New(sess)

	_, err := CreateSubscriptionsTable(client, opts.Table)

	if err != nil {
		return nil, err
	}

	db := DynamoDBSubscriptionsDatabase{
		client: client,
	}

	return &db, nil
}

func (db *DynamoDBSubscriptionsDatabase) AddSubscription(sub *subscription.Subscription) error {
	return errors.New("Please write me")
}

func (db *DynamoDBSubscriptionsDatabase) RemoveSubscription(sub *subscription.Subscription) error {
	return errors.New("Please write me")
}

func (db *DynamoDBSubscriptionsDatabase) UpdateSubscription(sub *subscription.Subscription) error {
	return errors.New("Please write me")
}

func (db *DynamoDBSubscriptionsDatabase) GetSubscriptionWithAddress(addr string) (*subscription.Subscription, error) {
	return nil, errors.New("Please write me")
}

func (db *DynamoDBSubscriptionsDatabase) ListSubscriptionsConfirmed(ctx context.Context, callback database.ListSubscriptionsFunc) error {
	return errors.New("Please write me")
}

func (db *DynamoDBSubscriptionsDatabase) ListSubscriptionsUnconfirmed(ctx context.Context, callback database.ListSubscriptionsFunc) error {
	return errors.New("Please write me")
}
