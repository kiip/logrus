package logrus

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestPrettyFormatting(t *testing.T) {
	pf := &PrettyFormatter{DisableColors: true}

	testCases := []struct {
		value    string
		expected string
	}{
		{`foo`, "time=0001-01-01T00:00:00Z level=panic test=foo\n"},
	}

	for _, tc := range testCases {
		b, _ := pf.Format(WithField("test", tc.value))

		if string(b) != tc.expected {
			t.Errorf("formatting expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestPrettyQuoting(t *testing.T) {
	pf := &PrettyFormatter{DisableColors: true}

	checkQuoting := func(q bool, value interface{}) {
		b, _ := pf.Format(WithField("test", value))
		idx := bytes.Index(b, ([]byte)("test="))
		cont := bytes.Contains(b[idx+5:], []byte("\""))
		if cont != q {
			if q {
				t.Errorf("quoting expected for: %#v", value)
			} else {
				t.Errorf("quoting not expected for: %#v", value)
			}
		}
	}

	checkQuoting(false, "")
	checkQuoting(false, "abcd")
	checkQuoting(false, "v1.0")
	checkQuoting(false, "1234567890")
	checkQuoting(false, "/foobar")
	checkQuoting(false, "foo_bar")
	checkQuoting(false, "foo@bar")
	checkQuoting(false, "foobar^")
	checkQuoting(false, "+/-_^@f.oobar")
	checkQuoting(false, "foobar$")
	checkQuoting(false, "&foobar")
	checkQuoting(false, "x y")
	checkQuoting(false, "x,y")
	checkQuoting(false, errors.New("invalid"))
	checkQuoting(false, errors.New("invalid argument"))

	// Test for quoting empty fields.
	pf.QuoteEmptyFields = true
	checkQuoting(false, "")
	checkQuoting(false, "abcd")
	checkQuoting(false, errors.New("invalid argument"))
}

func TestPrettyEscaping(t *testing.T) {
	pf := &PrettyFormatter{DisableColors: true}

	testCases := []struct {
		value    string
		expected string
	}{
		{`ba"r`, `ba"r`},
		{`ba'r`, `ba'r`},
	}

	for _, tc := range testCases {
		b, _ := pf.Format(WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestPrettyDisableTimestampWithColoredOutput(t *testing.T) {
	pf := &PrettyFormatter{DisableTimestamp: true, ForceColors: true}

	b, _ := pf.Format(WithField("test", "test"))
	if strings.Contains(string(b), "[0000]") {
		t.Error("timestamp not expected when DisableTimestamp is true")
	}
}

// TODO add tests for sorting etc., this requires a parser for the text
// formatter output.
