package dberrors

import "errors"

var (
	ErrUnsupportedProviderType = errors.New("unsupported database provider type")
	ErrNotFound                = errors.New("not found")
)
