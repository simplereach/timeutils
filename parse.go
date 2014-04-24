package timeutils

// #include "approxidate.h"
// #cgo LDFLAGS: -lm
import "C"
import "time"
import "fmt"

// Takes a string and passes it through Approxidate
// Parses into a time.Time
func ParseDateString(dt string) (t time.Time, err error) {
	date := C.struct_timeval{}

	ok := C.approxidate(C.CString(dt), &date)
	if int(ok) != 0 {
		err = fmt.Errorf("Invlid Date Format %s", dt)
		return
	}

	t = time.Unix(int64(date.tv_sec), int64(date.tv_usec)*1000)
	return
}

// Parses a milliseconds-since-epoch time stamp to a time.Time
func ParseMillis(dt int64) (t time.Time, err error) {
	t = time.Unix(0, dt*1000*1000)
	return
}
