# Release v5.0.1: Major Dependency Updates and Security Improvements

## 🚀 What's New

This release focuses on major dependency updates, security improvements, and modernization of the codebase.

## 📦 Major Updates

### Go Runtime
- **Updated Go from 1.19 to 1.23.0** with toolchain go1.23.2
- Improved performance and latest security patches
- Enhanced build system compatibility

### Docker Integration
- **Updated Docker client from v20.10.22 to v28.3.0** (breaking change)
- Fixed Docker API compatibility issues
- Updated import paths and API calls for modern Docker client
- Enhanced container management capabilities

### Dependencies
- **Updated fatih/color** from v1.13.0 to v1.18.0
- **Updated golang.org/x/crypto** to v0.39.0 for latest security patches
- **Updated sirupsen/logrus** from v1.9.0 to v1.9.3
- **Updated all golang.org/x/* packages** to latest versions
- Added new OpenTelemetry dependencies for enhanced observability

## 🔒 Security Improvements

- ✅ Latest security patches for all dependencies
- ✅ Updated Go runtime with security fixes
- ✅ Modern Docker client with improved security features
- ✅ Updated crypto libraries with latest algorithms
- ✅ Addressed multiple vulnerability reports from Dependabot

## 🛠️ Technical Changes

### Code Updates
- Fixed Docker API imports and types in `pkg/docker/docker.go`
- Enhanced error handling and test stability
- Improved map iteration handling in tests
- Updated release workflow to use Go 1.23

### Documentation
- Added comprehensive dependency update documentation
- Moved documentation to organized `docs/` structure
- Enhanced release notes and change tracking

## ✅ Quality Assurance

- **All tests pass** with the updated dependencies
- **Project builds successfully** on all platforms
- **Backward compatibility maintained** for all YAML configurations and CLI options
- **No breaking changes** to public API

## 🔄 Compatibility

- **Minimum Go version**: Now requires Go 1.21+ (previously 1.19+)
- **Docker API**: Compatible with latest Docker Engine versions
- **YAML configurations**: All existing configurations remain unchanged
- **CLI options**: All command-line options work as before

## 📋 Files Changed

- `go.mod` and `go.sum`: Updated all dependencies
- `pkg/docker/docker.go`: Fixed Docker API compatibility
- `pkg/errors/errors_test.go`: Enhanced test stability
- `.github/workflows/release.yaml`: Updated to Go 1.23
- `docs/DEPENDENCY_UPDATES.md`: Added comprehensive documentation

## 🚨 Breaking Changes

The Docker client update includes breaking changes in the API, but these are handled internally and don't affect end users. All existing functionality remains the same.

## 📖 Documentation

For detailed information about the dependency updates, see:
- [docs/DEPENDENCY_UPDATES.md](docs/DEPENDENCY_UPDATES.md)
- [docs/README.md](docs/README.md)

## 🙏 Contributors

Thanks to all contributors and the community for reporting issues and suggesting improvements!

---

**Full Changelog**: https://github.com/LANIXX-HH/runfromyaml/compare/v5.0.0...v5.0.1
