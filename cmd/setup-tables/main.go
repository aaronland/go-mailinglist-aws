package main

import (
	_ "flag"
	_ "log"

	_ "github.com/aaronland/go-mailinglist-database-dynamodb"
)

func main() {

	var client *aws_dynamodb.DynamoDB
	
	t, ok := auth_dynamodb.DynamoDBTables["accesstokens"]

	if !ok {
		log.Fatalf("Failed to derive access tokens definition")
	}

	opts := &aa_dynamodb.CreateTablesOptions{
		Tables: map[string]*aws_dynamodb.CreateTableInput{
			"accesstokens": t,
		},

		// Add all the go-mailinglist-database-dynamodb tables here
		Refresh: true,
	}

	err := aa_dynamodb.CreateTables(client, opts)

	if err != nil {
		log.Fatalf("Failed to create access tokens database, %v", err)
	}
	
}
