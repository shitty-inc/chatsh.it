package main

import (
	"context"
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	svc dynamodbiface.DynamoDBAPI
	session dynamodbiface.DynamoDBAPI
)

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc = dynamodb.New(session)

	callbackURL := url.URL{
		Scheme: "https",
		Host:   request.RequestContext.DomainName,
		Path:   request.RequestContext.Stage,
	}

	api := apigatewaymanagementapi.New(func(o *apigatewaymanagementapi.Options) {
		o.EndpointResolver = apigatewaymanagementapi.EndpointResolverFromURL(callbackURL.String())
	})

	_, err := api.PostToConnection(context.Background(), &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte("Test Post to Connection"),
	})
	if err != nil {
		log.Fatalf("Error calling PostToConnection: %s", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: "Sent",
	}, nil
}

func init() {
	session = session.Must(session.NewSession())

	svc = dynamodb.New(session)
}

func main() {
	lambda.Start(handler)
}
