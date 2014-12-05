package timeutils

// #include "approxidate.h"
// #include <stdlib.h>
// #cgo LDFLAGS: -lm
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

// ParseDateString takes a string and passes it through Approxidate
// Parses into a time.Time
func ParseDateString(dt string) (time.Time, error) {
	date := C.struct_timeval{}

	cStr := C.CString(dt)
	ok := C.approxidate(cStr, &date)
	C.free(unsafe.Pointer(cStr))
	if int(ok) != 0 {
		return time.Time{}, fmt.Errorf("Invlid Date Format %s", dt)
	}

	return time.Unix(int64(date.tv_sec), int64(date.tv_usec)*1000), nil
}

// ParseMillis parses a milliseconds-since-epoch time stamp to a time.Time
func ParseMillis(dt int64) (time.Time, error) {
	return time.Unix(0, dt*1000*1000), nil
}
