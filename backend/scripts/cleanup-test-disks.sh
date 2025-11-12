#!/bin/bash
# Script to clean up test loop devices
# Usage: sudo ./cleanup-test-disks.sh

set -e

TEST_DIR="/tmp/stumpfworks-test-disks"

echo "Cleaning up test disks..."

# Find all loop devices pointing to our test directory
for LOOP_DEV in $(losetup -a | grep "${TEST_DIR}" | cut -d: -f1); do
    echo "Detaching ${LOOP_DEV}"
    losetup -d "${LOOP_DEV}"
done

# Remove disk images
if [ -d "${TEST_DIR}" ]; then
    echo "Removing test disk images from ${TEST_DIR}"
    rm -rf "${TEST_DIR}"
fi

echo "Cleanup complete!"
