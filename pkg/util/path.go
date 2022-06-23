package util

import (
	"path"
	"strings"
)

func CleanBase(dir string) string {
	base := path.Base(dir)
	i := strings.IndexRune(base, '.')
	if i != -1 {
		base = base[:i]
	}
	return base
}
