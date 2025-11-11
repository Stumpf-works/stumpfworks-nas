# UI/UX Design Specification

> macOS-inspired interface design for Stumpf.Works NAS Solution

---

## Table of Contents

- [Design Philosophy](#design-philosophy)
- [Design System](#design-system)
- [Layout Components](#layout-components)
- [Window System](#window-system)
- [Applications](#applications)
- [Interactions & Animations](#interactions--animations)
- [Responsive Design](#responsive-design)
- [Accessibility](#accessibility)

---

## Design Philosophy

### Core Principles

1. **Familiarity**: Users who love macOS should feel at home
2. **Clarity**: Information hierarchy is obvious
3. **Depth**: Visual layers create spatial understanding
4. **Fluidity**: Smooth, natural animations
5. **Efficiency**: Fast access to common tasks
6. **Beauty**: Aesthetic excellence drives engagement

### Visual Language

Inspired by macOS Big Sur, Monterey, and Sonoma:
- **Glassmorphism**: Translucent surfaces with blur
- **Vibrancy**: Colors that adapt to content behind them
- **Shadows**: Realistic depth cues
- **Rounded Corners**: Soft, friendly shapes
- **Subtle Gradients**: Gentle color transitions
- **Iconography**: SF Symbols-inspired icons (but custom/open-source)

---

## Design System

### Color Palette

#### Light Mode

```css
/* Primary */
--primary-blue: #007AFF;
--primary-green: #34C759;
--primary-red: #FF3B30;
--primary-orange: #FF9500;
--primary-purple: #AF52DE;

/* Neutrals */
--gray-50: #F9FAFB;
--gray-100: #F3F4F6;
--gray-200: #E5E7EB;
--gray-300: #D1D5DB;
--gray-400: #9CA3AF;
--gray-500: #6B7280;
--gray-600: #4B5563;
--gray-700: #374151;
--gray-800: #1F2937;
--gray-900: #111827;

/* Glassmorphism */
--glass-light: rgba(255, 255, 255, 0.8);
--glass-blur: 40px;
```

#### Dark Mode

```css
/* Primary (same as light, but adjusted opacity) */
--primary-blue-dark: #0A84FF;

/* Neutrals */
--gray-dark-50: #1E1E1E;
--gray-dark-100: #2C2C2C;
--gray-dark-200: #3A3A3A;
--gray-dark-300: #4A4A4A;
--gray-dark-400: #6E6E6E;
--gray-dark-500: #8E8E93;

/* Glassmorphism */
--glass-dark: rgba(30, 30, 30, 0.7);
--glass-blur: 40px;
```

---

### Typography

**Font Family:**
- Primary: `system-ui, -apple-system, "Segoe UI", sans-serif`
- Monospace: `"SF Mono", Monaco, "Cascadia Code", monospace`

**Scale:**
```css
--text-xs: 0.75rem;    /* 12px */
--text-sm: 0.875rem;   /* 14px */
--text-base: 1rem;     /* 16px */
--text-lg: 1.125rem;   /* 18px */
--text-xl: 1.25rem;    /* 20px */
--text-2xl: 1.5rem;    /* 24px */
--text-3xl: 1.875rem;  /* 30px */
--text-4xl: 2.25rem;   /* 36px */
```

**Weights:**
- Regular: 400
- Medium: 500
- Semibold: 600
- Bold: 700

---

### Spacing

Based on 4px grid:

```css
--spacing-1: 0.25rem;  /* 4px */
--spacing-2: 0.5rem;   /* 8px */
--spacing-3: 0.75rem;  /* 12px */
--spacing-4: 1rem;     /* 16px */
--spacing-6: 1.5rem;   /* 24px */
--spacing-8: 2rem;     /* 32px */
--spacing-12: 3rem;    /* 48px */
--spacing-16: 4rem;    /* 64px */
```

---

### Shadows

```css
/* macOS-style shadows */
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
--shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
--shadow-lg: 0 10px 25px rgba(0, 0, 0, 0.15);
--shadow-xl: 0 20px 40px rgba(0, 0, 0, 0.2);

/* Window shadow */
--shadow-window: 0 20px 60px rgba(0, 0, 0, 0.25);
```

---

### Border Radius

```css
--radius-sm: 0.375rem;  /* 6px */
--radius-md: 0.5rem;    /* 8px */
--radius-lg: 0.75rem;   /* 12px */
--radius-xl: 1rem;      /* 16px */
--radius-2xl: 1.5rem;   /* 24px */
--radius-full: 9999px;  /* Circular */
```

---

## Layout Components

### 1. Desktop Environment

The root container for the entire UI.

**Structure:**
```
<Desktop>
  <Wallpaper />
  <WindowManager>
    {open windows}
  </WindowManager>
  <TopBar />
  <Dock />
  <Launchpad />
  <ControlCenter />
  <NotificationCenter />
</Desktop>
```

**Features:**
- Dynamic wallpaper (changes with time of day)
- Right-click context menu (desktop actions)
- Drag-and-drop support (files to desktop)

---

### 2. Top Bar (Menu Bar)

Fixed bar at the top of the screen.

**Left Section:**
- **Apple Logo** → System menu (About, Restart, Shutdown)
- **App Name** → Current focused app name

**Right Section:**
- **System Indicators**: CPU, RAM, disk, network (mini graphs)
- **Clock**: Time and date
- **User Avatar** → Quick user menu
- **Control Center Icon**
- **Notification Icon** (with badge count)

**Height:** 32px

**Style:**
- Glassmorphism background
- Blur effect
- Subtle bottom border
- Dark text (light mode), light text (dark mode)

**Behavior:**
- Always visible (fixed position)
- Auto-hide on fullscreen windows (optional)

---

### 3. Dock

Animated app launcher at the bottom of the screen.

**Features:**
- **App Icons**: Pinned and running apps
- **Hover Magnification**: Icons scale up on hover (macOS style)
- **Running Indicators**: Small dot below running apps
- **Badge Notifications**: Red badge with count on app icon
- **Bounce Animation**: Apps can request attention
- **Drag to Reorder**: Rearrange icons
- **Trash/Recent Files**: Right side (optional)

**Dimensions:**
- Height: 60px (default), up to 80px (magnified)
- Width: Auto (based on number of apps)
- Padding: 8px

**Style:**
- Glassmorphism background
- Rounded corners (--radius-2xl)
- Subtle shadow (--shadow-lg)
- Centered horizontally

**Animations:**
```tsx
// Dock Icon Magnification
const dockVariants = {
  hover: { scale: 1.5, y: -10 },
  tap: { scale: 0.95 },
};

<motion.div
  variants={dockVariants}
  whileHover="hover"
  whileTap="tap"
>
  <AppIcon />
</motion.div>
```

---

### 4. Launchpad

Full-screen app grid overlay.

**Trigger:**
- Click Launchpad icon in Dock
- Keyboard shortcut (F4)
- Pinch gesture (future)

**Features:**
- **Grid Layout**: 6 columns x 4 rows (desktop), responsive on smaller screens
- **Pagination**: Dots at bottom if more apps than fit on one screen
- **Search Bar**: Top center, filters apps as you type
- **Folders**: Group related apps
- **Animations**: Apps zoom in from Dock position

**Style:**
- Blurred background (desktop fades out)
- App icons with labels below
- Smooth entrance/exit animations

**Animation Example:**
```tsx
<AnimatePresence>
  {isOpen && (
    <motion.div
      initial={{ opacity: 0, scale: 0.8 }}
      animate={{ opacity: 1, scale: 1 }}
      exit={{ opacity: 0, scale: 0.8 }}
      transition={{ duration: 0.3 }}
    >
      <AppGrid />
    </motion.div>
  )}
</AnimatePresence>
```

---

### 5. Control Center

Quick settings panel (slides in from top-right).

**Sections:**
- **Network**: Wi-Fi, Ethernet status, toggle
- **Volume**: Slider
- **Brightness**: Slider (if applicable)
- **Night Mode**: Toggle dark mode
- **Do Not Disturb**: Toggle notifications
- **System Resources**: CPU, RAM gauges

**Style:**
- Glassmorphism card
- Rounded corners
- Compact, organized layout
- Toggle switches (iOS style)

**Behavior:**
- Slide in from top-right corner
- Click outside to close
- Quick actions (no navigation away)

---

### 6. Notification Center

Notification and widget panel (slides in from top-right).

**Sections:**
- **Today**: Date, weather widget (optional)
- **Notifications**: List of recent notifications
  - Grouped by app
  - Expandable (show more info)
  - Clear all button
- **Widgets**: System status, calendar, etc.

**Style:**
- Similar to Control Center
- Scrollable list
- Swipe to dismiss (individual notifications)

---

## Window System

### Window Component

Every app runs in a window (like macOS).

**Structure:**
```tsx
<Window
  title="App Name"
  icon={<Icon />}
  initialPosition={{ x: 100, y: 100 }}
  initialSize={{ width: 800, height: 600 }}
  minSize={{ width: 400, height: 300 }}
  onClose={() => {}}
  onMinimize={() => {}}
  onMaximize={() => {}}
>
  <AppContent />
</Window>
```

**Features:**
- **Title Bar**: App icon, title, traffic lights (close, minimize, maximize)
- **Draggable**: Click and drag title bar to move
- **Resizable**: Drag edges/corners to resize
- **Min/Max**: Respect min/max dimensions
- **Focus Management**: Bring to front on click
- **Multiple Windows**: Same app can have multiple windows

---

### Window Traffic Lights

Iconic macOS window controls (top-left).

**Buttons:**
1. **Red (Close)**: Close window
2. **Yellow (Minimize)**: Minimize to Dock
3. **Green (Maximize)**: Toggle fullscreen/restore

**Behavior:**
- Hidden by default (macOS style)
- Visible on hover over title bar
- Hover on each button shows a symbol (×, −, ⤢)

**Style:**
```css
.traffic-light {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 8px;
}

.traffic-light.close { background: #FF5F57; }
.traffic-light.minimize { background: #FFBD2E; }
.traffic-light.maximize { background: #28CA42; }
```

---

### Window States

1. **Normal**: Regular window on desktop
2. **Minimized**: Hidden, icon animates to Dock
3. **Maximized**: Fills entire screen (except Dock/TopBar)
4. **Fullscreen**: Fills entire screen (hides Dock/TopBar)
5. **Focused**: Current active window (higher z-index, shadow)
6. **Background**: Inactive windows (slightly dimmed)

---

### Window Manager

Manages all open windows.

**Responsibilities:**
- Track open windows (state)
- Handle focus (bring to front)
- Handle minimize/maximize/close
- Persist window positions (localStorage)
- Animate window transitions
- Prevent windows from going off-screen

**State Example:**
```ts
interface Window {
  id: string;
  appId: string;
  title: string;
  position: { x: number; y: number };
  size: { width: number; height: number };
  state: 'normal' | 'minimized' | 'maximized' | 'fullscreen';
  zIndex: number;
  isFocused: boolean;
}
```

---

## Applications

### Core Apps

#### 1. Dashboard

**Purpose:** System overview and quick actions

**Layout:**
- **Header**: "Dashboard", date, time
- **Widgets Grid**:
  - System Resources (CPU, RAM, Disk)
  - Network Activity (upload/download graphs)
  - Recent Activity (logs)
  - Quick Actions (restart, shutdown, update)
- **Charts**: Line charts (Recharts or Chart.js)

**Style:**
- Card-based layout
- Real-time updates (WebSocket)
- Responsive grid (1-3 columns)

---

#### 2. Storage Manager

**Purpose:** Manage disks, volumes, snapshots

**Layout:**
- **Sidebar**: List of disks/volumes
- **Main Panel**: Details, SMART data, usage graphs
- **Toolbar**: Create, Delete, Snapshot buttons

**Features:**
- Visual disk layout (pie charts)
- Create LVM/mdadm volumes (wizard)
- SMART monitoring (alerts)
- Snapshot manager (list, create, restore, delete)

---

#### 3. File Station

**Purpose:** Web-based file manager

**Layout:**
- **Sidebar**: Folder tree
- **Main Panel**: File grid/list view
- **Toolbar**: Upload, Download, New Folder, Delete
- **Breadcrumb**: Current path

**Features:**
- Drag-and-drop upload
- File previews (images, videos, PDFs)
- Context menu (right-click)
- Permissions editor
- Sharing (generate public links)

---

#### 4. User Manager

**Purpose:** Manage users and groups

**Layout:**
- **List View**: Users with avatar, name, role
- **Detail Panel**: Edit user (name, email, password, groups, quotas)
- **Toolbar**: Add User, Delete User

**Features:**
- Role-based permissions
- Quota management
- Activity logs per user

---

#### 5. Network Manager

**Purpose:** Configure network interfaces

**Layout:**
- **Interface List**: Ethernet, Wi-Fi, VLANs
- **Detail Panel**: IP config, gateway, DNS
- **Firewall Tab**: Rules list

**Features:**
- Static/DHCP configuration
- VLAN creation
- Bonding/teaming
- Firewall rules (ufw integration)

---

#### 6. Share Manager

**Purpose:** Configure file shares

**Layout:**
- **Share List**: SMB, NFS, FTP shares
- **Detail Panel**: Path, permissions, users
- **Toolbar**: Add Share, Delete Share

**Features:**
- SMB/CIFS configuration
- NFS exports
- FTP/SFTP users
- WebDAV support

---

#### 7. Plugin Center

**Purpose:** Browse and install plugins

**Layout:**
- **Grid/List View**: Available plugins (App Store style)
- **Detail Panel**: Plugin info, screenshots, reviews
- **Installed Tab**: Manage installed plugins

**Features:**
- Search and filter
- Install/uninstall
- Enable/disable
- Update plugins
- Plugin settings

---

### App Icon Design

**Guidelines:**
- Rounded square shape (iOS/macOS style)
- Gradient backgrounds
- Symbolic icons (simple, recognizable)
- Consistent size (512x512 base, scaled to 64x64 for Dock)
- PNG or SVG format

**Example Color Schemes:**
- Dashboard: Blue gradient (#007AFF → #5AC8FA)
- Storage: Orange gradient (#FF9500 → #FF6B00)
- File Station: Purple gradient (#AF52DE → #BF5AF2)
- Network: Green gradient (#34C759 → #30D158)

---

## Interactions & Animations

### Micro-Interactions

#### 1. Button Hover
```tsx
<motion.button
  whileHover={{ scale: 1.05 }}
  whileTap={{ scale: 0.95 }}
>
  Click Me
</motion.button>
```

#### 2. Card Hover
```tsx
<motion.div
  whileHover={{ y: -4, boxShadow: 'var(--shadow-lg)' }}
  transition={{ duration: 0.2 }}
>
  <Card />
</motion.div>
```

#### 3. Dock Icon Magnification
- Scale: 1 → 1.5
- Y-axis: 0 → -10px
- Smooth easing

#### 4. Window Open/Close
- Open: Scale from 0.5 → 1, Opacity 0 → 1
- Close: Scale 1 → 0.5, Opacity 1 → 0
- Duration: 300ms

---

### Page Transitions

```tsx
<motion.div
  initial={{ opacity: 0, x: 20 }}
  animate={{ opacity: 1, x: 0 }}
  exit={{ opacity: 0, x: -20 }}
  transition={{ duration: 0.3 }}
>
  <Page />
</motion.div>
```

---

### Loading States

#### Skeleton Loaders
Placeholder while content loads (gray boxes that pulse).

```tsx
<motion.div
  className="skeleton"
  animate={{ opacity: [0.5, 1, 0.5] }}
  transition={{ repeat: Infinity, duration: 1.5 }}
/>
```

#### Spinner
macOS-style spinner (circular, rotating).

---

## Responsive Design

### Breakpoints

```css
/* Tailwind default breakpoints */
sm: 640px   /* Small devices (tablets) */
md: 768px   /* Medium devices (landscape tablets) */
lg: 1024px  /* Large devices (laptops) */
xl: 1280px  /* Extra large devices (desktops) */
2xl: 1536px /* 2X large devices (large desktops) */
```

### Responsive Behavior

#### Desktop (1024px+)
- Full Desktop layout
- Dock at bottom
- Multiple windows
- Sidebar navigation in apps

#### Tablet (768px - 1023px)
- Single window at a time
- Dock becomes smaller
- Sidebar collapses to hamburger menu

#### Mobile (< 768px)
- No Desktop environment
- Native mobile app layout
- Bottom tab navigation
- Swipe gestures
- No windows (fullscreen views)

---

## Accessibility

### WCAG 2.1 AA Compliance

✅ **Color Contrast**
- Text: 4.5:1 minimum
- Large text: 3:1 minimum
- UI components: 3:1 minimum

✅ **Keyboard Navigation**
- All interactive elements focusable
- Tab order logical
- Focus indicators visible
- Keyboard shortcuts (with ⌘/Ctrl modifiers)

✅ **Screen Readers**
- Semantic HTML
- ARIA labels where needed
- Alt text for images
- Descriptive button labels

✅ **Motion**
- Respect `prefers-reduced-motion`
- Option to disable animations
- No auto-playing videos

---

### Keyboard Shortcuts

| Action | Shortcut |
|--------|----------|
| Open Launchpad | `F4` |
| Open Dashboard | `⌘ + D` |
| Close Window | `⌘ + W` |
| Minimize Window | `⌘ + M` |
| Maximize Window | `⌘ + F` |
| Toggle Dark Mode | `⌘ + Shift + D` |
| Open Settings | `⌘ + ,` |
| Search | `⌘ + K` |

---

## Component Library

### Reusable Components

#### Button
```tsx
<Button variant="primary|secondary|danger" size="sm|md|lg">
  Click Me
</Button>
```

#### Input
```tsx
<Input
  type="text|password|email"
  label="Username"
  placeholder="Enter username"
/>
```

#### Card
```tsx
<Card>
  <CardHeader>Title</CardHeader>
  <CardBody>Content</CardBody>
  <CardFooter>Actions</CardFooter>
</Card>
```

#### Modal
```tsx
<Modal isOpen={isOpen} onClose={onClose}>
  <ModalHeader>Title</ModalHeader>
  <ModalBody>Content</ModalBody>
  <ModalFooter>
    <Button onClick={onClose}>Cancel</Button>
    <Button variant="primary">Confirm</Button>
  </ModalFooter>
</Modal>
```

#### Table
```tsx
<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Name</TableHead>
      <TableHead>Status</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    <TableRow>
      <TableCell>Item 1</TableCell>
      <TableCell>Active</TableCell>
    </TableRow>
  </TableBody>
</Table>
```

---

## Design Deliverables Checklist

- [ ] Color palette defined (light + dark mode)
- [ ] Typography scale and fonts
- [ ] Spacing system
- [ ] Component library (Storybook)
- [ ] Icon set (custom or Heroicons/Lucide)
- [ ] Animation guidelines
- [ ] Desktop layout mockups
- [ ] App-specific mockups
- [ ] Responsive layouts (mobile, tablet)
- [ ] Accessibility audit
- [ ] User testing (usability)

---

## Inspiration & References

- macOS Big Sur / Monterey UI
- [Apple Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/)
- Glassmorphism examples: [neumorphism.io](https://neumorphism.io/)
- Animation inspiration: [Dribbble](https://dribbble.com/)
- Component libraries: [Radix UI](https://www.radix-ui.com/), [Headless UI](https://headlessui.com/)

---

**Document Version:** 1.0.0
**Last Updated:** 2025-11-11
**Figma Mockups:** Coming in Phase 3
