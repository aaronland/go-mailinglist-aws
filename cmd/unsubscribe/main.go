package main

import (
	"flag"
	"github.com/aaronland/go-mailinglist-database-dynamodb"
	"log"
	"os"
)

func main() {

	dsn := flag.String("dsn", "", "...")
	addr := flag.String("address", "", "...")

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

	log.Println("SUB", sub)

	os.Exit(0)
}
