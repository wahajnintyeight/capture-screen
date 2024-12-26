package aws

import (
	"bytes"
	"capture-screen/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

)

type S3Service struct {
    client *s3.Client
    bucket string
}

func NewS3Service(ctx context.Context) (*S3Service, error) {
    b := config.GetEnv("S3_BUCKET_NAME")
	accessKey := config.GetEnv("S3_ACCESS_KEY_ID")
	secretKey := config.GetEnv("S3_SECRET_ACCESS_KEY")
	region := config.GetEnv("S3_REGION")
	customResolver := awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))

	cfg, err := awsconfig.LoadDefaultConfig(ctx, customResolver, awsconfig.WithRegion(region))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}
	
    client := s3.NewFromConfig(cfg)
    return &S3Service{client: client, bucket: b}, nil
}

func (s *S3Service) UploadImage(ctx context.Context, imageBytes []byte, deviceName string) (string, error) {
    
	key := config.GetEnv("S3_FOLDER_NAME") + "/" + deviceName + "/" + time.Now().Format("2006-01-02-15-04-05") + ".png"
	log.Println("AWS KEY", key)

	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(imageBytes),
		// ACL:    types.ObjectCannedACLPublicRead, // Make object publicly readable
	}
	log.Println("AWS BUCKET", s.bucket)

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	} else {
		url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
		log.Println("AWS URL", url)
		return url, nil
	}

}
