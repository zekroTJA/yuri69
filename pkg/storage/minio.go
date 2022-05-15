package storage

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zekrotja/yuri69/pkg/util"
)

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Location        string
	Secure          bool
}

// Minio implements the Storage interface for
// the MinIO SDK to connect to a MinIO instance,
// Amazon S3 or Google Cloud.
type Minio struct {
	client   *minio.Client
	location string
}

var _ IStorage = (*Minio)(nil)

func NewMinio(c MinioConfig) (*Minio, error) {
	var t Minio
	var err error

	t.location = c.Location
	t.client, err = minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, ""),
		Secure: c.Secure,
	})

	return &t, err
}

func (t *Minio) BucketExists(name string) (bool, error) {
	ctx, cancel := timeoutContext()
	defer cancel()
	return t.client.BucketExists(ctx, name)
}

func (t *Minio) CreateBucket(name string, location ...string) error {
	ctx, cancel := timeoutContext()
	defer cancel()
	return t.client.MakeBucket(ctx, name, minio.MakeBucketOptions{
		Region: t.getLocation(location),
	})
}

func (t *Minio) CreateBucketIfNotExists(name string, location ...string) error {
	ok, err := t.BucketExists(name)
	if err == nil && !ok {
		err = t.CreateBucket(name, location...)
	}

	return err
}

func (t *Minio) PutObject(
	bucketName string,
	objectName string,
	reader io.Reader,
	objectSize int64,
	mimeType string,
) error {
	if err := t.CreateBucketIfNotExists(bucketName, t.location); err != nil {
		return err
	}

	ctx, cancel := timeoutContext(5 * time.Minute)
	defer cancel()
	_, err := t.client.PutObject(ctx, bucketName, objectName, reader, objectSize,
		minio.PutObjectOptions{
			ContentType: mimeType,
		})

	return err
}

func (t *Minio) GetObject(bucketName, objectName string) (io.ReadCloser, int64, error) {
	ctx, _ := timeoutContext(5 * time.Minute)
	obj, err := t.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, err
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, 0, err
	}

	return obj, stat.Size, err
}

func (t *Minio) DeleteObject(bucketName, objectName string) error {
	ctx, cancel := timeoutContext()
	defer cancel()
	return t.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (t *Minio) getLocation(loc []string) string {
	if len(loc) > 0 {
		return loc[0]
	}
	return t.location
}

func timeoutContext(d ...time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), util.Opt(d, 1*time.Second))
}
