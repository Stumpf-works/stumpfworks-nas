# Stumpf.Works NAS - Cluster Integration

**Vision:** Enterprise-Grade Multi-Node Cluster für horizontale Skalierung, Load Balancing und verteilte Speicher
**Version:** 1.3.0 (geplant für v1.3+)
**Status:** 🎯 Konzept-Phase

---

## Executive Summary

Eine **generelle Cluster-Integration** verwandelt Stumpf.Works NAS von einem Single-Server-System in ein **skalierbares, verteiltes Storage-Cluster**:

- ✅ **Multi-Node Support** (3+ Nodes statt nur 2)
- ✅ **Distributed Storage** (GlusterFS, Ceph)
- ✅ **Load Balancing** (HAProxy, Service Distribution)
- ✅ **Container Orchestration** (Docker Swarm, K3s)
- ✅ **Auto-Scaling** (Horizontal Skalierung)
- ✅ **Centralized Management** (Single Pane of Glass)
- ✅ **Distributed Monitoring** (Cluster-wide Metrics)

**Unterschied zu Phase 6 HA:**
- **Phase 6 HA:** 2-Node Failover (Active-Passive)
- **Cluster Integration:** Multi-Node Scale-Out (Active-Active)

---

## 1. Cluster-Architektur

### 1.1 Cluster-Topologie

```
┌──────────────────────────────────────────────────────────────────┐
│                    Stumpf.Works NAS Cluster                      │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Client Access Layer                                            │
│  ┌────────────────────────────────────────────────────────┐     │
│  │  Virtual IP: 10.0.0.100                                │     │
│  │  HAProxy Load Balancer (Web UI, API, SMB, NFS)        │     │
│  └────────────────────────────────────────────────────────┘     │
│                           │                                      │
│           ┌───────────────┼───────────────┐                     │
│           │               │               │                     │
│  ┌────────▼──────┐ ┌─────▼──────┐ ┌─────▼──────┐              │
│  │  Node 1       │ │  Node 2    │ │  Node 3    │              │
│  │  (Manager)    │ │  (Manager) │ │  (Worker)  │              │
│  ├───────────────┤ ├────────────┤ ├────────────┤              │
│  │ • Web UI      │ │ • Web UI   │ │ • Web UI   │              │
│  │ • API Server  │ │ • API      │ │ • API      │              │
│  │ • etcd        │ │ • etcd     │ │            │              │
│  │ • Docker      │ │ • Docker   │ │ • Docker   │              │
│  │ • GlusterFS   │ │ • GlusterFS│ │ • GlusterFS│              │
│  └───────────────┘ └────────────┘ └────────────┘              │
│           │               │               │                     │
│           └───────────────┼───────────────┘                     │
│                           │                                      │
│  ┌────────────────────────▼────────────────────────┐            │
│  │  Distributed Storage Layer                       │            │
│  │  • GlusterFS Replicated Volume (3x)             │            │
│  │  • Ceph Object Storage (Optional)                │            │
│  │  • DRBD for Critical Data (Active-Active)       │            │
│  └──────────────────────────────────────────────────┘            │
│                                                                  │
│  ┌──────────────────────────────────────────────────┐            │
│  │  Cluster Coordination Layer                      │            │
│  │  • etcd (Distributed Config Store)               │            │
│  │  • Consul (Service Discovery)                    │            │
│  │  • Raft Consensus (Leader Election)              │            │
│  └──────────────────────────────────────────────────┘            │
└──────────────────────────────────────────────────────────────────┘
```

### 1.2 Node-Typen

#### **Manager Nodes** (Control Plane)
- Führen Cluster-Management aus
- etcd Key-Value Store
- API Endpoint
- Web UI
- Minimum: 3 Nodes für Quorum (Raft)

#### **Worker Nodes** (Data Plane)
- Führen hauptsächlich Storage-Services aus
- Können skaliert werden (horizontal)
- Kein etcd
- Nur API Client

#### **Edge Nodes** (Optional)
- Spezielle Nodes für I/O-intensive Workloads
- z.B. Video-Transkodierung
- Docker Container Execution

### 1.3 Cluster-Modi

#### Mode 1: **Replicated Mode** (High Availability)
```
Node 1: [Data A] [Data B] [Data C]
Node 2: [Data A] [Data B] [Data C]  ← Full Copy
Node 3: [Data A] [Data B] [Data C]  ← Full Copy
```
- **Vorteil:** Maximale Redundanz, kein Datenverlust
- **Nachteil:** 3x Storage-Overhead
- **Use Case:** Kritische Daten, Datenbanken

#### Mode 2: **Distributed Mode** (Scale-Out)
```
Node 1: [Data A] [Data D] [Data G]
Node 2: [Data B] [Data E] [Data H]
Node 3: [Data C] [Data F] [Data I]
```
- **Vorteil:** Mehr Gesamtkapazität (3x)
- **Nachteil:** Node-Ausfall = Datenverlust
- **Use Case:** Große Dateien, Archive

#### Mode 3: **Dispersed Mode** (Erasure Coding)
```
Node 1: [Data A] [Parity 1]
Node 2: [Data B] [Parity 2]
Node 3: [Data C] [Parity 3]
Node 4: [Data D] [Parity 4]  ← Kann 1 Node Ausfall überleben
```
- **Vorteil:** Guter Balance zwischen Redundanz und Kapazität
- **Nachteil:** Komplexer, langsamer
- **Use Case:** Archiv-Storage, Cold Storage

---

## 2. Distributed Storage Layer

### 2.1 GlusterFS Integration

**Was ist GlusterFS?**
- Scale-out Network Attached Storage (NAS)
- FUSE-basiertes Dateisystem
- Keine Metadata-Server (fully distributed)
- Transparent für Clients

#### Backend Implementation

**Neue Datei:** `/backend/internal/system/cluster/glusterfs.go`

```go
package cluster

type GlusterFSManager struct {
    shell *executor.ShellExecutor
}

type GlusterVolume struct {
    Name          string   `json:"name"`
    Type          string   `json:"type"`        // Replicate, Distribute, Disperse
    Status        string   `json:"status"`      // Started, Stopped
    Bricks        []Brick  `json:"bricks"`
    ReplicaCount  int      `json:"replica_count"`
    DisperseCount int      `json:"disperse_count"`
    Options       map[string]string `json:"options"`
}

type Brick struct {
    Node string `json:"node"`
    Path string `json:"path"`
}

// Volume Operations
func (g *GlusterFSManager) CreateVolume(name string, volumeType string, bricks []Brick, replicaCount int) error
func (g *GlusterFSManager) DeleteVolume(name string) error
func (g *GlusterFSManager) StartVolume(name string) error
func (g *GlusterFSManager) StopVolume(name string) error
func (g *GlusterFSManager) GetVolumeInfo(name string) (*GlusterVolume, error)
func (g *GlusterFSManager) ListVolumes() ([]GlusterVolume, error)

// Brick Operations
func (g *GlusterFSManager) AddBrick(volumeName string, brick Brick) error
func (g *GlusterFSManager) RemoveBrick(volumeName string, brick Brick) error
func (g *GlusterFSManager) ReplaceBrick(volumeName string, oldBrick, newBrick Brick) error

// Rebalancing
func (g *GlusterFSManager) RebalanceVolume(volumeName string) error
func (g *GlusterFSManager) GetRebalanceStatus(volumeName string) (int, error) // Progress %

// Healing (Self-Heal)
func (g *GlusterFSManager) HealVolume(volumeName string) error
func (g *GlusterFSManager) GetHealInfo(volumeName string) (map[string]int, error) // Node → Files to heal

// Peer Management
func (g *GlusterFSManager) AddPeer(hostname string) error
func (g *GlusterFSManager) RemovePeer(hostname string) error
func (g *GlusterFSManager) ListPeers() ([]string, error)
```

**Implementation Example:**

```go
// CreateVolume - Erstellt GlusterFS Volume
func (g *GlusterFSManager) CreateVolume(name string, volumeType string, bricks []Brick, replicaCount int) error {
    // Build brick list: node1:/data/brick1 node2:/data/brick1 node3:/data/brick1
    var brickArgs []string
    for _, brick := range bricks {
        brickArgs = append(brickArgs, fmt.Sprintf("%s:%s", brick.Node, brick.Path))
    }

    args := []string{"volume", "create", name}

    // Add replica/disperse count
    switch volumeType {
    case "replicate":
        args = append(args, "replica", fmt.Sprintf("%d", replicaCount))
    case "disperse":
        args = append(args, "disperse", fmt.Sprintf("%d", replicaCount))
    case "distribute":
        // No extra args for distributed
    }

    args = append(args, brickArgs...)

    // gluster volume create myvolume replica 3 node1:/brick1 node2:/brick1 node3:/brick1
    _, err := g.shell.Execute("gluster", args...)
    if err != nil {
        return fmt.Errorf("failed to create volume: %w", err)
    }

    // Auto-start volume
    return g.StartVolume(name)
}

// GetVolumeInfo - Liest Volume-Informationen
func (g *GlusterFSManager) GetVolumeInfo(name string) (*GlusterVolume, error) {
    // gluster volume info <name> --xml
    output, err := g.shell.Execute("gluster", "volume", "info", name, "--xml")
    if err != nil {
        return nil, err
    }

    // Parse XML output
    var vol GlusterVolume
    // ... XML parsing ...

    return &vol, nil
}
```

#### Frontend UI

**Komponente:** `/frontend/src/apps/ClusterManager/tabs/DistributedStorage.tsx`

```tsx
export function DistributedStorage() {
  const [volumes, setVolumes] = useState<GlusterVolume[]>([]);
  const [showCreateDialog, setShowCreateDialog] = useState(false);

  return (
    <div>
      <h2>Distributed Storage (GlusterFS)</h2>

      {/* Volume List */}
      <VolumeList volumes={volumes} onEdit={handleEdit} onDelete={handleDelete} />

      {/* Create Volume Dialog */}
      <CreateVolumeDialog
        isOpen={showCreateDialog}
        onClose={() => setShowCreateDialog(false)}
        onSubmit={handleCreateVolume}
      />
    </div>
  );
}

// CreateVolumeDialog
interface CreateVolumeDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreateVolumeRequest) => void;
}

function CreateVolumeDialog({ isOpen, onClose, onSubmit }: CreateVolumeDialogProps) {
  const [volumeType, setVolumeType] = useState<'replicate' | 'distribute' | 'disperse'>('replicate');
  const [bricks, setBricks] = useState<Brick[]>([]);

  return (
    <Dialog open={isOpen} onClose={onClose}>
      <h3>Create GlusterFS Volume</h3>

      {/* Volume Type */}
      <Select label="Volume Type" value={volumeType} onChange={setVolumeType}>
        <option value="replicate">Replicated (High Availability)</option>
        <option value="distribute">Distributed (Scale-Out)</option>
        <option value="disperse">Dispersed (Erasure Coding)</option>
      </Select>

      {/* Brick Selection */}
      <BrickSelector
        nodes={clusterNodes}
        bricks={bricks}
        onChange={setBricks}
        minBricks={volumeType === 'replicate' ? 2 : 1}
      />

      {/* Replica Count (only for replicate) */}
      {volumeType === 'replicate' && (
        <Input
          type="number"
          label="Replica Count"
          min={2}
          max={bricks.length}
          value={replicaCount}
          onChange={setReplicaCount}
        />
      )}

      <Button onClick={handleSubmit}>Create Volume</Button>
    </Dialog>
  );
}
```

**UI Design:**

```
┌─────────────────────────────────────────────────────────┐
│  Distributed Storage (GlusterFS)                        │
├─────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────┐   │
│  │ Volume         Type       Status    Bricks  Size │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ data-volume    Replicate  Started   3x      5TB  │   │
│  │  └─ node1:/brick1 (5TB used)                     │   │
│  │  └─ node2:/brick1 (5TB used)                     │   │
│  │  └─ node3:/brick1 (5TB used)                     │   │
│  │                          [Rebalance] [Stop] [Del] │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ archive        Distribute Started   3x      15TB │   │
│  │  └─ node1:/brick2 (5TB)                          │   │
│  │  └─ node2:/brick2 (5TB)                          │   │
│  │  └─ node3:/brick2 (5TB)                          │   │
│  │                          [Rebalance] [Stop] [Del] │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  [+ Create Volume]  [Add Brick]  [Peer Management]      │
└─────────────────────────────────────────────────────────┘
```

### 2.2 Ceph Integration (Optional)

**Was ist Ceph?**
- Software-defined Storage
- Object Storage (S3-kompatibel)
- Block Storage (RBD)
- File Storage (CephFS)

**Backend:** `/backend/internal/system/cluster/ceph.go`

```go
type CephManager struct {
    shell *executor.ShellExecutor
}

// Cluster Operations
func (c *CephManager) GetClusterStatus() (*CephClusterStatus, error)
func (c *CephManager) GetOSDTree() ([]CephOSD, error)

// Pool Operations
func (c *CephManager) CreatePool(name string, pgNum int) error
func (c *CephManager) DeletePool(name string) error
func (c *CephManager) ListPools() ([]CephPool, error)

// RBD (Block Device) Operations
func (c *CephManager) CreateImage(poolName, imageName string, size uint64) error
func (c *CephManager) DeleteImage(poolName, imageName string) error
func (c *CephManager) MapImage(poolName, imageName string) (string, error) // Returns /dev/rbdX
func (c *CephManager) UnmapImage(device string) error
```

---

## 3. Load Balancing & Service Distribution

### 3.1 HAProxy Integration

**Was ist HAProxy?**
- High-Performance Load Balancer
- TCP/HTTP Load Balancing
- Health Checks
- SSL Termination

#### Use Cases im Cluster

1. **Web UI Load Balancing**
```
Client → HAProxy (10.0.0.100:8080)
          ├─→ Node 1:8080 (Round-Robin)
          ├─→ Node 2:8080
          └─→ Node 3:8080
```

2. **API Load Balancing**
```
API Calls → HAProxy (10.0.0.100:8080/api)
             ├─→ Node 1:8080/api
             ├─→ Node 2:8080/api
             └─→ Node 3:8080/api
```

3. **SMB/NFS Load Balancing**
```
SMB Client → HAProxy (10.0.0.100:445)
              ├─→ Node 1:445
              ├─→ Node 2:445
              └─→ Node 3:445
```

#### Backend Implementation

**Neue Datei:** `/backend/internal/system/cluster/haproxy.go`

```go
package cluster

type HAProxyManager struct {
    shell      *executor.ShellExecutor
    configPath string // /etc/haproxy/haproxy.cfg
}

type LoadBalancer struct {
    Name      string   `json:"name"`
    Frontend  Frontend `json:"frontend"`
    Backend   Backend  `json:"backend"`
}

type Frontend struct {
    Bind     string `json:"bind"`      // *:8080
    Mode     string `json:"mode"`      // tcp, http
    DefaultBackend string `json:"default_backend"`
}

type Backend struct {
    Name    string         `json:"name"`
    Mode    string         `json:"mode"`
    Balance string         `json:"balance"` // roundrobin, leastconn, source
    Servers []BackendServer `json:"servers"`
}

type BackendServer struct {
    Name   string `json:"name"`
    Host   string `json:"host"`
    Port   int    `json:"port"`
    Check  bool   `json:"check"`   // Enable health checks
    Backup bool   `json:"backup"`  // Backup server
}

// HAProxy Operations
func (h *HAProxyManager) CreateLoadBalancer(lb LoadBalancer) error
func (h *HAProxyManager) DeleteLoadBalancer(name string) error
func (h *HAProxyManager) GetLoadBalancerStatus(name string) (*LBStatus, error)
func (h *HAProxyManager) ListLoadBalancers() ([]LoadBalancer, error)
func (h *HAProxyManager) AddBackendServer(lbName string, server BackendServer) error
func (h *HAProxyManager) RemoveBackendServer(lbName, serverName string) error
func (h *HAProxyManager) GetStats() (*HAProxyStats, error)
func (h *HAProxyManager) Reload() error // Reload config without downtime
```

**Implementation:**

```go
// CreateLoadBalancer - Erstellt HAProxy Load Balancer
func (h *HAProxyManager) CreateLoadBalancer(lb LoadBalancer) error {
    // Read existing config
    config, err := h.readConfig()
    if err != nil {
        return err
    }

    // Add frontend section
    frontendConfig := fmt.Sprintf(`
frontend %s
    bind %s
    mode %s
    default_backend %s
`, lb.Name, lb.Frontend.Bind, lb.Frontend.Mode, lb.Frontend.DefaultBackend)

    // Add backend section
    backendConfig := fmt.Sprintf(`
backend %s
    mode %s
    balance %s
`, lb.Backend.Name, lb.Backend.Mode, lb.Backend.Balance)

    for _, server := range lb.Backend.Servers {
        checkStr := ""
        if server.Check {
            checkStr = " check"
        }
        backupStr := ""
        if server.Backup {
            backupStr = " backup"
        }
        backendConfig += fmt.Sprintf("    server %s %s:%d%s%s\n",
            server.Name, server.Host, server.Port, checkStr, backupStr)
    }

    // Append to config
    config += frontendConfig + backendConfig

    // Write config
    if err := h.writeConfig(config); err != nil {
        return err
    }

    // Reload HAProxy
    return h.Reload()
}

// Reload - Reload HAProxy without downtime
func (h *HAProxyManager) Reload() error {
    // systemctl reload haproxy (graceful reload)
    _, err := h.shell.Execute("systemctl", "reload", "haproxy")
    return err
}
```

#### Frontend UI

**Komponente:** `/frontend/src/apps/ClusterManager/tabs/LoadBalancer.tsx`

```tsx
export function LoadBalancer() {
  const [loadBalancers, setLoadBalancers] = useState<LoadBalancer[]>([]);

  return (
    <div>
      <h2>Load Balancer (HAProxy)</h2>

      {/* Load Balancer List */}
      <LBList loadBalancers={loadBalancers} />

      {/* Create LB */}
      <CreateLBDialog />
    </div>
  );
}
```

**UI Design:**

```
┌─────────────────────────────────────────────────────────┐
│  Load Balancer (HAProxy)                                │
├─────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────┐   │
│  │ Name        Frontend      Backend Servers  Stats │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ webui-lb    *:8080 (HTTP) 3 servers       [View] │   │
│  │  └─ node1:8080  ✅ UP    (50 req/s)              │   │
│  │  └─ node2:8080  ✅ UP    (45 req/s)              │   │
│  │  └─ node3:8080  ❌ DOWN  (0 req/s)               │   │
│  │                             [Edit] [Delete]       │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ smb-lb      *:445 (TCP)   3 servers       [View] │   │
│  │  └─ node1:445   ✅ UP    (1.2 Gbps)              │   │
│  │  └─ node2:445   ✅ UP    (1.1 Gbps)              │   │
│  │  └─ node3:445   ✅ UP    (0.9 Gbps)              │   │
│  │                             [Edit] [Delete]       │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  [+ Create Load Balancer]  [View Global Stats]          │
└─────────────────────────────────────────────────────────┘
```

---

## 4. Container Orchestration

### 4.1 Docker Swarm Mode

**Was ist Docker Swarm?**
- Native Docker Cluster
- Service Orchestration
- Load Balancing
- Rolling Updates

**Vorteil:** Bereits in Docker integriert, keine zusätzliche Software

#### Backend Implementation

**Erweitern:** `/backend/internal/docker/swarm.go`

```go
package docker

type SwarmManager struct {
    client *client.Client
}

// Swarm Operations
func (s *SwarmManager) InitSwarm(advertiseAddr string) error
func (s *SwarmManager) JoinSwarm(token, managerAddr string) error
func (s *SwarmManager) LeaveSwarm(force bool) error
func (s *SwarmManager) GetSwarmInfo() (*SwarmInfo, error)
func (s *SwarmManager) ListNodes() ([]SwarmNode, error)

// Service Operations
func (s *SwarmManager) CreateService(spec ServiceSpec) (string, error)
func (s *SwarmManager) UpdateService(serviceID string, spec ServiceSpec) error
func (s *SwarmManager) DeleteService(serviceID string) error
func (s *SwarmManager) ListServices() ([]SwarmService, error)
func (s *SwarmManager) GetServiceLogs(serviceID string) (string, error)
func (s *SwarmManager) ScaleService(serviceID string, replicas uint64) error

// Stack Operations (Docker Compose on Swarm)
func (s *SwarmManager) DeployStack(name string, composeFile []byte) error
func (s *SwarmManager) RemoveStack(name string) error
func (s *SwarmManager) ListStacks() ([]Stack, error)
```

**UI:**

```
┌─────────────────────────────────────────────────────────┐
│  Container Orchestration (Docker Swarm)                 │
├─────────────────────────────────────────────────────────┤
│  Swarm Status: ✅ Active (3 nodes, 1 manager)           │
│                                                         │
│  Nodes:                                                 │
│  ┌──────────────────────────────────────────────────┐   │
│  │ Node       Role     Status  Availability  CPU    │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ node1      Manager  Ready   Active        45%   │   │
│  │ node2      Worker   Ready   Active        72%   │   │
│  │ node3      Worker   Ready   Active        58%   │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  Services:                                              │
│  ┌──────────────────────────────────────────────────┐   │
│  │ Name       Replicas   Mode       Image           │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ nginx      3/3        Replicated nginx:latest    │   │
│  │  └─ nginx.1  node1  ✅ Running                   │   │
│  │  └─ nginx.2  node2  ✅ Running                   │   │
│  │  └─ nginx.3  node3  ✅ Running                   │   │
│  │                                  [Scale] [Update]│   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  [+ Deploy Service]  [Deploy Stack]                     │
└─────────────────────────────────────────────────────────┘
```

### 4.2 K3s Integration (Optional)

**Was ist K3s?**
- Lightweight Kubernetes Distribution
- Perfect für Edge/IoT
- < 100 MB Binary
- Full Kubernetes API

**Use Case:** Wenn du vollständige Kubernetes-Features brauchst

---

## 5. Cluster Coordination & Service Discovery

### 5.1 etcd Integration

**Was ist etcd?**
- Distributed Key-Value Store
- Raft Consensus
- Strong Consistency
- Configuration Management

#### Backend Implementation

**Neue Datei:** `/backend/internal/system/cluster/etcd.go`

```go
package cluster

import clientv3 "go.etcd.io/etcd/client/v3"

type EtcdManager struct {
    client *clientv3.Client
}

// Key-Value Operations
func (e *EtcdManager) Put(key, value string) error
func (e *EtcdManager) Get(key string) (string, error)
func (e *EtcdManager) Delete(key string) error
func (e *EtcdManager) List(prefix string) (map[string]string, error)

// Watch Operations
func (e *EtcdManager) Watch(key string, callback func(string, string)) error

// Leader Election
func (e *EtcdManager) ElectLeader(sessionName string) (bool, error) // Returns true if elected
func (e *EtcdManager) ResignLeader(sessionName string) error

// Distributed Lock
func (e *EtcdManager) AcquireLock(lockName string, ttl int) error
func (e *EtcdManager) ReleaseLock(lockName string) error
```

**Use Cases:**

1. **Cluster-weite Konfiguration**
```go
// Schreibe Config
etcdManager.Put("/cluster/config/smb/workgroup", "WORKGROUP")

// Lese Config auf anderem Node
workgroup, _ := etcdManager.Get("/cluster/config/smb/workgroup")
```

2. **Service Discovery**
```go
// Node registriert sich selbst
etcdManager.Put("/cluster/nodes/node1/ip", "10.0.0.11")
etcdManager.Put("/cluster/nodes/node1/status", "online")

// Andere Nodes finden alle Nodes
nodes := etcdManager.List("/cluster/nodes/")
```

3. **Leader Election**
```go
// Nur ein Node führt Backup aus
isLeader, _ := etcdManager.ElectLeader("backup-scheduler")
if isLeader {
    // Führe Backup aus
}
```

### 5.2 Consul Integration (Alternative zu etcd)

**Was ist Consul?**
- Service Discovery
- Health Checking
- KV Store
- Multi-Datacenter Support

**Vorteil über etcd:** Built-in Service Discovery + Health Checks

---

## 6. Cluster Manager UI

### 6.1 Neue App: Cluster Manager

**Datei:** `/frontend/src/apps/ClusterManager/ClusterManager.tsx`

```tsx
type ClusterTab = 'overview' | 'nodes' | 'storage' | 'loadbalancer' | 'orchestration' | 'config';

export function ClusterManager() {
  const [activeTab, setActiveTab] = useState<ClusterTab>('overview');

  return (
    <div className="cluster-manager">
      <Header title="Cluster Management" />

      <Tabs>
        <Tab id="overview" icon="🌐" label="Overview" />
        <Tab id="nodes" icon="🖥️" label="Nodes" />
        <Tab id="storage" icon="💾" label="Distributed Storage" />
        <Tab id="loadbalancer" icon="⚖️" label="Load Balancer" />
        <Tab id="orchestration" icon="🐳" label="Container Orchestration" />
        <Tab id="config" icon="⚙️" label="Configuration" />
      </Tabs>

      <TabContent>
        {activeTab === 'overview' && <ClusterOverview />}
        {activeTab === 'nodes' && <NodeManagement />}
        {activeTab === 'storage' && <DistributedStorage />}
        {activeTab === 'loadbalancer' && <LoadBalancer />}
        {activeTab === 'orchestration' && <ContainerOrchestration />}
        {activeTab === 'config' && <ClusterConfiguration />}
      </TabContent>
    </div>
  );
}
```

### 6.2 Cluster Overview Tab

**UI Design:**

```
┌─────────────────────────────────────────────────────────┐
│  Cluster Overview                                       │
├─────────────────────────────────────────────────────────┤
│  Cluster Name: production-nas                           │
│  Status: ✅ Healthy (3/3 nodes online)                  │
│  Leader: node1                                          │
│                                                         │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐    │
│  │ Total Nodes  │ │ Total Storage│ │ Active Svcs  │    │
│  │      3       │ │     45 TB    │ │     12       │    │
│  └──────────────┘ └──────────────┘ └──────────────┘    │
│                                                         │
│  Cluster Topology:                                      │
│  ┌────────────────────────────────────────────────┐     │
│  │       ┌─────────┐                              │     │
│  │       │ node1   │  Manager, Leader             │     │
│  │       │ 10.0.11 │  ✅ Online                   │     │
│  │       └────┬────┘                              │     │
│  │   ┌────────┴────────┐                          │     │
│  │   │                 │                          │     │
│  │ ┌─▼─────┐       ┌──▼──────┐                   │     │
│  │ │ node2 │       │ node3   │                    │     │
│  │ │10.0.12│       │ 10.0.13 │                    │     │
│  │ │✅ Onl.│       │ ✅ Onl. │                    │     │
│  │ └───────┘       └─────────┘                    │     │
│  │ Manager         Worker                         │     │
│  └────────────────────────────────────────────────┘     │
│                                                         │
│  Recent Events:                                         │
│  • 10:45 - node3 joined cluster                        │
│  • 10:30 - GlusterFS volume 'data' healed 15 files    │
│  • 10:15 - Service 'nginx' scaled from 2 to 3 replicas│
└─────────────────────────────────────────────────────────┘
```

### 6.3 Node Management Tab

```
┌─────────────────────────────────────────────────────────┐
│  Node Management                                        │
├─────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────┐   │
│  │ Node    IP        Role    Status  CPU  RAM  Disk │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ 👑 node1 10.0.0.11 Manager Online  45%  72%  68% │   │
│  │    • etcd Leader                                 │   │
│  │    • 5 Services Running                          │   │
│  │    • GlusterFS Brick: /data/brick1 (15TB)        │   │
│  │                       [SSH] [Logs] [Drain]       │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ ✅ node2 10.0.0.12 Manager Online  58%  65%  71% │   │
│  │    • etcd Member                                 │   │
│  │    • 4 Services Running                          │   │
│  │    • GlusterFS Brick: /data/brick1 (15TB)        │   │
│  │                       [SSH] [Logs] [Drain]       │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ ⚙️ node3 10.0.0.13 Worker  Online  72%  81%  65% │   │
│  │    • Docker Swarm Worker                         │   │
│  │    • 3 Services Running                          │   │
│  │    • GlusterFS Brick: /data/brick1 (15TB)        │   │
│  │                       [SSH] [Logs] [Remove]      │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  [+ Add Node]  [Promote to Manager]  [Demote to Worker]│
└─────────────────────────────────────────────────────────┘
```

---

## 7. Cluster-wide Features

### 7.1 Distributed Monitoring

**Prometheus Federation:**
- Jeder Node hat eigenen Prometheus
- Zentrale Prometheus-Instanz federated alle Nodes
- Grafana Dashboard zeigt Cluster-wide Metriken

**UI:**

```
┌─────────────────────────────────────────────────────────┐
│  Cluster Monitoring (Aggregated)                        │
├─────────────────────────────────────────────────────────┤
│  Total CPU Usage:  65% ████████▓░░░                     │
│  Total RAM Usage:  72% █████████░░░                     │
│  Total Disk Usage: 68% ████████▓░░░                     │
│                                                         │
│  Per-Node Breakdown:                                    │
│  [Line Chart showing CPU/RAM/Disk for all 3 nodes]     │
│                                                         │
│  Active Alerts:                                         │
│  ⚠️ node3 - High CPU (>80%) for 5 minutes              │
│  ⚠️ GlusterFS - Self-heal in progress (15 files)       │
└─────────────────────────────────────────────────────────┘
```

### 7.2 Distributed Logging

**Loki/Promtail:**
- Promtail auf jedem Node
- Loki zentral
- Logs von allen Nodes aggregiert

### 7.3 Centralized Configuration Management

**etcd-basierte Config:**
```
/cluster/
├── config/
│   ├── smb/workgroup
│   ├── smb/netbios_name
│   ├── nfs/version
│   └── global/timezone
├── nodes/
│   ├── node1/
│   │   ├── ip
│   │   ├── status
│   │   └── last_heartbeat
│   ├── node2/
│   └── node3/
└── services/
    ├── web-ui/
    │   ├── replicas
    │   └── version
    └── nginx/
```

---

## 8. Implementation Roadmap

### Phase 1: Foundation (Woche 1-3)

**Sprint 1.1: etcd Integration (5 Tage)**
- [ ] etcd Setup (3-Node Cluster)
- [ ] Backend: etcd Manager
- [ ] API Endpoints für KV Operations
- [ ] Frontend: Config Management UI

**Sprint 1.2: Node Discovery & Management (5 Tage)**
- [ ] Node Registration in etcd
- [ ] Heartbeat Mechanism
- [ ] Frontend: Node Management Tab
- [ ] Node Metrics Collection

**Sprint 1.3: Cluster Manager App (5 Tage)**
- [ ] Frontend: Cluster Manager App Shell
- [ ] Overview Tab mit Topology
- [ ] Node Management Tab
- [ ] Basic Health Checks

### Phase 2: Distributed Storage (Woche 4-6)

**Sprint 2.1: GlusterFS Backend (7 Tage)**
- [ ] GlusterFS Manager Implementation
- [ ] Volume Operations
- [ ] Brick Operations
- [ ] Peer Management
- [ ] API Endpoints (10+)

**Sprint 2.2: GlusterFS Frontend (5 Tage)**
- [ ] Frontend API Client
- [ ] Distributed Storage Tab
- [ ] Create Volume Dialog
- [ ] Volume Management UI
- [ ] Rebalance & Heal UI

**Sprint 2.3: Testing & Optimization (3 Tage)**
- [ ] 3-Node GlusterFS Setup Testing
- [ ] Replicated Volume Testing
- [ ] Distributed Volume Testing
- [ ] Performance Benchmarking

### Phase 3: Load Balancing (Woche 7-8)

**Sprint 3.1: HAProxy Integration (5 Tage)**
- [ ] HAProxy Manager Backend
- [ ] Load Balancer Creation
- [ ] Backend Server Management
- [ ] Stats Collection
- [ ] API Endpoints (8+)

**Sprint 3.2: HAProxy Frontend (3 Tage)**
- [ ] Load Balancer Tab
- [ ] Create LB Dialog
- [ ] Stats Dashboard
- [ ] Health Check Visualization

**Sprint 3.3: Pre-configured LBs (2 Tage)**
- [ ] Auto-create Web UI LB
- [ ] Auto-create API LB
- [ ] Auto-create SMB LB (optional)
- [ ] Auto-create NFS LB (optional)

### Phase 4: Container Orchestration (Woche 9-11)

**Sprint 4.1: Docker Swarm Backend (7 Tage)**
- [ ] Swarm Manager Implementation
- [ ] Service Operations
- [ ] Stack Operations
- [ ] Node Management
- [ ] API Endpoints (15+)

**Sprint 4.2: Docker Swarm Frontend (5 Tage)**
- [ ] Container Orchestration Tab
- [ ] Service List & Management
- [ ] Stack Deployment UI
- [ ] Swarm Node Visualization

**Sprint 4.3: Integration with Docker Manager (3 Tage)**
- [ ] Merge Swarm Tab into Docker Manager
- [ ] Unified Docker UI (Standalone + Swarm)
- [ ] Migration Path (Standalone → Swarm)

### Phase 5: Advanced Features (Woche 12-14)

**Sprint 5.1: Distributed Monitoring (5 Tage)**
- [ ] Prometheus Federation Setup
- [ ] Per-Node Exporters
- [ ] Centralized Grafana
- [ ] Cluster-wide Dashboards

**Sprint 5.2: Distributed Logging (5 Tage)**
- [ ] Loki Setup
- [ ] Promtail auf allen Nodes
- [ ] Log Aggregation UI
- [ ] Log Search & Filter

**Sprint 5.3: Auto-Scaling (Optional, 5 Tage)**
- [ ] Resource Monitoring
- [ ] Auto-Scaling Rules Engine
- [ ] Scale-Up/Scale-Down Logic
- [ ] UI für Scaling Policies

### Phase 6: Testing & Documentation (Woche 15)

**Sprint 6: Integration & Docs (5 Tage)**
- [ ] End-to-End Testing
- [ ] Cluster Setup Guide
- [ ] Best Practices Documentation
- [ ] Performance Tuning Guide
- [ ] Troubleshooting Guide

---

## 9. Gesamtaufwand

| Phase | Dauer | Features |
|-------|-------|----------|
| **Phase 1** | 3 Wochen | etcd, Node Management, Cluster UI |
| **Phase 2** | 3 Wochen | GlusterFS Distributed Storage |
| **Phase 3** | 2 Wochen | HAProxy Load Balancing |
| **Phase 4** | 3 Wochen | Docker Swarm Orchestration |
| **Phase 5** | 3 Wochen | Monitoring, Logging, Auto-Scaling |
| **Phase 6** | 1 Woche | Testing & Docs |

**Gesamt:** **15 Wochen** (~3.5-4 Monate)

---

## 10. Deliverables v1.3.0

✅ **Multi-Node Cluster Support** (3+ Nodes)
✅ **etcd Distributed Configuration**
✅ **GlusterFS Distributed Storage** (Replicate/Distribute/Disperse)
✅ **HAProxy Load Balancing** (Web UI, API, SMB, NFS)
✅ **Docker Swarm Orchestration**
✅ **Cluster Manager UI** (Single Pane of Glass)
✅ **Distributed Monitoring** (Prometheus Federation)
✅ **Distributed Logging** (Loki)
✅ **Node Management** (Add, Remove, Drain)
✅ **Service Discovery** (etcd/Consul)
✅ **Auto-Scaling** (Optional)

---

## 11. Vergleich: HA vs Cluster Integration

| Feature | Phase 6 HA | Cluster Integration |
|---------|------------|---------------------|
| **Nodes** | 2 Nodes (Active-Passive) | 3+ Nodes (Active-Active) |
| **Failover** | ✅ Automatic Failover | ✅ Automatic + Load Balancing |
| **Storage** | DRBD (1:1 Mirror) | GlusterFS (N-Way Replicate/Distribute) |
| **Scaling** | ❌ Vertical Only | ✅ Horizontal Scale-Out |
| **Load Balancing** | ⚠️ VIP Only | ✅ HAProxy Multi-Backend |
| **Orchestration** | ❌ Manual | ✅ Docker Swarm |
| **Configuration** | Per-Node | ✅ Centralized (etcd) |
| **Monitoring** | Per-Node | ✅ Cluster-wide Aggregated |
| **Use Case** | Small Business, Homelab | Enterprise, Large Scale |

---

## 12. Hardware-Anforderungen

### Minimum Cluster (3 Nodes)

**Node 1 (Manager):**
- CPU: 4 Cores
- RAM: 8 GB
- Disk: 500 GB (OS + etcd)
- Network: 2x 1 Gbps (1 Data + 1 Heartbeat)

**Node 2 (Manager):**
- CPU: 4 Cores
- RAM: 8 GB
- Disk: 500 GB

**Node 3 (Worker):**
- CPU: 4 Cores
- RAM: 8 GB
- Disk: 500 GB

**Shared Storage (GlusterFS Bricks):**
- Node 1: 5 TB
- Node 2: 5 TB
- Node 3: 5 TB
- **Total:** 15 TB (Replicated) oder 15 TB (Distributed)

### Empfohlen (5+ Nodes)

Für Production: **Minimum 5 Nodes**
- 3 Manager Nodes (Quorum)
- 2+ Worker Nodes

---

**Next Step:** Soll ich mit der Implementation starten oder noch Details ausarbeiten? 🚀
