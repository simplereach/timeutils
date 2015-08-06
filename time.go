package timeutils

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Time wraps time.Time overriddin the json marshal/unmarshal to pass
// timestamp as integer
type Time struct {
	time.Time `bson:",inline"`
	nano      int32
}

// Nano sets the nano flag and return the object.
// This will produce nano timestamp when marshaling.
func (t *Time) Nano(set bool) Time {
	return *t.NanoPtr(set)
}

// NanoPtr is like Nano but returns a pointer to time.
func (t *Time) NanoPtr(set bool) *Time {
	if set {
		atomic.StoreInt32(&t.nano, 1)
	} else {
		atomic.StoreInt32(&t.nano, 0)
	}
	return t
}

// MarshalJSON implements json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	if atomic.LoadInt32(&t.nano) == 1 {
		return []byte(strconv.FormatInt(t.Time.UTC().UnixNano(), 10)), nil
	}
	return []byte(strconv.FormatInt(t.Time.UTC().Unix(), 10)), nil
}

// UnmarshalJSON implements json.Unmarshaler inferface.
func (t *Time) UnmarshalJSON(buf []byte) error {
	// Try to parse the timestamp integer
	ts, err := strconv.ParseInt(string(buf), 10, 64)
	if err == nil {
		if len(buf) == 19 {
			t.Time = time.Unix(ts/1e9, ts%1e9)
		} else {
			t.Time = time.Unix(ts, 0)
		}
		return nil
	}
	// Try the default unmarshal
	if err := json.Unmarshal(buf, &t.Time); err == nil {
		return nil
	}
	str := strings.Trim(string(buf), `"`)
	if str == "null" || str == "" {
		return nil
	}
	// Try to manually parse the data
	tt, err := ParseDateString(str)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

// GetBSON implements mgo/bson.Getter interface.
func (t Time) GetBSON() (interface{}, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	if atomic.LoadInt32(&t.nano) == 1 {
		return strconv.FormatInt(t.Time.UTC().UnixNano(), 10), nil
	}
	return strconv.FormatInt(t.Time.UTC().Unix(), 10), nil
}

// SetBSON implements mgo/bson.Setter interface.
func (t *Time) SetBSON(raw bson.Raw) error {
	// Try the default unmarshal
	if err := raw.Unmarshal(&t.Time); err == nil {
		return nil
	}

	// Try to pull the timestamp as an int
	var tsInt int64
	if err := raw.Unmarshal(&tsInt); err == nil {
		if tsInt > 5e9 {
			t.Time = time.Unix(tsInt/1e9, tsInt%1e9)
		} else {
			t.Time = time.Unix(tsInt, 0)
		}
	}

	// Try to pull the timestamp as a string
	var tsStr string
	if err := raw.Unmarshal(&tsStr); err != nil {
		return err
	}
	if tsStr == "" {
		return nil
	}
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err == nil {
		if len(tsStr) == 19 {
			t.Time = time.Unix(ts/1e9, ts%1e9)
		} else {
			t.Time = time.Unix(ts, 0)
		}
		return nil
	}

	// Try the json umarshal
	if err := json.Unmarshal([]byte(`"`+tsStr+`"`), &t.Time); err == nil {
		return nil
	}

	// Try to manually parse the data
	tt, err := ParseDateString(tsStr)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}
