#!/usr/bin/env python3
"""
Generate registry.json from all plugin.json files in plugins/ directory
"""

import json
import os
from datetime import datetime
from pathlib import Path

PLUGINS_DIR = Path("plugins")
REGISTRY_FILE = Path("registry.json")
REPO_URL = "https://github.com/Stumpf-works/stumpfworks-nas-apps"


def find_plugin_manifests():
    """Find all plugin.json files"""
    manifests = []

    for plugin_dir in PLUGINS_DIR.iterdir():
        if not plugin_dir.is_dir():
            continue

        manifest_file = plugin_dir / "plugin.json"
        if manifest_file.exists():
            manifests.append(manifest_file)

    return manifests


def load_plugin_manifest(manifest_file):
    """Load and parse a plugin.json file"""
    with open(manifest_file, 'r') as f:
        data = json.load(f)

    plugin_id = data.get('id')
    plugin_dir = manifest_file.parent.name
    version = data.get('version', '1.0.0')

    # Construct URLs
    repository_url = f"{REPO_URL}/tree/main/plugins/{plugin_dir}"

    # Try to find release asset
    download_url = f"{REPO_URL}/releases/download/{plugin_dir}-v{version}/{plugin_dir}-v{version}.tar.gz"

    # Build registry entry
    entry = {
        "id": plugin_id,
        "name": data.get('name'),
        "version": version,
        "author": data.get('author'),
        "description": data.get('description'),
        "icon": data.get('icon', '🔌'),
        "category": data.get('category', 'utilities'),
        "repository_url": repository_url,
        "download_url": download_url,
        "homepage": data.get('links', {}).get('homepage', repository_url),
        "min_nas_version": data.get('requires', {}).get('minNasVersion', '0.1.0'),
        "require_docker": data.get('requires', {}).get('docker', False),
        "required_ports": data.get('requires', {}).get('ports', []),
        "screenshots": [],
        "tags": []
    }

    # Add tags from category and name
    tags = [entry['category']]
    if 'tags' in data:
        tags.extend(data['tags'])
    entry['tags'] = tags

    return entry


def generate_registry():
    """Generate the complete registry.json"""
    manifests = find_plugin_manifests()

    print(f"Found {len(manifests)} plugin manifests")

    plugins = []
    for manifest in manifests:
        try:
            entry = load_plugin_manifest(manifest)
            plugins.append(entry)
            print(f"✅ Loaded {entry['name']} v{entry['version']}")
        except Exception as e:
            print(f"❌ Error loading {manifest}: {e}")
            continue

    # Sort plugins by name
    plugins.sort(key=lambda p: p['name'])

    # Build registry
    registry = {
        "version": "1.0.0",
        "updated": datetime.utcnow().isoformat() + "Z",
        "repository": REPO_URL,
        "plugins": plugins
    }

    # Write registry.json
    with open(REGISTRY_FILE, 'w') as f:
        json.dump(registry, f, indent=2, ensure_ascii=False)
        f.write('\n')  # Add trailing newline

    print(f"\n✅ Generated registry.json with {len(plugins)} plugins")
    return registry


if __name__ == "__main__":
    generate_registry()
