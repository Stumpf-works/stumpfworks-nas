# Stumpf.Works NAS - Testing Guide

> Complete guide for testing the application stack

## 🎉 Phase 4 Complete!

**What's been implemented:**
- ✅ Phase 1: Architecture & Documentation
- ✅ Phase 2: Backend Core Infrastructure
- ✅ Phase 3: Frontend Framework & UI System
- ✅ Phase 4: Core Applications

---

## 🚀 Quick Start (Testing)

### Prerequisites

Before testing, ensure you have:
- **Go 1.21+** installed
- **Node.js 18+** installed
- **Terminal access** (2 terminal windows recommended)

---

## Step 1: Start the Backend

Open **Terminal 1**:

```bash
# Navigate to backend directory
cd /home/user/stumpfworks-nas/backend

# Download Go dependencies (first time only)
go mod download

# Start the backend server
go run cmd/stumpfworks-server/main.go
```

**Expected Output:**
```
Stumpf.Works NAS v0.1.0-alpha
Starting server...
2025-11-11T20:00:00Z	INFO	Configuration loaded	{"environment": "development", "version": "0.1.0-alpha"}
2025-11-11T20:00:00Z	INFO	Database connected successfully	{"driver": "sqlite", "path": "./data/stumpfworks.db"}
2025-11-11T20:00:00Z	INFO	Database migrations completed successfully
2025-11-11T20:00:00Z	INFO	Default admin user created	{"username": "admin", "password": "admin (PLEASE CHANGE THIS!)"}
2025-11-11T20:00:00Z	INFO	HTTP server starting	{"address": "0.0.0.0:8080", "environment": "development"}
2025-11-11T20:00:00Z	INFO	Server started successfully	{"address": "0.0.0.0:8080", "health": "http://0.0.0.0:8080/health", "api": "http://0.0.0.0:8080/api/v1"}
```

✅ **Backend is running on:** `http://localhost:8080`

---

## Step 2: Test Backend API

Keep Terminal 1 running, open **Terminal 2**:

```bash
# Test health endpoint
curl http://localhost:8080/health

# Expected response:
# {"success":true,"data":{"status":"ok","service":"Stumpf.Works NAS","version":"0.1.0-alpha"}}
```

```bash
# Test login (default admin credentials)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Expected response (you'll get a JWT token):
# {"success":true,"data":{"accessToken":"eyJhbGc...","refreshToken":"eyJhbGc...","user":{...}}}
```

---

## Step 3: Start the Frontend

In **Terminal 2** (or a new Terminal 3):

```bash
# Navigate to frontend directory
cd /home/user/stumpfworks-nas/frontend

# Install dependencies (first time only)
npm install

# Start the dev server
npm run dev
```

**Expected Output:**
```
  VITE v5.0.0  ready in 500 ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
```

✅ **Frontend is running on:** `http://localhost:3000`

---

## Step 4: Open in Browser

Open your web browser and navigate to:

```
http://localhost:3000
```

---

## 🧪 Testing Checklist

### 1. Login Screen

**What to test:**
- ✅ Login form appears with glassmorphism effect
- ✅ Username field is focused automatically
- ✅ Default credentials are displayed

**Test Login:**
- Username: `admin`
- Password: `admin`
- Click "Sign In"

**Expected:**
- ✅ Loading spinner appears
- ✅ Redirects to Desktop after successful login
- ✅ Token is stored in localStorage

---

### 2. Desktop Environment

**What to test:**
- ✅ Gradient wallpaper background visible
- ✅ Top bar appears at top (with logo, metrics, time)
- ✅ Dock appears at bottom with 6 app icons
- ✅ Dock icons have magnification effect on hover

**Top Bar Elements:**
- ✅ "Stumpf.Works NAS" title visible
- ✅ CPU and Memory metrics updating
- ✅ Theme toggle button (sun/moon icon)
- ✅ User avatar (first letter of username)
- ✅ Time and date updating every second

**Dock Apps:**
- 📊 Dashboard
- 💾 Storage
- 📁 Files
- 👥 Users
- 🌐 Network
- ⚙️ Settings

---

### 3. Dashboard App

**How to test:**
- Click the **📊 Dashboard** icon in the Dock

**Expected:**
- ✅ Window opens with smooth animation
- ✅ Traffic lights visible (🔴 🟡 🟢) on title bar hover
- ✅ Window is draggable by title bar
- ✅ Dashboard content loads

**Dashboard Features:**
- ✅ Real-time CPU usage (percentage + progress bar)
- ✅ Real-time Memory usage (used/total)
- ✅ Disk usage per partition
- ✅ Network statistics (bytes sent/received)
- ✅ Data updates every 3 seconds

**Window Controls:**
- ✅ 🔴 Red: Closes window
- ✅ 🟡 Yellow: Minimizes window to Dock
- ✅ 🟢 Green: Maximizes/restores window
- ✅ Drag title bar to move window
- ✅ Click anywhere on window to bring to front

---

### 4. User Manager App

**How to test:**
- Click the **👥 Users** icon in the Dock

**Expected:**
- ✅ Window opens with User Manager
- ✅ Default admin user is visible in grid
- ✅ User cards show: avatar, username, email, role, status

**Test Create User:**
1. Click "**+ Create User**" button
2. Fill in the form:
   - Username: `testuser`
   - Email: `test@example.com`
   - Password: `password123`
   - Full Name: `Test User`
   - Role: Select **User**
3. Click "**Create**"

**Expected:**
- ✅ User is created
- ✅ Modal closes
- ✅ New user appears in grid
- ✅ User card shows: testuser, test@example.com, "User" role badge

**Test Edit User:**
1. Click "**Edit**" on testuser card
2. Change email to: `newemail@example.com`
3. Click "**Update**"

**Expected:**
- ✅ User is updated
- ✅ Email changes to newemail@example.com

**Test Delete User:**
1. Click "**Delete**" on testuser card
2. Confirm deletion in popup

**Expected:**
- ✅ Confirmation dialog appears
- ✅ User is deleted after confirmation
- ✅ User card disappears from grid

---

### 5. Settings App

**How to test:**
- Click the **⚙️ Settings** icon in the Dock

**Expected:**
- ✅ Window opens with Settings
- ✅ User Information section shows current user
- ✅ Appearance section has dark mode toggle
- ✅ System Information section shows OS details
- ✅ About section shows app name and version
- ✅ Logout button visible

**Test Dark Mode:**
1. Click the **toggle switch** in Appearance section
2. Watch the UI change

**Expected:**
- ✅ Background changes to dark
- ✅ Text colors invert
- ✅ Cards adapt to dark theme
- ✅ Dock adapts to dark theme
- ✅ Theme persists on page reload

**Test System Information:**
- ✅ Hostname displayed
- ✅ Platform displayed (linux, darwin, etc.)
- ✅ OS displayed (debian, ubuntu, etc.)
- ✅ Architecture displayed (amd64, arm64, etc.)
- ✅ CPU Cores count displayed
- ✅ Uptime formatted (e.g., "0d 2h 15m")

---

### 6. Multi-Window Management

**How to test:**
1. Open **Dashboard** (📊)
2. Open **Users** (👥)
3. Open **Settings** (⚙️)

**Expected:**
- ✅ All 3 windows are open simultaneously
- ✅ Each window is independently draggable
- ✅ Clicking a window brings it to front
- ✅ Running indicators (dots) appear below icons in Dock
- ✅ Windows can overlap
- ✅ Z-index stacking works correctly

**Test Window States:**
1. Minimize Dashboard (🟡 yellow button)
2. Maximize Users (🟢 green button)
3. Close Settings (🔴 red button)

**Expected:**
- ✅ Dashboard disappears but icon shows running in Dock
- ✅ Users window fills the screen (except Dock/TopBar)
- ✅ Settings window closes and icon shows not running
- ✅ Click minimized Dashboard icon to restore

---

### 7. Placeholder Apps

**How to test:**
- Click **💾 Storage**, **📁 Files**, or **🌐 Network**

**Expected:**
- ✅ Window opens with "Coming Soon" placeholder
- ✅ Shows 🚧 construction icon
- ✅ App name displayed
- ✅ Window controls still work

---

### 8. Logout

**How to test:**
1. Open **⚙️ Settings**
2. Scroll to bottom
3. Click "**Logout**" button

**Expected:**
- ✅ API logout call is made
- ✅ Tokens removed from localStorage
- ✅ Redirects to Login screen
- ✅ Must login again to access Desktop

---

## 🐛 Troubleshooting

### Backend won't start

**Error:** `port 8080 already in use`
```bash
# Kill the process using port 8080
lsof -ti:8080 | xargs kill -9

# Or use a different port
cd backend
PORT=8081 go run cmd/stumpfworks-server/main.go
```

### Frontend won't start

**Error:** `EADDRINUSE: port 3000 already in use`
```bash
# Kill the process using port 3000
lsof -ti:3000 | xargs kill -9

# Or use a different port
PORT=3001 npm run dev
```

### Can't login

**Issue:** Invalid credentials error

**Solution:**
- Default credentials: `admin` / `admin`
- Check backend logs for "Default admin user created"
- Database might be missing, delete `backend/data/` and restart backend

### CORS errors in browser console

**Issue:** CORS policy blocking requests

**Solution:**
- Ensure Vite proxy is configured (should be automatic)
- Check `frontend/vite.config.ts` has proxy settings
- Restart frontend dev server

### Metrics not updating

**Issue:** Top bar or Dashboard shows 0% or no data

**Solution:**
- Check backend is running
- Open browser DevTools → Network tab
- Check for 401 errors (token expired, logout and login again)
- Check for 500 errors (backend issue)

### Dark mode not working

**Issue:** Toggle doesn't change theme

**Solution:**
- Clear browser localStorage
- Refresh page (Ctrl+R or Cmd+R)
- Check browser console for errors

---

## 📊 Performance Expectations

**Backend:**
- Startup time: < 2 seconds
- API response time: < 100ms
- Memory usage: ~ 30-50 MB
- CPU usage: < 1% idle, < 5% active

**Frontend:**
- Initial load: < 1 second (dev mode)
- API calls: < 100ms
- Window animations: 60 FPS
- Memory usage: ~ 100-150 MB

---

## 🎯 Feature Summary

### Implemented ✅

**Backend:**
- REST API (Chi router)
- JWT Authentication
- User CRUD (admin only)
- System Metrics API
- WebSocket server (basic)
- SQLite database
- Auto-migrations
- Default admin user

**Frontend:**
- macOS-like Desktop
- Animated Dock (magnification)
- Top Menu Bar (metrics, time, theme)
- Window Management (drag, minimize, maximize, close)
- Login screen
- Dashboard (real-time metrics)
- User Manager (full CRUD)
- Settings (user info, theme, system info, logout)
- Dark mode
- Responsive design

### Coming Soon 🚧

- Storage Manager (disks, volumes, SMART)
- File Station (file browser, upload/download)
- Network Manager (interfaces, firewall)
- Plugin System
- Launchpad (app grid)
- Control Center
- Notification Center
- WebSocket real-time updates

---

## 🎉 Success Criteria

If all the following work, the application is functioning correctly:

- ✅ Backend starts without errors
- ✅ Frontend starts without errors
- ✅ Login works with admin/admin
- ✅ Desktop appears with wallpaper, Dock, TopBar
- ✅ Dashboard opens and shows real-time metrics
- ✅ User Manager can create, edit, and delete users
- ✅ Settings shows system info and allows logout
- ✅ Multiple windows can be open simultaneously
- ✅ Windows are draggable and manageable
- ✅ Dark mode toggle works
- ✅ Logout returns to login screen

---

## 📝 Test Results Log

Use this section to log your test results:

```
Date: ___________
Tester: ___________

Backend Startup: [ ] Pass [ ] Fail
Frontend Startup: [ ] Pass [ ] Fail
Login: [ ] Pass [ ] Fail
Desktop Load: [ ] Pass [ ] Fail
Dashboard: [ ] Pass [ ] Fail
User Manager: [ ] Pass [ ] Fail
Settings: [ ] Pass [ ] Fail
Multi-Window: [ ] Pass [ ] Fail
Dark Mode: [ ] Pass [ ] Fail
Logout: [ ] Pass [ ] Fail

Notes:
_________________________________
_________________________________
```

---

**Happy Testing! 🚀**

If you find any bugs, please note them for future fixes!
