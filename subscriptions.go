package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	"github.com/aaronland/go-mailinglist/database"
	"github.com/aaronland/go-mailinglist/subscription"
	aws "github.com/aws/aws-sdk-go/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const SUBSCRIPTIONS_DEFAULT_TABLENAME string = "subscriptions"
const SUBSCRIPTIONS_DEFAULT_BILLINGMODE string = "PAY_PER_REQUEST"

type DynamoDBSubscriptionsDatabase struct {
	database.SubscriptionsDatabase
	client *aws_dynamodb.DynamoDB
	table  string
}

func init() {
	ctx := context.Background()
	database.RegisterSubscriptionsDatabase(ctx, "awsdynamodb", NewDynamoDBSubscriptionsDatabase)
}

func NewDynamoDBSubscriptionsDatabase(ctx context.Context, uri string) (database.SubscriptionsDatabase, error) {

	client, err := aa_dynamodb.NewClientWithURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create session, %w", err)
	}

	table := SUBSCRIPTIONS_DEFAULT_TABLENAME

	db := DynamoDBSubscriptionsDatabase{
		client: client,
		table:  table,
	}

	return &db, nil
}

func (db *DynamoDBSubscriptionsDatabase) GetSubscriptionWithAddress(ctx context.Context, addr string) (*subscription.Subscription, error) {

	req := &aws_dynamodb.GetItemInput{
		TableName: aws.String(db.table),
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

func (db *DynamoDBSubscriptionsDatabase) AddSubscription(ctx context.Context, sub *subscription.Subscription) error {

	existing_sub, err := db.GetSubscriptionWithAddress(ctx, sub.Address)

	if err != nil && !database.IsNotExist(err) {
		return err
	}

	if existing_sub != nil {
		return errors.New("Subscription already exists")
	}

	return db.putSubscription(ctx, sub)
}

func (db *DynamoDBSubscriptionsDatabase) RemoveSubscription(ctx context.Context, sub *subscription.Subscription) error {

	req := &aws_dynamodb.DeleteItemInput{
		TableName: aws.String(db.table),
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

func (db *DynamoDBSubscriptionsDatabase) UpdateSubscription(ctx context.Context, sub *subscription.Subscription) error {

	return db.putSubscription(ctx, sub)
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
			TableName: aws.String(db.table),
		}

		return querySubscriptions(ctx, db.client, req, callback)
	*/

	req := &aws_dynamodb.ScanInput{
		TableName: aws.String(db.table),
	}

	return db.scanSubscriptions(ctx, req, callback)
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
		TableName:            aws.String(db.table),
	}

	return db.scanSubscriptions(ctx, req, callback)
}

func (db *DynamoDBSubscriptionsDatabase) putSubscription(ctx context.Context, sub *subscription.Subscription) error {

	item, err := aws_dynamodbattribute.MarshalMap(sub)

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

			err = callback(ctx, sub)

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

func (db *DynamoDBSubscriptionsDatabase) scanSubscriptions(ctx context.Context, req *aws_dynamodb.ScanInput, callback database.ListSubscriptionsFunc) error {

	for {

		rsp, err := db.client.Scan(req)

		if err != nil {
			return err
		}

		for _, item := range rsp.Items {

			sub, err := itemToSubscription(item)

			if err != nil {
				return err
			}

			err = callback(ctx, sub)

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
