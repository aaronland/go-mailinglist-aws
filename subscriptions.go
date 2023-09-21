package dynamodb

import (
	"context"
	"errors"
	"strconv"

	"github.com/aaronland/go-aws-session"
	"github.com/aaronland/go-mailinglist/database"
	"github.com/aaronland/go-mailinglist/subscription"
	aws "github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const SUBSCRIPTIONS_DEFAULT_TABLENAME string = "subscriptions"
const SUBSCRIPTIONS_DEFAULT_BILLINGMODE string = "PAY_PER_REQUEST"

type DynamoDBSubscriptionsDatabaseOptions struct {
	TableName   string
	BillingMode string
	CreateTable bool
}

func DefaultDynamoDBSubscriptionsDatabaseOptions() *DynamoDBSubscriptionsDatabaseOptions {

	opts := DynamoDBSubscriptionsDatabaseOptions{
		TableName:   SUBSCRIPTIONS_DEFAULT_TABLENAME,
		BillingMode: SUBSCRIPTIONS_DEFAULT_BILLINGMODE,
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

	rsp, err := db.client.GetItem(req)

	if err != nil {
		return nil, err
	}

	return itemToSubscription(rsp.Item)
}

func (db *DynamoDBSubscriptionsDatabase) AddSubscription(sub *subscription.Subscription) error {

	existing_sub, err := db.GetSubscriptionWithAddress(sub.Address)

	if err != nil && !database.IsNotExist(err) {
		return err
	}

	if existing_sub != nil {
		return errors.New("Subscription already exists")
	}

	return putSubscription(db.client, db.options, sub)
}

func (db *DynamoDBSubscriptionsDatabase) RemoveSubscription(sub *subscription.Subscription) error {

	req := &aws_dynamodb.DeleteItemInput{
		TableName: aws.String(db.options.TableName),
		Key: map[string]*aws_dynamodb.AttributeValue{
			"address": {
				S: aws.String(sub.Address),
			},
		},
	}

	_, err := db.client.DeleteItem(req)

	if err != nil {
		return err
	}

	return nil
}

func (db *DynamoDBSubscriptionsDatabase) UpdateSubscription(sub *subscription.Subscription) error {

	return putSubscription(db.client, db.options, sub)
}

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBMapper.QueryScanExample.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Scan.html#Scan.FilterExpression
// https://github.com/markuscraig/dynamodb-examples/blob/master/go/movies_scan.go
// https://github.com/markuscraig/dynamodb-examples/blob/master/go/movies_query_year.go

func (db *DynamoDBSubscriptionsDatabase) ListSubscriptions(ctx context.Context, callback database.ListSubscriptionsFunc) error {

	/*
		req := &aws_dynamodb.QueryInput{
			ProjectionExpression: aws.String("#confirmed, address"),
			KeyConditionExpression: aws.String("#confirmed > :zero"),
			ExpressionAttributeNames: map[string]*string{
				"#confirmed": aws.String("confirmed"),
			},
			ExpressionAttributeValues: map[string]*aws_dynamodb.AttributeValue{
				":zero": {
					N: aws.String("0"),
				},
			},
			TableName: aws.String(db.options.TableName),
		}

		return querySubscriptions(ctx, db.client, req, callback)
	*/

	req := &aws_dynamodb.ScanInput{
		TableName: aws.String(db.options.TableName),
	}

	return scanSubscriptions(ctx, db.client, req, callback)
}

func (db *DynamoDBSubscriptionsDatabase) ListSubscriptionsWithStatus(ctx context.Context, callback database.ListSubscriptionsFunc, status ...int) error {

	if len(status) == 0 {
		return errors.New("Missing status(es)")
	}

	if len(status) > 1 {
		return errors.New("Multiple status(es) are not supported yet.")
	}

	// only supporting one status is not a feature - it just hasn't been implemented yet...
	// (20190627/thisisaaronland)

	state := status[0]

	str_state := strconv.Itoa(state)

	req := &aws_dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"),
		},
		ExpressionAttributeValues: map[string]*aws_dynamodb.AttributeValue{
			":state": {
				N: aws.String(str_state),
			},
		},
		FilterExpression:     aws.String("#status = :state"),
		ProjectionExpression: aws.String("#status, address"),
		TableName:            aws.String(db.options.TableName),
	}

	return scanSubscriptions(ctx, db.client, req, callback)
}

func putSubscription(client *aws_dynamodb.DynamoDB, opts *DynamoDBSubscriptionsDatabaseOptions, sub *subscription.Subscription) error {

	item, err := aws_dynamodbattribute.MarshalMap(sub)

	if err != nil {
		return err
	}

	req := &aws_dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(opts.TableName),
	}

	_, err = client.PutItem(req)

	if err != nil {
		return err
	}

	return nil
}

func itemToSubscription(item map[string]*aws_dynamodb.AttributeValue) (*subscription.Subscription, error) {

	var sub *subscription.Subscription

	err := aws_dynamodbattribute.UnmarshalMap(item, &sub)

	if err != nil {
		return nil, err
	}

	if sub.Address == "" {
		return nil, new(database.NoRecordError)
	}

	return sub, nil
}

func querySubscriptions(ctx context.Context, client *aws_dynamodb.DynamoDB, req *aws_dynamodb.QueryInput, callback database.ListSubscriptionsFunc) error {

	for {

		rsp, err := client.Query(req)

		if err != nil {
			return err
		}

		for _, item := range rsp.Items {

			sub, err := itemToSubscription(item)

			if err != nil {
				return err
			}

			err = callback(sub)

			if err != nil {
				return err
			}
		}

		req.ExclusiveStartKey = rsp.LastEvaluatedKey

		if rsp.LastEvaluatedKey == nil {
			break
		}
	}

	return nil
}

func scanSubscriptions(ctx context.Context, client *aws_dynamodb.DynamoDB, req *aws_dynamodb.ScanInput, callback database.ListSubscriptionsFunc) error {

	for {

		rsp, err := client.Scan(req)

		if err != nil {
			return err
		}

		for _, item := range rsp.Items {

			sub, err := itemToSubscription(item)

			if err != nil {
				return err
			}

			err = callback(sub)

			if err != nil {
				return err
			}
		}

		req.ExclusiveStartKey = rsp.LastEvaluatedKey

		if rsp.LastEvaluatedKey == nil {
			break
		}
	}

	return nil
}
