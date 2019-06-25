package dynamodb

// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-list-tables.html
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-create-table.html

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateSubscriptionsTable(client *dynamodb.DynamoDB, table_name string) (bool, error) {

	has_table, err := hasTable(client, table_name)

	if err != nil {
		return false, err
	}

	if has_table {
		return true, nil
	}

	req := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
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
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Address"),
				KeyType:       aws.String("S"),
			},
		},

		TableName: aws.String(table_name),
	}

	_, err = client.CreateTable(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateConfirmationsTable(client *dynamodb.DynamoDB, table_name string) (bool, error) {

	has_table, err := hasTable(client, table_name)

	if err != nil {
		return false, err
	}

	if has_table {
		return true, nil
	}

	req := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
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
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Code"),
				KeyType:       aws.String("S"),
			},
			{
				AttributeName: aws.String("Address"),
				KeyType:       aws.String("S"),
			},
		},

		TableName: aws.String(table_name),
	}

	_, err = client.CreateTable(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

func hasTable(client *dynamodb.DynamoDB, table string) (bool, error) {

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

func listTables(client *dynamodb.DynamoDB) ([]string, error) {

	tables := make([]string, 0)

	input := &dynamodb.ListTablesInput{}

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
