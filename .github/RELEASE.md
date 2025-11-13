# Release Process

This document describes how to create a new release of Stumpf.Works NAS.

## Automated Release via GitHub Actions

The project uses GitHub Actions to automatically build and publish releases when you push a version tag.

### Creating a New Release

1. **Update version numbers** (if needed):
   - Backend: `backend/internal/updates/update_service.go` → `CurrentVersion`
   - Frontend: `frontend/package.json` → `version`
   - Main: `backend/cmd/stumpfworks-server/main.go` → `AppVersion`

2. **Commit your changes**:
   ```bash
   git add -A
   git commit -m "chore: bump version to v0.3.0"
   git push
   ```

3. **Create and push a tag**:
   ```bash
   git tag v0.3.0
   git push origin v0.3.0
   ```

4. **GitHub Actions will automatically**:
   - Build backend binaries for Linux AMD64 and ARM64
   - Build and package the frontend
   - Generate a changelog from commits since the last release
   - Create a GitHub Release with all binaries as attachments
   - Generate SHA256 checksums

5. **Check the release**:
   - Go to: https://github.com/Stumpf-works/stumpfworks-nas/releases
   - The new release should appear within 5-10 minutes

## Version Naming

Use semantic versioning (SemVer):
- **Major** (v1.0.0): Breaking changes
- **Minor** (v0.3.0): New features, backwards compatible
- **Patch** (v0.2.1): Bug fixes, backwards compatible

Examples:
- `v0.3.0` - Minor release with new features
- `v0.3.1` - Patch release with bug fixes
- `v1.0.0` - Major release (first stable version)

## Pre-releases

For beta or release candidate versions:
```bash
git tag v0.3.0-beta.1
git push origin v0.3.0-beta.1
```

The workflow will mark these as "pre-release" on GitHub.

## Manual Release (if needed)

If the automated workflow fails, you can manually create a release:

1. **Build Backend**:
   ```bash
   cd backend
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stumpfworks-nas-linux-amd64 ./cmd/stumpfworks-server
   CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o stumpfworks-nas-linux-arm64 ./cmd/stumpfworks-server
   ```

2. **Build Frontend**:
   ```bash
   cd frontend
   npm ci
   npm run build
   tar -czf stumpfworks-nas-frontend.tar.gz -C dist .
   ```

3. **Create Release on GitHub**:
   - Go to Releases → New Release
   - Tag: v0.3.0
   - Upload the binaries
   - Write release notes

## Troubleshooting

### Workflow fails with "permission denied"

Make sure the repository settings allow GitHub Actions to create releases:
- Settings → Actions → General → Workflow permissions
- Select "Read and write permissions"

### Tag already exists

Delete and recreate the tag:
```bash
git tag -d v0.3.0
git push origin :refs/tags/v0.3.0
git tag v0.3.0
git push origin v0.3.0
```

### Build fails

Check the GitHub Actions logs:
- Go to Actions tab
- Click on the failed workflow run
- Check the logs for errors

Common issues:
- Missing dependencies in package.json
- Go module issues (run `go mod tidy`)
- TypeScript compilation errors
