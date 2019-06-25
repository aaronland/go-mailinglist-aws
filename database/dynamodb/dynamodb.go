package dynamodb

import (
	"errors"
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

	return false, errors.New("Please write me")
}

func CreateConfirmationsTable(client *dynamodb.DynamoDB, table_name string) (bool, error) {

	has_table, err := hasTable(client, table_name)

	if err != nil {
		return false, err
	}

	if has_table {
		return true, nil
	}

	return false, errors.New("Please write me")
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
