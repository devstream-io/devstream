package s3

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	awsutil "github.com/devstream-io/devstream/internal/pkg/aws/util"
)

type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3File struct {
	Bucket string
	Region string
	Key    string
	ctx    context.Context
	api    S3API
}

func NewS3File(ctx context.Context, api S3API, bucket, region, key string) (*S3File, error) {
	file := S3File{
		Bucket: bucket,
		Region: region,
		Key:    key,
		ctx:    ctx,
		api:    api,
	}
	return &file, nil
}

func (f *S3File) Put(data []byte) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(f.Bucket),
		Key:    aws.String(f.Key),
		Body:   bytes.NewReader(data),
	}
	_, err := f.api.PutObject(f.ctx, params)
	if err != nil {
		awsutil.LogAWSError(err)
		return fmt.Errorf("failed to upload %s to bucket %s", f.Key, f.Bucket)
	}
	return nil
}

func (f *S3File) Get() ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(f.Bucket),
		Key:    aws.String(f.Key),
	}
	out, err := f.api.GetObject(f.ctx, params)
	if err != nil {
		awsutil.LogAWSError(err)
		return nil, fmt.Errorf("failed to download %s from bucket %s", f.Key, f.Bucket)
	}

	defer out.Body.Close()
	data, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
