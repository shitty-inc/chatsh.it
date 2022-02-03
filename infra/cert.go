package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/acm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createCertificate(ctx *pulumi.Context) (*acm.CertificateValidation, error) {
	awsUsEast1, err := aws.NewProvider(ctx, "aws-provider-us-east-1", &aws.ProviderArgs{Region: pulumi.String("us-east-1")})
	if err != nil {
		return nil, err
	}

	sslCertificate, err := acm.NewCertificate(ctx,
		"ssl-cert",
		&acm.CertificateArgs{
			DomainName:       pulumi.String("chatsh.it"),
			ValidationMethod: pulumi.String("DNS"),
		},
		pulumi.Provider(awsUsEast1),
	)
	if err != nil {
		return nil, err
	}

	validatedSslCertificate, err := acm.NewCertificateValidation(ctx,
    "ssl-cert-validation",
    &acm.CertificateValidationArgs{
			CertificateArn: sslCertificate.Arn,
    },
    pulumi.Provider(awsUsEast1),
	)
	if err != nil {
		return nil, err
	}

	return validatedSslCertificate, nil
}
