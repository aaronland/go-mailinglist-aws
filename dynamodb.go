package dynamodb

// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-list-tables.html
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-create-table.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html#HowItWorks.CoreComponents.TablesItemsAttributes
// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#CreateTableInput
// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#AttributeDefinition

import (
	"github.com/aws/aws-sdk-go/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
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
				AttributeName: aws.String("Address"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Created"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Confirmed"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Status"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*aws_dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Address"),
				KeyType:       aws.String("HASH"),
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
				AttributeName: aws.String("Code"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Action"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Created"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Status"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*aws_dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Code"),
				KeyType:       aws.String("HASH"),
			},
		},
		/*
			GlobalSecondaryIndexes: []*aws_dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String("Address"),
					KeySchema:       aws.String("HASH"),
				},
			},
		*/
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
