package timeutils

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Format enum type.
type Format int32

// Format enum values.
const (
	Timestamp Format = iota
	TimestampNano
	ANSIC
	UnixDate
	RubyDate
	RFC822
	RFC822Z
	RFC850
	RFC1123
	RFC1123Z
	RFC3339
	RFC3339Nano
	Kitchen
)

// Common errors.
var (
	ErrInvalidFormat = errors.New("invalid format")
)

// Time wraps time.Time overriddin the json marshal/unmarshal to pass
// timestamp as integer
type Time struct {
	time.Time `bson:",inline"`
	format    Format
}

// NewTime create a new Time object with the given format.
func NewTime(t time.Time, format Format) Time {
	return Time{
		Time:   t,
		format: format,
	}
}

// FormatMode create a copy of the time object
// and sets the format method to be used by
// the marhsal functions.
func (t Time) FormatMode(format Format) Time {
	t.format = format
	return t
}

func (t Time) formatTime(mode int) ([]byte, error) {
	var ret string

	switch t.format {
	case ANSIC:
		ret = t.Time.Format(time.ANSIC)
	case UnixDate:
		ret = t.Time.Format(time.UnixDate)
	case RubyDate:
		ret = t.Time.Format(time.RubyDate)
	case RFC822:
		ret = t.Time.Format(time.RFC822)
	case RFC822Z:
		ret = t.Time.Format(time.RFC822Z)
	case RFC850:
		ret = t.Time.Format(time.RFC850)
	case RFC1123:
		ret = t.Time.Format(time.RFC1123)
	case RFC1123Z:
		ret = t.Time.Format(time.RFC1123Z)
	case RFC3339:
		ret = t.Time.Format(time.RFC3339)
	case RFC3339Nano:
		ret = t.Time.Format(time.RFC3339Nano)
	case Kitchen:
		ret = t.Time.Format(time.Kitchen)
	case Timestamp:
		return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
	case TimestampNano:
		return []byte(strconv.FormatInt(t.Time.UnixNano(), 10)), nil
	default:
		return nil, ErrInvalidFormat
	}
	switch mode {
	default:
		fallthrough
	case 0: // json
		return []byte(`"` + ret + `"`), nil
	case 1: // bson
		return []byte(ret), nil
	}
}

// MarshalJSON implements json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return t.formatTime(0)
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
	return t.formatTime(1)
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
