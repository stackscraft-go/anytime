# Contributing

Thanks for your interest in contributing to `anytime`.

## Development

1. Fork the repository and create a feature branch.
2. Keep changes focused and add/update tests.
3. Run checks locally:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go vet ./...
```

The project currently enforces full statement coverage in CI.

## Pull requests

- Describe the motivation and behavior change.
- Include examples for new parsing layouts or behaviors.
- Avoid breaking public APIs unless discussed first.

## Reporting issues

Please include:

- Go version
- Input value(s)
- Expected behavior
- Actual behavior
