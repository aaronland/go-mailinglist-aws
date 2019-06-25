package main

import (
	"flag"
	"github.com/aaronland/go-mailinglist-aws/database/dynamodb"
	"github.com/aaronland/go-mailinglist/subscription"
	"log"
)

func main() {

	dsn := flag.String("dsn", "", "...")

	addr := flag.String("address", "", "...")
	enabled := flag.Bool("enabled", false, "...")

	flag.Parse()

	opts := dynamodb.DefaultDynamoDBSubscriptionsDatabaseOptions()

	db, err := dynamodb.NewDynamoDBSubscriptionsDatabaseWithDSN(*dsn, opts)

	if err != nil {
		log.Fatal(err)
	}

	sub, err := subscription.NewSubscription(*addr)

	if err != nil {
		log.Fatal(err)
	}

	if *enabled {
		sub.Confirm()
		sub.Enable()
	}

	err = db.AddSubscription(sub)

	if err != nil {
		log.Fatal(err)
	}

	if !sub.IsConfirmed() {
		// send confirmation code...
	}
}
