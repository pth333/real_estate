package initialize

import (
	"context"
	"real_estate_be/internal/global"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3() {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(global.Config.R2.Region),
		config.WithCredentialsProvider(
			aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     global.Config.R2.AccessKeyID,
					SecretAccessKey: global.Config.R2.SecretAccessKey,
				}, nil
			}),
		),
	)
	if err != nil {
		panic("Failed to load S3 config: " + err.Error())
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(global.Config.R2.Endpoint)
		o.UsePathStyle = true
	})

	global.S3Client = client
}
