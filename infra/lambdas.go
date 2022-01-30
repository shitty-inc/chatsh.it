package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createLambdas(ctx *pulumi.Context) (*lambda.Function, *lambda.Function, *lambda.Function, error) {
	role, err := iam.NewRole(ctx, "auth-exec-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Sid": "",
				"Effect": "Allow",
				"Principal": {
					"Service": "lambda.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}`),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	logPolicy, err := iam.NewRolePolicy(ctx, "lambda-log-policy", &iam.RolePolicyArgs{
		Role: role.Name,
		Policy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": [
					"logs:CreateLogGroup",
					"logs:CreateLogStream",
					"logs:PutLogEvents"
				],
				"Resource": "arn:aws:logs:*:*:*"
			}]
		}`),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	onConnect, err := lambda.NewFunction(ctx, "chatshit-onConnect", &lambda.FunctionArgs{
		Runtime: lambda.RuntimeGo1dx,
		Code: pulumi.NewAssetArchive(map[string]interface{}{
			".": pulumi.NewFileArchive("./functions/onConnect"),
		}),
		Handler: pulumi.String("handler"),
		Role:    role.Arn,
	}, pulumi.DependsOn([]pulumi.Resource{logPolicy}))

	if err != nil {
		return nil, nil, nil, err
	}

	onDisconnect, err := lambda.NewFunction(ctx, "chatshit-onDisconnect", &lambda.FunctionArgs{
		Runtime: lambda.RuntimeGo1dx,
		Code: pulumi.NewAssetArchive(map[string]interface{}{
			".": pulumi.NewFileArchive("./functions/onDisconnect"),
		}),
		Handler: pulumi.String("handler"),
		Role:    role.Arn,
	}, pulumi.DependsOn([]pulumi.Resource{logPolicy}))

	if err != nil {
		return nil, nil, nil, err
	}

	sendMessage, err := lambda.NewFunction(ctx, "chatshit-sendMessage", &lambda.FunctionArgs{
		Runtime: lambda.RuntimeGo1dx,
		Code: pulumi.NewAssetArchive(map[string]interface{}{
			".": pulumi.NewFileArchive("./functions/sendMessage"),
		}),
		Handler: pulumi.String("handler"),
		Role:    role.Arn,
	}, pulumi.DependsOn([]pulumi.Resource{logPolicy}))

	if err != nil {
		return nil, nil, nil, err
	}
	return onConnect, onDisconnect, sendMessage, nil
}
