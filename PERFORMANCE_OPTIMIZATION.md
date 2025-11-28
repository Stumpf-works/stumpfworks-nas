# Performance Optimization Guide

**Ziel:** Maximale Performance für Stumpf.Works NAS
**Status:** Performance-Audit & Optimierungs-Roadmap

---

## Executive Summary

Stumpf.Works NAS läuft bereits gut, aber es gibt Potenzial für **erhebliche Performance-Verbesserungen**:

- 🎯 **Backend:** API Response Time < 100ms (aktuell: ?ms)
- 🎯 **Frontend:** Initial Load < 2s (aktuell: ?s)
- 🎯 **Database:** Query Time < 50ms (aktuell: ?ms)
- 🎯 **WebSocket:** Real-time Updates < 100ms Latenz

**Geschätzte Verbesserung:** 30-50% schneller

---

## 1. Backend Performance (Go)

### 1.1 API Response Time Optimization

#### Problem: Langsame API Endpoints

**Audit durchführen:**
```bash
# Add timing middleware
# backend/internal/api/middleware/timing.go
```

```go
package middleware

import (
    "net/http"
    "time"
    "github.com/rs/zerolog/log"
)

func TimingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        next.ServeHTTP(w, r)

        duration := time.Since(start)

        // Log slow requests (> 500ms)
        if duration > 500*time.Millisecond {
            log.Warn().
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Dur("duration", duration).
                Msg("Slow API request detected")
        }

        // Add header for debugging
        w.Header().Set("X-Response-Time", duration.String())
    })
}
```

**Optimierungen:**

1. **Reduce Shell Executor Calls**
```go
// BAD: Multiple shell calls
func GetDiskInfo() {
    shell.Execute("lsblk", "-J")           // Call 1
    shell.Execute("df", "-h")              // Call 2
    shell.Execute("smartctl", "-a", "/dev/sda") // Call 3
}

// GOOD: Batch commands or cache results
func GetDiskInfo() {
    // Use single command with multiple outputs
    output := shell.Execute("lsblk -J && df -h && smartctl -a /dev/sda")
    // Parse all at once
}
```

2. **Add Caching for Expensive Operations**
```go
// backend/pkg/cache/cache.go
package cache

import (
    "sync"
    "time"
)

type Cache struct {
    data map[string]cacheEntry
    mu   sync.RWMutex
    ttl  time.Duration
}

type cacheEntry struct {
    value      interface{}
    expiration time.Time
}

func NewCache(ttl time.Duration) *Cache {
    return &Cache{
        data: make(map[string]cacheEntry),
        ttl:  ttl,
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    entry, exists := c.data[key]
    if !exists || time.Now().After(entry.expiration) {
        return nil, false
    }

    return entry.value, true
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[key] = cacheEntry{
        value:      value,
        expiration: time.Now().Add(c.ttl),
    }
}
```

**Usage:**
```go
var diskCache = cache.NewCache(30 * time.Second)

func (h *DiskHandler) ListDisks(w http.ResponseWriter, r *http.Request) {
    // Try cache first
    if cached, ok := diskCache.Get("disks"); ok {
        utils.RespondSuccess(w, cached)
        return
    }

    // Expensive operation
    disks, err := disk.ListDisks()
    if err != nil {
        utils.RespondError(w, err)
        return
    }

    // Cache for 30 seconds
    diskCache.Set("disks", disks)

    utils.RespondSuccess(w, disks)
}
```

3. **Use Goroutines for Parallel Operations**
```go
// BAD: Sequential
func GetSystemInfo() (*SystemInfo, error) {
    cpu := getCPUInfo()        // 100ms
    memory := getMemoryInfo()  // 50ms
    disk := getDiskInfo()      // 200ms
    // Total: 350ms
}

// GOOD: Parallel
func GetSystemInfo() (*SystemInfo, error) {
    var cpu CPUInfo
    var memory MemoryInfo
    var disk DiskInfo
    var wg sync.WaitGroup

    wg.Add(3)

    go func() {
        defer wg.Done()
        cpu = getCPUInfo()
    }()

    go func() {
        defer wg.Done()
        memory = getMemoryInfo()
    }()

    go func() {
        defer wg.Done()
        disk = getDiskInfo()
    }()

    wg.Wait() // Total: 200ms (longest task)

    return &SystemInfo{CPU: cpu, Memory: memory, Disk: disk}, nil
}
```

4. **Database Query Optimization**
```go
// BAD: N+1 Query Problem
func GetUsersWithGroups() []User {
    users := db.Find(&User{})
    for _, user := range users {
        user.Groups = db.Find(&Group{UserID: user.ID}) // N queries!
    }
}

// GOOD: Eager Loading
func GetUsersWithGroups() []User {
    var users []User
    db.Preload("Groups").Find(&users) // Single query with JOIN
    return users
}
```

5. **Add Database Indexes**
```go
// backend/internal/db/migrations/add_indexes.go
func AddPerformanceIndexes(db *gorm.DB) error {
    // Index frequently queried fields
    db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)")
    db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp)")
    db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id)")
    db.Exec("CREATE INDEX IF NOT EXISTS idx_shares_name ON shares(name)")

    return nil
}
```

6. **Connection Pool Tuning**
```go
// backend/internal/db/database.go
func InitDB() (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    sqlDB, _ := db.DB()

    // Performance tuning
    sqlDB.SetMaxOpenConns(25)          // Limit concurrent connections
    sqlDB.SetMaxIdleConns(5)           // Keep idle connections
    sqlDB.SetConnMaxLifetime(5 * time.Minute) // Recycle connections

    return db, nil
}
```

### 1.2 File Operations Optimization

**Problem:** File upload/download zu langsam

**Optimierungen:**

1. **Chunk Size Optimization**
```go
// backend/internal/api/handlers/files.go
const (
    OptimalChunkSize = 8 * 1024 * 1024 // 8 MB chunks (war vermutlich kleiner)
)

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    // Use optimal buffer size
    buf := make([]byte, OptimalChunkSize)

    for {
        n, err := r.Body.Read(buf)
        if n > 0 {
            file.Write(buf[:n])
        }
        if err != nil {
            break
        }
    }
}
```

2. **Streaming für große Dateien**
```go
func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
    file, _ := os.Open(filepath)
    defer file.Close()

    // Set headers für streaming
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

    // Stream direkt (kein Buffer im RAM)
    io.Copy(w, file)
}
```

3. **Parallel File Hashing**
```go
func CalculateChecksums(files []string) map[string]string {
    results := make(map[string]string)
    var mu sync.Mutex
    var wg sync.WaitGroup

    // Limit concurrency to CPU count
    semaphore := make(chan struct{}, runtime.NumCPU())

    for _, file := range files {
        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release

            hash := calculateHash(f)

            mu.Lock()
            results[f] = hash
            mu.Unlock()
        }(file)
    }

    wg.Wait()
    return results
}
```

### 1.3 Docker API Optimization

**Problem:** Docker-Operationen blockieren

**Optimierung:**
```go
// Use context with timeout
func (d *DockerService) ListContainers(ctx context.Context) ([]Container, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
    if err != nil {
        return nil, err
    }

    return containers, nil
}
```

---

## 2. Frontend Performance (React)

### 2.1 Initial Load Time Optimization

#### Current Issues:
- Großes Bundle (alle Apps laden sofort)
- Keine Code-Splitting
- Unnötige Re-Renders

**Optimierungen:**

1. **Code Splitting (Lazy Loading)**
```typescript
// frontend/src/apps/index.tsx

// BAD: All apps loaded upfront
import { Dashboard } from './Dashboard/Dashboard';
import { FileManager } from './FileManager/FileManager';
// ... all apps

// GOOD: Lazy load apps
const Dashboard = lazy(() => import('./Dashboard/Dashboard'));
const FileManager = lazy(() => import('./FileManager/FileManager'));
const DockerManager = lazy(() => import('./DockerManager/DockerManager'));
// ... all apps

export const registeredApps: App[] = [
  {
    id: 'dashboard',
    name: 'Dashboard',
    icon: '📊',
    component: Dashboard, // Loaded on demand
    ...
  },
  // ...
];
```

2. **Suspense Wrapper**
```typescript
// frontend/src/App.tsx
import { Suspense } from 'react';

function App() {
  return (
    <Suspense fallback={<LoadingScreen />}>
      <Router>
        {/* Apps load on demand */}
      </Router>
    </Suspense>
  );
}
```

3. **Bundle Analysis**
```bash
# Analyze bundle size
npm run build
npx vite-bundle-visualizer

# Look for:
# - Large dependencies (can we tree-shake?)
# - Duplicate code (can we dedupe?)
# - Unused imports (can we remove?)
```

4. **Optimize Dependencies**
```typescript
// BAD: Import entire library
import _ from 'lodash';
const result = _.debounce(fn, 300);

// GOOD: Import only what you need
import debounce from 'lodash/debounce';
const result = debounce(fn, 300);
```

5. **Remove Unused Dependencies**
```bash
# Check for unused deps
npx depcheck

# Remove if not used
npm uninstall <package>
```

### 2.2 React Re-render Optimization

**Problem:** Unnecessary re-renders

**Optimierungen:**

1. **Use React.memo for Components**
```typescript
// frontend/src/components/FileItem.tsx

// BAD: Re-renders on every parent update
export function FileItem({ file, onSelect }) {
  return <div onClick={() => onSelect(file)}>{file.name}</div>;
}

// GOOD: Only re-renders when props change
export const FileItem = React.memo(({ file, onSelect }) => {
  return <div onClick={() => onSelect(file)}>{file.name}</div>;
});
```

2. **Use useMemo for Expensive Calculations**
```typescript
// BAD: Recalculates on every render
function FileList({ files }) {
  const sortedFiles = files.sort((a, b) => a.name.localeCompare(b.name));
  const filteredFiles = sortedFiles.filter(f => f.size > 1000);

  return <div>{filteredFiles.map(...)}</div>;
}

// GOOD: Only recalculates when files change
function FileList({ files }) {
  const processedFiles = useMemo(() => {
    const sorted = [...files].sort((a, b) => a.name.localeCompare(b.name));
    return sorted.filter(f => f.size > 1000);
  }, [files]);

  return <div>{processedFiles.map(...)}</div>;
}
```

3. **Use useCallback for Event Handlers**
```typescript
// BAD: New function on every render
function FileManager() {
  const handleFileSelect = (file) => {
    console.log(file);
  };

  return files.map(f => <FileItem file={f} onSelect={handleFileSelect} />);
}

// GOOD: Stable function reference
function FileManager() {
  const handleFileSelect = useCallback((file) => {
    console.log(file);
  }, []); // Dependencies

  return files.map(f => <FileItem file={f} onSelect={handleFileSelect} />);
}
```

4. **Virtualization for Long Lists**
```typescript
// frontend/src/apps/FileManager/FileManager.tsx
import { FixedSizeList } from 'react-window';

// BAD: Renders ALL files (1000s of DOM nodes)
function FileList({ files }) {
  return (
    <div>
      {files.map(file => <FileItem key={file.id} file={file} />)}
    </div>
  );
}

// GOOD: Only renders visible items (~20 DOM nodes)
function FileList({ files }) {
  return (
    <FixedSizeList
      height={600}
      itemCount={files.length}
      itemSize={50}
      width="100%"
    >
      {({ index, style }) => (
        <div style={style}>
          <FileItem file={files[index]} />
        </div>
      )}
    </FixedSizeList>
  );
}
```

### 2.3 API Request Optimization

**Problem:** Zu viele API-Requests

**Optimierungen:**

1. **Request Deduplication**
```typescript
// frontend/src/api/client.ts
const requestCache = new Map<string, Promise<any>>();

async function fetchWithDedup(url: string) {
  // Check if request is already in flight
  if (requestCache.has(url)) {
    return requestCache.get(url);
  }

  // Start request and cache promise
  const promise = fetch(url).then(r => r.json());
  requestCache.set(url, promise);

  // Clear cache after completion
  promise.finally(() => {
    setTimeout(() => requestCache.delete(url), 100);
  });

  return promise;
}
```

2. **SWR (Stale-While-Revalidate)**
```bash
npm install swr
```

```typescript
// frontend/src/hooks/useFiles.ts
import useSWR from 'swr';

export function useFiles(path: string) {
  const { data, error, mutate } = useSWR(
    `/api/v1/files?path=${path}`,
    fetcher,
    {
      revalidateOnFocus: false, // Don't refetch on window focus
      dedupingInterval: 2000,   // Dedupe requests within 2s
    }
  );

  return {
    files: data,
    isLoading: !error && !data,
    error,
    refresh: mutate,
  };
}
```

3. **Batch Multiple Requests**
```typescript
// BAD: 3 separate requests
async function loadDashboard() {
  const metrics = await fetch('/api/v1/metrics/current');
  const health = await fetch('/api/v1/system/health');
  const alerts = await fetch('/api/v1/alerts');
}

// GOOD: Single batch request
async function loadDashboard() {
  const data = await fetch('/api/v1/dashboard'); // Returns all data
  // Backend combines all 3 queries
}
```

### 2.4 Image & Asset Optimization

**Optimierungen:**

1. **Optimize Images**
```bash
# Install imagemin
npm install -D vite-plugin-imagemin

# Add to vite.config.ts
import viteImagemin from 'vite-plugin-imagemin';

export default defineConfig({
  plugins: [
    viteImagemin({
      gifsicle: { optimizationLevel: 7 },
      optipng: { optimizationLevel: 7 },
      mozjpeg: { quality: 80 },
      pngquant: { quality: [0.8, 0.9] },
      svgo: { plugins: [{ name: 'removeViewBox' }] },
    }),
  ],
});
```

2. **Use SVG Icons instead of Emojis**
```typescript
// BAD: Emoji icons (can render differently)
icon: '📊'

// GOOD: SVG icons (consistent, scalable)
import { ChartBarIcon } from '@heroicons/react/24/outline';
icon: <ChartBarIcon className="w-6 h-6" />
```

3. **Lazy Load Images**
```typescript
<img
  src={thumbnailUrl}
  loading="lazy" // Native lazy loading
  alt="File preview"
/>
```

### 2.5 Framer Motion Optimization

**Problem:** Animationen verlangsamen UI

**Optimierungen:**

1. **Reduce Motion Complexity**
```typescript
// BAD: Complex animation on every item
{items.map(item => (
  <motion.div
    initial={{ opacity: 0, x: -100, rotate: 90 }}
    animate={{ opacity: 1, x: 0, rotate: 0 }}
    transition={{ type: 'spring', stiffness: 100 }}
  >
    {item}
  </motion.div>
))}

// GOOD: Simple, GPU-accelerated animation
{items.map(item => (
  <motion.div
    initial={{ opacity: 0 }}
    animate={{ opacity: 1 }}
    transition={{ duration: 0.2 }}
  >
    {item}
  </motion.div>
))}
```

2. **Use layoutId Sparingly**
```typescript
// Only use layoutId for shared element transitions
// Not needed for simple animations
```

3. **Disable Animations for Long Lists**
```typescript
function FileList({ files }) {
  // No animation for lists > 100 items
  const shouldAnimate = files.length < 100;

  return files.map((file, index) => (
    shouldAnimate ? (
      <motion.div animate={{ opacity: 1 }}>{file.name}</motion.div>
    ) : (
      <div>{file.name}</div>
    )
  ));
}
```

---

## 3. Database Performance (PostgreSQL)

### 3.1 Query Optimization

**Optimierungen:**

1. **Add Missing Indexes**
```sql
-- Analyze slow queries
EXPLAIN ANALYZE SELECT * FROM audit_logs WHERE user_id = 1;

-- Add indexes for foreign keys
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_shares_user_id ON shares(user_id);
CREATE INDEX idx_files_parent_id ON files(parent_id);

-- Composite index for common queries
CREATE INDEX idx_audit_logs_user_timestamp ON audit_logs(user_id, created_at DESC);
```

2. **Use LIMIT for Large Result Sets**
```go
// Always paginate large results
func GetAuditLogs(page, pageSize int) []AuditLog {
    var logs []AuditLog
    db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&logs)
    return logs
}
```

3. **Batch Inserts**
```go
// BAD: Insert one-by-one
for _, log := range logs {
    db.Create(&log) // N queries
}

// GOOD: Batch insert
db.CreateInBatches(logs, 100) // Single query
```

### 3.2 Connection Pool Tuning

```yaml
# config.yaml
database:
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m
```

---

## 4. Network Performance

### 4.1 HTTP/2 & Compression

**Optimierungen:**

1. **Enable Gzip Compression**
```go
// backend/internal/api/middleware/compression.go
import "github.com/go-chi/chi/v5/middleware"

router.Use(middleware.Compress(5)) // Gzip level 5
```

2. **HTTP/2 Support**
```go
// backend/cmd/stumpfworks-server/main.go
srv := &http.Server{
    Addr:    ":8080",
    Handler: router,
}

// HTTP/2 automatically enabled for HTTPS
srv.ListenAndServeTLS("cert.pem", "key.pem")
```

### 4.2 WebSocket Optimization

**Problem:** WebSocket-Verbindungen verbrauchen Ressourcen

**Optimierungen:**

1. **Heartbeat/Ping-Pong**
```go
// Keep connections alive efficiently
ticker := time.NewTicker(30 * time.Second)
for {
    select {
    case <-ticker.C:
        conn.WriteMessage(websocket.PingMessage, nil)
    }
}
```

2. **Message Batching**
```go
// Don't send every tiny update immediately
// Batch updates every 100ms
ticker := time.NewTicker(100 * time.Millisecond)
var batchedUpdates []Update

for {
    select {
    case update := <-updateChan:
        batchedUpdates = append(batchedUpdates, update)
    case <-ticker.C:
        if len(batchedUpdates) > 0 {
            conn.WriteJSON(batchedUpdates)
            batchedUpdates = nil
        }
    }
}
```

---

## 5. Caching Strategy

### 5.1 Multi-Level Cache

```
┌─────────────┐
│  Browser    │ ← Cache-Control headers (static assets)
├─────────────┤
│  Redis      │ ← API responses (optional, future)
├─────────────┤
│  In-Memory  │ ← Frequently accessed data (current)
├─────────────┤
│  Database   │ ← Persistent storage
└─────────────┘
```

### 5.2 Cache Implementation

```go
// backend/pkg/cache/multi_level.go
type MultiLevelCache struct {
    memory *sync.Map
    ttl    time.Duration
}

func (c *MultiLevelCache) Get(key string) (interface{}, bool) {
    // Try memory first
    if val, ok := c.memory.Load(key); ok {
        return val, true
    }

    // TODO: Try Redis if available

    return nil, false
}
```

### 5.3 Cache-Control Headers

```go
// Static assets
w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year

// API responses
w.Header().Set("Cache-Control", "private, max-age=60") // 1 minute

// No cache for sensitive data
w.Header().Set("Cache-Control", "no-store")
```

---

## 6. Monitoring & Profiling

### 6.1 Go Profiling

**Enable pprof:**
```go
// backend/cmd/stumpfworks-server/main.go
import _ "net/http/pprof"

go func() {
    log.Println("Starting pprof server on :6060")
    http.ListenAndServe(":6060", nil)
}()
```

**Profile CPU:**
```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

**Profile Memory:**
```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

### 6.2 Frontend Performance Monitoring

**React DevTools Profiler:**
```bash
# Use React DevTools in browser
# Record performance during interactions
# Identify slow components
```

**Lighthouse CI:**
```bash
npm install -g @lhci/cli

# Run Lighthouse
lhci autorun --upload.target=temporary-public-storage
```

---

## 7. Performance Checklist

### Backend ✅
- [ ] Add timing middleware
- [ ] Implement caching for expensive operations
- [ ] Use goroutines for parallel operations
- [ ] Optimize database queries (indexes, eager loading)
- [ ] Tune database connection pool
- [ ] Enable Gzip compression
- [ ] Optimize file operations (chunk size, streaming)
- [ ] Add pprof for profiling

### Frontend ✅
- [ ] Implement code splitting (lazy loading)
- [ ] Use React.memo for components
- [ ] Use useMemo for expensive calculations
- [ ] Use useCallback for event handlers
- [ ] Virtualize long lists (react-window)
- [ ] Implement SWR or React Query
- [ ] Optimize images and assets
- [ ] Reduce animation complexity
- [ ] Bundle size analysis

### Database ✅
- [ ] Add missing indexes
- [ ] Use pagination for large results
- [ ] Batch inserts
- [ ] Analyze slow queries (EXPLAIN ANALYZE)

### Network ✅
- [ ] Enable HTTP/2
- [ ] Gzip compression
- [ ] Cache-Control headers
- [ ] WebSocket message batching

---

## 8. Expected Results

| Metric | Before | Target | Improvement |
|--------|--------|--------|-------------|
| **API Response** | ~500ms | < 100ms | 5x faster |
| **Initial Load** | ~5s | < 2s | 2.5x faster |
| **File Upload** | ~10 MB/s | ~50 MB/s | 5x faster |
| **Database Query** | ~200ms | < 50ms | 4x faster |
| **Bundle Size** | ~2 MB | < 500 KB | 4x smaller |

---

## 9. Implementation Priority

### Phase 1: Quick Wins (1 Woche) 🔴
- [ ] Add caching for disk/docker/metrics endpoints
- [ ] Code splitting (lazy load apps)
- [ ] Database indexes
- [ ] Gzip compression
- [ ] React.memo for heavy components

**Expected:** 20-30% performance improvement

### Phase 2: Deep Optimization (2 Wochen) 🟡
- [ ] Goroutines for parallel operations
- [ ] Virtualization for long lists
- [ ] SWR implementation
- [ ] Bundle size optimization
- [ ] WebSocket message batching

**Expected:** Additional 15-20% improvement

### Phase 3: Advanced (Optional) 🟢
- [ ] Redis caching layer
- [ ] CDN for static assets
- [ ] Service Workers
- [ ] Edge caching

**Expected:** Additional 10-15% improvement

---

## 10. Performance Testing

### Backend Load Testing
```bash
# Install hey
go install github.com/rakyll/hey@latest

# Test API endpoint
hey -n 1000 -c 10 http://localhost:8080/api/v1/disks

# Results:
# - Requests/sec
# - Average response time
# - 95th percentile
```

### Frontend Performance Testing
```bash
# Lighthouse
npx lighthouse http://localhost:8080 --view

# Metrics to check:
# - First Contentful Paint (FCP)
# - Largest Contentful Paint (LCP)
# - Time to Interactive (TTI)
# - Total Blocking Time (TBT)
```

---

## Zusammenfassung

**Gesamtaufwand:** 3-4 Wochen für vollständige Optimierung

**Erwartete Verbesserung:** 30-50% schneller

**Nächster Schritt:** Start mit Phase 1 (Quick Wins) für sofortige Verbesserung!
