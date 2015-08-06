package timeutils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var timestampTestData = []struct {
	str string
	t   time.Time
}{
	{"1417773080000001234", time.Unix(1417773080, 1234)},
	{"1417791080003300000", time.Unix(1417791080, 3300000)},
	{"1420019480100000000", time.Unix(1420019480, 100000000)},
	{"1420062689999999999", time.Unix(1420062680, 9999999999)},
}

type tsTest struct {
	TS Time `json:"ts"`
}

func TestJSONRoundTripTimestamp(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	for i, elem := range timestampTestData {
		str := `{"ts":` + elem.str + `}`
		ts := &tsTest{}
		if err := json.Unmarshal([]byte(str), &ts); err != nil {
			t.Fatal(err)
		}
		if expect, got := elem.t, ts.TS; expect.Sub(got.Time) != 0 {
			t.Errorf("[%d] Unexpected result.\nExpect:\t%s\nGot:\t%s\n", i, expect, got)
		}
		ts.TS.NanoPtr(true)
		buf, err := json.Marshal(ts)
		if err != nil {
			t.Fatal(err)
		}
		if expect, got := str, string(buf); expect != got {
			t.Errorf("[%d] Unexpected result.\nExpect:\t%s\nGot:\t%s\n", i, expect, got)
		}
	}
}

func TestJSONNull(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	str := `{"ts":null}`
	ts := &tsTest{}
	if err := json.Unmarshal([]byte(str), &ts); err != nil {
		t.Fatal(err)
	}
	if !ts.TS.IsZero() {
		t.Fatalf("Unexpected result. Time not zero")
	}
	buf, err := json.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := str, string(buf); expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestJSONRegular(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	str := `{"ts":"2015-08-06T12:17:25.881396749Z"}`
	ts := &tsTest{}
	if err := json.Unmarshal([]byte(str), &ts); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(1438863445, 881396749), ts.TS; expect.Sub(got.Time) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
	ts.TS.Nano(true)
	buf, err := json.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := `{"ts":1438863445881396749}`, string(buf); expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestJSONTimestampSimple(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	str := `{"ts":"141779108"}`
	ts := &tsTest{}
	if err := json.Unmarshal([]byte(str), &ts); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(141779108, 0), ts.TS; expect.Sub(got.Time) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
	ts.TS.Nano(true)
	buf, err := json.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := `{"ts":141779108000000000}`, string(buf); expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestJSONTimestampSimpleInt(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	str := `{"ts":141779108}`
	ts := &tsTest{}
	if err := json.Unmarshal([]byte(str), &ts); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(141779108, 0), ts.TS; expect.Sub(got.Time) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
	ts.TS.Nano(false)
	buf, err := json.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := `{"ts":141779108}`, string(buf); expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestJSONTimestampManual(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	for i, elem := range testData {
		tt, err := ParseDateString(elem.out)
		if err != nil {
			t.Fatal(err)
		}
		str := `{"ts":"` + elem.in + `"}`
		ts := &tsTest{}
		if err := json.Unmarshal([]byte(str), &ts); err != nil {
			t.Fatal(err)
		}
		if expect, got := tt, ts.TS; expect.Sub(got.Time) != 0 {
			t.Errorf("[%d] Unexpected result.\nExpect:\t%s\nGot:\t%s\n", i, expect, got)
		}
		ts.TS.Nano(true)
		buf, err := json.Marshal(ts)
		if err != nil {
			t.Fatal(err)
		}
		if expect, got := fmt.Sprintf(`{"ts":%d}`, tt.UnixNano()), string(buf); expect != got {
			t.Errorf("[%d] Unexpected result.\nExpect:\t%s\nGot:\t%s\n", i, expect, got)
		}
	}
}

func TestBSONRoundTripTimestamp(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	for i, elem := range timestampTestData {
		b := bson.M{"ts": elem.str}
		buf, err := bson.Marshal(b)
		if err != nil {
			t.Fatal(err)
		}
		ts := &tsTest{}
		if err := bson.Unmarshal(buf, &ts); err != nil {
			t.Fatal(err)
		}
		if expect, got := elem.t, ts.TS; expect.Sub(got.Time) != 0 {
			t.Errorf("[%d] Unexpected result.\nExpect:\t%s\nGot:\t%s\n", i, expect, got)
		}
	}
}

func TestBSONNull(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	{
		buf, err := bson.Marshal(bson.M{"ts": ""})
		if err != nil {
			t.Fatal(err)
		}
		ts := &tsTest{}
		if err := bson.Unmarshal(buf, &ts); err != nil {
			t.Fatal(err)
		}
		if !ts.TS.IsZero() {
			t.Fatalf("Unexpected result. Time not zero")
		}
		buf2, err := bson.Marshal(ts)
		if err != nil {
			t.Fatal(err)
		}
		if err := bson.Unmarshal(buf2, &ts); err != nil {
			t.Fatal(err)
		}
	}
	{
		buf, err := bson.Marshal(bson.M{"ts": nil})
		if err != nil {
			t.Fatal(err)
		}
		ts := &tsTest{}
		if err := bson.Unmarshal(buf, &ts); err != nil {
			t.Fatal(err)
		}
		if !ts.TS.IsZero() {
			t.Fatalf("Unexpected result. Time not zero")
		}
	}
}

func TestBSONRegular(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	buf, err := bson.Marshal(bson.M{"ts": "2015-08-06T12:17:25.881396749Z"})
	if err != nil {
		t.Fatal(err)
	}
	ts := &tsTest{}
	if err := bson.Unmarshal(buf, &ts); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(1438863445, 881396749), ts.TS; expect.Sub(got.Time) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestBSONTimestampSimple(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	buf, err := bson.Marshal(bson.M{"ts": "141779108"})
	if err != nil {
		t.Fatal(err)
	}
	ts := &tsTest{}
	if err := bson.Unmarshal(buf, &ts); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(141779108, 0), ts.TS; expect.Sub(got.Time) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestBSONTimestampSimpleIntNano(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	buf, err := bson.Marshal(bson.M{"ts": 141779108000000999})
	if err != nil {
		t.Fatal(err)
	}
	ts := &tsTest{}
	if err := bson.Unmarshal(buf, &ts); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(141779108, 999), ts.TS; expect.Sub(got.Time) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
	ts.TS.Nano(true)
	buf2, err := bson.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := "141779108000000999", string(buf2); !strings.Contains(got, expect) {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestBSONTimestampSimpleInt(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	buf, err := bson.Marshal(bson.M{"ts": 141779108})
	if err != nil {
		t.Fatal(err)
	}
	ts := &tsTest{}
	if err := bson.Unmarshal(buf, &ts); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(141779108, 0), ts.TS; expect.Sub(got.Time) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestBSONTimestampManual(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	for i, elem := range testData {
		tt, err := ParseDateString(elem.out)
		if err != nil {
			t.Fatal(err)
		}
		buf, err := bson.Marshal(bson.M{"ts": elem.in})
		if err != nil {
			t.Fatal(err)
		}

		ts := &tsTest{}
		if err := bson.Unmarshal(buf, &ts); err != nil {
			t.Fatal(err)
		}
		if expect, got := tt, ts.TS; expect.Sub(got.Time) != 0 {
			t.Errorf("[%d] Unexpected result.\nExpect:\t%s\nGot:\t%s\n", i, expect, got)
		}
		buf2, err := bson.Marshal(ts)
		if err != nil {
			t.Fatal(err)
		}
		if expect, got := fmt.Sprintf("%d", tt.Unix()), string(buf2); !strings.Contains(got, expect) {
			t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
		}
	}
}
