# Technology Stack

> Detailed analysis of technology choices and trade-offs for Stumpf.Works NAS Solution

---

## Table of Contents

- [Overview](#overview)
- [Backend Technology](#backend-technology)
- [Frontend Technology](#frontend-technology)
- [Database](#database)
- [Infrastructure](#infrastructure)
- [Development Tools](#development-tools)
- [Trade-off Analysis](#trade-off-analysis)

---

## Overview

Technology decisions for Stumpf.Works NAS prioritize:
1. **Performance** - System-level operations require efficiency
2. **Developer Experience** - Modern tooling and clear patterns
3. **Maintainability** - Code should be readable and testable
4. **Community** - Active ecosystems and available libraries
5. **Production Readiness** - Battle-tested in production environments

---

## Backend Technology

### Primary Choice: **Go (Golang)**

**Version:** Go 1.21+

**Rationale:**

✅ **Pros:**
- **Performance**: Compiled binary, fast execution, low memory footprint
- **Concurrency**: Goroutines and channels make concurrent operations trivial
- **System Integration**: Excellent for system-level programming (storage, network, processes)
- **Single Binary**: Easy deployment - no runtime dependencies
- **Standard Library**: Comprehensive stdlib (HTTP server, JSON, crypto, etc.)
- **Cross-Platform**: Compile for Linux, macOS, Windows from any OS
- **Ecosystem**: Rich libraries for storage (LVM, ZFS), Docker API, libvirt
- **Error Handling**: Explicit error handling reduces runtime surprises
- **Testing**: Built-in testing framework (`go test`)
- **Tooling**: `go fmt`, `go vet`, `go mod` - excellent developer experience

❌ **Cons:**
- **Verbosity**: More boilerplate compared to Python
- **Generic Support**: Generics added in 1.18, still maturing
- **Learning Curve**: Slightly steeper for developers from dynamic languages

**Use Cases in Stumpf.Works:**
- API server (REST + WebSocket)
- Storage management (LVM, mdadm, ZFS)
- Network configuration
- Plugin system
- System monitoring
- Container/VM management

**Key Libraries:**
- **Web Framework**: [Chi](https://github.com/go-chi/chi) or [Gin](https://github.com/gin-gonic/gin)
- **WebSocket**: [gorilla/websocket](https://github.com/gorilla/websocket)
- **Database**: [GORM](https://gorm.io/) or [sqlx](https://github.com/jmoiron/sqlx)
- **Configuration**: [viper](https://github.com/spf13/viper)
- **Logging**: [zap](https://github.com/uber-go/zap) or [zerolog](https://github.com/rs/zerolog)
- **Docker**: [docker/docker](https://github.com/moby/moby) (official Go client)
- **Testing**: [testify](https://github.com/stretchr/testify)

---

### Alternative: Python (FastAPI)

**Considered but not chosen for core backend.**

✅ **Pros:**
- Rapid development
- Huge ecosystem (numpy, pandas for data processing)
- Excellent for scripting and prototyping
- Easy to learn

❌ **Cons:**
- **Performance**: Slower than Go for system operations
- **Concurrency**: GIL limits true parallelism
- **Deployment**: Requires Python runtime and dependencies
- **Type Safety**: Optional typing (mypy) not enforced at runtime

**Verdict:** Python may be used for **plugin development** or **utility scripts**, but Go is preferred for the core backend.

---

### Alternative: Rust

**Considered but not chosen for initial version.**

✅ **Pros:**
- Maximum performance
- Memory safety without garbage collection
- Growing ecosystem

❌ **Cons:**
- **Steep Learning Curve**: Ownership model is challenging
- **Development Speed**: Slower to prototype
- **Ecosystem**: Smaller library selection for system management
- **Compile Times**: Slower than Go

**Verdict:** Rust is an excellent choice for performance-critical plugins or future rewrites, but Go offers a better balance for initial development.

---

## Frontend Technology

### Primary Stack: **React + TypeScript + Vite + TailwindCSS + Framer Motion**

---

### React 18

**Why React?**

✅ **Pros:**
- **Ecosystem**: Largest component library ecosystem
- **Maturity**: Battle-tested, production-ready
- **Community**: Huge community, extensive documentation
- **Performance**: Virtual DOM, concurrent features
- **Tooling**: Excellent dev tools (React DevTools)
- **Hiring**: Easier to find React developers

**Key Features Used:**
- Hooks (`useState`, `useEffect`, `useContext`, custom hooks)
- Context API (lightweight state management)
- Suspense and lazy loading (code splitting)
- Concurrent rendering (React 18 features)

---

### TypeScript 5

**Why TypeScript?**

✅ **Pros:**
- **Type Safety**: Catch errors at compile time
- **Intellisense**: Better IDE autocomplete
- **Refactoring**: Safer code changes
- **Documentation**: Types serve as inline documentation
- **Ecosystem**: Most libraries have TypeScript definitions

**Configuration:**
- Strict mode enabled
- Path aliases (`@/components`, `@/utils`)
- ESLint + Prettier integration

---

### Vite

**Why Vite over Create React App or Webpack?**

✅ **Pros:**
- **Speed**: Lightning-fast HMR (Hot Module Replacement)
- **Modern**: Built for ES modules
- **Plugins**: Rich plugin ecosystem
- **Build Performance**: Rollup-based production builds
- **Developer Experience**: Instant server start

**Configuration:**
- React plugin (`@vitejs/plugin-react`)
- Path aliases
- Environment variables
- Proxy for API during development

---

### TailwindCSS 3

**Why TailwindCSS?**

✅ **Pros:**
- **Utility-First**: Rapid UI development
- **Consistency**: Design tokens in config
- **Performance**: PurgeCSS removes unused styles
- **Responsive**: Mobile-first breakpoints
- **Customization**: Fully customizable theme
- **macOS Aesthetic**: Easy to implement glassmorphism, blur effects

**Configuration:**
- Custom color palette (macOS-inspired)
- Custom spacing and border radius
- Dark mode support
- Custom animations

**Example macOS Theme:**
```js
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        'macos-gray': '#1e1e1e',
        'macos-blue': '#007aff',
        'macos-blur': 'rgba(255, 255, 255, 0.8)',
      },
      backdropBlur: {
        'macos': '40px',
      },
      boxShadow: {
        'macos': '0 10px 40px rgba(0, 0, 0, 0.15)',
      },
    },
  },
}
```

---

### Framer Motion

**Why Framer Motion?**

✅ **Pros:**
- **Declarative Animations**: Easy to write smooth animations
- **Layout Animations**: Automatic layout transitions
- **Gestures**: Drag, hover, tap interactions
- **Variants**: Coordinated animations across components
- **Performance**: Optimized for 60fps

**Use Cases:**
- Dock animations (bounce, magnification)
- Window transitions (open, close, minimize, maximize)
- Launchpad grid animations
- Page transitions
- Micro-interactions (button hovers, etc.)

**Example:**
```tsx
import { motion } from 'framer-motion';

<motion.div
  initial={{ opacity: 0, y: 20 }}
  animate={{ opacity: 1, y: 0 }}
  exit={{ opacity: 0, y: -20 }}
  transition={{ duration: 0.3 }}
>
  {content}
</motion.div>
```

---

### State Management: Zustand

**Why Zustand over Redux/MobX?**

✅ **Pros:**
- **Simplicity**: Minimal boilerplate
- **Performance**: No context provider overhead
- **TypeScript**: Excellent TypeScript support
- **Devtools**: Redux DevTools integration
- **Size**: ~1KB (vs Redux ~10KB)

**Example Store:**
```ts
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
```

---

### HTTP Client: Axios

**Why Axios?**

✅ **Pros:**
- **Interceptors**: Easy to add auth headers
- **Error Handling**: Consistent error format
- **Timeouts**: Built-in timeout support
- **Automatic JSON**: Parses JSON responses automatically
- **Browser Compatibility**: Works everywhere

**Example:**
```ts
// api/client.ts
import axios from 'axios';

const client = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default client;
```

---

## Database

### Primary: **SQLite**

**Why SQLite for a NAS?**

✅ **Pros:**
- **Embedded**: No separate database server
- **Zero Configuration**: No setup required
- **ACID Compliant**: Reliable transactions
- **Fast**: Excellent read performance
- **File-Based**: Easy backups (copy file)
- **Cross-Platform**: Works everywhere
- **Low Memory**: Minimal resource usage

❌ **Cons:**
- **Concurrency**: Limited write concurrency
- **Scalability**: Not suitable for high-traffic web apps

**Verdict:** Perfect for a single-node NAS. Most operations are read-heavy (viewing storage, users, settings). Writes are infrequent (creating shares, adding users).

**Schema:**
- Users and permissions
- Storage volumes and snapshots
- Network configuration
- Plugin metadata
- System settings
- Logs and audit trail

**ORM:** GORM (Go) or sqlx for simpler queries

---

### Alternative: PostgreSQL

**When to switch to PostgreSQL:**
- Multi-node clustering (future)
- High write concurrency requirements
- Advanced querying (JSON, full-text search)

**Migration Path:** GORM supports multiple databases, so switching later is straightforward.

---

## Infrastructure

### Operating System: **Debian Bookworm (Stable)**

**Why Debian?**

✅ **Pros:**
- **Stability**: Rock-solid, conservative updates
- **Package Availability**: Huge repository (apt)
- **Long-Term Support**: Security updates for years
- **Community**: Large, helpful community
- **Predictability**: No surprises in stable releases
- **Universal**: Works on x86, ARM, etc.

**Why not Ubuntu?**
- Debian is upstream of Ubuntu
- No Snap packages forced on users
- More conservative (fewer breaking changes)

**Why not Arch/Fedora?**
- Too bleeding-edge for a NAS (stability > latest features)

---

### Web Server: **nginx**

**Why nginx?**

✅ **Pros:**
- **Performance**: Handles thousands of concurrent connections
- **Reverse Proxy**: Perfect for backend API
- **Static Files**: Serves frontend efficiently
- **HTTPS**: Let's Encrypt integration
- **Configuration**: Simple, declarative config

**Example Config:**
```nginx
server {
  listen 443 ssl http2;
  server_name nas.local;

  ssl_certificate /etc/letsencrypt/live/nas.local/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/nas.local/privkey.pem;

  location /api/ {
    proxy_pass http://localhost:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
  }

  location /ws {
    proxy_pass http://localhost:8080;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
  }

  location / {
    root /var/www/stumpfworks/frontend;
    try_files $uri /index.html;
  }
}
```

---

### Containerization: **Docker + Podman**

**Why support both?**

- **Docker**: Most popular, huge ecosystem
- **Podman**: Rootless, daemonless, more secure

**Use Cases:**
- Plugin isolation
- Development environment
- Third-party services (Nextcloud, Plex, etc.)

---

### Virtualization: **libvirt + KVM**

**Why libvirt?**

✅ **Pros:**
- **Standard**: Industry standard for Linux VMs
- **Performance**: Near-native performance
- **Management**: CLI (`virsh`) and API
- **Tooling**: Compatible with virt-manager

---

## Development Tools

### Version Control: **Git + GitHub**

**Workflow:**
- Main branch: stable releases
- Develop branch: integration
- Feature branches: `feature/name`
- Release branches: `release/v1.0.0`

---

### CI/CD: **GitHub Actions**

**Pipeline:**
1. Lint (golangci-lint, ESLint)
2. Test (Go tests, Jest/Vitest)
3. Build (compile backend, build frontend)
4. Security scan (Trivy, Dependabot)
5. Release (GitHub Releases, Docker Hub)

---

### Package Manager:

**Backend:** `go mod`
**Frontend:** `npm` or `pnpm` (faster)

---

### Code Quality:

**Backend:**
- `golangci-lint` (linting)
- `go fmt` (formatting)
- `go vet` (static analysis)

**Frontend:**
- ESLint (linting)
- Prettier (formatting)
- Husky (pre-commit hooks)

---

### Documentation:

- **Markdown**: All docs in `/docs`
- **OpenAPI/Swagger**: API documentation
- **Storybook**: Component documentation (optional)

---

## Trade-off Analysis

### Backend: Go vs Python vs Rust

| Factor | Go | Python | Rust |
|--------|----|---------| -----|
| **Performance** | ⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Development Speed** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **Concurrency** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Ecosystem (System)** | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Learning Curve** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐ |
| **Deployment** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

**Verdict:** Go wins for backend.

---

### Frontend: React vs Vue vs Svelte

| Factor | React | Vue | Svelte |
|--------|-------|-----|--------|
| **Ecosystem** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Performance** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Learning Curve** | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Community** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Animation Support** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **TypeScript** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

**Verdict:** React wins for frontend (ecosystem + animation libraries).

---

### Database: SQLite vs PostgreSQL

| Factor | SQLite | PostgreSQL |
|--------|--------|------------|
| **Setup Complexity** | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **Performance (Reads)** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Performance (Writes)** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Concurrency** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Resource Usage** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Backup** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

**Verdict:** SQLite for initial release. Migrate to PostgreSQL if clustering is needed.

---

## Summary

### Chosen Stack

**Backend:**
- Go 1.21+
- Chi/Gin router
- GORM + SQLite
- Zap logging

**Frontend:**
- React 18 + TypeScript 5
- Vite
- TailwindCSS 3
- Framer Motion
- Zustand
- Axios

**Infrastructure:**
- Debian Bookworm
- nginx
- Docker + Podman
- libvirt + KVM

**Development:**
- Git + GitHub
- GitHub Actions
- golangci-lint + ESLint
- Docker Compose (dev environment)

---

**Document Version:** 1.0.0
**Last Updated:** 2025-11-11
**Next Review:** After Phase 1 completion
