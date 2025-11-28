# Cluster Activation Wizard - Easy Cluster Setup

**Konzept:** Cluster-Support ist **OPTIONAL** - Single-Node by Default
**Philosophie:** "Starte einfach, skaliere wenn nötig"

---

## Design-Prinzipien

### 1. **Single-Node First** ✅
- Standard-Installation ist IMMER Single-Node
- Volle Funktionalität ohne Cluster
- Keine Performance-Overhead

### 2. **Easy Activation** 🚀
- One-Click Cluster Wizard
- Guided Setup mit Validierung
- Auto-Detection von Nodes
- Rollback bei Fehler

### 3. **Gradual Migration** 📈
```
Single-Node → 2-Node HA → 3+ Node Cluster
    ↓            ↓              ↓
  Default    Failover    Scale-Out & Load Balancing
```

### 4. **Optional Features** ⚙️
- Distributed Storage: Optional aktivieren
- Load Balancing: Optional aktivieren
- Container Orchestration: Optional aktivieren
- Jedes Feature einzeln toggle-bar

---

## Cluster Activation Flow

### Schritt 1: Single-Node Installation (DEFAULT)

```
┌─────────────────────────────────────────────┐
│  Welcome to Stumpf.Works NAS!               │
├─────────────────────────────────────────────┤
│                                             │
│  Setup Mode:                                │
│  ◉ Standalone Server (Recommended)          │
│  ○ High Availability Cluster (2 Nodes)      │
│  ○ Scale-Out Cluster (3+ Nodes)             │
│                                             │
│  [Continue with Standalone]                 │
└─────────────────────────────────────────────┘
```

**Standard:** User klickt "Standalone" → Normale Installation

### Schritt 2: Cluster-Hinweis in Dashboard

**Nach Installation** zeigt Dashboard diskret:

```
┌─────────────────────────────────────────────┐
│  Dashboard                                  │
├─────────────────────────────────────────────┤
│  System Health: 92% ✅                      │
│  Storage: 5.2 TB / 10 TB                    │
│  Docker: 12 containers                      │
│                                             │
│  💡 Pro Tip:                                │
│  ┌──────────────────────────────────────┐   │
│  │ Need High Availability or more       │   │
│  │ storage? Enable Cluster Mode!        │   │
│  │                                       │   │
│  │ [Learn More] [Enable Cluster]        │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

**Nicht aufdringlich** - nur ein Hinweis, kein Zwang!

### Schritt 3: Cluster Activation Wizard

User klickt "Enable Cluster" → Wizard startet:

#### Screen 1: Cluster Mode Selection

```
┌─────────────────────────────────────────────────────────┐
│  Cluster Activation Wizard (Step 1/4)                   │
├─────────────────────────────────────────────────────────┤
│  Choose Cluster Mode:                                   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │ ◉ High Availability (2 Nodes)                    │   │
│  │   ✅ Automatic failover                          │   │
│  │   ✅ Zero downtime for maintenance               │   │
│  │   ✅ Data mirroring (DRBD)                       │   │
│  │   ⚠️  Storage: 2x (mirrored)                     │   │
│  │                                                   │   │
│  │   Recommended for: Small Business, Critical Data │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │ ○ Scale-Out Cluster (3+ Nodes)                   │   │
│  │   ✅ Horizontal scaling                          │   │
│  │   ✅ Load balancing                              │   │
│  │   ✅ Distributed storage (GlusterFS)             │   │
│  │   ✅ Container orchestration (Swarm)             │   │
│  │   ℹ️  Storage: Configurable (1x-3x)              │   │
│  │                                                   │   │
│  │   Recommended for: Enterprise, Large Scale       │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  [Cancel]                              [Next: Add Nodes]│
└─────────────────────────────────────────────────────────┘
```

#### Screen 2: Node Discovery

```
┌─────────────────────────────────────────────────────────┐
│  Cluster Activation Wizard (Step 2/4)                   │
├─────────────────────────────────────────────────────────┤
│  Add Cluster Nodes:                                     │
│                                                         │
│  Current Node: node1 (10.0.0.11) - This Server          │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │ Discovered Nodes:                                │   │
│  │                                                   │   │
│  │ ☐ 10.0.0.12  (nas-backup)    SSH: ✅  Ping: 2ms │   │
│  │    • Version: 1.1.0 Compatible                   │   │
│  │    • Disk Space: 8 TB available                  │   │
│  │    • RAM: 16 GB                                   │   │
│  │                                   [Add to Cluster]│   │
│  │                                                   │   │
│  │ ☐ 10.0.0.13  (nas-worker)    SSH: ✅  Ping: 3ms │   │
│  │    • Version: 1.1.0 Compatible                   │   │
│  │    • Disk Space: 12 TB available                 │   │
│  │    • RAM: 32 GB                                   │   │
│  │                                   [Add to Cluster]│   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  Or add manually:                                       │
│  IP Address: [_____________]  [Scan]                    │
│                                                         │
│  Selected Nodes: 0                                      │
│  Minimum Required: 1 (HA) or 2 (Scale-Out)              │
│                                                         │
│  [Back]                              [Next: Configure]  │
└─────────────────────────────────────────────────────────┘
```

**Auto-Discovery:**
- Scannt lokales Netzwerk (Subnet)
- Findet andere Stumpf.Works NAS Instanzen
- Prüft Version-Kompatibilität
- Zeigt Hardware-Specs

#### Screen 3: Feature Selection (Scale-Out only)

```
┌─────────────────────────────────────────────────────────┐
│  Cluster Activation Wizard (Step 3/4)                   │
├─────────────────────────────────────────────────────────┤
│  Select Cluster Features:                               │
│                                                         │
│  Core Features (Always Enabled):                        │
│  ✅ Multi-Node Management                               │
│  ✅ Node Health Monitoring                              │
│  ✅ Centralized Configuration (etcd)                    │
│                                                         │
│  Optional Features:                                     │
│                                                         │
│  ☑️ Distributed Storage (GlusterFS)                     │
│     Replicate data across nodes for redundancy          │
│     Required packages: glusterfs-server, glusterfs-client│
│     Storage overhead: 2x-3x (depending on replica count)│
│                                                         │
│  ☑️ Load Balancing (HAProxy)                            │
│     Distribute web/API/SMB traffic across nodes         │
│     Required packages: haproxy                          │
│                                                         │
│  ☑️ Container Orchestration (Docker Swarm)              │
│     Deploy containers across cluster nodes              │
│     Required: Docker 20.10+                             │
│                                                         │
│  ☐ Advanced Monitoring (Prometheus Federation)          │
│     Cluster-wide metrics aggregation (Optional)         │
│     Requires: Additional 2GB RAM per node               │
│                                                         │
│  [Back]                                     [Next: Review]│
└─────────────────────────────────────────────────────────┘
```

**User wählt Features:**
- Nur was sie brauchen
- Klare Beschreibung + Requirements
- Kann später aktiviert/deaktiviert werden

#### Screen 4: Review & Confirm

```
┌─────────────────────────────────────────────────────────┐
│  Cluster Activation Wizard (Step 4/4)                   │
├─────────────────────────────────────────────────────────┤
│  Review Configuration:                                  │
│                                                         │
│  Cluster Mode: Scale-Out (3 Nodes)                      │
│                                                         │
│  Nodes:                                                 │
│  • node1 (10.0.0.11) - Manager, This Server             │
│  • nas-backup (10.0.0.12) - Manager                     │
│  • nas-worker (10.0.0.13) - Worker                      │
│                                                         │
│  Features:                                              │
│  ✅ Distributed Storage (GlusterFS)                     │
│     Volume Type: Replicated (3x)                        │
│     Estimated Total Capacity: 8 TB (24 TB raw)          │
│                                                         │
│  ✅ Load Balancing (HAProxy)                            │
│     Virtual IP: 10.0.0.100 (auto-assigned)              │
│     Backends: Web UI, API, SMB, NFS                     │
│                                                         │
│  ✅ Container Orchestration (Docker Swarm)              │
│     Manager Nodes: 2                                    │
│     Worker Nodes: 1                                     │
│                                                         │
│  ⚠️  This will:                                         │
│  • Install packages on all nodes                        │
│  • Reconfigure network settings                         │
│  • Require brief downtime (~5 minutes)                  │
│  • Cannot be easily reversed (backup recommended!)      │
│                                                         │
│  ☑️ I have backed up my data                            │
│  ☑️ I understand this is a major change                 │
│                                                         │
│  [Back]         [Cancel]         [Activate Cluster 🚀]  │
└─────────────────────────────────────────────────────────┘
```

### Schritt 4: Activation Progress

```
┌─────────────────────────────────────────────────────────┐
│  Activating Cluster...                                  │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ✅ Validating node connectivity                        │
│  ✅ Installing dependencies on nodes                    │
│  🔄 Setting up etcd cluster (node1, node2)              │
│     ████████░░░░░░░░░░░░░░░░ 40%                       │
│  ⏳ Configuring GlusterFS                               │
│  ⏳ Setting up HAProxy                                  │
│  ⏳ Initializing Docker Swarm                           │
│  ⏳ Finalizing configuration                            │
│                                                         │
│  Current Step: etcd initialization...                   │
│  Estimated Time Remaining: 3 minutes                    │
│                                                         │
│  [View Detailed Logs]                                   │
└─────────────────────────────────────────────────────────┘
```

**Automatische Schritte:**
1. SSH Keys austauschen (passwordless)
2. Packages installieren (apt-get install ...)
3. etcd Cluster aufsetzen
4. GlusterFS Peers hinzufügen
5. HAProxy konfigurieren
6. Docker Swarm initialisieren
7. Firewall-Regeln setzen
8. Health Checks

**Bei Fehler:**
- Rollback auf vorherigen Zustand
- Detaillierte Fehlermeldung
- "Retry" oder "Cancel" Option

### Schritt 5: Success!

```
┌─────────────────────────────────────────────────────────┐
│  ✅ Cluster Activated Successfully!                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Your cluster is now ready:                             │
│                                                         │
│  🌐 Web Access:                                         │
│     http://10.0.0.100:8080 (Load-Balanced VIP)          │
│     http://10.0.0.11:8080 (Direct: node1)               │
│     http://10.0.0.12:8080 (Direct: nas-backup)          │
│     http://10.0.0.13:8080 (Direct: nas-worker)          │
│                                                         │
│  📊 Cluster Status:                                     │
│     • 3/3 nodes online and healthy                      │
│     • GlusterFS volume 'data' replicated 3x             │
│     • HAProxy load balancer active                      │
│     • Docker Swarm ready (2 managers, 1 worker)         │
│                                                         │
│  Next Steps:                                            │
│  1. Test cluster failover (Settings → Cluster → Test)   │
│  2. Deploy your first distributed service               │
│  3. Configure monitoring alerts                         │
│                                                         │
│  [Go to Cluster Manager]  [View Dashboard]              │
└─────────────────────────────────────────────────────────┘
```

---

## Easy Management nach Activation

### Dashboard zeigt Cluster-Status

```
┌─────────────────────────────────────────────┐
│  Dashboard                                  │
├─────────────────────────────────────────────┤
│  🌐 Cluster Status: ✅ Healthy              │
│  Nodes: 3/3 online                          │
│  VIP: 10.0.0.100                            │
│                                             │
│  Total Storage: 8 TB / 24 TB (replicated 3x)│
│  Docker Services: 5 (distributed)           │
│                                             │
│  [Cluster Manager] [Add Node] [Settings]    │
└─────────────────────────────────────────────┘
```

### Settings → Cluster

```
┌─────────────────────────────────────────────────────────┐
│  Settings → Cluster                                     │
├─────────────────────────────────────────────────────────┤
│  Cluster Mode: ✅ Enabled (Scale-Out, 3 Nodes)          │
│                                                         │
│  Features:                                              │
│  ┌──────────────────────────────────────────────────┐   │
│  │ ✅ Distributed Storage     [Configure] [Disable] │   │
│  │ ✅ Load Balancing          [Configure] [Disable] │   │
│  │ ✅ Container Orchestration [Configure] [Disable] │   │
│  │ ❌ Advanced Monitoring     [Enable]              │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  Quick Actions:                                         │
│  [Add Node]  [Test Failover]  [Cluster Health Check]    │
│                                                         │
│  Danger Zone:                                           │
│  [Convert to HA Mode]  [Disable Cluster (Standalone)]   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**Einfaches Toggle:**
- User kann Features einzeln aktivieren/deaktivieren
- Cluster selbst bleibt aktiv
- z.B. GlusterFS deaktivieren → Normale lokale Disks

### Cluster Manager App (nur wenn Cluster aktiviert)

```
┌─────────────────────────────────────────────────────────┐
│  Cluster Manager                                        │
├─────────────────────────────────────────────────────────┤
│  [Overview] [Nodes] [Storage] [Load Balancer] [Config]  │
│                                                         │
│  Overview Tab:                                          │
│  • Cluster Health: ✅ All systems operational           │
│  • Leader: node1                                        │
│  • Uptime: 15 days                                      │
│  • Last Failover: None                                  │
│                                                         │
│  [Manage Features] [Add Node] [Test Cluster]            │
└─────────────────────────────────────────────────────────┘
```

**Cluster Manager erscheint NUR wenn Cluster aktiv ist!**
- Kein Clutter für Single-Node User
- Sauber und fokussiert

---

## Migration Paths

### Path 1: Standalone → HA (2 Nodes)

```
1. User: "Enable Cluster"
2. Wizard: "Choose HA Mode"
3. Wizard: "Add 1 Node (nas-backup)"
4. Wizard: Auto-configure DRBD + Pacemaker + VIP
5. Done: 2-Node HA Cluster
```

**Was passiert:**
- DRBD mirroring zwischen Node 1 & 2
- Pacemaker/Corosync für Failover
- Keepalived VIP
- Keine GlusterFS (nicht nötig für 2 Nodes)

### Path 2: Standalone → Scale-Out (3+ Nodes)

```
1. User: "Enable Cluster"
2. Wizard: "Choose Scale-Out Mode"
3. Wizard: "Add 2+ Nodes"
4. Wizard: "Select Features (GlusterFS, HAProxy, Swarm)"
5. Wizard: "Configure Storage (Replicate/Distribute)"
6. Done: Scale-Out Cluster
```

### Path 3: HA → Scale-Out (Upgrade)

```
1. Current: 2-Node HA (DRBD)
2. User: "Settings → Cluster → Upgrade to Scale-Out"
3. Wizard: "Add 1+ Nodes (minimum 3 total)"
4. Wizard: "Migrate DRBD → GlusterFS?"
5. Done: Migrated to Scale-Out
```

**Migration ohne Datenverlust:**
- DRBD Daten → GlusterFS Volume kopieren
- DRBD deaktivieren
- GlusterFS replizieren

---

## Feature Toggles (nach Activation)

### Settings → Cluster → Features

```
┌─────────────────────────────────────────────────────────┐
│  Cluster Features                                       │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │ Distributed Storage (GlusterFS)                  │   │
│  │ Status: ✅ Enabled                               │   │
│  │                                                   │   │
│  │ Current Config:                                   │   │
│  │ • Volume 'data': Replicated 3x (8TB usable)      │   │
│  │ • Nodes: node1, nas-backup, nas-worker           │   │
│  │ • Self-Heal: Enabled                             │   │
│  │                                                   │   │
│  │ [Configure]  [Disable Feature]                   │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │ Load Balancing (HAProxy)                         │   │
│  │ Status: ✅ Enabled                               │   │
│  │                                                   │   │
│  │ Active Load Balancers:                           │   │
│  │ • Web UI (http://10.0.0.100:8080)                │   │
│  │ • API (http://10.0.0.100:8080/api)               │   │
│  │ • SMB (\\10.0.0.100)                             │   │
│  │                                                   │   │
│  │ [Configure]  [Disable Feature]                   │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │ Container Orchestration (Docker Swarm)           │   │
│  │ Status: ✅ Enabled                               │   │
│  │                                                   │   │
│  │ Swarm Status:                                    │   │
│  │ • Managers: 2 (node1, nas-backup)                │   │
│  │ • Workers: 1 (nas-worker)                        │   │
│  │ • Services: 5                                    │   │
│  │                                                   │   │
│  │ [Manage Swarm]  [Disable Feature]                │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │ Advanced Monitoring (Prometheus Federation)      │   │
│  │ Status: ❌ Disabled                              │   │
│  │                                                   │   │
│  │ Enable cluster-wide metrics aggregation?         │   │
│  │ • Requires: +2GB RAM per node                    │   │
│  │ • Provides: Centralized Grafana dashboards       │   │
│  │                                                   │   │
│  │ [Enable Feature]                                 │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**User kann jedes Feature individual togglen:**
- GlusterFS disable → Zurück zu lokalen Disks
- HAProxy disable → Direkter Node-Zugriff
- Swarm disable → Standalone Docker

---

## Backend: Cluster Detection

### Automatische Feature-Detection im Frontend

```typescript
// frontend/src/hooks/useCluster.ts
export function useCluster() {
  const [clusterEnabled, setClusterEnabled] = useState(false);
  const [clusterMode, setClusterMode] = useState<'standalone' | 'ha' | 'scaleout'>('standalone');
  const [features, setFeatures] = useState({
    glusterfs: false,
    haproxy: false,
    swarm: false,
    monitoring: false,
  });

  useEffect(() => {
    // GET /api/v1/cluster/status
    clusterApi.getStatus().then(resp => {
      if (resp.success && resp.data) {
        setClusterEnabled(resp.data.enabled);
        setClusterMode(resp.data.mode);
        setFeatures(resp.data.features);
      }
    });
  }, []);

  return { clusterEnabled, clusterMode, features };
}

// In App.tsx
const { clusterEnabled } = useCluster();

// Conditional Rendering
{clusterEnabled && (
  <App id="cluster-manager" name="Cluster Manager" icon="🌐" />
)}
```

**Cluster Manager App erscheint NUR wenn Cluster aktiviert!**

### Backend API

```go
// GET /api/v1/cluster/status
{
  "success": true,
  "data": {
    "enabled": true,
    "mode": "scaleout",
    "nodes": 3,
    "leader": "node1",
    "features": {
      "glusterfs": true,
      "haproxy": true,
      "swarm": true,
      "monitoring": false
    },
    "health": "healthy"
  }
}
```

---

## Vorteile dieses Ansatzes

### ✅ Für Single-Node User
- **Keine Complexity**: Sehen nichts von Cluster
- **Volle Performance**: Kein Overhead
- **Einfache Updates**: Keine Cluster-Abhängigkeiten

### ✅ Für Cluster User
- **Easy Activation**: One-Click Wizard
- **Gradual Migration**: Kann schrittweise skalieren
- **Optional Features**: Nur was sie brauchen
- **Easy Management**: Dedicated Cluster Manager App

### ✅ Für Stumpf.Works
- **Marketing**: "Skaliert von Homelab bis Enterprise"
- **Competitive Advantage**: Einfacher als TrueNAS SCALE
- **Flexibility**: Unterstützt alle Use Cases
- **Clean Code**: Features sind modular

---

## Implementation Priority

### Phase 0: Detection & Infrastructure (Woche 1)
- [ ] Cluster Status API (`/api/v1/cluster/status`)
- [ ] Frontend: `useCluster` Hook
- [ ] Conditional Rendering (Cluster Manager nur wenn enabled)

### Phase 1: Cluster Activation Wizard (Woche 2-3)
- [ ] Wizard UI (4 Screens)
- [ ] Node Discovery (Auto-Scan)
- [ ] Backend: Cluster Activation API
- [ ] SSH Key Exchange
- [ ] Auto-Installation von Packages

### Phase 2: Feature Toggles (Woche 4)
- [ ] Settings → Cluster → Features
- [ ] Enable/Disable GlusterFS
- [ ] Enable/Disable HAProxy
- [ ] Enable/Disable Swarm

### Phase 3: Cluster Features (see CLUSTER_INTEGRATION.md)
- Parallel zu Phase 6 (HA) implementieren
- GlusterFS, HAProxy, Swarm wie geplant

---

## Fazit

**"Cluster ist optional, aber wenn du es brauchst, ist es super einfach!"**

- ✅ Standarduser sehen nichts von Cluster-Complexity
- ✅ Enterprise-Kunden aktivieren Cluster in 5 Minuten
- ✅ Gradual Migration: Start small, scale big
- ✅ Feature Toggles: Nur was du brauchst

**Nächster Schritt:** Soll ich mit dem Activation Wizard UI prototypen? 🚀
