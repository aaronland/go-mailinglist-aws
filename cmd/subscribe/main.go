package main

import (
	"flag"
	"log"
	"os"

	"github.com/aaronland/go-mailinglist-database-dynamodb"
	"github.com/aaronland/go-mailinglist/subscription"
)

func main() {

	dsn := flag.String("dsn", "", "...")
	addr := flag.String("address", "", "...")
	enabled := flag.Bool("enabled", false, "...")

	subs_table := flag.String("subscriptions-table", dynamodb.SUBSCRIPTIONS_DEFAULT_TABLENAME, "...")

	flag.Parse()

	opts := dynamodb.DefaultDynamoDBSubscriptionsDatabaseOptions()
	opts.TableName = *subs_table

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
		log.Println("SEND CONFIRMATION HERE...")
	}

	os.Exit(0)
}
