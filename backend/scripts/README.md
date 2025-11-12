# Storage Testing Scripts

Scripts for testing storage management functionality with loop devices.

## Quick Start

### 1. Create Test Disks

```bash
sudo ./create-test-disks.sh 3 1024
```

This creates 3 loop devices of 1024MB each.

### 2. Use in Stumpf.Works NAS

1. Start the backend server
2. Open the web UI
3. Navigate to Storage Manager
4. Create volumes using the loop devices (`/dev/loop1`, `/dev/loop2`, etc.)

### 3. Clean Up

```bash
sudo ./cleanup-test-disks.sh
```

## Prerequisites

- Linux system with loop device support
- Root/sudo access
- `losetup` utility (usually pre-installed)

## Testing Different Filesystems

The storage manager supports:
- **ext4** - Requires `e2fsprogs` (usually pre-installed)
- **xfs** - Requires `apt-get install xfsprogs`
- **btrfs** - Requires `apt-get install btrfs-progs`

## Advanced Usage

### Create More Disks

```bash
# Create 5 disks of 2GB each
sudo ./create-test-disks.sh 5 2048
```

### List Current Loop Devices

```bash
losetup -a
```

### Manual Cleanup of Single Device

```bash
sudo losetup -d /dev/loop1
```

## Safety Notes

- Loop devices are safe to use for testing
- No real hardware disks are affected
- Data is stored in `/tmp/stumpfworks-test-disks/`
- Cleanup script removes all test data

## Troubleshooting

### "No loop devices available"

```bash
# Load loop module
sudo modprobe loop

# Or increase max loop devices
sudo modprobe loop max_loop=16
```

### "Permission denied"

All scripts require root/sudo access:

```bash
sudo ./create-test-disks.sh
```

### "Device busy"

Unmount any mounted test volumes first:

```bash
sudo umount /mnt/test-volume
```

Then run cleanup.
