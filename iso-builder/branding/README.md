# StumpfWorks NAS OS - Branding

This directory contains visual branding assets for the ISO.

## Files

- **logo.png**: StumpfWorks logo (displayed in installer and boot menu)
- **splash.png**: Boot splash screen (1920x1080 recommended)
- **theme.conf**: Color scheme and theme configuration

## Customization

### Logo Requirements
- Format: PNG with transparency
- Size: 200x200px recommended
- Used in: Boot menu, installer, first-boot screen

### Splash Screen Requirements
- Format: PNG or JPG
- Size: 1920x1080px (Full HD)
- Used in: Plymouth boot splash

### Theme Configuration

Edit `theme.conf` to customize colors and appearance:

```ini
[Colors]
primary=#0066cc
secondary=#004499
background=#1a1a1a
text=#ffffff

[Boot]
show_logo=true
show_progress=true
animation=fade
```

## Adding Custom Branding

1. Replace `logo.png` with your logo
2. Replace `splash.png` with your splash screen
3. Update `theme.conf` with your colors
4. Rebuild ISO: `sudo ./build-iso.sh`

## Default Branding

If no custom branding is provided, the system will use:
- StumpfWorks default logo
- Debian boot splash
- Dark theme with blue accents
