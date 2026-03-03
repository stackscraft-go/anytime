// Package anytime provides flexible parsing for date and datetime strings.
//
// The main entrypoint is Parse, which attempts a curated set of common layouts,
// plus unix timestamps (seconds, milliseconds, microseconds, nanoseconds, and
// fractional unix seconds). The returned Time value is immutable and stores the
// successful layout used for parsing.
//
// String, MarshalText, and MarshalJSON always render using the stored layout.
// For unix timestamp inputs, the stored layout defaults to time.RFC3339Nano.
package anytime
