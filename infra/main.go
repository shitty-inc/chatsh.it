package main

import (
	"github.com/pulumi/pulumi-aws/sdk/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/acm"
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

		domain := pulumi.String("chatsh.it")
		awsUsEast1, err := aws.NewProvider(ctx, "aws-provider-us-east-1", &aws.ProviderArgs{Region: pulumi.String("us-east-1")})
		sslCertificate, err := acm.NewCertificate(ctx,
			"ssl-cert",
			&acm.CertificateArgs{
				DomainName:       domain,
				ValidationMethod: pulumi.String("DNS"),
			},
			pulumi.Provider(awsUsEast1),
		)

		ctx.Export("url", gateway.Url)

		return nil
	})
}
