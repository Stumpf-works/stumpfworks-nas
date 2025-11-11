# Stumpf.Works NAS - Frontend

macOS-inspired web interface for Stumpf.Works NAS Solution.

## Features

✅ **Phase 3 Complete** - Frontend Framework & UI System

- macOS-like Desktop Environment
- Animated Dock with magnification effect
- Top Menu Bar with real-time system metrics
- Window Manager (draggable, resizable windows)
- Traffic Lights (close, minimize, maximize)
- Login screen with authentication
- Dashboard app with real-time metrics
- Dark mode support
- Responsive design
- Glassmorphism UI effects
- Framer Motion animations

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool & dev server
- **TailwindCSS 3** - Utility-first CSS framework
- **Framer Motion** - Animation library
- **Zustand** - State management
- **Axios** - HTTP client

## Project Structure

```
frontend/src/
├── api/                    # API client and services
│   ├── client.ts          # Axios instance with interceptors
│   ├── auth.ts            # Authentication API
│   └── system.ts          # System metrics API
├── apps/                  # Applications
│   └── Dashboard/         # Dashboard app
├── components/            # Shared components
│   ├── ui/               # UI primitives
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   └── Card.tsx
│   └── Window.tsx        # Window component
├── layout/                # Core layout components
│   ├── Desktop.tsx       # Desktop environment
│   ├── TopBar.tsx        # Top menu bar
│   ├── Dock.tsx          # Bottom dock
│   └── WindowManager.tsx # Window management
├── pages/                 # Pages
│   └── Login.tsx         # Login page
├── store/                 # Zustand stores
│   └── index.ts          # All stores
├── styles/                # Global styles
│   └── index.css         # Tailwind imports
├── types/                 # TypeScript types
│   └── index.ts
├── App.tsx               # Root component
└── main.tsx              # Entry point
```

## Getting Started

### Prerequisites

- Node.js 18+
- npm or pnpm

### Installation

```bash
cd frontend

# Install dependencies
npm install

# Or with pnpm
pnpm install
```

### Development

```bash
# Start dev server
npm run dev

# Or use Make (from project root)
make dev-frontend
```

The frontend will start on `http://localhost:3000` and proxy API requests to `http://localhost:8080`.

### Building

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

## Configuration

Environment variables can be set in `.env`:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## Features

### Desktop Environment

- macOS-like desktop with wallpaper
- Gradient background
- Window management
- Dock navigation

### Top Bar

- System name and logo
- Real-time metrics (CPU, Memory)
- Theme toggle (dark/light)
- User avatar
- Time and date

### Dock

- Animated app icons
- Hover magnification effect
- Running indicators
- Tooltips
- Smooth spring animations

### Windows

- Draggable title bar
- Resizable (coming soon)
- Traffic lights (close, minimize, maximize)
- Focus management
- Z-index stacking
- Smooth animations

### Login

- Clean authentication UI
- Form validation
- Error handling
- Default credentials displayed

### Dashboard

- Real-time system metrics
- CPU usage with per-core stats
- Memory usage
- Disk usage (per partition)
- Network statistics
- Auto-refresh every 3 seconds

## State Management

Uses Zustand with multiple stores:

- **AuthStore** - User authentication state
- **SystemStore** - System metrics
- **WindowStore** - Window management
- **ThemeStore** - Dark/light mode

## Styling

### TailwindCSS Configuration

- Custom macOS color palette
- Glassmorphism utilities
- Custom shadows
- Dark mode support
- Custom animations

### Design Tokens

```js
colors: {
  'macos-blue': '#007AFF',
  'macos-green': '#34C759',
  'macos-red': '#FF3B30',
  // ... more colors
}
```

## API Integration

### Authentication

```ts
import { authApi } from '@/api/auth';

// Login
const response = await authApi.login({ username, password });

// Get current user
const user = await authApi.getCurrentUser();
```

### System Metrics

```ts
import { systemApi } from '@/api/system';

// Get real-time metrics
const metrics = await systemApi.getMetrics();
```

## Default Credentials

**⚠️ Change in production!**

- **Username:** admin
- **Password:** admin

## Development Tips

### Path Aliases

TypeScript path aliases are configured for cleaner imports:

```ts
import Button from '@/components/ui/Button';
import { useAuthStore } from '@/store';
import { authApi } from '@/api/auth';
```

### Adding New Apps

1. Create app component in `src/apps/YourApp/`
2. Register in `src/layout/Dock.tsx`:

```ts
const registeredApps = [
  // ...
  {
    id: 'yourapp',
    name: 'Your App',
    icon: '🎨',
    component: YourApp,
    defaultSize: { width: 800, height: 600 },
  },
];
```

### Creating Windows

Windows are managed by Zustand store:

```ts
import { useWindowStore } from '@/store';

const openWindow = useWindowStore((state) => state.openWindow);

// Open a new window
openWindow(appConfig);
```

## Linting & Formatting

```bash
# Lint code
npm run lint

# Format code
npm run format
```

## Testing

```bash
# Run tests (when added)
npm run test
```

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## Next Steps (Phase 4)

- More applications (Storage Manager, File Station, etc.)
- Launchpad (app grid)
- Control Center (quick settings)
- Notification Center
- Context menus
- Keyboard shortcuts
- WebSocket integration for real-time updates
- Mobile responsive improvements

## License

MIT License - see [LICENSE](../LICENSE) for details.
