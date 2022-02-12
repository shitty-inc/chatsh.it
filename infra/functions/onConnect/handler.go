package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	svc dynamodbiface.DynamoDBAPI
)

type Item struct {
	ConnectionId string
}

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	item := Item {
    ConnectionId: request.RequestContext.ConnectionID,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Error marshalling new item: %s", err)
	}

	input := &dynamodb.PutItemInput{
    Item:      av,
    TableName: aws.String("chatshit-9388b75"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Error calling PutItem: %s", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: "Connected",
	}, nil
}

func init() {
	session := session.Must(session.NewSession())

	svc = dynamodb.New(session)
}

func main() {
	lambda.Start(handler)
}
