package anytime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeMethods(t *testing.T) {
	assert.False(t, Now().IsZero())

	baseStd := time.Date(2026, 3, 1, 14, 15, 16, 123456789, time.FixedZone("X", 3600))
	base := FromTime(baseStd)
	other := FromTime(baseStd.Add(2 * time.Hour))

	assert.Equal(t, baseStd, base.Time())
	assert.Equal(t, base.String(), base.GoString())
	assert.False(t, base.IsZero())
	assert.True(t, (Time{}).IsZero())
	assert.True(t, base.Before(other))
	assert.True(t, other.After(base))
	assert.True(t, base.Equal(base))
	assert.Equal(t, 0, base.Compare(base))
	assert.Equal(t, 2*time.Hour, other.Sub(base))

	assert.Equal(t, baseStd.Add(time.Minute), base.Add(time.Minute).Time())
	assert.Equal(t, baseStd.AddDate(1, 0, 0), base.AddDate(1, 0, 0).Time())
	assert.Equal(t, baseStd.Round(time.Second), base.Round(time.Second).Time())
	assert.Equal(t, baseStd.Truncate(time.Second), base.Truncate(time.Second).Time())
	assert.Equal(t, baseStd.UTC(), base.UTC().Time())
	assert.Equal(t, baseStd.Local(), base.Local().Time())

	ny, err := time.LoadLocation("America/New_York")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, baseStd.In(ny), base.In(ny).Time())

	gn, goff := base.Zone()
	wn, woff := baseStd.Zone()
	assert.Equal(t, wn, gn)
	assert.Equal(t, woff, goff)
	assert.Equal(t, baseStd.Location(), base.Location())

	gh, gm, gs := base.Clock()
	wh, wm, ws := baseStd.Clock()
	assert.Equal(t, wh, gh)
	assert.Equal(t, wm, gm)
	assert.Equal(t, ws, gs)

	gy, gmon, gd := base.Date()
	wy, wmon, wd := baseStd.Date()
	assert.Equal(t, wy, gy)
	assert.Equal(t, wmon, gmon)
	assert.Equal(t, wd, gd)

	assert.Equal(t, baseStd.Year(), base.Year())
	assert.Equal(t, baseStd.Month(), base.Month())
	assert.Equal(t, baseStd.Day(), base.Day())
	assert.Equal(t, baseStd.Weekday(), base.Weekday())
	giy, giw := base.ISOWeek()
	wiy, wiw := baseStd.ISOWeek()
	assert.Equal(t, wiy, giy)
	assert.Equal(t, wiw, giw)
	assert.Equal(t, baseStd.YearDay(), base.YearDay())
	assert.Equal(t, baseStd.Hour(), base.Hour())
	assert.Equal(t, baseStd.Minute(), base.Minute())
	assert.Equal(t, baseStd.Second(), base.Second())
	assert.Equal(t, baseStd.Nanosecond(), base.Nanosecond())
	assert.Equal(t, baseStd.Unix(), base.Unix())
	assert.Equal(t, baseStd.UnixMilli(), base.UnixMilli())
	assert.Equal(t, baseStd.UnixMicro(), base.UnixMicro())
	assert.Equal(t, baseStd.UnixNano(), base.UnixNano())

	_, err = base.MarshalBinary()
	assert.NoError(t, err)
}

func TestLayoutFallbackAndWithLayoutEmpty(t *testing.T) {
	var zero Time
	assert.Equal(t, defaultLayout, zero.Layout())

	parsed, err := Parse("2026-03-01")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, defaultLayout, parsed.WithLayout("").Layout())
}

func TestUnmarshalJSONBranches(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantNil bool
	}{
		{name: "null", input: "null", wantNil: true},
		{name: "invalid json", input: "{", wantErr: true},
		{name: "invalid date string", input: `"bad-date"`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Time
			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			if tt.wantNil {
				assert.True(t, got.IsZero())
			}
		})
	}
}

func TestUnmarshalJSONDirect(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{name: "success", input: []byte(`"2026-03-01"`)},
		{name: "trimmed null", input: []byte("  null\n")},
		{name: "invalid json syntax", input: []byte("{"), wantErr: true},
		{name: "invalid string payload", input: []byte(`"not-a-date"`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Time
			err := got.UnmarshalJSON(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestUnmarshalJSONNullResetsValue(t *testing.T) {
	parsed, err := Parse("2026-03-01")
	if !assert.NoError(t, err) {
		return
	}

	got := *parsed
	assert.NoError(t, got.UnmarshalJSON([]byte("null")))
	assert.True(t, got.IsZero())
}
