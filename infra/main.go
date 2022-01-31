package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		onConnect, onDisconnect, sendMessage, err := createLambdas(ctx)
		if err != nil {
			return err
		}

		gateway, err := createApiGateway(ctx, onConnect, onDisconnect, sendMessage)
		if err != nil {
			return err
		}

		ctx.Export("url", gateway.Url)

		return nil
	})
}