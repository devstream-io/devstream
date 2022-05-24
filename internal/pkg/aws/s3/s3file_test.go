package s3_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	s3api "github.com/devstream-io/devstream/internal/pkg/aws/s3"
)

const (
	TestBucket  = "test-bucket"
	TestRegion  = "ap-northeast-1"
	TestKey     = "foo.txt"
	TestContent = "hello, world!"
)

type MockS3Client struct {
	t       *testing.T
	bucket  string
	key     string
	content []byte
}

func NewMockS3Client(t *testing.T, bucket, key string, content []byte) *MockS3Client {
	return &MockS3Client{
		t:       t,
		bucket:  bucket,
		key:     key,
		content: content,
	}
}

func (mock *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	mock.t.Helper()

	checkStringParam(mock.t, "bucket", mock.bucket, params.Bucket)
	checkStringParam(mock.t, "key", mock.key, params.Key)

	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(TestContent))),
	}, nil
}

func (mock *MockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	mock.t.Helper()

	checkStringParam(mock.t, "bucket", mock.bucket, params.Bucket)
	checkStringParam(mock.t, "key", mock.key, params.Key)
	checkBodyParam(mock.t, mock.content, params.Body)

	return &s3.PutObjectOutput{}, nil
}

func checkStringParam(t *testing.T, paramName, expected string, actual *string) {
	t.Helper()

	if actual == nil {
		t.Fatalf("expect %s to not be nil", paramName)
	}
	if expected != *actual {
		t.Errorf("expect %v, got %v", expected, actual)
	}
}

func checkBodyParam(t *testing.T, expected []byte, body io.Reader) {
	t.Helper()

	actual, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatalf("failed to get data from body: %s", err)
	}
	if !bytes.Equal(expected, actual) {
		t.Errorf("expect %v, got %v", expected, actual)
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		bucket string
		key    string
		expect []byte
	}{
		{
			bucket: TestBucket,
			key:    TestKey,
			expect: []byte(TestContent),
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			file, err := s3api.NewS3File(context.TODO(), NewMockS3Client(t, tt.bucket, tt.key, nil), TestBucket, TestRegion, TestKey)
			if err != nil {
				t.Fatalf("failed to create s3 file: %s", err)
			}

			content, err := file.Get()
			if err != nil {
				t.Fatalf("failed to get content from s3 file: %s", err)
			}

			if e, a := tt.expect, content; !bytes.Equal(e, a) {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestPut(t *testing.T) {
	cases := []struct {
		bucket  string
		key     string
		content []byte
	}{
		{
			bucket:  TestBucket,
			key:     TestKey,
			content: []byte(TestContent),
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			file, err := s3api.NewS3File(context.TODO(), NewMockS3Client(t, tt.bucket, tt.key, tt.content), TestBucket, TestRegion, TestKey)
			if err != nil {
				t.Fatalf("failed to create s3 file: %s", err)
			}

			err = file.Put([]byte(TestContent))
			if err != nil {
				t.Fatalf("failed to get content from s3 file: %s", err)
			}
		})
	}
}
