package main

import (
	"flag"
	"github.com/aaronland/go-mailinglist-aws/database/dynamodb"
	"log"
)

func main() {

	dsn := flag.String("dsn", "", "...")

	// table names here... or in dsn

	flag.Parse()

	subscribe_opts := dynamodb.DefaultDynamoDBSubscriptionsDatabaseOptions()
	confirm_opts := dynamodb.DefaultDynamoDBConfirmationsDatabaseOptions()

	var err error

	_, err = dynamodb.NewDynamoDBSubscriptionsDatabaseWithDSN(*dsn, subscribe_opts)

	if err != nil {
		log.Printf("Failed to set up %s table, %s\n", subscribe_opts.Table, err)
	}

	_, err = dynamodb.NewDynamoDBConfirmationsDatabaseWithDSN(*dsn, confirm_opts)

	if err != nil {
		log.Printf("Failed to set up %s table, %s\n", confirm_opts.Table, err)
	}

}
