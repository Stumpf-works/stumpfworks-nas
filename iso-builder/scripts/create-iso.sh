#!/bin/bash
# Create bootable ISO from the built system

set -e

BUILD_DIR="$1"
OUTPUT_DIR="$2"
ISO_NAME="$3"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Creating ISO image..."

# Create squashfs filesystem from chroot
echo "Creating squashfs filesystem..."
mksquashfs "${BUILD_DIR}/chroot" "${BUILD_DIR}/iso/live/filesystem.squashfs" \
    -comp xz \
    -e boot

echo "✓ Squashfs created"

# Copy kernel and initrd from chroot
echo "Copying kernel and initrd..."
cp "${BUILD_DIR}/chroot/boot/vmlinuz-"* "${BUILD_DIR}/iso/live/vmlinuz"
cp "${BUILD_DIR}/chroot/boot/initrd.img-"* "${BUILD_DIR}/iso/live/initrd.img"

echo "✓ Kernel and initrd copied"

# Create GRUB configuration
cat > "${BUILD_DIR}/iso/boot/grub/grub.cfg" << 'GRUBEOF'
set timeout=10
set default=0

menuentry "Install StumpfWorks NAS OS" {
    linux /live/vmlinuz boot=live components quiet splash
    initrd /live/initrd.img
}

menuentry "Install StumpfWorks NAS OS (Safe Mode)" {
    linux /live/vmlinuz boot=live components nomodeset
    initrd /live/initrd.img
}

menuentry "Live Mode (Try without installing)" {
    linux /live/vmlinuz boot=live components persistence
    initrd /live/initrd.img
}

menuentry "Rescue Mode" {
    linux /live/vmlinuz boot=live components rescue/enable=true
    initrd /live/initrd.img
}
GRUBEOF

echo "✓ GRUB configuration created"

# Create isolinux configuration (for legacy BIOS)
if [ -f /usr/lib/ISOLINUX/isolinux.bin ]; then
    cp /usr/lib/ISOLINUX/isolinux.bin "${BUILD_DIR}/iso/isolinux/"
    cp /usr/lib/syslinux/modules/bios/*.c32 "${BUILD_DIR}/iso/isolinux/" 2>/dev/null || true

    cat > "${BUILD_DIR}/iso/isolinux/isolinux.cfg" << 'ISOLINUXEOF'
DEFAULT install
TIMEOUT 100
PROMPT 0

LABEL install
    MENU LABEL Install StumpfWorks NAS OS
    KERNEL /live/vmlinuz
    APPEND initrd=/live/initrd.img boot=live components quiet splash

LABEL safe
    MENU LABEL Install StumpfWorks NAS OS (Safe Mode)
    KERNEL /live/vmlinuz
    APPEND initrd=/live/initrd.img boot=live components nomodeset

LABEL live
    MENU LABEL Live Mode (Try without installing)
    KERNEL /live/vmlinuz
    APPEND initrd=/live/initrd.img boot=live components persistence

LABEL rescue
    MENU LABEL Rescue Mode
    KERNEL /live/vmlinuz
    APPEND initrd=/live/initrd.img boot=live components rescue/enable=true
ISOLINUXEOF

    echo "✓ ISOLINUX configuration created"
fi

# Create ISO using xorriso (supports both BIOS and UEFI)
echo "Building ISO with xorriso..."

xorriso -as mkisofs \
    -iso-level 3 \
    -full-iso9660-filenames \
    -volid "STUMPFWORKS_NAS" \
    -appid "StumpfWorks NAS OS v1.0.0" \
    -publisher "Stumpf.Works Team" \
    -preparer "stumpfworks-nas-iso-builder" \
    -eltorito-boot isolinux/isolinux.bin \
    -eltorito-catalog isolinux/boot.cat \
    -no-emul-boot \
    -boot-load-size 4 \
    -boot-info-table \
    -isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin \
    -eltorito-alt-boot \
    -e boot/grub/efi.img \
    -no-emul-boot \
    -isohybrid-gpt-basdat \
    -output "${OUTPUT_DIR}/${ISO_NAME}" \
    "${BUILD_DIR}/iso" \
    2>/dev/null || {
        # Fallback: simpler ISO creation if xorriso advanced features fail
        echo "Falling back to basic ISO creation..."
        xorriso -as mkisofs \
            -iso-level 3 \
            -volid "STUMPFWORKS_NAS" \
            -output "${OUTPUT_DIR}/${ISO_NAME}" \
            "${BUILD_DIR}/iso"
    }

# Make ISO hybrid (bootable from USB)
if command -v isohybrid &> /dev/null && [ -f "${OUTPUT_DIR}/${ISO_NAME}" ]; then
    isohybrid "${OUTPUT_DIR}/${ISO_NAME}" 2>/dev/null || true
fi

echo "✓ ISO created: ${OUTPUT_DIR}/${ISO_NAME}"

# Calculate checksums
echo "Calculating checksums..."
cd "${OUTPUT_DIR}"
sha256sum "${ISO_NAME}" > "${ISO_NAME}.sha256"
md5sum "${ISO_NAME}" > "${ISO_NAME}.md5"

echo "✓ Checksums created"
echo "  SHA256: $(cat "${ISO_NAME}.sha256")"
echo "  MD5:    $(cat "${ISO_NAME}.md5")"
