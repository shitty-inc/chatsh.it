package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		onConnect, onDisconnect, onMessage, err := createLambdas(ctx)
		if err != nil {
			return err
		}

		cert, err := createCertificate(ctx)
		if err != nil {
			return err
		}

		gateway, err := createApiGateway(ctx, onConnect, onDisconnect, onMessage, cert)
		if err != nil {
			return err
		}

		err = createDynamoDBTable(ctx)
		if err != nil {
			return err
		}

		err = createWebsiteBucket(ctx)
		if err != nil {
			return err
		}

		ctx.Export("url", gateway)

		return nil
	})
}
