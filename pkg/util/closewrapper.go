package util

import "io"

type readCloserWrapper struct {
	io.ReadCloser

	afterClose func(error) error
}

var _ io.ReadCloser = (*readCloserWrapper)(nil)

func WrapReadCloser(rc io.ReadCloser, afterClose func(error) error) readCloserWrapper {
	return readCloserWrapper{
		ReadCloser: rc,
		afterClose: afterClose,
	}
}

func (t readCloserWrapper) Close() error {
	err := t.ReadCloser.Close()
	return t.afterClose(err)
}
