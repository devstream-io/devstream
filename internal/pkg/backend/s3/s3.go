package s3

import (
	"context"

	"github.com/devstream-io/devstream/pkg/util/cloud/aws/s3"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type S3Backend struct {
	file *s3.S3File
}

func NewS3Backend(bucket, region, key string) (*S3Backend, error) {
	if err := validate(bucket, region, key); err != nil {
		return nil, err
	}

	log.Infof("Using s3 backend. Bucket: %s, region: %s, key: %s.", bucket, region, key)

	ctx := context.Background()
	client, err := s3.NewClient(ctx, region)
	if err != nil {
		log.Fatalf("Creating s3 client failed: %s.", err)
	}

	file, err := s3.NewS3File(ctx, client, bucket, region, key)
	if err != nil {
		log.Fatalf("Creating remote state file %s failed.", key)
	}

	return &S3Backend{
		file: file,
	}, nil
}

func (b *S3Backend) Read() ([]byte, error) {
	return b.file.Get()
}

func (b *S3Backend) Write(data []byte) error {
	return b.file.Put(data)
}
