package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aaronland/go-mailinglist-database-dynamodb"
	"github.com/aaronland/go-mailinglist/subscription"
)

func main() {

	subs_uri := flag.String("subscriptions-uri", "", "...")
	addr := flag.String("address", "", "...")
	enabled := flag.Bool("enabled", false, "...")

	flag.Parse()

	ctx := context.Background()

	db, err := dynamodb.NewDynamoDBSubscriptionsDatabase(ctx, *subs_uri)

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

	err = db.AddSubscription(ctx, sub)

	if err != nil {
		log.Fatal(err)
	}

	if !sub.IsConfirmed() {
		log.Println("SEND CONFIRMATION HERE...")
	}

	os.Exit(0)
}
