package storage

import (
	"errors"
	"io"
	"os"
	"path"
)

type FileConfig struct {
	BasePath string
}

// File implements the Storage interface for a
// local file storage provider.
type File struct {
	basePath string
}

var _ IStorage = (*File)(nil)

func NewFile(c FileConfig) (*File, error) {
	return &File{basePath: c.BasePath}, nil
}

func (t *File) BucketExists(name string) (bool, error) {
	stat, err := os.Stat(path.Join(t.basePath, name))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !stat.IsDir() {
		return false, errors.New("basePath is a file")
	}
	return true, nil
}

func (t *File) CreateBucket(name string, location ...string) error {
	return os.MkdirAll(path.Join(t.basePath, name), os.ModeDir)
}

func (t *File) CreateBucketIfNotExists(name string, location ...string) error {
	ok, err := t.BucketExists(name)
	if err == nil && !ok {
		err = t.CreateBucket(name, location...)
	}

	return err
}

func (t *File) PutObject(
	bucketName string,
	objectName string,
	reader io.Reader,
	objectSize int64,
	mimeType string,
) error {
	if err := t.CreateBucketIfNotExists(bucketName); err != nil {
		return err
	}

	fd := path.Join(t.basePath, bucketName, objectName)

	stat, err := os.Stat(fd)
	var fh *os.File

	if os.IsNotExist(err) {
		fh, err = os.Create(fd)
	}
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.New("given file dir is a location")
	}

	fh, err = os.Open(fd)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = io.CopyN(fh, reader, objectSize)
	return err
}

func (t *File) GetObject(bucketName string, objectName string) (io.ReadCloser, int64, error) {
	fd := path.Join(t.basePath, bucketName, objectName)
	stat, err := os.Stat(fd)
	var fh *os.File

	if os.IsNotExist(err) {
		return nil, 0, errors.New("file does not exist")
	} else if err != nil {
		return nil, 0, err
	} else if stat.IsDir() {
		return nil, 0, errors.New("given file dir is a location")
	} else {
		fh, err = os.Open(fd)
	}

	return fh, stat.Size(), err
}

func (t *File) DeleteObject(bucketName, objectName string) error {
	fd := path.Join(t.basePath, bucketName, objectName)
	return os.Remove(fd)
}
