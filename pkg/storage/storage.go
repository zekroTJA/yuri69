package storage

import (
	"io"
	"strings"
)

type IStorage interface {
	BucketExists(name string) (bool, error)
	CreateBucket(name string, location ...string) error
	CreateBucketIfNotExists(name string, location ...string) error

	PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, mimeType string) error
	GetObject(bucketName, objectName string) (io.ReadCloser, int64, error)
	DeleteObject(bucketName, objectName string) error
}

type StorageConfig struct {
	Type  string
	Minio MinioConfig
	File  FileConfig
}

func New(c StorageConfig) (IStorage, error) {
	switch strings.ToLower(c.Type) {
	case "local", "file", "files":
		return NewFile(c.File)
	case "minio", "s3":
		return NewMinio(c.Minio)
	default:
		return nil, ErrUnsupportedProviderType
	}
}
