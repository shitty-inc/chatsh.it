package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createDynamoDBTable(ctx *pulumi.Context) (error) {
	_, err := dynamodb.NewTable(ctx, "chatshit", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("ConnectionId"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode:  pulumi.String("PAY_PER_REQUEST"),
		HashKey:      pulumi.String("ConnectionId"),
		Ttl: &dynamodb.TableTtlArgs{
			AttributeName: pulumi.String("TimeToExist"),
			Enabled:       pulumi.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
