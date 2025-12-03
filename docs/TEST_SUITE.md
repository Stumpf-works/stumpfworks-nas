# Testing Guide - Stumpf.Works NAS

**Last Updated:** December 3, 2025
**Test Coverage Goal:** 80%+

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Test Infrastructure](#test-infrastructure)
3. [Running Tests](#running-tests)
4. [Writing Tests](#writing-tests)
5. [Test Coverage](#test-coverage)
6. [CI/CD Integration](#cicd-integration)
7. [Best Practices](#best-practices)

---

## Overview

This document describes the testing strategy and infrastructure for Stumpf.Works NAS. Our goal is to maintain **80%+ test coverage** across the codebase to ensure production readiness.

### Test Types

1. **Unit Tests** - Test individual functions and handlers
2. **Integration Tests** - Test API endpoints with real database
3. **E2E Tests** - Test complete user workflows (Playwright)
4. **Load Tests** - Performance and stress testing (k6)
5. **Benchmark Tests** - Performance profiling

---

## Test Infrastructure

### Backend (Go)

**Testing Framework:**
- `testing` - Standard Go testing package
- `testify` - Assertions and test suites
- `httptest` - HTTP handler testing

**Test Utilities:**
- `backend/internal/testutil/http.go` - HTTP test helpers
- `backend/internal/testutil/database.go` - Test database setup
- `backend/internal/testutil/fixtures.go` - Test fixtures and mock data

### Frontend (TypeScript/React)

**Testing Framework:**
- `vitest` - Fast unit test runner
- `@testing-library/react` - React component testing
- `playwright` - E2E testing (planned)

---

## Running Tests

### Quick Commands

```bash
# Run all tests with coverage
make test

# Run quick tests (no coverage)
make test-quick

# Backend tests only
make test-backend

# Frontend tests only
make test-frontend

# Generate coverage report
make test-coverage

# Run with race detector
make test-race
```

### Direct Script Usage

```bash
# Comprehensive test suite
./scripts/run-tests.sh

# Backend only
./scripts/run-tests.sh --backend

# Frontend only
./scripts/run-tests.sh --frontend

# Verbose output
./scripts/run-tests.sh --verbose

# With race detector
./scripts/run-tests.sh --race
```

### Coverage Reports

After running tests, coverage reports are available at:
- **Backend:** `coverage/backend-coverage.html`
- **Frontend:** `coverage/frontend/index.html`

---

## Writing Tests

### Backend Unit Tests

#### Basic Handler Test

```go
package handlers

import (
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMyHandler_Success(t *testing.T) {
    // Setup
    req := httptest.NewRequest("GET", "/api/test", nil)
    rr := httptest.NewRecorder()

    // Execute
    MyHandler(rr, req)

    // Assert
    assert.Equal(t, 200, rr.Code)

    var response map[string]interface{}
    err := json.Unmarshal(rr.Body.Bytes(), &response)
    require.NoError(t, err)
    assert.Contains(t, response, "data")
}
```

#### Using Test Utilities

```go
func TestWithHTTPHelper(t *testing.T) {
    h := testutil.NewHTTPTest(t)

    // Create request
    req := h.MakeRequest("POST", "/api/users", map[string]string{
        "username": "test",
    })

    // Execute
    rr := h.ExecuteRequest(CreateUser, req)

    // Assert
    h.AssertStatusCode(rr, 201)
    h.AssertJSONResponse(rr, &response)
}
```

#### Database Tests

```go
func TestWithDatabase(t *testing.T) {
    // Setup test database
    db := testutil.SetupTestDBWithModels(t, &models.User{})
    defer testutil.CleanupTestDB(t, db)

    // Create test data
    user := testutil.CreateTestUser("alice", "user")
    db.Create(user)

    // Test your function
    result, err := FindUser(db, "alice")
    assert.NoError(t, err)
    assert.Equal(t, "alice", result.Username)
}
```

### Frontend Unit Tests

#### Component Test

```typescript
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import MyComponent from './MyComponent';

describe('MyComponent', () => {
  it('renders correctly', () => {
    render(<MyComponent title="Test" />);
    expect(screen.getByText('Test')).toBeInTheDocument();
  });
});
```

---

## Test Coverage

### Current Status (December 2025)

#### Backend Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| `internal/users` | 85% | ✅ Complete |
| `internal/files` | 78% | ✅ Complete |
| `internal/api/handlers` | **25%** | 🚧 In Progress |
| `internal/storage` | 0% | ❌ Todo |
| `internal/docker` | 0% | ❌ Todo |
| `internal/cloudbackup` | 0% | ❌ Todo |

**Overall Backend Coverage:** ~15% → **Goal: 80%+**

#### Handlers Coverage Status

✅ **Tested (2/38):**
- `auth.go` - Authentication endpoints
- `health.go` - Health check endpoints

🚧 **In Progress (5/38):**
- `storage.go`
- `docker.go`
- `ups.go`
- `cloudbackup.go`
- `alertrules.go`

❌ **Todo (31/38):** All remaining handlers

#### Frontend Coverage

**Overall Frontend Coverage:** ~0% → **Goal: 70%+**

---

## CI/CD Integration

### Local Test Scripts

Since we're not using GitHub Actions, tests are run via local scripts:

**Pre-commit Hook (Recommended):**
```bash
#!/bin/bash
# .git/hooks/pre-commit

./scripts/test-quick.sh
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi
```

**Pre-push Hook:**
```bash
#!/bin/bash
# .git/hooks/pre-push

./scripts/run-tests.sh
if [ $? -ne 0 ]; then
    echo "Tests failed. Push aborted."
    exit 1
fi
```

---

## Best Practices

### 1. Test Naming

- **Pattern:** `Test<FunctionName>_<Scenario>`
- **Examples:**
  - `TestLogin_Success`
  - `TestLogin_InvalidCredentials`
  - `TestLogin_EmptyPassword`

### 2. Test Structure

Follow the **Arrange-Act-Assert** pattern:

```go
func TestExample(t *testing.T) {
    // Arrange - Setup test data and dependencies
    user := createTestUser()

    // Act - Execute the function under test
    result, err := DoSomething(user)

    // Assert - Verify the results
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### 3. Table-Driven Tests

Use table-driven tests for multiple scenarios:

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectedErr bool
    }{
        {"valid input", "test", false},
        {"empty input", "", true},
        {"too long", strings.Repeat("a", 1000), true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if tt.expectedErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 4. Test Independence

- Each test should be independent
- Use `t.Parallel()` for parallel execution when safe
- Clean up resources with `defer`

### 5. Mock External Dependencies

- Mock database calls
- Mock HTTP clients
- Mock file system operations
- Use interfaces for dependency injection

### 6. Coverage Goals

- **Critical paths:** 90%+
- **Business logic:** 80%+
- **Handlers:** 80%+
- **Utilities:** 70%+

### 7. Don't Test

- External libraries (already tested)
- Generated code
- Trivial getters/setters
- Simple constructors

---

## Continuous Improvement

### Weekly Goals

1. **Week 1:** Auth, Health, Users packages → 80%
2. **Week 2:** Storage, Docker, Networking → 80%
3. **Week 3:** Cloud Backup, UPS, Monitoring → 80%
4. **Week 4:** Integration tests, E2E tests

### Monitoring

Track coverage trends:
```bash
# Generate coverage report
make test-coverage

# View HTML report
open coverage/backend-coverage.html
```

---

## Troubleshooting

### Common Issues

#### Build Errors

**Problem:** `undefined: syscall.Statfs_t`
**Solution:** Platform-specific code needs build tags. Add:
```go
// +build linux

package files
```

#### Slow Tests

**Problem:** Tests taking too long
**Solution:** Use `testing.Short()` for quick tests:
```go
if testing.Short() {
    t.Skip("Skipping slow test in short mode")
}
```

Run with: `go test -short ./...`

#### Flaky Tests

**Problem:** Tests pass/fail randomly
**Solution:**
- Avoid time.Sleep() - use channels/waitgroups
- Don't rely on execution order
- Clean up goroutines properly

---

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Vitest Documentation](https://vitest.dev/)
- [Testing Best Practices](https://github.com/golang/go/wiki/TestComments)

---

**Questions?** Open an issue on GitHub or check the [CONTRIBUTING.md](CONTRIBUTING.md)
