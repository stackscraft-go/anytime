package anytime

import (
	"bytes"
	"encoding/json"
	"time"
)

const defaultLayout = time.RFC3339

// Time is an immutable wrapper around time.Time that keeps track of the
// successful input format used by Parse.
type Time struct {
	v      time.Time
	layout string
}

// FromTime creates a Time from a standard time.Time using RFC3339 as the
// default string format.
func FromTime(t time.Time) Time {
	return Time{
		v:      t,
		layout: defaultLayout,
	}
}

// Now returns the current local time wrapped in Time.
func Now() Time {
	return FromTime(time.Now())
}

// Time returns a copy of the wrapped standard library time.Time.
func (t Time) Time() time.Time {
	return t.v
}

// Layout returns the Go time layout used when formatting layout-based values.
func (t Time) Layout() string {
	if t.layout == "" {
		return defaultLayout
	}
	return t.layout
}

// WithLayout returns a new Time with the same instant but a different layout
// for String formatting.
func (t Time) WithLayout(layout string) Time {
	if layout == "" {
		layout = defaultLayout
	}
	t.layout = layout
	return t
}

// String formats using the current tracked parse format.
func (t Time) String() string {
	return t.v.Format(t.Layout())
}

func (t Time) GoString() string { return t.String() }

func (t Time) IsZero() bool                  { return t.v.IsZero() }
func (t Time) Before(u Time) bool            { return t.v.Before(u.v) }
func (t Time) After(u Time) bool             { return t.v.After(u.v) }
func (t Time) Equal(u Time) bool             { return t.v.Equal(u.v) }
func (t Time) Compare(u Time) int            { return t.v.Compare(u.v) }
func (t Time) Sub(u Time) time.Duration      { return t.v.Sub(u.v) }
func (t Time) Add(d time.Duration) Time      { t.v = t.v.Add(d); return t }
func (t Time) AddDate(y, m, d int) Time      { t.v = t.v.AddDate(y, m, d); return t }
func (t Time) Round(d time.Duration) Time    { t.v = t.v.Round(d); return t }
func (t Time) Truncate(d time.Duration) Time { t.v = t.v.Truncate(d); return t }
func (t Time) UTC() Time                     { t.v = t.v.UTC(); return t }
func (t Time) Local() Time                   { t.v = t.v.Local(); return t }
func (t Time) In(loc *time.Location) Time    { t.v = t.v.In(loc); return t }
func (t Time) Zone() (string, int)           { return t.v.Zone() }
func (t Time) Location() *time.Location      { return t.v.Location() }
func (t Time) Clock() (hour, min, sec int)   { return t.v.Clock() }
func (t Time) Date() (year int, month time.Month, day int) {
	return t.v.Date()
}
func (t Time) Year() int                      { return t.v.Year() }
func (t Time) Month() time.Month              { return t.v.Month() }
func (t Time) Day() int                       { return t.v.Day() }
func (t Time) Weekday() time.Weekday          { return t.v.Weekday() }
func (t Time) ISOWeek() (int, int)            { return t.v.ISOWeek() }
func (t Time) YearDay() int                   { return t.v.YearDay() }
func (t Time) Hour() int                      { return t.v.Hour() }
func (t Time) Minute() int                    { return t.v.Minute() }
func (t Time) Second() int                    { return t.v.Second() }
func (t Time) Nanosecond() int                { return t.v.Nanosecond() }
func (t Time) Unix() int64                    { return t.v.Unix() }
func (t Time) UnixMilli() int64               { return t.v.UnixMilli() }
func (t Time) UnixMicro() int64               { return t.v.UnixMicro() }
func (t Time) UnixNano() int64                { return t.v.UnixNano() }
func (t Time) MarshalBinary() ([]byte, error) { return t.v.MarshalBinary() }
func (t Time) MarshalJSON() ([]byte, error)   { return json.Marshal(t.String()) }
func (t Time) MarshalText() ([]byte, error)   { return []byte(t.String()), nil }

func (t *Time) UnmarshalText(data []byte) error {
	parsed, err := Parse(string(data))
	if err != nil {
		return err
	}
	*t = *parsed
	return nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*t = Time{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return t.UnmarshalText([]byte(s))
}
