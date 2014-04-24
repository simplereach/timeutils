package timeutils

// #include "approxidate.h"
// #cgo LDFLAGS: -lm
import "C"
import "time"
import "fmt"

func parseDateStr(dt string) (t time.Time, err error) {
	date := C.struct_timeval{}

	ok := C.approxidate(C.CString(dt), &date)
	if int(ok) != 0 {
		err = fmt.Errorf("Invlid Date Format %s", dt)
		return
	}

	t = time.Unix(int64(date.tv_sec), int64(date.tv_usec)*1000)
	return
}

func parseDateInt(dt int64) (t time.Time, err error) {
	t = time.Unix(0, dt*1000*1000)
	return
}

func ParseDate(dt interface{}) (time.Time, error) {
	idt, ok := dt.(int64)
	if ok {
		return parseDateInt(idt)
	}

	sdt, ok := dt.(string)
	if ok {
		return parseDateStr(sdt)
	}

	panic("Invlid type supplied for ParseDate")
}
