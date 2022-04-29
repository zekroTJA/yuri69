package errs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zekrotja/yuri69/pkg/util"
)

type unwrappable interface {
	Unwrap() error
}

type StatusError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (t StatusError) Error() string {
	return fmt.Sprintf("%d: %s", t.Status, t.Message)
}

type UserError struct {
	inner error
}

func WrapUserError(v any, status ...int) UserError {
	var inner error
	switch vt := v.(type) {
	case error:
		inner = vt
	case string:
		inner = errors.New(vt)
	default:
		inner = fmt.Errorf("%+v", vt)
	}

	return UserError{
		inner: StatusError{
			Status:  util.Opt(status, http.StatusBadRequest),
			Message: inner.Error(),
		},
	}
}

func (t UserError) Error() string {
	return t.inner.Error()
}

func (t UserError) Unwrap() error {
	return t.inner
}

func As[T error](err any) (v T, ok bool) {
	if err == nil {
		return v, false
	}

	if v, ok = err.(T); ok {
		return v, ok
	}

	if uerr, ok := err.(unwrappable); ok {
		return As[T](uerr.Unwrap())
	}

	return v, false
}
