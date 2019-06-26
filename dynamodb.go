package dynamodb

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GSI.html

// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-list-tables.html
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-create-table.html

// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#CreateTableInput
// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#AttributeDefinition
// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#KeySchemaElement

import (
	"github.com/aws/aws-sdk-go/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	aws_dynamodbattribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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
				AttributeName: aws.String("created"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("confirmed"),
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
				IndexName: aws.String("Confirmed"),
				KeySchema: []*aws_dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("confirmed"),
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
			/*
				{
					AttributeName: aws.String("created"),
					KeyType:       aws.String("RANGE"),
				},
			*/
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

func PutItem(client *aws_dynamodb.DynamoDB, opts *DynamoDBSubscriptionsDatabaseOptions, item interface{}) error {

	enc_item, err := aws_dynamodbattribute.MarshalMap(item)

	if err != nil {
		return err
	}

	req := &aws_dynamodb.PutItemInput{
		Item:      enc_item,
		TableName: aws.String(opts.TableName),
	}

	_, err = client.PutItem(req)

	if err != nil {
		return err
	}

	return nil
}
