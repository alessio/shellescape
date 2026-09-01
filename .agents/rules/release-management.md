# Rule: Release Management & Versioning

## Semantic Versioning
- This repository follows [Semantic Versioning 2.0.0](https://semver.org/):
  - **PATCH** (`v1.x.Y`): Bug fixes, internal performance enhancements, test improvements, or dependency updates with no API changes.
  - **MINOR** (`v1.X.0`): New backward-compatible functions, CLI flags, or enhancements.
  - **MAJOR** (`vX.0.0`): Breaking changes to the public API signature or behavior.

## GoReleaser Configuration
- Builds the `escargs` binary across platforms:
  - OS: Linux, Windows, macOS (`darwin`), FreeBSD.
  - Arch: `amd64`, `arm64`, `arm` (v6, v7), `amd64` (v2, v3).
  - CGO is disabled (`CGO_ENABLED=0`).
- Generates `checksums.txt` and GitHub release notes automatically.
- Changelog excludes commit titles starting with `docs:` or `test:`.

## Release Workflow Execution
1. Ensure all tests, linters, and build checks pass on `master`.
2. Check `.goreleaser.yml` configuration:
   ```bash
   goreleaser check
   ```
3. Test a snapshot release locally:
   ```bash
   goreleaser release --snapshot --clean
   ```
4. Create and push a semantic git tag (e.g. `v1.2.0`):
   ```bash
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin v1.2.0
   ```
5. GitHub Actions (`.github/workflows/release.yaml`) triggers GoReleaser to publish the artifacts.
