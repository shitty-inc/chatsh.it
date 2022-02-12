package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createApiGateway(ctx *pulumi.Context, onConnect *lambda.Function, onDisconnect *lambda.Function, sendMessage *lambda.Function, cert *acm.CertificateValidation) (*apigatewayv2.Api, error) {
	api, err := apigatewayv2.NewApi(ctx, "api", &apigatewayv2.ApiArgs{
		Name:                     pulumi.String("chatshit"),
		ProtocolType:             pulumi.String("WEBSOCKET"),
		RouteSelectionExpression: pulumi.String(fmt.Sprintf("%v%v", "$", "request.body.action")),
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewIntegration(ctx, "connect-integration", &apigatewayv2.IntegrationArgs{
		ApiId:                   api.ID(),
		IntegrationType:         pulumi.String("AWS"),
		ContentHandlingStrategy: pulumi.String("CONVERT_TO_TEXT"),
		IntegrationMethod:       pulumi.String("POST"),
		IntegrationUri:          onConnect.InvokeArn,
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, "connect-route", &apigatewayv2.RouteArgs{
		ApiId:         api.ID(),
		RouteKey:      pulumi.String(fmt.Sprintf("%v%v", "$", "connect")),
	})
	if err != nil {
		return nil, err
	}

	_, err = lambda.NewPermission(ctx, "connect-permission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  onConnect.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewIntegration(ctx, "disconnect-integration", &apigatewayv2.IntegrationArgs{
		ApiId:                   api.ID(),
		IntegrationType:         pulumi.String("AWS"),
		ContentHandlingStrategy: pulumi.String("CONVERT_TO_TEXT"),
		IntegrationMethod:       pulumi.String("POST"),
		IntegrationUri:          onDisconnect.InvokeArn,
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, "disconnect-route", &apigatewayv2.RouteArgs{
		ApiId:         api.ID(),
		RouteKey:      pulumi.String(fmt.Sprintf("%v%v", "$", "disconnect")),
	})
	if err != nil {
		return nil, err
	}

	_, err = lambda.NewPermission(ctx, "disconnect-permission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  onDisconnect.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewIntegration(ctx, "sendmessage-integration", &apigatewayv2.IntegrationArgs{
		ApiId:                   api.ID(),
		IntegrationType:         pulumi.String("AWS"),
		ContentHandlingStrategy: pulumi.String("CONVERT_TO_TEXT"),
		IntegrationMethod:       pulumi.String("POST"),
		IntegrationUri:          sendMessage.InvokeArn,
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, "sendmessage-route", &apigatewayv2.RouteArgs{
		ApiId:         api.ID(),
		RouteKey:      pulumi.String("sendmessage"),
	})
	if err != nil {
		return nil, err
	}

	_, err = lambda.NewPermission(ctx, "sendmessage-permission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  sendMessage.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
	})
	if err != nil {
		return nil, err
	}

	domainName, err := apigatewayv2.NewDomainName(ctx, "api-domain", &apigatewayv2.DomainNameArgs{
		DomainName: pulumi.String("api.chatsh.it"),
		DomainNameConfiguration: &apigatewayv2.DomainNameDomainNameConfigurationArgs{
			CertificateArn: cert.CertificateArn,
			EndpointType:   pulumi.String("REGIONAL"),
			SecurityPolicy: pulumi.String("TLS_1_2"),
		},
	})
	if err != nil {
		return nil, err
	}

	stage, err := apigatewayv2.NewStage(ctx, "api-stage", &apigatewayv2.StageArgs{
		ApiId: api.ID(),
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewApiMapping(ctx, "api-mapping", &apigatewayv2.ApiMappingArgs{
		ApiId:      api.ID(),
		DomainName: domainName.DomainName,
		Stage:      stage.ID(),
	})
	if err != nil {
		return nil, err
	}

	return api, nil
}
