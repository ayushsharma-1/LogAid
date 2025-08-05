# Release Checklist

This checklist ensures all necessary steps are completed before releasing a new version of LogAid.

## Pre-Release Checklist

### 📋 Code Quality
- [ ] All tests pass (`make test`)
- [ ] Race condition tests pass (`make test-race`)
- [ ] Code coverage above 80% (`make test-cover`)
- [ ] Linting passes (`make lint`)
- [ ] Code formatted (`make fmt`)
- [ ] Go vet passes (`make vet`)

### 📚 Documentation
- [ ] README.md is up to date
- [ ] CHANGELOG.md has new version entry
- [ ] ARCHITECTURE.md reflects current state
- [ ] CONTRIBUTING.md is current
- [ ] All public functions have godoc comments
- [ ] Examples in documentation work

### 🔧 Configuration
- [ ] .env.example includes all options
- [ ] Default configuration is sensible
- [ ] Version information is correct
- [ ] Build flags are appropriate

### 🧪 Testing
- [ ] Manual testing on Linux
- [ ] Manual testing on macOS
- [ ] Manual testing on Windows
- [ ] Docker image builds and runs
- [ ] Installation script works
- [ ] All plugins tested
- [ ] AI integration tested (if keys available)

### 📦 Build System
- [ ] Makefile targets work
- [ ] GitHub Actions workflows validate
- [ ] Multi-platform builds succeed
- [ ] Docker builds succeed
- [ ] Archives are created correctly

## Release Process

### 1. Version Preparation
```bash
# Update version in relevant files
# Update CHANGELOG.md
# Commit all changes
git add .
git commit -m "chore: prepare release v1.0.0"
```

### 2. Create Release Tag
```bash
# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 3. Verify Automated Release
- [ ] GitHub Actions release workflow runs
- [ ] All platform binaries are built
- [ ] GitHub release is created
- [ ] Docker image is published
- [ ] Release notes are generated

### 4. Post-Release Verification
- [ ] Installation script downloads new version
- [ ] Docker image works with new tag
- [ ] Release artifacts are accessible
- [ ] Documentation links work

### 5. Announce Release
- [ ] Update project documentation
- [ ] Social media announcement (if applicable)
- [ ] Notify users/contributors

## Emergency Hotfix Process

### If critical issues are found:

1. **Create hotfix branch**
```bash
git checkout -b hotfix/v1.0.1 v1.0.0
```

2. **Apply minimal fix**
```bash
# Make necessary changes
git commit -m "fix: critical issue description"
```

3. **Test thoroughly**
```bash
make test
make test-race
# Manual testing
```

4. **Release hotfix**
```bash
git tag -a v1.0.1 -m "Hotfix v1.0.1"
git push origin v1.0.1
```

5. **Merge back to main**
```bash
git checkout main
git merge hotfix/v1.0.1
git push origin main
```

## Version Strategy

- **MAJOR** (1.0.0 → 2.0.0): Breaking changes
- **MINOR** (1.0.0 → 1.1.0): New features, backward compatible
- **PATCH** (1.0.0 → 1.0.1): Bug fixes, backward compatible

## Supported Platforms

### Primary Platforms (Tested)
- Linux AMD64
- Linux ARM64
- macOS AMD64 (Intel)
- macOS ARM64 (Apple Silicon)
- Windows AMD64

### Distribution Channels
- GitHub Releases (binaries)
- Docker Hub / GitHub Container Registry
- Go package registry
- Installation script

## Rollback Procedure

If a release has critical issues:

1. **Remove problematic release**
```bash
gh release delete v1.0.0
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```

2. **Revert to previous working version**
3. **Communicate issue to users**
4. **Prepare fixed release**

## Quality Gates

### Automated Checks (CI/CD)
- Unit tests must pass
- Integration tests must pass
- Security scans must pass
- Multi-platform builds must succeed
- Docker builds must succeed

### Manual Checks
- Installation script tested on clean system
- Key features manually verified
- Documentation accuracy confirmed
- Performance regression check

---

**Note**: This checklist should be updated as the project evolves to ensure it remains comprehensive and accurate.
