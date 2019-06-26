package main

import (
	"errors"
	"flag"
	"github.com/aaronland/go-mailinglist-database-dynamodb"
	"log"
	"os"
)

func main() {

	dsn := flag.String("dsn", "", "...")
	addr := flag.String("address", "", "...")
	action := flag.String("action", "", "...")

	subs_table := flag.String("subscriptions-table", dynamodb.SUBSCRIPTIONS_DEFAULT_TABLENAME, "...")

	flag.Parse()

	opts := dynamodb.DefaultDynamoDBSubscriptionsDatabaseOptions()
	opts.TableName = *subs_table

	db, err := dynamodb.NewDynamoDBSubscriptionsDatabaseWithDSN(*dsn, opts)

	if err != nil {
		log.Fatal(err)
	}

	sub, err := db.GetSubscriptionWithAddress(*addr)

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

	err = db.UpdateSubscription(sub)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
