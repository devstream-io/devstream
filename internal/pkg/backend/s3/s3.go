package s3

import (
	"context"

	"github.com/devstream-io/devstream/internal/pkg/aws/s3"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type S3Backend struct {
	file *s3.S3File
}

func NewS3Backend(bucket, region, key string) *S3Backend {
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
	}
}

func (b *S3Backend) Read() ([]byte, error) {
	return b.file.Get()
}

func (b *S3Backend) Write(data []byte) error {
	return b.file.Put(data)
}
