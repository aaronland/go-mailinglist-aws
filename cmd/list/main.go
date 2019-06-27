package main

import (
	"context"
	"errors"
	"flag"
	"github.com/aaronland/go-mailinglist-database-dynamodb"
	"github.com/aaronland/go-mailinglist/subscription"
	"log"
	"os"
)

func main() {

	dsn := flag.String("dsn", "", "...")
	str_status := flag.String("status", "", "...")

	subs_table := flag.String("subscriptions-table", dynamodb.SUBSCRIPTIONS_DEFAULT_TABLENAME, "...")

	flag.Parse()

	opts := dynamodb.DefaultDynamoDBSubscriptionsDatabaseOptions()
	opts.TableName = *subs_table

	db, err := dynamodb.NewDynamoDBSubscriptionsDatabaseWithDSN(*dsn, opts)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cb := func(sub *subscription.Subscription) error {
		log.Println(sub.Address)
		return nil
	}

	err = db.ListSubscriptionsWithStatus(ctx, cb, status)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
