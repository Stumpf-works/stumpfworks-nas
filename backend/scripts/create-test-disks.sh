#!/bin/bash
# Script to create loop devices for testing storage management
# Usage: sudo ./create-test-disks.sh [count] [size_in_mb]

set -e

COUNT=${1:-3}
SIZE_MB=${2:-1024}
TEST_DIR="/tmp/stumpfworks-test-disks"

echo "Creating ${COUNT} test disk(s) of ${SIZE_MB}MB each..."

# Create directory for test disk images
mkdir -p "${TEST_DIR}"

# Create loop devices
for i in $(seq 1 $COUNT); do
    DISK_FILE="${TEST_DIR}/test-disk-${i}.img"

    # Create sparse file
    echo "Creating disk image: ${DISK_FILE}"
    dd if=/dev/zero of="${DISK_FILE}" bs=1M count=${SIZE_MB} 2>/dev/null || true

    # Setup loop device
    LOOP_DEV=$(losetup -f)
    echo "Setting up loop device: ${LOOP_DEV}"
    losetup "${LOOP_DEV}" "${DISK_FILE}"

    echo "Created loop device: ${LOOP_DEV} -> ${DISK_FILE}"
done

echo ""
echo "Test disks created successfully!"
echo ""
echo "To list loop devices:"
echo "  losetup -a"
echo ""
echo "To clean up later:"
echo "  sudo ./cleanup-test-disks.sh"
echo ""
echo "Example usage in NAS UI:"
echo "  1. Go to Storage Manager"
echo "  2. Click 'Create Volume'"
echo "  3. Select a loop device (e.g., /dev/loop1)"
echo "  4. Choose filesystem (ext4)"
echo "  5. Set mount point (e.g., /mnt/test-volume)"
