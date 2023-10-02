GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

# https://aws.amazon.com/about-aws/whats-new/2018/08/use-amazon-dynamodb-local-more-easily-with-the-new-docker-image/
# https://hub.docker.com/r/amazon/dynamodb-local/

dynamo-local:
	docker run --rm -it -p 8000:8000 amazon/dynamodb-local

# Set up accesstokens and mailing list-related table(s)

dynamo-tables-local:
	go run -mod $(GOMOD) \
		cmd/setup-tables/main.go \
		-refresh \
		-client-uri 'awsdynamodb://?local=true'
