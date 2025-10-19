# Release Guide

This guide explains how to create releases for the TMDB MCP Server project.

## Automated Release Process

This project uses GitHub Actions to automate the release process. When you push a git tag that matches the pattern `v*`, the following happens automatically:

1. **GitHub Release Creation**: A new release is created with the tag name
2. **Multi-Platform Binaries**: Pre-compiled binaries are built for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)
3. **Docker Images**: Multi-architecture Docker images are built and pushed to GitHub Container Registry

## Creating a Release

### Prerequisites

- Ensure you have the latest changes merged to `main` branch
- All tests should be passing
- Update version numbers in code if needed
- Update CHANGELOG.md (if exists)

### Steps

1. **Create and push a version tag**:
   ```bash
   # Make sure you're on the main branch and up to date
   git checkout main
   git pull origin main

   # Create a semantic version tag
   git tag v1.0.0

   # Push the tag to trigger the release workflow
   git push origin v1.0.0
   ```

2. **Monitor the release**:
   - Go to the [Actions tab](https://github.com/XDwanj/tmdb-mcp/actions) in GitHub
   - Watch the "Release" workflow run
   - The workflow will create the release and build all assets

3. **Verify the release**:
   - Check the [Releases page](https://github.com/XDwanj/tmdb-mcp/releases)
   - Verify all binaries are attached
   - Check that Docker images are available at `ghcr.io/XDwanj/tmdb-mcp`

## Version Tagging

### Semantic Versioning

Use semantic versioning for tags:
- **Major**: `v1.0.0` - Breaking changes
- **Minor**: `v1.1.0` - New features, backward compatible
- **Patch**: `v1.0.1` - Bug fixes, backward compatible

### Pre-releases

For pre-release versions:
```bash
# Alpha release
git tag v1.0.0-alpha.1

# Beta release
git tag v1.0.0-beta.1

# Release candidate
git tag v1.0.0-rc.1
```

## Release Assets

### Binary Downloads

Each release includes the following binaries:
- `tmdb-mcp-linux-amd64.tar.gz`
- `tmdb-mcp-linux-arm64.tar.gz`
- `tmdb-mcp-darwin-amd64.tar.gz`
- `tmdb-mcp-darwin-arm64.tar.gz`
- `tmdb-mcp-windows-amd64.zip`

### Docker Images

Docker images are available at:
- `ghcr.io/XDwanj/tmdb-mcp:latest` (always points to latest release)
- `ghcr.io/XDwanj/tmdb-mcp:v1.0.0` (version-specific)
- `ghcr.io/XDwanj/tmdb-mcp:v1.0.0-alpha.1` (pre-release versions)

All images support multiple architectures:
- `linux/amd64`
- `linux/arm64`

## Manual Release (If Needed)

If the automated process fails, you can manually create a release:

1. **Build binaries locally**:
   ```bash
   # Install goreleaser if not already installed
   go install github.com/goreleaser/goreleaser@latest

   # Build release
   goreleaser build --clean
   ```

2. **Create GitHub Release manually**:
   - Go to [Releases page](https://github.com/XDwanj/tmdb-mcp/releases)
   - Click "Create a new release"
   - Choose your tag
   - Upload the built binaries
   - Write release notes

3. **Build and push Docker image manually**:
   ```bash
   # Build for multiple platforms
   docker buildx build --platform linux/amd64,linux/arm64 \
     -t ghcr.io/XDwanj/tmdb-mcp:v1.0.0 \
     -t ghcr.io/XDwanj/tmdb-mcp:latest \
     --push .
   ```

## Troubleshooting

### Common Issues

1. **Workflow fails on tag push**:
   - Check the workflow logs for errors
   - Ensure the tag follows the `v*` pattern
   - Verify GitHub token permissions

2. **Docker build fails**:
   - Check if Dockerfile syntax is valid
   - Verify all build dependencies are available
   - Check for size limits (GitHub has limits)

3. **Missing assets**:
   - Verify all build jobs completed successfully
   - Check if any build steps timed out
   - Ensure file paths are correct

### Getting Help

- Check the [GitHub Actions logs](https://github.com/XDwanj/tmdb-mcp/actions)
- Review the [release workflow file](../.github/workflows/release.yml)
- Open an issue for any problems encountered

## Release Checklist

Before creating a release:

- [ ] All tests are passing
- [ ] Documentation is updated
- [ ] Version numbers are updated
- [ ] CHANGELOG is updated (if applicable)
- [ ] Code is reviewed and approved
- [ ] Security scan passes (if configured)

After release:

- [ ] Verify all assets are available
- [ ] Test Docker images
- [ ] Update documentation links if needed
- [ ] Announce the release (if applicable)