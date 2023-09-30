package dynamodb

import (
	"context"
	"errors"
	"fmt"

	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	"github.com/aaronland/go-mailinglist/database"
	"github.com/aaronland/go-mailinglist/delivery"
	aws "github.com/aws/aws-sdk-go/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const DELIVERIES_DEFAULT_TABLENAME string = "deliveries"
const DELIVERIES_DEFAULT_BILLINGMODE string = "PAY_PER_REQUEST"

type DynamoDBDeliveriesDatabase struct {
	database.DeliveriesDatabase
	client *aws_dynamodb.DynamoDB
	table  string
}

func NewDynamoDBDeliveriesDatabase(ctx context.Context, uri string) (database.DeliveriesDatabase, error) {

	client, err := aa_dynamodb.NewClientWithURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create session, %w", err)
	}

	table := DELIVERIES_DEFAULT_TABLENAME

	db := DynamoDBDeliveriesDatabase{
		client: client,
		table:  table,
	}

	return &db, nil
}

func (db *DynamoDBDeliveriesDatabase) GetDeliveryWithAddressAndMessageId(ctx context.Context, addr string, message_id string) (*delivery.Delivery, error) {

	req := &aws_dynamodb.GetItemInput{
		TableName: aws.String(db.table),
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

func (db *DynamoDBDeliveriesDatabase) AddDelivery(ctx context.Context, d *delivery.Delivery) error {

	existing_d, err := db.GetDeliveryWithAddressAndMessageId(ctx, d.Address, d.MessageId)

	if err != nil && !database.IsNotExist(err) {
		return err
	}

	if existing_d != nil {
		return errors.New("Delivery already exists")
	}

	return db.putDelivery(ctx, d)
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
			TableName: aws.String(db.table),
		}

		return queryDeliveries(ctx, db.client, req, callback)

	req := &aws_dynamodb.ScanInput{
		TableName: aws.String(db.table),
	}

	return scanDeliveries(ctx, db.client, req, callback)
}
*/

func (db *DynamoDBDeliveriesDatabase) putDelivery(ctx context.Context, sub *delivery.Delivery) error {

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
