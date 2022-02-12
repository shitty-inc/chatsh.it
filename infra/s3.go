package main

import (
	"io/fs"
	"mime"
	"path"
	"path/filepath"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createWebsiteBucket(ctx *pulumi.Context) (error) {
	siteDir := "../build"

	siteBucket, err := s3.NewBucket(ctx, "bucket", &s3.BucketArgs{
		Bucket: pulumi.String("chatsh.it"),
		Website: s3.BucketWebsiteArgs{
			IndexDocument: pulumi.String("index.html"),
		},
	})
	if err != nil {
		return err
	}

	err = filepath.Walk(siteDir, func(name string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(siteDir, name)
			if err != nil {
				return err
			}

			_, err = s3.NewBucketObject(ctx, rel, &s3.BucketObjectArgs{
				Bucket:      siteBucket.ID(),                                     // reference to the s3.Bucket object
				Source:      pulumi.NewFileAsset(name),                           // use FileAsset to point to a file
				ContentType: pulumi.String(mime.TypeByExtension(path.Ext(name))), // set the MIME type of the file
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if _, err := s3.NewBucketPolicy(ctx, "bucketPolicy", &s3.BucketPolicyArgs{
		Bucket: siteBucket.ID(),
		Policy: pulumi.Any(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect":    "Allow",
					"Principal": "*",
					"Action": []interface{}{
						"s3:GetObject",
					},
					"Resource": []interface{}{
						pulumi.Sprintf("arn:aws:s3:::%s/*", siteBucket.ID()),
					},
				},
			},
		}),
	}); err != nil {
		return err
	}

	return nil
}
