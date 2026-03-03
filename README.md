# anytime

[![CI](https://github.com/stackscraft-go/anytime/actions/workflows/ci.yml/badge.svg)](https://github.com/stackscraft-go/anytime/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/stackscraft-go/anytime.svg)](https://pkg.go.dev/github.com/stackscraft-go/anytime)
[![Go Report Card](https://goreportcard.com/badge/github.com/stackscraft-go/anytime)](https://goreportcard.com/report/github.com/stackscraft-go/anytime)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`anytime` is a small Go package for parsing real-world date and datetime
strings into an immutable value that remembers how it was parsed.

## Install

```bash
go get github.com/stackscraft-go/anytime
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/stackscraft-go/anytime"
)

func main() {
	t, err := anytime.Parse("2026-03-01 14:15:16")
	if err != nil {
		panic(err)
	}

	// String uses the successful parse format.
	fmt.Println(t.String())

	// Immutable operations return a new value.
	nextWeek := t.AddDate(0, 0, 7)
	fmt.Println(t.String())
	fmt.Println(nextWeek.String())
}
```

## Parse behavior

`Parse` supports many common layouts, including:

- RFC3339 / RFC3339Nano and other Go RFC layouts
- Common numeric date layouts (`YYYY-MM-DD`, `YYYY/MM/DD`, `DD/MM/YYYY`, etc)
- Common datetime layouts with and without fractional seconds
- Unix timestamps in seconds, milliseconds, microseconds, and nanoseconds

The chosen layout is preserved and can be inspected with:

- `Layout()`

When parsing unix timestamps, the stored layout defaults to `time.RFC3339Nano`.

## JSON/Text behavior

- `String()` formats with the current layout
- `MarshalText()` and `MarshalJSON()` also format with the current layout
- `UnmarshalText()` and `UnmarshalJSON()` delegate to `Parse`

## Immutability

`anytime.Time` has no exported fields and all transformation methods (such as `Add`, `AddDate`, `UTC`, `Round`, `WithLayout`) return a new value.

Use `Time()` when you need a plain `time.Time` copy for interoperability.
