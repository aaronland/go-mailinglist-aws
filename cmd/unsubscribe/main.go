package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aaronland/go-mailinglist-database-dynamodb"
)

func main() {

	subs_uri := flag.String("subscriptions-uri", "", "...")
	addr := flag.String("address", "", "...")

	flag.Parse()

	ctx := context.Background()

	db, err := dynamodb.NewDynamoDBSubscriptionsDatabase(ctx, *subs_uri)

	if err != nil {
		log.Fatal(err)
	}

	sub, err := db.GetSubscriptionWithAddress(ctx, *addr)

	if err != nil {
		log.Fatal(err)
	}

	err = db.RemoveSubscription(ctx, sub)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
