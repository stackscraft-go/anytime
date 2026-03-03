# Release Checklist

## Before tagging

- [ ] `go test ./...`
- [ ] `go vet ./...`
- [ ] Coverage remains at 100%
- [ ] README examples still compile logically
- [ ] CHANGELOG updated

## Tag and publish

```bash
git tag v0.1.0
git push origin v0.1.0
```

## After release

- [ ] Create GitHub release notes
- [ ] Verify package appears on pkg.go.dev
