# Update Summary for platform (data service)

## Changes Made

### 1. Go Version Update
- ✅ Updated Go version to 1.26 in go.mod
- Command used: `go mod edit -go=1.26 && go mod tidy`

### 2. CGO Disabled
- ✅ Updated Makefile to set CGO_ENABLED=0 for all builds
- Changed from CGO_ENABLED=1 to CGO_ENABLED=0 in the else branch

### 3. Dockerfile Updates
- ✅ Updated Dockerfile.data to use distroless base image (gcr.io/distroless/static:nonroot)
- ✅ Changed FROM --platform=$BUILDPLATFORM alpine:latest to FROM gcr.io/distroless/static:nonroot
- ✅ Updated user from tidepool to nonroot
- ✅ Removed all apk commands (apk update, upgrade, add ca-certificates tzdata, adduser)
- ✅ No gcc or musl-dev dependencies present

### 4. CHANGELOG.md
- ✅ Added new version entry 0.31.0 with increment to minor version
- Documented all changes made

## Rationale

### Go 1.26
- Latest stable version with performance improvements and security fixes
- Better build times and runtime performance

### Distroless Images
- Smaller attack surface (no shell, package manager, or unnecessary utilities)
- Reduced image size
- Better security posture
- Meets production best practices

### CGO Disabled
- Fully static binaries
- Better portability
- No C library dependencies
- Simplified deployment

## Date
2026-03-27

## Notes
- Platform uses a Makefile-based build system
- The main Dockerfile is Dockerfile.data for the data service
- CGO was previously enabled but has now been disabled for all builds

