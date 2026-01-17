# Dependency Inventory

This document provides a comprehensive inventory of all project dependencies, their current versions, health status, and update recommendations.

Last updated: 2026-01-08

## Table of Contents

- [Summary](#summary)
- [Security Status](#security-status)
- [Go Dependencies](#go-dependencies)
- [NPM Dependencies](#npm-dependencies)
- [Update Strategy](#update-strategy)
- [Related Issues](#related-issues)

---

## Summary

| Category | Total | Current | Outdated | Vulnerable |
|----------|-------|---------|----------|------------|
| Go Direct | 18 | 3 | 15 | TBD |
| Go Indirect | ~60 | - | - | TBD |
| NPM Production | 21 | 15 | 6 | 1 |
| NPM Development | 25 | 22 | 3 | 0 |

---

## Security Status

### Known Vulnerabilities

| Package | Type | Severity | CVE/Advisory | Status |
|---------|------|----------|--------------|--------|
| quill | npm | Moderate | GHSA-4943-9vgg-gr5r (XSS) | Deferred - requires breaking upgrade |

### Recently Fixed (2026-01-08)

| Package | Severity | Fix Applied |
|---------|----------|-------------|
| @modelcontextprotocol/sdk | High | npm audit fix |
| qs | High | npm audit fix |
| semver | High | @angular-eslint update |
| js-yaml | Moderate | @angular-eslint update |

---

## Go Dependencies

### Go Version

- **Current**: 1.20
- **Recommended**: 1.23+ (latest stable)

### Direct Dependencies

| Package | Version | Latest | Gap | Priority |
|---------|---------|--------|-----|----------|
| github.com/analogj/go-util | v0.0.0-20190301 | v0.0.0-20210417 | 2 years | Low |
| github.com/containrrr/shoutrrr | v0.8.0 | v0.8.0 | Current | - |
| github.com/fatih/color | v1.15.0 | v1.18.0 | 3 minor | Low |
| github.com/gin-gonic/gin | v1.6.3 | v1.11.0 | 4 minor | High |
| github.com/glebarez/sqlite | v1.4.5 | v1.11.0 | 6 minor | Medium |
| github.com/go-gormigrate/gormigrate/v2 | v2.0.0 | v2.1.5 | 1 minor | Low |
| github.com/golang/mock | v1.6.0 | v1.6.0 | Current | - |
| github.com/influxdata/influxdb-client-go/v2 | v2.9.0 | v2.14.0 | 5 minor | Medium |
| github.com/jaypipes/ghw | v0.6.1 | v0.21.2 | 15 minor | High |
| github.com/mitchellh/mapstructure | v1.5.0 | v1.5.0 | Current | - |
| github.com/prometheus/client_golang | v1.17.0 | v1.23.2 | 6 minor | Medium |
| github.com/samber/lo | v1.25.0 | v1.52.0 | 27 minor | Medium |
| github.com/sirupsen/logrus | v1.6.0 | v1.9.3 | 3 minor | Low |
| github.com/spf13/viper | v1.15.0 | v1.21.0 | 6 minor | Medium |
| github.com/stretchr/testify | v1.8.1 | v1.11.1 | 3 minor | Low |
| github.com/urfave/cli/v2 | v2.2.0 | v2.27.7 | 25 minor | High |
| golang.org/x/sync | v0.3.0 | v0.19.0 | 16 minor | Low |
| gorm.io/gorm | v1.23.5 | v1.31.1 | 8 minor | High |

### Deprecated Indirect Dependencies

| Package | Status | Replacement |
|---------|--------|-------------|
| github.com/golang/protobuf | Deprecated | google.golang.org/protobuf |
| github.com/deepmap/oapi-codegen | Deprecated | Consider alternatives |
| github.com/cncf/udpa/go | Deprecated | No longer maintained |

---

## NPM Dependencies

### Angular Framework

| Package | Version | Status |
|---------|---------|--------|
| @angular/core | ^21.0.5 | Current |
| @angular/common | ^21.0.5 | Current |
| @angular/compiler | ^21.0.5 | Current |
| @angular/forms | ^21.0.5 | Current |
| @angular/router | ^21.0.5 | Current |
| @angular/animations | ^21.0.5 | Current |
| @angular/platform-browser | ^21.0.5 | Current |
| @angular/platform-browser-dynamic | ^21.0.5 | Current |

### Angular Material (Version Mismatch)

| Package | Version | Expected | Status |
|---------|---------|----------|--------|
| @angular/material | ^16.2.14 | ^21.x | 5 major versions behind |
| @angular/cdk | ^16.2.14 | ^21.x | 5 major versions behind |
| @angular/material-moment-adapter | ^16.2.14 | ^21.x | 5 major versions behind |

### Production Dependencies

| Package | Version | Status | Notes |
|---------|---------|--------|-------|
| crypto-js | ^4.1.1 | Current | |
| highlight.js | ^11.6.0 | Current | |
| humanize-duration | ^3.27.3 | Current | |
| lodash | 4.17.21 | Current | Locked version |
| marked | ^17.0.1 | Current | |
| moment | ^2.29.4 | Deprecated | Maintainers recommend alternatives |
| ng-apexcharts | ^1.17.1 | Current | |
| perfect-scrollbar | ^1.5.5 | Current | |
| quill | ^1.3.7 | Vulnerable | XSS vulnerability, upgrade to v2.0.3 |
| rrule | ^2.7.1 | Current | |
| rxjs | ^7.5.7 | Current | |
| tslib | ^2.4.1 | Current | |
| web-animations-js | ^2.3.2 | Current | |
| zone.js | ^0.15.1 | Current | |

### Development Dependencies

| Package | Version | Status |
|---------|---------|--------|
| @angular/cli | ^21.0.3 | Current |
| @angular/build | ^21.0.3 | Current |
| @angular/compiler-cli | ^21.0.5 | Current |
| @angular/language-service | ^21.0.5 | Current |
| @angular-eslint/builder | ^21.1.0 | Current |
| @angular-eslint/eslint-plugin | ^21.1.0 | Current |
| @angular-eslint/eslint-plugin-template | ^21.1.0 | Current |
| @angular-eslint/template-parser | ^21.1.0 | Current |
| @angular-eslint/schematics | ^21.1.0 | Current |
| @typescript-eslint/eslint-plugin | ^5.62.0 | Current |
| @typescript-eslint/parser | ^5.62.0 | Current |
| apexcharts | ~3.35.0 | Outdated | |
| eslint | ^8.57.1 | Current | |
| eslint-config-prettier | ^8.10.2 | Current | |
| eslint-plugin-prettier | ^4.2.5 | Current | |
| jasmine-core | ^4.5.0 | Current | |
| jasmine-spec-reporter | ^7.0.0 | Current | |
| karma | ^6.4.1 | Current | |
| karma-chrome-launcher | ^3.1.1 | Current | |
| karma-coverage | ^2.2.0 | Current | |
| karma-jasmine | ^5.1.0 | Current | |
| karma-jasmine-html-reporter | ^2.0.0 | Current | |
| ngx-markdown | ^21.0.1 | Current | |
| prettier | ^2.8.8 | Current | |
| tailwindcss | ^3.2.3 | Outdated | v4.x available |
| ts-node | ^10.9.1 | Current | |
| typescript | ~5.9 | Current | |

---

## Update Strategy

### Immediate (This PR)

- [x] npm audit fix (non-breaking)
- [x] Update @angular-eslint to v21.x

### Phase 2: Low-Risk Go Updates

```bash
go get -u github.com/sirupsen/logrus@latest
go get -u github.com/fatih/color@latest
go get -u github.com/stretchr/testify@latest
go get -u golang.org/x/sync@latest
```

### Phase 3: Angular Material Alignment

```bash
ng update @angular/material @angular/cdk
```

### Phase 4: Medium-Risk Go Updates

```bash
go get -u github.com/samber/lo@latest
go get -u github.com/prometheus/client_golang@latest
go get -u github.com/influxdata/influxdb-client-go/v2@latest
go get -u github.com/spf13/viper@latest
```

### Phase 5: High-Risk Go Updates

```bash
go get -u github.com/gin-gonic/gin@latest
go get -u github.com/urfave/cli/v2@latest
go get -u github.com/jaypipes/ghw@latest
go get -u gorm.io/gorm@latest
```

---

## Related Issues

- GitHub #36: Dependency Health Check (this audit)
- GitHub #69: Quill v2.0 upgrade (XSS vulnerability fix)
- GitHub #70: moment.js migration to date-fns
- TBD: Angular Material v21 upgrade
- TBD: Go dependency updates (multiple phases)

---

## Maintenance

To check for vulnerabilities:

```bash
# NPM
cd webapp/frontend
npm audit

# Go (requires govulncheck)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

To check for outdated packages:

```bash
# NPM
cd webapp/frontend
npm outdated

# Go
go list -m -u all
```
