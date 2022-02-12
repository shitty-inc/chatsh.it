package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/acm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createCertificate(ctx *pulumi.Context) (*acm.CertificateValidation, error) {
	sslCertificate, err := acm.NewCertificate(ctx,
		"ssl-cert",
		&acm.CertificateArgs{
			DomainName:       pulumi.String("api.chatsh.it"),
			ValidationMethod: pulumi.String("DNS"),
		},
	)
	if err != nil {
		return nil, err
	}

	validatedSslCertificate, err := acm.NewCertificateValidation(ctx,
    "ssl-cert-validation",
    &acm.CertificateValidationArgs{
			CertificateArn: sslCertificate.Arn,
    },
	)
	if err != nil {
		return nil, err
	}

	return validatedSslCertificate, nil
}
