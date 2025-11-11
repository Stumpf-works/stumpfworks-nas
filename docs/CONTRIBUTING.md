# Contributing to Stumpf.Works NAS Solution

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the project.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Documentation](#documentation)

---

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors, regardless of:
- Experience level
- Gender identity and expression
- Sexual orientation
- Disability
- Personal appearance
- Body size
- Race
- Ethnicity
- Age
- Religion
- Nationality

### Expected Behavior

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Trolling, insulting/derogatory comments, and personal attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate

---

## How Can I Contribute?

### 1. Reporting Bugs

Before submitting a bug report:
- Check the [issue tracker](https://github.com/stumpfworks/nas/issues) for existing reports
- Try to reproduce the issue with the latest version

When submitting a bug report, include:
- **Clear title** - Descriptive, concise summary
- **Steps to reproduce** - Detailed steps to trigger the bug
- **Expected behavior** - What you expected to happen
- **Actual behavior** - What actually happened
- **Environment** - OS, version, hardware specs
- **Logs** - Relevant log output (sanitize sensitive data)
- **Screenshots** - If applicable

**Template:**

```markdown
### Bug Description
Brief description of the issue.

### Steps to Reproduce
1. Go to '...'
2. Click on '...'
3. See error

### Expected Behavior
What should happen.

### Actual Behavior
What actually happens.

### Environment
- OS: Debian 12 (Bookworm)
- Version: 1.0.0
- Browser: Chrome 120

### Logs
```
[paste logs here]
```

### Screenshots
[if applicable]
```

---

### 2. Suggesting Features

Before submitting a feature request:
- Check if the feature already exists
- Search existing feature requests

When submitting a feature request, include:
- **Use case** - Why is this feature needed?
- **Proposed solution** - How should it work?
- **Alternatives** - Other solutions you've considered
- **Mockups** - Wireframes or screenshots (if UI-related)

---

### 3. Contributing Code

Areas where contributions are welcome:
- **Bug fixes** - Fix reported issues
- **New features** - Implement planned features (check ROADMAP.md)
- **Performance** - Optimize slow operations
- **Tests** - Add missing test coverage
- **Documentation** - Improve or expand docs
- **Refactoring** - Improve code quality

Before starting:
- **Check existing issues** - See if someone is already working on it
- **Open an issue** - Discuss your approach before writing code
- **Get feedback** - Ensure your approach aligns with project goals

---

### 4. Improving Documentation

Documentation contributions are highly valued:
- Fix typos or grammatical errors
- Clarify confusing sections
- Add examples
- Translate documentation
- Write tutorials or guides

---

## Development Setup

### Prerequisites

- **Go 1.21+** (for backend)
- **Node.js 18+** (for frontend)
- **Docker** (for containerized development)
- **Git**

### Clone Repository

```bash
git clone https://github.com/stumpfworks/nas.git
cd nas
```

### Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Run backend
go run cmd/stumpfworks-server/main.go

# Backend runs on http://localhost:8080
```

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev

# Frontend runs on http://localhost:3000
```

### Docker Setup (Alternative)

```bash
# Run entire stack with Docker Compose
docker-compose up -d

# Backend: http://localhost:8080
# Frontend: http://localhost:3000
```

---

## Coding Standards

### Go (Backend)

#### Style Guide

Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

**Key Points:**
- Use `gofmt` for formatting (automatic)
- Use `golangci-lint` for linting
- Exported names must have doc comments
- Error strings should not be capitalized
- Use meaningful variable names (avoid single letters except in small scopes)

**Example:**

```go
// Good
func GetStorageInfo(volumeID string) (*StorageInfo, error) {
    if volumeID == "" {
        return nil, errors.New("volume ID is required")
    }

    info, err := fetchStorageInfo(volumeID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch storage info: %w", err)
    }

    return info, nil
}

// Bad
func get_storage_info(v string) (*StorageInfo, error) {
    Info, Err := FetchStorageInfo(v)
    if Err != nil {
        return nil, Err
    }
    return Info, nil
}
```

#### Project Structure

```
backend/
├── cmd/                  # Main applications
│   └── stumpfworks-server/
│       └── main.go
├── internal/             # Private application code
│   ├── api/
│   ├── storage/
│   └── ...
├── pkg/                  # Public libraries
│   ├── logger/
│   └── ...
└── go.mod
```

#### Error Handling

Always handle errors explicitly:

```go
// Good
data, err := readFile(path)
if err != nil {
    return fmt.Errorf("read file: %w", err)
}

// Bad
data, _ := readFile(path)  // Never ignore errors!
```

#### Testing

Write table-driven tests:

```go
func TestCalculateSize(t *testing.T) {
    tests := []struct {
        name     string
        input    int64
        expected string
    }{
        {"bytes", 500, "500 B"},
        {"kilobytes", 1024, "1.0 KB"},
        {"megabytes", 1048576, "1.0 MB"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateSize(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

---

### TypeScript/React (Frontend)

#### Style Guide

Follow [Airbnb React/JSX Style Guide](https://github.com/airbnb/javascript/tree/master/react).

**Key Points:**
- Use TypeScript for all new code
- Use functional components and hooks (no class components)
- Use ESLint and Prettier (automatic)
- Prefer named exports over default exports
- Use meaningful component and variable names

**Example:**

```tsx
// Good
interface DashboardProps {
  userId: string;
  onRefresh: () => void;
}

export function Dashboard({ userId, onRefresh }: DashboardProps) {
  const [data, setData] = useState<DashboardData | null>(null);

  useEffect(() => {
    fetchDashboardData(userId).then(setData);
  }, [userId]);

  if (!data) {
    return <LoadingSpinner />;
  }

  return (
    <div className="dashboard">
      <h1>Dashboard</h1>
      <button onClick={onRefresh}>Refresh</button>
      <DataView data={data} />
    </div>
  );
}

// Bad
export default function dashboard(props: any) {
  const [d, setD] = useState(null);
  // ... (no type safety, unclear naming)
}
```

#### Component Structure

```
frontend/src/
├── layout/              # Core layout components
├── apps/                # Applications
│   └── Dashboard/
│       ├── Dashboard.tsx
│       ├── components/  # App-specific components
│       ├── hooks/       # App-specific hooks
│       └── utils/       # App-specific utilities
├── components/          # Shared components
├── hooks/               # Shared hooks
├── store/               # State management
├── api/                 # API client
└── utils/               # Shared utilities
```

#### State Management

Use Zustand for global state:

```tsx
// store/systemStore.ts
import create from 'zustand';

interface SystemStore {
  cpuUsage: number;
  memoryUsage: number;
  updateMetrics: (cpu: number, memory: number) => void;
}

export const useSystemStore = create<SystemStore>((set) => ({
  cpuUsage: 0,
  memoryUsage: 0,
  updateMetrics: (cpu, memory) => set({ cpuUsage: cpu, memoryUsage: memory }),
}));

// Usage in component
import { useSystemStore } from '@/store/systemStore';

export function SystemMonitor() {
  const { cpuUsage, memoryUsage } = useSystemStore();

  return (
    <div>
      <p>CPU: {cpuUsage}%</p>
      <p>Memory: {memoryUsage}%</p>
    </div>
  );
}
```

---

## Commit Guidelines

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring (no feature change)
- `perf`: Performance improvements
- `test`: Add or update tests
- `chore`: Maintenance tasks (dependencies, build config)

**Examples:**

```
feat(storage): add ZFS support

Implemented ZFS pool creation and management. Users can now
create and manage ZFS pools through the Storage Manager app.

Closes #123
```

```
fix(api): handle empty volume ID in GetStorageInfo

Previously, empty volume IDs caused a panic. Now returns a
proper error message.

Fixes #456
```

```
docs(plugin): clarify permission system

Added examples and improved explanation of the plugin
permission system in PLUGIN_DEV.md.
```

### Commit Frequency

- Commit often (atomic commits)
- Each commit should be a logical unit
- Don't mix unrelated changes in one commit

---

## Pull Request Process

### 1. Fork and Branch

```bash
# Fork the repository on GitHub

# Clone your fork
git clone https://github.com/YOUR_USERNAME/nas.git
cd nas

# Add upstream remote
git remote add upstream https://github.com/stumpfworks/nas.git

# Create feature branch
git checkout -b feature/my-feature
```

### 2. Make Changes

- Write code following style guides
- Write tests for new features
- Update documentation if needed
- Ensure all tests pass

### 3. Commit

```bash
git add .
git commit -m "feat(module): add new feature"
```

### 4. Push

```bash
git push origin feature/my-feature
```

### 5. Open Pull Request

- Go to GitHub and open a pull request
- Fill out the PR template
- Link related issues (e.g., "Closes #123")
- Request review from maintainers

### PR Template

```markdown
## Description
Brief description of changes.

## Related Issue
Closes #123

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manually tested

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests pass locally
```

### 6. Code Review

- Respond to feedback promptly
- Make requested changes
- Push updates to the same branch
- Re-request review when ready

### 7. Merge

Once approved, a maintainer will merge your PR.

---

## Testing

### Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/storage
```

### Frontend Tests

```bash
cd frontend

# Run unit tests (Vitest)
npm run test

# Run with coverage
npm run test:coverage

# Run E2E tests (Playwright)
npm run test:e2e
```

### Integration Tests

```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
make test-integration
```

---

## Documentation

### Where to Document

- **Code comments** - Complex logic, public APIs
- **README.md** - Project overview, quick start
- **docs/** - Comprehensive guides (architecture, plugins, etc.)
- **API docs** - OpenAPI/Swagger for REST API

### Documentation Standards

- Use clear, concise language
- Include examples
- Keep documentation up-to-date with code changes
- Use proper Markdown formatting

---

## Review Process

### What We Look For

✅ **Code Quality**
- Follows style guidelines
- Well-structured and readable
- No unnecessary complexity

✅ **Functionality**
- Works as intended
- Handles edge cases
- No regressions

✅ **Tests**
- Adequate test coverage
- Tests are meaningful (not just for coverage)

✅ **Documentation**
- Code is documented
- User-facing changes have updated docs

✅ **Performance**
- No obvious performance issues
- Efficient algorithms

✅ **Security**
- No security vulnerabilities
- Input validation
- Proper error handling

---

## Getting Help

- **GitHub Issues** - Ask questions, report bugs
- **Discussions** - General questions, ideas
- **Discord** - Real-time chat (link in README)
- **Email** - maintainers@stumpf.works

---

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Credited in release notes
- Acknowledged in the community

---

**Thank you for contributing to Stumpf.Works NAS Solution!** 🎉

---

**Document Version:** 1.0.0
**Last Updated:** 2025-11-11
