# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- Initial release of `anytime`.
- `Parse(string) (*Time, error)` with broad layout support and unix timestamps.
- Immutable `Time` wrapper with layout-aware `String`, `MarshalText`,
  `MarshalJSON`, `UnmarshalText`, and `UnmarshalJSON`.
- Comprehensive table-driven tests with 100% statement coverage.
