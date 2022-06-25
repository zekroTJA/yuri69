package util

func Opt[T any](v []T, def ...T) T {
	if len(v) == 0 {
		var altDef T
		return Opt(def, altDef)
	}
	return v[0]
}
