// Package docker provides Docker management functionality
package docker

// ComposeTemplate represents a pre-configured Docker Compose template
type ComposeTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Category    string            `json:"category"` // media, download, automation, monitoring, etc.
	Author      string            `json:"author"`
	Version     string            `json:"version"`
	Compose     string            `json:"compose"`      // Docker Compose YAML content
	Variables   map[string]string `json:"variables"`    // User-customizable variables
	Requirements struct {
		MinMemoryMB int      `json:"min_memory_mb"`
		MinDiskGB   int      `json:"min_disk_gb"`
		Ports       []int    `json:"ports"`
		Notes       []string `json:"notes"`
	} `json:"requirements"`
}

// BuiltinTemplates contains all pre-configured templates
var BuiltinTemplates = []ComposeTemplate{
	{
		ID:          "plex",
		Name:        "Plex Media Server",
		Description: "Stream your media collection to any device. Supports hardware transcoding and offline sync.",
		Icon:        "🎬",
		Category:    "media",
		Author:      "StumpfWorks",
		Version:     "1.0.0",
		Variables: map[string]string{
			"MEDIA_PATH":   "/mnt/media",
			"CONFIG_PATH":  "/var/lib/stumpfworks/plex/config",
			"TRANSCODE_PATH": "/var/lib/stumpfworks/plex/transcode",
			"PUID":         "1000",
			"PGID":         "1000",
			"TZ":           "Europe/Berlin",
		},
		Compose: `version: '3.8'

services:
  plex:
    image: lscr.io/linuxserver/plex:latest
    container_name: plex
    network_mode: host
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
      - VERSION=docker
      - PLEX_CLAIM= # Optional: Get from https://www.plex.tv/claim
    volumes:
      - {{CONFIG_PATH}}:/config
      - {{TRANSCODE_PATH}}:/transcode
      - {{MEDIA_PATH}}:/media
    devices:
      - /dev/dri:/dev/dri # Hardware transcoding (Intel QuickSync)
    restart: unless-stopped
    labels:
      - "com.stumpfworks.app=plex"
      - "com.stumpfworks.category=media"`,
	},
	{
		ID:          "jellyfin",
		Name:        "Jellyfin Media Server",
		Description: "Open-source media server with no premium features or tracking. Free alternative to Plex.",
		Icon:        "🍿",
		Category:    "media",
		Author:      "StumpfWorks",
		Version:     "1.0.0",
		Variables: map[string]string{
			"MEDIA_PATH":  "/mnt/media",
			"CONFIG_PATH": "/var/lib/stumpfworks/jellyfin/config",
			"CACHE_PATH":  "/var/lib/stumpfworks/jellyfin/cache",
			"PUID":        "1000",
			"PGID":        "1000",
			"TZ":          "Europe/Berlin",
		},
		Compose: `version: '3.8'

services:
  jellyfin:
    image: jellyfin/jellyfin:latest
    container_name: jellyfin
    network_mode: host
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
    volumes:
      - {{CONFIG_PATH}}:/config
      - {{CACHE_PATH}}:/cache
      - {{MEDIA_PATH}}:/media:ro
    devices:
      - /dev/dri:/dev/dri # Hardware transcoding
    restart: unless-stopped
    labels:
      - "com.stumpfworks.app=jellyfin"
      - "com.stumpfworks.category=media"`,
	},
	{
		ID:          "sonarr",
		Name:        "Sonarr",
		Description: "Automated TV show downloading and organization. Monitors RSS feeds and downloads via torrent/usenet.",
		Icon:        "📺",
		Category:    "automation",
		Author:      "StumpfWorks",
		Version:     "1.0.0",
		Variables: map[string]string{
			"CONFIG_PATH":    "/var/lib/stumpfworks/sonarr/config",
			"TV_PATH":        "/mnt/media/tv",
			"DOWNLOADS_PATH": "/mnt/downloads",
			"PUID":           "1000",
			"PGID":           "1000",
			"TZ":             "Europe/Berlin",
		},
		Compose: `version: '3.8'

services:
  sonarr:
    image: lscr.io/linuxserver/sonarr:latest
    container_name: sonarr
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
    volumes:
      - {{CONFIG_PATH}}:/config
      - {{TV_PATH}}:/tv
      - {{DOWNLOADS_PATH}}:/downloads
    ports:
      - 8989:8989
    restart: unless-stopped
    labels:
      - "com.stumpfworks.app=sonarr"
      - "com.stumpfworks.category=automation"`,
	},
	{
		ID:          "radarr",
		Name:        "Radarr",
		Description: "Automated movie downloading and organization. Monitors RSS feeds and downloads via torrent/usenet.",
		Icon:        "🎥",
		Category:    "automation",
		Author:      "StumpfWorks",
		Version:     "1.0.0",
		Variables: map[string]string{
			"CONFIG_PATH":    "/var/lib/stumpfworks/radarr/config",
			"MOVIES_PATH":    "/mnt/media/movies",
			"DOWNLOADS_PATH": "/mnt/downloads",
			"PUID":           "1000",
			"PGID":           "1000",
			"TZ":             "Europe/Berlin",
		},
		Compose: `version: '3.8'

services:
  radarr:
    image: lscr.io/linuxserver/radarr:latest
    container_name: radarr
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
    volumes:
      - {{CONFIG_PATH}}:/config
      - {{MOVIES_PATH}}:/movies
      - {{DOWNLOADS_PATH}}:/downloads
    ports:
      - 7878:7878
    restart: unless-stopped
    labels:
      - "com.stumpfworks.app=radarr"
      - "com.stumpfworks.category=automation"`,
	},
	{
		ID:          "prowlarr",
		Name:        "Prowlarr",
		Description: "Indexer manager for Sonarr, Radarr, Lidarr, and Readarr. Centralized indexer management.",
		Icon:        "🔍",
		Category:    "automation",
		Author:      "StumpfWorks",
		Version:     "1.0.0",
		Variables: map[string]string{
			"CONFIG_PATH": "/var/lib/stumpfworks/prowlarr/config",
			"PUID":        "1000",
			"PGID":        "1000",
			"TZ":          "Europe/Berlin",
		},
		Compose: `version: '3.8'

services:
  prowlarr:
    image: lscr.io/linuxserver/prowlarr:latest
    container_name: prowlarr
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
    volumes:
      - {{CONFIG_PATH}}:/config
    ports:
      - 9696:9696
    restart: unless-stopped
    labels:
      - "com.stumpfworks.app=prowlarr"
      - "com.stumpfworks.category=automation"`,
	},
	{
		ID:          "transmission",
		Name:        "Transmission",
		Description: "Fast and lightweight BitTorrent client with web interface. Low resource usage.",
		Icon:        "⬇️",
		Category:    "download",
		Author:      "StumpfWorks",
		Version:     "1.0.0",
		Variables: map[string]string{
			"CONFIG_PATH":    "/var/lib/stumpfworks/transmission/config",
			"DOWNLOADS_PATH": "/mnt/downloads",
			"WATCH_PATH":     "/mnt/downloads/watch",
			"PUID":           "1000",
			"PGID":           "1000",
			"TZ":             "Europe/Berlin",
		},
		Compose: `version: '3.8'

services:
  transmission:
    image: lscr.io/linuxserver/transmission:latest
    container_name: transmission
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
      - USER=admin
      - PASS=changeme
    volumes:
      - {{CONFIG_PATH}}:/config
      - {{DOWNLOADS_PATH}}:/downloads
      - {{WATCH_PATH}}:/watch
    ports:
      - 9091:9091
      - 51413:51413
      - 51413:51413/udp
    restart: unless-stopped
    labels:
      - "com.stumpfworks.app=transmission"
      - "com.stumpfworks.category=download"`,
	},
	{
		ID:          "qbittorrent",
		Name:        "qBittorrent",
		Description: "Feature-rich BitTorrent client with advanced features and modern web UI.",
		Icon:        "📦",
		Category:    "download",
		Author:      "StumpfWorks",
		Version:     "1.0.0",
		Variables: map[string]string{
			"CONFIG_PATH":    "/var/lib/stumpfworks/qbittorrent/config",
			"DOWNLOADS_PATH": "/mnt/downloads",
			"PUID":           "1000",
			"PGID":           "1000",
			"TZ":             "Europe/Berlin",
			"WEBUI_PORT":     "8080",
		},
		Compose: `version: '3.8'

services:
  qbittorrent:
    image: lscr.io/linuxserver/qbittorrent:latest
    container_name: qbittorrent
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
      - WEBUI_PORT={{WEBUI_PORT}}
    volumes:
      - {{CONFIG_PATH}}:/config
      - {{DOWNLOADS_PATH}}:/downloads
    ports:
      - {{WEBUI_PORT}}:{{WEBUI_PORT}}
      - 6881:6881
      - 6881:6881/udp
    restart: unless-stopped
    labels:
      - "com.stumpfworks.app=qbittorrent"
      - "com.stumpfworks.category=download"`,
	},
	{
		ID:          "media-stack-complete",
		Name:        "Complete Media Stack",
		Description: "All-in-one media solution: Plex + Sonarr + Radarr + Prowlarr + Transmission. Perfect for beginners.",
		Icon:        "🎯",
		Category:    "media",
		Author:      "StumpfWorks",
		Version:     "1.0.0",
		Variables: map[string]string{
			"MEDIA_PATH":     "/mnt/media",
			"DOWNLOADS_PATH": "/mnt/downloads",
			"CONFIG_BASE":    "/var/lib/stumpfworks",
			"PUID":           "1000",
			"PGID":           "1000",
			"TZ":             "Europe/Berlin",
		},
		Compose: `version: '3.8'

services:
  plex:
    image: lscr.io/linuxserver/plex:latest
    container_name: plex
    network_mode: host
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
      - VERSION=docker
    volumes:
      - {{CONFIG_BASE}}/plex/config:/config
      - {{CONFIG_BASE}}/plex/transcode:/transcode
      - {{MEDIA_PATH}}:/media
    devices:
      - /dev/dri:/dev/dri
    restart: unless-stopped

  sonarr:
    image: lscr.io/linuxserver/sonarr:latest
    container_name: sonarr
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
    volumes:
      - {{CONFIG_BASE}}/sonarr:/config
      - {{MEDIA_PATH}}/tv:/tv
      - {{DOWNLOADS_PATH}}:/downloads
    ports:
      - 8989:8989
    restart: unless-stopped

  radarr:
    image: lscr.io/linuxserver/radarr:latest
    container_name: radarr
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
    volumes:
      - {{CONFIG_BASE}}/radarr:/config
      - {{MEDIA_PATH}}/movies:/movies
      - {{DOWNLOADS_PATH}}:/downloads
    ports:
      - 7878:7878
    restart: unless-stopped

  prowlarr:
    image: lscr.io/linuxserver/prowlarr:latest
    container_name: prowlarr
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
    volumes:
      - {{CONFIG_BASE}}/prowlarr:/config
    ports:
      - 9696:9696
    restart: unless-stopped

  transmission:
    image: lscr.io/linuxserver/transmission:latest
    container_name: transmission
    environment:
      - PUID={{PUID}}
      - PGID={{PGID}}
      - TZ={{TZ}}
      - USER=admin
      - PASS=changeme
    volumes:
      - {{CONFIG_BASE}}/transmission:/config
      - {{DOWNLOADS_PATH}}:/downloads
      - {{DOWNLOADS_PATH}}/watch:/watch
    ports:
      - 9091:9091
      - 51413:51413
      - 51413:51413/udp
    restart: unless-stopped`,
	},
}

// GetTemplateByID returns a template by its ID
func GetTemplateByID(id string) *ComposeTemplate {
	for _, tpl := range BuiltinTemplates {
		if tpl.ID == id {
			return &tpl
		}
	}
	return nil
}

// GetTemplatesByCategory returns all templates in a category
func GetTemplatesByCategory(category string) []ComposeTemplate {
	var result []ComposeTemplate
	for _, tpl := range BuiltinTemplates {
		if tpl.Category == category {
			result = append(result, tpl)
		}
	}
	return result
}

// GetAllCategories returns all unique template categories
func GetAllCategories() []string {
	categoryMap := make(map[string]bool)
	for _, tpl := range BuiltinTemplates {
		categoryMap[tpl.Category] = true
	}

	var categories []string
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	return categories
}

// RenderTemplate replaces variables in the compose template with provided values
func RenderTemplate(template *ComposeTemplate, variables map[string]string) string {
	result := template.Compose

	// Merge default variables with user-provided ones
	finalVars := make(map[string]string)
	for k, v := range template.Variables {
		finalVars[k] = v
	}
	for k, v := range variables {
		finalVars[k] = v
	}

	// Replace all variables
	for key, value := range finalVars {
		placeholder := "{{" + key + "}}"
		result = replaceAll(result, placeholder, value)
	}

	return result
}

// Simple string replace helper
func replaceAll(s, old, new string) string {
	for i := 0; i < len(s); i++ {
		if hasPrefix(s[i:], old) {
			s = s[:i] + new + s[i+len(old):]
			i += len(new) - 1
		}
	}
	return s
}

func hasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}
