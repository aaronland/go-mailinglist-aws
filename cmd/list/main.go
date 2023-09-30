package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/aaronland/go-mailinglist-database-dynamodb"
	"github.com/aaronland/go-mailinglist/subscription"
)

func main() {

	subs_uri := flag.String("subscriptions-uri", "", "...")
	str_status := flag.String("status", "", "...")

	flag.Parse()

	ctx := context.Background()

	db, err := dynamodb.NewDynamoDBSubscriptionsDatabase(ctx, *subs_uri)

	if err != nil {
		log.Fatal(err)
	}

	status := -1

	switch *str_status {
	case "pending":
		status = subscription.SUBSCRIPTION_STATUS_PENDING
	case "enabled":
		status = subscription.SUBSCRIPTION_STATUS_ENABLED
	case "disabled":
		status = subscription.SUBSCRIPTION_STATUS_DISABLED
	case "blocked":
		status = subscription.SUBSCRIPTION_STATUS_BLOCKED
	default:
		err = errors.New("Invalid status")
	}

	if err != nil {
		log.Fatal(err)
	}

	cb := func(ctx context.Context, sub *subscription.Subscription) error {
		log.Println(sub.Address)
		return nil
	}

	err = db.ListSubscriptionsWithStatus(ctx, cb, status)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
