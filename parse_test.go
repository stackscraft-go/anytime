package anytime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantLayout string
		wantString string
		wantErr    bool
	}{
		{name: "rfc3339", input: "2026-03-01T14:15:16Z", wantLayout: time.RFC3339Nano, wantString: "2026-03-01T14:15:16Z"},
		{name: "date only", input: "2026-03-01", wantLayout: "2006-01-02", wantString: "2026-03-01"},
		{name: "custom datetime", input: "2026/03/01 09:10:11", wantLayout: "2006/01/02 15:04:05", wantString: "2026/03/01 09:10:11"},
		{name: "month name", input: "Mar 1, 2026", wantLayout: "Jan 2, 2006", wantString: "Mar 1, 2026"},
		{name: "dot datetime", input: "02.1.2026 15:04:05", wantLayout: "02.1.2006 15:04:05", wantString: "02.1.2026 15:04:05"},
		{name: "datetime with tz name", input: "2026-03-01 14:15:16 UTC", wantLayout: "2006-01-02 15:04:05 MST", wantString: "2026-03-01 14:15:16 UTC"},
		{name: "non padded datetime", input: "2026-3-1 5:4:3", wantLayout: "2006-1-2 15:4:5", wantString: "2026-3-1 05:4:3"},
		{name: "compact iso tz", input: "20260301T141516Z", wantLayout: "20060102T150405Z", wantString: "20260301T141516Z"},
		{name: "year month text", input: "Jan 2026", wantLayout: "Jan 2006", wantString: "Jan 2026"},
		{name: "time only", input: "14:15:16", wantLayout: time.TimeOnly, wantString: "14:15:16"},
		{name: "edge case before year", input: "Before 2026", wantLayout: "Before 2006", wantString: "Before 2026"},
		{name: "unix millis", input: "1700000000123", wantLayout: time.RFC3339Nano, wantString: "2023-11-14T22:13:20.123Z"},
		{name: "unix float seconds", input: "1700000000.123", wantLayout: time.RFC3339Nano, wantString: "2023-11-14T22:13:20.123Z"},
		{name: "invalid", input: "not-a-date", wantErr: true},
		{name: "empty", input: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.wantLayout, parsed.Layout())
			assert.Equal(t, tt.wantString, parsed.String())
		})
	}
}

func TestTimeImmutability(t *testing.T) {
	tests := []struct {
		name       string
		baseInput  string
		transform  func(Time) Time
		wantBase   string
		wantResult string
	}{
		{name: "add date", baseInput: "2026-03-01", transform: func(v Time) Time { return v.AddDate(0, 0, 5) }, wantBase: "2026-03-01", wantResult: "2026-03-06"},
		{name: "with layout", baseInput: "2026-03-01", transform: func(v Time) Time { return v.WithLayout("02/01/2006") }, wantBase: "2026-03-01", wantResult: "01/03/2026"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, err := Parse(tt.baseInput)
			if !assert.NoError(t, err) {
				return
			}

			result := tt.transform(*base)
			assert.Equal(t, tt.wantBase, base.String())
			assert.Equal(t, tt.wantResult, result.String())
		})
	}
}

func TestMarshalRespectsLayout(t *testing.T) {
	tests := []struct {
		name   string
		encode func(Time) (string, error)
		want   string
	}{
		{name: "marshal text", encode: func(v Time) (string, error) { b, err := v.MarshalText(); return string(b), err }, want: "01/03/2026"},
		{name: "marshal json", encode: func(v Time) (string, error) { b, err := v.MarshalJSON(); return string(b), err }, want: `"01/03/2026"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := Parse("2026-03-01")
			if !assert.NoError(t, err) {
				return
			}

			got, err := tt.encode(parsed.WithLayout("02/01/2006"))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUnmarshalUsesParseBehavior(t *testing.T) {
	tests := []struct {
		name       string
		decode     func(*Time) error
		wantLayout string
		wantString string
	}{
		{
			name: "unmarshal text",
			decode: func(v *Time) error {
				return v.UnmarshalText([]byte("2026/03/01 09:10:11"))
			},
			wantLayout: "2006/01/02 15:04:05",
			wantString: "2026/03/01 09:10:11",
		},
		{
			name: "unmarshal json",
			decode: func(v *Time) error {
				return json.Unmarshal([]byte(`"1700000000123"`), v)
			},
			wantLayout: time.RFC3339Nano,
			wantString: "2023-11-14T22:13:20.123Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Time
			if !assert.NoError(t, tt.decode(&got)) {
				return
			}

			assert.Equal(t, tt.wantLayout, got.Layout())
			assert.Equal(t, tt.wantString, got.String())
		})
	}
}

func TestParseUnixInternal(t *testing.T) {
	tests := []struct {
		name  string
		input string
		ok    bool
	}{
		{name: "seconds", input: "1700000000", ok: true},
		{name: "seconds plus sign", input: "+1700000000", ok: true},
		{name: "seconds minus sign", input: "-1700000000", ok: true},
		{name: "milliseconds", input: "1700000000123", ok: true},
		{name: "microseconds", input: "1700000000123456", ok: true},
		{name: "nanoseconds", input: "1700000000123456789", ok: true},
		{name: "invalid length", input: "17000000001", ok: false},
		{name: "invalid integer", input: "abc", ok: false},
		{name: "float seconds", input: "1700000000.123", ok: true},
		{name: "float plus sign", input: "+1700000000.1", ok: true},
		{name: "float minus sign", input: "-1700000000.1", ok: true},
		{name: "float missing frac", input: "1700000000.", ok: false},
		{name: "float missing secs", input: ".123", ok: false},
		{name: "float non-digit frac", input: "1700000000.a", ok: false},
		{name: "float overflow secs", input: "92233720368547758070.1", ok: false},
		{name: "float long frac", input: "1700000000.1234567899", ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := parseUnix(tt.input)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestAllDigits(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "empty", input: "", want: false},
		{name: "digits", input: "123456", want: true},
		{name: "non digits", input: "12a", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, allDigits(tt.input))
		})
	}
}
