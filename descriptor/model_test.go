package descriptor

import (
	"testing"
)

type ParseTest struct {
	in  string
	out Format
}

var parseTests = []ParseTest{
	// Test noop
	{"arw", Format("arw")},
	// Test lower casing
	{"Cr2", Format("cr2")},
	{"CRW", Format("crw")},
	{"NeF", Format("nef")},
	// Test stripping
	{"arw ", Format("arw")},
	{" arw", Format("arw")},
	// Test multiple operations
	{" 	\nNeF   \t ", Format("nef")},
}

func TestParseFormat(t *testing.T) {
	for _, test := range parseTests {
		actual := ParseFormat(test.in)
		if actual != test.out {
			t.Errorf("ParseFormat(%v) = %v; want %v", test.in, actual, test.out)
		}
	}
}
