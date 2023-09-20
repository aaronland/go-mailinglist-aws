package dynamodb

import (
	"context"
	"errors"

	"github.com/aaronland/go-aws-session"
	"github.com/aaronland/go-mailinglist/database"
	"github.com/aaronland/go-mailinglist/delivery"
	aws "github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"	
)

const DELIVERIES_DEFAULT_TABLENAME string = "deliveries"

type DynamoDBDeliveriesDatabaseOptions struct {
	TableName   string
	BillingMode string
	CreateTable bool
}

func DefaultDynamoDBDeliveriesDatabaseOptions() *DynamoDBDeliveriesDatabaseOptions {

	opts := DynamoDBDeliveriesDatabaseOptions{
		TableName:   DELIVERIES_DEFAULT_TABLENAME,
		BillingMode: "PAY_PER_REQUEST",
		CreateTable: false,
	}

	return &opts
}

type DynamoDBDeliveriesDatabase struct {
	database.DeliveriesDatabase
	client  *aws_dynamodb.DynamoDB
	options *DynamoDBDeliveriesDatabaseOptions
}

func NewDynamoDBDeliveriesDatabaseWithDSN(dsn string, opts *DynamoDBDeliveriesDatabaseOptions) (database.DeliveriesDatabase, error) {

	sess, err := session.NewSessionWithDSN(dsn)

	if err != nil {
		return nil, err
	}

	return NewDynamoDBDeliveriesDatabaseWithSession(sess, opts)
}

func NewDynamoDBDeliveriesDatabaseWithSession(sess *aws_session.Session, opts *DynamoDBDeliveriesDatabaseOptions) (database.DeliveriesDatabase, error) {

	client := aws_dynamodb.New(sess)

	if opts.CreateTable {

		_, err := CreateDeliveriesTable(client, opts)

		if err != nil {
			return nil, err
		}
	}

	db := DynamoDBDeliveriesDatabase{
		client:  client,
		options: opts,
	}

	return &db, nil
}

func (db *DynamoDBDeliveriesDatabase) GetDeliveryWithAddressAndMessageId(addr string, message_id string) (*delivery.Delivery, error) {

	req := &aws_dynamodb.GetItemInput{
		TableName: aws.String(db.options.TableName),
		Key: map[string]*aws_dynamodb.AttributeValue{
			"address": {
				S: aws.String(addr),
			},
			"message_id": {
				S: aws.String(message_id),
			},
		},
	}

	rsp, err := db.client.GetItem(req)

	if err != nil {
		return nil, err
	}

	return itemToDelivery(rsp.Item)
}

func (db *DynamoDBDeliveriesDatabase) AddDelivery(d *delivery.Delivery) error {

	existing_d, err := db.GetDeliveryWithAddressAndMessageId(d.Address, d.MessageId)

	if err != nil && !database.IsNotExist(err) {
		return err
	}

	if existing_d != nil {
		return errors.New("Delivery already exists")
	}

	return putDelivery(db.client, db.options, d)
}

/*
func (db *DynamoDBDeliveriesDatabase) ListDeliveries(ctx context.Context, callback database.ListDeliveriesFunc) error {

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

		return queryDeliveries(ctx, db.client, req, callback)

	req := &aws_dynamodb.ScanInput{
		TableName: aws.String(db.options.TableName),
	}

	return scanDeliveries(ctx, db.client, req, callback)
}
*/

func putDelivery(client *aws_dynamodb.DynamoDB, opts *DynamoDBDeliveriesDatabaseOptions, sub *delivery.Delivery) error {

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

func itemToDelivery(item map[string]*aws_dynamodb.AttributeValue) (*delivery.Delivery, error) {

	var sub *delivery.Delivery

	err := aws_dynamodbattribute.UnmarshalMap(item, &sub)

	if err != nil {
		return nil, err
	}

	if sub.Address == "" {
		return nil, new(database.NoRecordError)
	}

	return sub, nil
}

func scanDeliveries(ctx context.Context, client *aws_dynamodb.DynamoDB, req *aws_dynamodb.ScanInput, callback database.ListDeliveriesFunc) error {

	for {

		rsp, err := client.Scan(req)

		if err != nil {
			return err
		}

		for _, item := range rsp.Items {

			sub, err := itemToDelivery(item)

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
