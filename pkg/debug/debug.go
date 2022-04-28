package debug

import (
	"sync/atomic"
)

var debugMode uint32 = 0

func Enabled() bool {
	return atomic.LoadUint32(&debugMode) == 1
}

func SetEnabled(isDebug bool) {
	if isDebug {
		atomic.StoreUint32(&debugMode, 1)
	} else {
		atomic.StoreUint32(&debugMode, 0)
	}
}
