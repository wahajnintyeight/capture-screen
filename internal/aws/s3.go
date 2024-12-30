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
    
	// Delete previous screenshots for this device
	prefix := config.GetEnv("S3_FOLDER_NAME") + "/" + deviceName + "/"
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}
	
	result, err := s.client.ListObjectsV2(ctx, listInput)
	if err != nil {
		return "", fmt.Errorf("failed to list objects: %w", err)
	}

	for _, obj := range result.Contents {
		deleteInput := &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    obj.Key,
		}
		
		_, err = s.client.DeleteObject(ctx, deleteInput)
		if err != nil {
			return "", fmt.Errorf("failed to delete object %s: %w", *obj.Key, err)
		}
	}

	// Upload new screenshot
	key := prefix + time.Now().Format("2006-01-02-15-04-05") + ".png"
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(imageBytes),
	}

	_, err = s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
	return url, nil
    
}
