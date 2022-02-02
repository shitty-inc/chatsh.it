package main

import (
	awsapigateway "github.com/pulumi/pulumi-aws-apigateway/sdk/go/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createApiGateway(ctx *pulumi.Context, onConnect *lambda.Function, onDisconnect *lambda.Function, sendMessage *lambda.Function, cert *acm.CertificateValidation) (*awsapigateway.RestAPI, error) {
	localPath := "../build"
	getMethod := awsapigateway.MethodGET

	restAPI, err := awsapigateway.NewRestAPI(ctx, "chatshit", &awsapigateway.RestAPIArgs{
		Routes: []awsapigateway.RouteArgs{
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

	_, err = apigateway.NewDomainName(ctx, "api-domain-name",
		&apigateway.DomainNameArgs{
			CertificateArn: cert.CertificateArn,
			DomainName:     pulumi.String("chatsh.it"),
		},
	)
	if err != nil {
		return nil, err
	}

	return restAPI, nil
}
