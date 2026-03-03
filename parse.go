package anytime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var layoutCandidates = []string{
	// Common ISO and API formats (highest priority)
	time.RFC3339Nano,
	time.RFC3339,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006/01/02 15:04:05",
	"2006/01/02 15:04",
	"2006-01-02",
	"2006/01/02",

	// Common regional formats
	"02/01/2006 15:04:05",
	"02/01/2006",
	"01/02/2006",
	"02-01-2006",
	"01-02-2006",

	// Date, time & time zone formats
	"2006-01-02T15:04:05Z",
	"2006-01-02 15:04:05-07",
	"2006/01/02 15:04:05-07",
	"2006-01-02 15:04:05 -0700",
	"2006/01/02 15:04 -0700",
	"2006-01-02T15:04:05-0700",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 MST",
	"2006/01/02 15:04 MST",
	"2006/01/02 15:04:05 MST",
	"2006-01-02 15:04:05 (MST+3)",
	"2006/01/02 15:04:05 (MST+3)",
	"2006-01-02 15:04:05 (MST)",
	"2006/01/02 15:04:05 (MST)",
	"02-Jan-2006 15:04:05 MST",
	"02-01-2006 15:04:05 MST+3",
	"02-01-2006 15:04:05 MST",
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.UnixDate,
	time.RubyDate,

	// Fractional seconds
	"2006-01-02 15:04:05.999",
	"2006-01-02 15:04:05.999999",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02T15:04:05.999",
	"2006-01-02T15:04:05.999999",
	"2006-01-02T15:04:05.999999999",

	// Compact ISO-like formats
	"20060102T150405",
	"20060102T150405Z",
	"20060102T150405-0700",
	"20060102150405",

	// Numeric compact and partial date formats
	"20060102 15:04:05",
	"20060102",
	"200601",
	"2006-01",

	// Textual and month-name formats
	"Jan 2, 2006",
	"January 2, 2006",
	"Jan 2, 2006 15:04",
	"Jan 2, 2006 15:04:05",
	"02-Jan-2006",
	"02-Jan-2006 15:04:05",
	"02 Jan 2006",
	"Mon Jan 2 2006",
	"January _2 2006",
	"2006-Jan-02",
	"2006-Jan-02.",
	"Jan 2006",
	"January 2006",

	// Dotted and non-padded variants
	"2006.01.02 15:04:05",
	"02.01.2006 15:04:05",
	"02.1.2006 15:04:05",
	"2.1.2006 15:04:05",
	"2006.01.02",
	"2006. 01. 02.",
	"02.01.2006",
	"2006-1-2",
	"2006-01-02+15:04",
	"2006-1-2 15:4",
	"2006-1-2 15:4:5",
	"2/1/2006",
	"1/2/2006",
	time.ANSIC,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,

	// Edge cases
	"before Jan-2006",
	"Before 20060102",
	"Before 200601",
	"Before 2006",
	"before 20060102",
	"before 200601",
	"before 2006",
}

// Parse parses many common date/datetime formats and unix timestamps.
//
// The returned Time keeps track of the chosen successful parse format.
func Parse(s string) (*Time, error) {
	raw := strings.TrimSpace(s)
	if raw == "" {
		return nil, fmt.Errorf("anytime: cannot parse empty input")
	}

	if parsed, ok := parseUnix(raw); ok {
		return &parsed, nil
	}

	for _, candidate := range layoutCandidates {
		parsed, err := time.Parse(candidate, raw)
		if err == nil {
			t := Time{v: parsed, layout: candidate}
			return &t, nil
		}
	}

	return nil, fmt.Errorf("anytime: unsupported date/datetime format: %q", s)
}

func parseUnix(s string) (Time, bool) {
	if strings.Contains(s, ".") {
		return parseUnixFloatSeconds(s)
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Time{}, false
	}

	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}

	var t time.Time
	switch len(s) {
	case 10:
		t = time.Unix(v, 0).UTC()
	case 13:
		t = time.UnixMilli(v).UTC()
	case 16:
		t = time.UnixMicro(v).UTC()
	case 19:
		t = time.Unix(0, v).UTC()
	default:
		return Time{}, false
	}

	return Time{v: t, layout: time.RFC3339Nano}, true
}

func parseUnixFloatSeconds(s string) (Time, bool) {
	raw := s
	sign := int64(1)
	if strings.HasPrefix(raw, "-") {
		raw = raw[1:]
		sign = -1
	} else if strings.HasPrefix(raw, "+") {
		raw = raw[1:]
	}

	parts := strings.SplitN(raw, ".", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Time{}, false
	}
	if !allDigits(parts[0]) || !allDigits(parts[1]) {
		return Time{}, false
	}

	secAbs, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Time{}, false
	}

	frac := parts[1]
	if len(frac) > 9 {
		frac = frac[:9]
	}
	for len(frac) < 9 {
		frac += "0"
	}

	fracAbs, _ := strconv.ParseInt(frac, 10, 64)

	sec := sign * secAbs
	nsec := sign * fracAbs

	t := time.Unix(sec, nsec).UTC()
	return Time{v: t, layout: time.RFC3339Nano}, true
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
