package timeutils

import (
	"os"
	"testing"
)

var testData = []struct{ in, out string }{
	{"2014-12-05 09:51:20.939152 -0500", "2014-12-05 14:51:20.939152 +0000 UTC"},
	{"2014-12-05 09:51:20.939152 -0500 EST", "2014-12-05 14:51:20.939152 +0000 UTC"},
	{"2014-12-05 09:51:20.939152", "2014-12-05 09:51:20.939152 +0000 UTC"},
	{"2014/12/05 09:51:20.939152", "2014-12-05 09:51:20.939152 +0000 UTC"},
	{"2014.12.05 09:51:20.939152", "2014-12-05 09:51:20.939152 +0000 UTC"},
	{"09:51:20.939152 2014-31-12", "2014-12-31 09:51:20.939152 +0000 UTC"},
	{"09:51:20.939152am 2014-31-12", "2014-12-31 09:51:20.939152 +0000 UTC"},
	{"09:51:20.939152pm 2014-31-12", "2014-12-31 21:51:20.939152 +0000 UTC"},
}

func TestParseDateString(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	for i, elem := range testData {
		if tt, err := ParseDateString(elem.in); err != nil {
			t.Error(err)
		} else if elem.out != tt.String() {
			t.Errorf("[%d] Unexpected parsed time.\nExpect:\t%s\nGot:\t%s\n", i, elem.out, tt)
		}
	}
}
