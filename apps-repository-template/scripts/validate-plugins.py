#!/usr/bin/env python3
"""
Validate all plugin.json files
"""

import json
import sys
from pathlib import Path

PLUGINS_DIR = Path("plugins")

# Required fields
REQUIRED_FIELDS = [
    'id',
    'name',
    'version',
    'author',
    'description',
    'icon',
    'category'
]

# Valid categories
VALID_CATEGORIES = [
    'storage',
    'media',
    'communication',
    'development',
    'monitoring',
    'networking',
    'productivity',
    'security',
    'utilities'
]


def validate_semver(version):
    """Validate semantic versioning"""
    parts = version.split('.')
    if len(parts) < 3:
        return False

    try:
        # Check numeric parts
        for part in parts[:3]:
            int(part)
        return True
    except ValueError:
        return False


def validate_plugin_id(plugin_id):
    """Validate plugin ID format (reverse domain notation)"""
    if not plugin_id:
        return False

    parts = plugin_id.split('.')
    return len(parts) >= 3


def validate_plugin_manifest(manifest_file):
    """Validate a single plugin.json file"""
    errors = []
    warnings = []

    try:
        with open(manifest_file, 'r') as f:
            data = json.load(f)
    except json.JSONDecodeError as e:
        return [f"Invalid JSON: {e}"], []
    except Exception as e:
        return [f"Error reading file: {e}"], []

    # Check required fields
    for field in REQUIRED_FIELDS:
        if field not in data:
            errors.append(f"Missing required field: {field}")

    # Validate ID format
    if 'id' in data and not validate_plugin_id(data['id']):
        errors.append(f"Invalid plugin ID format: {data['id']} (should be reverse domain notation)")

    # Validate version format
    if 'version' in data and not validate_semver(data['version']):
        warnings.append(f"Version {data['version']} doesn't follow semantic versioning")

    # Validate category
    if 'category' in data and data['category'] not in VALID_CATEGORIES:
        errors.append(f"Invalid category: {data['category']}. Must be one of: {', '.join(VALID_CATEGORIES)}")

    # Check description length
    if 'description' in data:
        desc_len = len(data['description'])
        if desc_len < 20:
            warnings.append(f"Description is very short ({desc_len} characters)")
        elif desc_len > 200:
            warnings.append(f"Description is very long ({desc_len} characters)")

    # Check for common issues
    if 'name' in data and len(data['name']) < 3:
        errors.append("Plugin name is too short")

    if 'author' in data and '@' in data['author']:
        warnings.append("Author field should not contain email (use links section instead)")

    return errors, warnings


def main():
    """Validate all plugins"""
    all_errors = []
    all_warnings = []

    print("🔍 Validating plugin manifests...\n")

    for plugin_dir in PLUGINS_DIR.iterdir():
        if not plugin_dir.is_dir():
            continue

        manifest_file = plugin_dir / "plugin.json"
        if not manifest_file.exists():
            print(f"⚠️  {plugin_dir.name}: Missing plugin.json")
            all_warnings.append(f"{plugin_dir.name}: Missing plugin.json")
            continue

        errors, warnings = validate_plugin_manifest(manifest_file)

        if errors:
            print(f"❌ {plugin_dir.name}:")
            for error in errors:
                print(f"   ERROR: {error}")
            all_errors.extend(errors)

        if warnings:
            print(f"⚠️  {plugin_dir.name}:")
            for warning in warnings:
                print(f"   WARNING: {warning}")
            all_warnings.extend(warnings)

        if not errors and not warnings:
            print(f"✅ {plugin_dir.name}: Valid")

    print(f"\n{'='*60}")
    print(f"Total plugins checked: {len(list(PLUGINS_DIR.iterdir()))}")
    print(f"Errors: {len(all_errors)}")
    print(f"Warnings: {len(all_warnings)}")

    if all_errors:
        print("\n❌ Validation failed!")
        sys.exit(1)
    else:
        print("\n✅ All plugins are valid!")
        sys.exit(0)


if __name__ == "__main__":
    main()
