package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynamo dynamodbiface.DynamoDBAPI
)

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.DeleteItemInput{
    TableName: aws.String("chatshit-9388b75"),
		Key: map[string]*dynamodb.AttributeValue {
			"ConnectionId": {
				S: aws.String(request.RequestContext.ConnectionID),
			},
		},
	}

	_, err := dynamo.DeleteItem(input)
	if err != nil {
		log.Fatalf("Error calling DeleteItem: %s", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: "Disconnected",
	}, nil
}

func init() {
	session := session.Must(session.NewSession())

	dynamo = dynamodb.New(session)
}

func main() {
	lambda.Start(handler)
}
