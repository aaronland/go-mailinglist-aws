package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
)

func SubscriptionTables() map[string]*aws_dynamodb.CreateTableInput {

	tables := map[string]*aws_dynamodb.CreateTableInput{
		"deliveries": &aws_dynamodb.CreateTableInput{
			AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("address"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("message_id"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*aws_dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("address"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("message_id"),
					KeyType:       aws.String("RANGE"),
				},
			},
			GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String("status"),
					KeySchema: []*aws_dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("message_id"),
							KeyType:       aws.String("HASH"),
						},
					},
					Projection: &aws_dynamodb.Projection{
						// maybe just address...?
						ProjectionType: aws.String("ALL"),
					},
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
		},
		"confirmations": &aws_dynamodb.CreateTableInput{
			AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("code"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("address"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("created"),
					AttributeType: aws.String("N"),
				},
			},
			KeySchema: []*aws_dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("code"),
					KeyType:       aws.String("HASH"),
				},
			},
			GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String("address"),
					KeySchema: []*aws_dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("address"),
							KeyType:       aws.String("HASH"),
						},
					},
					Projection: &aws_dynamodb.Projection{
						ProjectionType: aws.String("ALL"),
					},
				},
				{
					IndexName: aws.String("created"),
					KeySchema: []*aws_dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("created"),
							KeyType:       aws.String("HASH"),
						},
					},
					Projection: &aws_dynamodb.Projection{
						ProjectionType: aws.String("INCLUDE"),
						NonKeyAttributes: []*string{
							aws.String("code"),
						},
					},
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
		},
		"eventlogs": &aws_dynamodb.CreateTableInput{
			AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("address"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("event"),
					AttributeType: aws.String("N"),
				},
				{
					AttributeName: aws.String("created"),
					AttributeType: aws.String("N"),
				},
			},
			KeySchema: []*aws_dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("address"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("created"),
					KeyType:       aws.String("RANGE"),
				},
			},
			GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String("address"),
					KeySchema: []*aws_dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("address"),
							KeyType:       aws.String("HASH"),
						},
					},
					Projection: &aws_dynamodb.Projection{
						// maybe just address...?
						ProjectionType: aws.String("ALL"),
					},
				},
				{
					IndexName: aws.String("event"),
					KeySchema: []*aws_dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("event"),
							KeyType:       aws.String("HASH"),
						},
					},
					Projection: &aws_dynamodb.Projection{
						// maybe just address...?
						ProjectionType: aws.String("ALL"),
					},
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
		},
		"subscriptions": &aws_dynamodb.CreateTableInput{
			AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("address"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("status"),
					AttributeType: aws.String("N"),
				},
			},
			KeySchema: []*aws_dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("address"),
					KeyType:       aws.String("HASH"),
				},
			},
			GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String("status"),
					KeySchema: []*aws_dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("status"),
							KeyType:       aws.String("HASH"),
						},
					},
					Projection: &aws_dynamodb.Projection{
						// maybe just address...?
						ProjectionType: aws.String("ALL"),
					},
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
		},
	}

	return tables
}

/*
func CreateSubscriptionsTable(client *aws_dynamodb.DynamoDB, opts *DynamoDBSubscriptionsDatabaseOptions) (bool, error) {

	has_table, err := hasTable(client, opts.TableName)

	if err != nil {
		return false, err
	}

	if has_table {
		return true, nil
	}

	req := &aws_dynamodb.CreateTableInput{
		AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("address"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("status"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*aws_dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("address"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("status"),
				KeySchema: []*aws_dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("status"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &aws_dynamodb.Projection{
					// maybe just address...?
					ProjectionType: aws.String("ALL"),
				},
			},
		},
		BillingMode: aws.String(opts.BillingMode),
		TableName:   aws.String(opts.TableName),
	}

	_, err = client.CreateTable(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateEventLogsTable(client *aws_dynamodb.DynamoDB, opts *DynamoDBEventLogsDatabaseOptions) (bool, error) {

	has_table, err := hasTable(client, opts.TableName)

	if err != nil {
		return false, err
	}

	if has_table {
		return true, nil
	}

	req := &aws_dynamodb.CreateTableInput{
		AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("address"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("event"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("created"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*aws_dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("address"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("created"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("address"),
				KeySchema: []*aws_dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("address"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &aws_dynamodb.Projection{
					// maybe just address...?
					ProjectionType: aws.String("ALL"),
				},
			},
			{
				IndexName: aws.String("event"),
				KeySchema: []*aws_dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("event"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &aws_dynamodb.Projection{
					// maybe just address...?
					ProjectionType: aws.String("ALL"),
				},
			},
		},
		BillingMode: aws.String(opts.BillingMode),
		TableName:   aws.String(opts.TableName),
	}

	_, err = client.CreateTable(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateConfirmationsTable(client *aws_dynamodb.DynamoDB, opts *DynamoDBConfirmationsDatabaseOptions) (bool, error) {

	has_table, err := hasTable(client, opts.TableName)

	if err != nil {
		return false, err
	}

	if has_table {
		return true, nil
	}

	req := &aws_dynamodb.CreateTableInput{
		AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("code"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("address"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("created"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*aws_dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("code"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("address"),
				KeySchema: []*aws_dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("address"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &aws_dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
			},
			{
				IndexName: aws.String("created"),
				KeySchema: []*aws_dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("created"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &aws_dynamodb.Projection{
					ProjectionType: aws.String("INCLUDE"),
					NonKeyAttributes: []*string{
						aws.String("code"),
					},
				},
			},
		},
		BillingMode: aws.String(opts.BillingMode),
		TableName:   aws.String(opts.TableName),
	}

	_, err = client.CreateTable(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateDeliveriesTable(client *aws_dynamodb.DynamoDB, opts *DynamoDBDeliveriesDatabaseOptions) (bool, error) {

	has_table, err := hasTable(client, opts.TableName)

	if err != nil {
		return false, err
	}

	if has_table {
		return true, nil
	}

	req := &aws_dynamodb.CreateTableInput{
		AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("address"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("message_id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*aws_dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("address"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("message_id"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("status"),
				KeySchema: []*aws_dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("message_id"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &aws_dynamodb.Projection{
					// maybe just address...?
					ProjectionType: aws.String("ALL"),
				},
			},
		},
		BillingMode: aws.String(opts.BillingMode),
		TableName:   aws.String(opts.TableName),
	}

	_, err = client.CreateTable(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

func hasTable(client *aws_dynamodb.DynamoDB, table string) (bool, error) {

	tables, err := listTables(client)

	if err != nil {
		return false, err
	}

	has_table := false

	for _, name := range tables {

		if name == table {
			has_table = true
			break
		}
	}

	return has_table, nil
}

func listTables(client *aws_dynamodb.DynamoDB) ([]string, error) {

	tables := make([]string, 0)

	input := &aws_dynamodb.ListTablesInput{}

	for {

		rsp, err := client.ListTables(input)

		if err != nil {
			return nil, err
		}

		for _, n := range rsp.TableNames {
			tables = append(tables, *n)
		}

		input.ExclusiveStartTableName = rsp.LastEvaluatedTableName

		if rsp.LastEvaluatedTableName == nil {
			break
		}
	}

	return tables, nil
}

*/
