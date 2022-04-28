package storage

import "errors"

var (
	ErrUnsupportedProviderType = errors.New("unsupported storage provider type")
)
