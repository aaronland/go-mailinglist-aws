package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/aaronland/go-mailinglist-database-dynamodb"
)

func main() {

	subs_uri := flag.String("subscriptions-uri", "", "...")
	addr := flag.String("address", "", "...")
	action := flag.String("action", "", "...")

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

	switch *action {

	case "enable":
		err = sub.Enable()
	case "disable":
		err = sub.Disable()
	case "block":
		err = sub.Block()
	case "unblock":
		err = sub.Unblock()
	case "confirm":
		err = sub.Confirm()
	case "unconfirm":
		err = sub.Unconfirm()
	default:
		err = errors.New("Invalid action")
	}

	if err != nil {
		log.Fatal(err)
	}

	err = db.UpdateSubscription(ctx, sub)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
