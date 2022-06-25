package util

import "strings"

func SplitIfNotEmpty(v string, sep string) []string {
	if len(v) == 0 {
		return nil
	}
	return strings.Split(v, sep)
}

func SplitAndClean(v string, sep string) []string {
	split := SplitIfNotEmpty(v, sep)
	if len(split) != 0 {
		for i, s := range split {
			split[i] = strings.TrimSpace(s)
		}
	}
	return split
}

func StartsWithAny(v string, prefixes []string) (bool, string) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(v, prefix) {
			return true, prefix
		}
	}
	return false, ""
}
