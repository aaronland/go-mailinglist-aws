package main

import (
	"context"
	"flag"
	"log"

	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	"github.com/aaronland/go-mailinglist-database-dynamodb"
)

func main() {

	client_uri := flag.String("client-uri", "", "...")
	refresh := flag.Bool("refresh", false, "...")
	flag.Parse()

	ctx := context.Background()

	client, err := aa_dynamodb.NewClientWithURI(ctx, *client_uri)

	if err != nil {
		log.Fatalf("Failed to create new client, %w", err)
	}

	opts := &aa_dynamodb.CreateTablesOptions{
		Tables:  dynamodb.DynamoDBTables,
		Refresh: *refresh,
	}

	err = aa_dynamodb.CreateTables(client, opts)

	if err != nil {
		log.Fatalf("Failed to create access tokens database, %v", err)
	}

}
