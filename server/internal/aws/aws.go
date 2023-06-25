package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

type AWS interface {
	CreateFile(ctx context.Context, input s3.PutObjectInput) (output *s3.PutObjectOutput, err error)
	DeleteFile(ctx context.Context, input s3.DeleteObjectInput) (output *s3.DeleteObjectOutput, err error)
}

type AWSservice struct {
	s3 *s3.Client
}

func Initialize(ctx context.Context, region string) AWS {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	return &AWSservice{
		s3: s3.NewFromConfig(cfg),
	}
}

func (aws *AWSservice) CreateFile(ctx context.Context, input s3.PutObjectInput) (output *s3.PutObjectOutput, err error) {
	output, err = aws.s3.PutObject(ctx, &input)
	return
}

func (aws *AWSservice) DeleteFile(ctx context.Context, input s3.DeleteObjectInput) (output *s3.DeleteObjectOutput, err error) {
	output, err = aws.s3.DeleteObject(ctx, &input)
	return
}
