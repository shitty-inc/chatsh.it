package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynamo dynamodbiface.DynamoDBAPI
	api *apigatewaymanagementapi.ApiGatewayManagementApi
)

type DynamoItem struct {
	Id string
	ConnectionId string
}

type Message struct {
	Action    string `json:"action"`
	Payload map[string]interface{} `json:"payload"`
}

// Register a new ConnectionId with the given ID
func register(ConnectionId string, Id string) error {
	item, err := dynamodbattribute.MarshalMap(DynamoItem {
		Id: Id,
    ConnectionId: ConnectionId,
	})
	if err != nil {
		log.Fatalf("Error marshalling new item: %s", err)
	}

	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
    Item:      item,
    TableName: aws.String("chatshit-ee079db"),
	})
	if err != nil {
		log.Fatalf("Error calling PutItem: %s", err)
	}

	json, err := json.Marshal(&Message{
		Action: "registered",
		Payload: map[string]interface{}{
			"Id": Id,
			"ConnectionId": ConnectionId,
		},
	})
	if err != nil {
		log.Fatalf("Error marshalling message: %s", err)
	}

	_, err = api.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(ConnectionId),
		Data:         []byte(json),
	})
	if err != nil {
		log.Fatalf("Error calling PostToConnection: %s", err)
	}

	return err
}

// Send a message to all connections from an ID
func send(Id string, message Message) error {
	input := &dynamodb.QueryInput{
		TableName: aws.String("chatshit-ee079db"),
		KeyConditionExpression: aws.String("Id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(Id),
			},
		},
	}

	result, err := dynamo.Query(input)
	if err != nil {
		log.Fatalf("Error calling Query: %s", err)
	}

	log.Println("Found connections", result.Items)

	for _, i := range result.Items {
		item := DynamoItem{}

		if err := dynamodbattribute.UnmarshalMap(i, &item); err != nil {
			log.Fatalf("Error unmarshalling item: %s", err)
		}

		json, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Error marshalling message: %s", err)
		}

		_, err = api.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(item.ConnectionId),
			Data:         []byte(json),
		})
		if err != nil {
			goneError := &apigatewaymanagementapi.GoneException{}
			if errors.As(err, &goneError) {
				log.Printf("Connection %s not found, removing from database", item.ConnectionId)
				dynamo.DeleteItem(&dynamodb.DeleteItemInput{
					TableName: aws.String("chatshit-ee079db"),
					Key: map[string]*dynamodb.AttributeValue{
						"Id": {
							S: aws.String(item.Id),
						},
					},
				})
			} else {
				log.Fatalf("Error calling PostToConnection: %s", err)
			}
		}
	}

	return err
}

// Main Lambda handler function
func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	message := Message{}
	if err := json.NewDecoder(strings.NewReader(request.Body)).Decode(&message); err != nil {
		log.Fatalf("Unable to decode body: %s", err)
	}

	switch message.Action {
	case "register":
		register(request.RequestContext.ConnectionID, message.Payload["id"].(string))
	case "exchange":
		send(message.Payload["theirID"].(string), Message{
			Action: "exchange",
			Payload: map[string]interface{}{
				"id": message.Payload["myID"].(string),
				"publicKey": message.Payload["publicKey"].(string),
			},
		})
	case "send":
		send(message.Payload["id"].(string), Message{
			Action: "message",
			Payload: map[string]interface{}{
				"message": message.Payload["message"].(string),
			},
		})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: "Sent",
	}, nil
}

func init() {
	session := session.Must(session.NewSession())

	dynamo = dynamodb.New(session)
	api = apigatewaymanagementapi.New(session, &aws.Config{
		Endpoint: aws.String("https://k91b8mc10c.execute-api.eu-west-1.amazonaws.com/api-stage-3989ad2"),
	})
}

func main() {
	lambda.Start(handler)
}
