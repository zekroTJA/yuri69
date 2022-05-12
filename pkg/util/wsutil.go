package util

import (
	"strconv"

	routing "github.com/zekrotja/ozzo-routing/v2"
)

func QueryInt(ctx *routing.Context, name string, def int) (int, error) {
	vStr := ctx.Query(name)
	if vStr == "" {
		return def, nil
	}

	v, err := strconv.Atoi(vStr)
	if err != nil {
		return 0, err
	}

	return v, nil
}
