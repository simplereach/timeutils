package timeutils_test

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/simplereach/timeutils"
)

func ExampleNewTime() {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	std := time.Unix(1438883568, 790859087)
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(timeutils.NewTime(std, timeutils.TimestampNano))
	_ = enc.Encode(timeutils.NewTime(std, timeutils.Timestamp))
	_ = enc.Encode(timeutils.NewTime(std, timeutils.RFC3339))
	_ = enc.Encode(timeutils.NewTime(std, timeutils.RFC3339Nano))

	// output:
	// 1438883568790859087
	// 1438883568
	// "2015-08-06T17:52:48Z"
	// "2015-08-06T17:52:48.790859087Z"
}

func ExampleTime() {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	var t timeutils.Time
	_ = json.Unmarshal([]byte(`1438883568`), &t)
	fmt.Println(t.UTC())
	t = timeutils.Time{}
	_ = json.Unmarshal([]byte(`1438883568790859087`), &t)
	fmt.Println(t.UTC())
	t = timeutils.Time{}
	_ = json.Unmarshal([]byte(`"2015-08-06T17:52:48Z"`), &t)
	fmt.Println(t.UTC())
	t = timeutils.Time{}
	_ = json.Unmarshal([]byte(`"09:51:20.939152pm 2014-31-12"`), &t)
	fmt.Println(t.UTC())

	// output:
	// 2015-08-06 17:52:48 +0000 UTC
	// 2015-08-06 17:52:48.790859087 +0000 UTC
	// 2015-08-06 17:52:48 +0000 UTC
	// 2014-12-31 21:51:20.939152 +0000 UTC

}
