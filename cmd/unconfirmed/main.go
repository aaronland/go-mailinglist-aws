package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-mailinglist-database-dynamodb"
	"github.com/aaronland/go-mailinglist/subscription"
	"log"
	"os"
)

func main() {

	dsn := flag.String("dsn", "", "...")

	subs_table := flag.String("subscriptions-table", dynamodb.SUBSCRIPTIONS_DEFAULT_TABLENAME, "...")

	flag.Parse()

	opts := dynamodb.DefaultDynamoDBSubscriptionsDatabaseOptions()
	opts.TableName = *subs_table

	db, err := dynamodb.NewDynamoDBSubscriptionsDatabaseWithDSN(*dsn, opts)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cb := func(sub *subscription.Subscription) error {
		log.Println(sub.Address)
		return nil
	}

	err = db.ListSubscriptionsUnconfirmed(ctx, cb)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
