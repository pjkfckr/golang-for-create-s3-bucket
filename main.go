package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"os"
)

const bucketName = "go-aws-demo-991"

func main() {
	var (
		s3Client *s3.Client
		err      error
	)
	ctx := context.Background()
	s3Client, err = initS3Client(ctx, "go-iam")
	if err != nil {
		fmt.Printf("initS3Client error: %s", err)
		os.Exit(1)
	}

	err = createS3Bucket(ctx, s3Client)
	if err != nil {
		fmt.Printf("createS3Bucket error: %s", err)
		os.Exit(1)
	}
}

func initS3Client(ctx context.Context, profile string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %s", err)
	}
	return s3.NewFromConfig(cfg), nil
}

func createS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	canMake, err := findBucketByName(ctx, s3Client)
	if err != nil {
		return err
	}

	if canMake {
		_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: "ap-northeast-2",
			},
		})
		if err != nil {
			return fmt.Errorf("CreateBucket error: %s", err)
		}
	}

	return nil
}

func findBucketByName(ctx context.Context, s3Client *s3.Client) (bool, error) {
	buckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return false, fmt.Errorf("FindByBucketName error: %s", err)
	}

	canMake := true
	for _, bucket := range buckets.Buckets {
		if *bucket.Name == bucketName {
			canMake = false
		}
	}
	return canMake, nil

}
