---
name: release-workflow
description: >-
  Use this skill when preparing, validating, or cutting a new version release of shellescape and escargs using GoReleaser and GitHub Actions.
---

# Release Management Runbook for `shellescape`

This skill provides step-by-step instructions for releasing a new version of `shellescape` and `escargs`.

## 1. Pre-Release Verification Checklist

Before tagging a release, run all verification steps:

- [ ] Ensure all working trees are clean: `git status`
- [ ] Run full test suite with race detector: `go test -v -race ./...`
- [ ] Run benchmarks to verify no performance regressions: `go test -v -bench=. -benchmem ./...`
- [ ] Run linter: `golangci-lint run ./...`
- [ ] Verify GoReleaser configuration syntax:
  ```bash
  goreleaser check
  ```

## 2. Test Local Snapshot Build

Build snapshot binaries for all supported target platforms without publishing:

```bash
goreleaser release --snapshot --clean
```

Inspect the generated binaries in `dist/`:
- Linux: `dist/escargs_linux_amd64_v1/escargs`, `dist/escargs_linux_arm64/escargs`
- macOS: `dist/escargs_darwin_amd64_v1/escargs`, `dist/escargs_darwin_arm64/escargs`
- Windows: `dist/escargs_windows_amd64_v1/escargs.exe`
- FreeBSD: `dist/escargs_freebsd_amd64_v1/escargs`

Test the local snapshot binary:
```bash
./dist/escargs_darwin_arm64/escargs -V
```

Clean up snapshot artifacts:
```bash
rm -rf dist/
```

## 3. Tagging & Releasing

1. Determine the next Semantic Version according to changes (e.g. `v1.2.1` for patch, `v1.3.0` for minor feature).
2. Create an annotated git tag:
   ```bash
   git tag -a v1.3.0 -m "Release v1.3.0"
   ```
3. Push the tag to GitHub:
   ```bash
   git push origin v1.3.0
   ```
4. Monitor the GitHub Actions release workflow (`.github/workflows/release.yaml`).
5. Verify the published GitHub Release and checksums.
