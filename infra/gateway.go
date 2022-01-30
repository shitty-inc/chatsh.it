package main

import (
	apigateway "github.com/pulumi/pulumi-aws-apigateway/sdk/go/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createApiGateway(ctx *pulumi.Context, onConnect *lambda.Function, onDisconnect *lambda.Function, sendMessage *lambda.Function) (*apigateway.RestAPI, error) {
	localPath := "../build"
	getMethod := apigateway.MethodGET
	restAPI, err := apigateway.NewRestAPI(ctx, "chatshit", &apigateway.RestAPIArgs{
		Routes: []apigateway.RouteArgs{
			{
				Path:      "/",
				LocalPath: &localPath,
			},
			{
				Path:         "onConnect",
				Method:       &getMethod,
				EventHandler: onConnect,
			},
			{
				Path:         "onDisconnect",
				Method:       &getMethod,
				EventHandler: onDisconnect,
			},
			{
				Path:         "sendMessage",
				Method:       &getMethod,
				EventHandler: sendMessage,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return restAPI, nil
}
