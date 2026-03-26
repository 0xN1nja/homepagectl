package widgets

import "strings"

type AuthType int

const (
	AuthNone AuthType = iota
	AuthAPIKey
	AuthToken
	AuthUserPass
)

type Auth struct {
	Type AuthType
}

type Info struct {
	WidgetType  string
	DisplayName string
	Icon        string
	Description string
	Group       string
	Auth        Auth
	ExtraFields map[string]string
}

func Lookup(containerName string) Info {
	name := normalize(containerName)
	for _, e := range entries {
		if strings.Contains(name, e.pattern) {
			return Info{
				WidgetType:  e.widgetType,
				DisplayName: e.displayName,
				Icon:        e.icon,
				Description: e.description,
				Group:       e.group,
				Auth:        e.auth,
				ExtraFields: e.extraFields,
			}
		}
	}
	return Info{
		DisplayName: toTitle(containerName),
		Icon:        strings.ToLower(containerName),
		Description: "Docker container",
		Group:       "Services",
		Auth:        Auth{Type: AuthNone},
	}
}

func normalize(s string) string {
	s = strings.ToLower(s)
	for _, suffix := range []string{"-app", "-web", "-server", "-backend"} {
		s = strings.TrimSuffix(s, suffix)
	}
	return s
}

func toTitle(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

var (
	apiKey   = Auth{Type: AuthAPIKey}
	token    = Auth{Type: AuthToken}
	userPass = Auth{Type: AuthUserPass}
	noAuth   = Auth{Type: AuthNone}
)

type entry struct {
	pattern     string
	widgetType  string
	displayName string
	icon        string
	description string
	group       string
	auth        Auth
	extraFields map[string]string
}

var entries = []entry{
	{"plex", "plex", "Plex", "plex", "Stream movies and TV shows", "Media", apiKey, nil},
	{"jellyfin", "jellyfin", "Jellyfin", "jellyfin", "Free media system", "Media", apiKey, nil},
	{"emby", "emby", "Emby", "emby", "Personal media server", "Media", apiKey, nil},
	{"navidrome", "navidrome", "Navidrome", "navidrome", "Music streaming server", "Media", userPass, nil},
	{"audiobookshelf", "audiobookshelf", "Audiobookshelf", "audiobookshelf", "Audiobook and podcast server", "Media", apiKey, nil},
	{"tubearchivist", "tubearchivist", "Tube Archivist", "tubearchivist", "YouTube media server", "Media", userPass, nil},
	{"tube-archivist", "tubearchivist", "Tube Archivist", "tubearchivist", "YouTube media server", "Media", userPass, nil},
	{"yourspotify", "yourspotify", "Your Spotify", "yourspotify", "Spotify listening history tracker", "Media", apiKey, nil},
	{"your-spotify", "yourspotify", "Your Spotify", "yourspotify", "Spotify listening history tracker", "Media", apiKey, nil},
	{"tautulli", "tautulli", "Tautulli", "tautulli", "Plex media server analytics", "Media", apiKey, nil},
	{"jellystat", "jellystat", "Jellystat", "jellystat", "Jellyfin statistics", "Media", apiKey, nil},
	{"kavita", "kavita", "Kavita", "kavita", "Digital library for manga and books", "Media", apiKey, nil},
	{"komga", "komga", "Komga", "komga", "Comics and manga server", "Media", userPass, nil},
	{"calibre-web", "calibre-web", "Calibre-Web", "calibre-web", "eBook library web frontend", "Media", noAuth, nil},
	{"calibre", "calibre-web", "Calibre", "calibre", "eBook library management", "Media", noAuth, nil},
	{"mylar", "mylar", "Mylar3", "mylar", "Comic book manager", "Media", apiKey, nil},
	{"readarr", "readarr", "Readarr", "readarr", "eBook collection manager", "Media", apiKey, nil},
	{"sonarr", "sonarr", "Sonarr", "sonarr", "TV show collection manager", "Media", apiKey, nil},
	{"radarr", "radarr", "Radarr", "radarr", "Movie collection manager", "Media", apiKey, nil},
	{"lidarr", "lidarr", "Lidarr", "lidarr", "Music collection manager", "Media", apiKey, nil},
	{"bazarr", "bazarr", "Bazarr", "bazarr", "Subtitle management", "Media", apiKey, nil},
	{"prowlarr", "prowlarr", "Prowlarr", "prowlarr", "Indexer manager for the arr stack", "Media", apiKey, nil},
	{"trailarr", "trailarr", "Trailarr", "trailarr", "Trailer downloader for the arr stack", "Media", apiKey, nil},
	{"requestrr", "", "Requestrr", "requestrr", "Media request chatbot", "Media", noAuth, nil},
	{"jackett", "jackett", "Jackett", "jackett", "Torrent indexer proxy", "Media", apiKey, nil},
	{"overseerr", "overseerr", "Overseerr", "overseerr", "Request movies and TV shows", "Media", apiKey, nil},
	{"jellyseerr", "jellyseerr", "Jellyseerr", "jellyseerr", "Request movies and TV shows", "Media", apiKey, nil},
	{"ombi", "ombi", "Ombi", "ombi", "Media request management", "Media", apiKey, nil},
	{"qbittorrent", "qbittorrent", "qBittorrent", "qbittorrent", "BitTorrent client", "Media", userPass, nil},
	{"transmission", "transmission", "Transmission", "transmission", "BitTorrent client", "Media", userPass, nil},
	{"deluge", "deluge", "Deluge", "deluge", "Lightweight BitTorrent client", "Media", userPass, nil},
	{"rutorrent", "rutorrent", "ruTorrent", "rutorrent", "BitTorrent client", "Media", userPass, nil},
	{"nzbget", "nzbget", "NZBGet", "nzbget", "Usenet downloader", "Media", userPass, nil},
	{"sabnzbd", "sabnzbd", "SABnzbd", "sabnzbd", "Usenet downloader", "Media", apiKey, nil},
	{"autobrr", "autobrr", "Autobrr", "autobrr", "Download automation for torrents", "Media", apiKey, nil},
	{"pyload", "pyload", "pyLoad", "pyload", "Download manager", "Media", userPass, nil},
	{"jdownloader", "jdownloader", "JDownloader", "jdownloader", "Automated download manager", "Media", userPass, nil},
	{"tdarr", "tdarr", "Tdarr", "tdarr", "Distributed transcoding automation", "Media", noAuth, nil},
	{"immich", "immich", "Immich", "immich", "Self-hosted photo and video backup", "Services", apiKey, nil},
	{"photoprism", "photoprism", "PhotoPrism", "photoprism", "AI-powered photo library", "Services", userPass, nil},
	{"frigate", "frigate", "Frigate", "frigate", "NVR with real-time AI object detection", "Services", noAuth, nil},
	{"mealie", "mealie", "Mealie", "mealie", "Self-hosted recipe manager", "Services", apiKey, nil},
	{"homeassistant", "homeassistant", "Home Assistant", "home-assistant", "Home automation platform", "Services", token, nil},
	{"home-assistant", "homeassistant", "Home Assistant", "home-assistant", "Home automation platform", "Services", token, nil},
	{"homebridge", "homebridge", "Homebridge", "homebridge", "HomeKit for non-HomeKit devices", "Services", userPass, nil},
	{"esphome", "esphome", "ESPHome", "esphome", "ESP8266/ESP32 firmware builder", "Services", noAuth, nil},
	{"node-red", "", "Node-RED", "node-red", "Flow-based programming", "Services", noAuth, nil},
	{"pihole", "pihole", "Pi-hole", "pi-hole", "Network-wide ad blocker", "Services", apiKey, nil},
	{"pi-hole", "pihole", "Pi-hole", "pi-hole", "Network-wide ad blocker", "Services", apiKey, nil},
	{"adguard", "adguard-home", "AdGuard Home", "adguard-home", "DNS-based ad and tracker blocker", "Services", userPass, nil},
	{"unifi", "unifi-controller", "Unifi Controller", "unifi", "Ubiquiti network management", "Services", userPass, nil},
	{"traefik", "traefik", "Traefik", "traefik", "Modern HTTP reverse proxy", "Services", noAuth, nil},
	{"nginx-proxy-manager", "nginx-proxy-manager", "Nginx Proxy Manager", "nginx-proxy-manager", "Easy Nginx reverse proxy with SSL", "Services", userPass, nil},
	{"wireguard", "", "WireGuard", "wireguard", "Fast VPN", "Services", noAuth, nil},
	{"wg-easy", "wgeasy", "Wg-Easy", "wgeasy", "WireGuard VPN with web UI", "Services", noAuth, nil},
	{"wgeasy", "wgeasy", "Wg-Easy", "wgeasy", "WireGuard VPN with web UI", "Services", noAuth, nil},
	{"gluetun", "gluetun", "Gluetun", "gluetun", "VPN client in a container", "Services", noAuth, nil},
	{"tailscale", "tailscale", "Tailscale", "tailscale", "Zero config VPN", "Services", apiKey, nil},
	{"crowdsec", "crowdsec", "CrowdSec", "crowdsec", "Collaborative security engine", "Services", apiKey, nil},
	{"opnsense", "opnsense", "OPNsense", "opnsense", "Open source firewall", "Services", userPass, nil},
	{"pfsense", "pfsense", "pfSense", "pfsense", "Open source firewall", "Services", userPass, nil},
	{"grafana", "grafana", "Grafana", "grafana", "Observability and metrics dashboards", "Services", apiKey, nil},
	{"prometheus", "prometheus", "Prometheus", "prometheus", "Metrics collection and alerting", "Services", noAuth, nil},
	{"netdata", "netdata", "Netdata", "netdata", "Real-time performance monitoring", "Services", noAuth, nil},
	{"glances", "glances", "Glances", "glances", "System monitoring dashboard", "Services", noAuth, map[string]string{"version": "4", "metric": "info", "chart": "false"}},
	{"uptime-kuma", "uptimekuma", "Uptime Kuma", "uptime-kuma", "Self-hosted uptime monitoring", "Services", apiKey, nil},
	{"uptimekuma", "uptimekuma", "Uptime Kuma", "uptime-kuma", "Self-hosted uptime monitoring", "Services", apiKey, nil},
	{"healthchecks", "healthchecks", "Health Checks", "healthchecks", "Cron job monitoring", "Services", apiKey, nil},
	{"scrutiny", "scrutiny", "Scrutiny", "scrutiny", "Hard drive health monitoring", "Services", noAuth, nil},
	{"speedtest", "speedtest-tracker", "Speedtest Tracker", "speedtest-tracker", "Internet speed test tracker", "Services", apiKey, nil},
	{"zabbix", "zabbix", "Zabbix", "zabbix", "Enterprise-grade monitoring", "Services", userPass, nil},
	{"portainer", "portainer", "Portainer", "portainer", "Docker management", "Services", apiKey, map[string]string{"env": "1"}},
	{"watchtower", "watchtower", "Watchtower", "watchtower", "Automatic Docker container updates", "Services", apiKey, nil},
	{"gitea", "gitea", "Gitea", "gitea", "Self-hosted Git service", "Services", apiKey, nil},
	{"gitlab", "gitlab", "GitLab", "gitlab", "DevOps platform", "Services", apiKey, nil},
	{"authentik", "authentik", "Authentik", "authentik", "Identity provider and SSO", "Services", apiKey, map[string]string{"slug": "YOUR_FLOW_SLUG"}},
	{"vaultwarden", "", "Vaultwarden", "vaultwarden", "Bitwarden-compatible password manager", "Services", noAuth, nil},
	{"bitwarden", "", "Bitwarden", "bitwarden", "Password manager", "Services", noAuth, nil},
	{"nextcloud", "nextcloud", "Nextcloud", "nextcloud", "Self-hosted cloud storage", "Services", userPass, nil},
	{"filebrowser", "filebrowser", "Filebrowser", "filebrowser", "Web-based file manager", "Services", userPass, nil},
	{"kopia", "kopia", "Kopia", "kopia", "Cross-platform backup tool", "Services", userPass, nil},
	{"homebox", "homebox", "Homebox", "homebox", "Home inventory management", "Services", noAuth, nil},
	{"vikunja", "vikunja", "Vikunja", "vikunja", "Open-source to-do app", "Services", apiKey, nil},
	{"wallos", "wallos", "Wallos", "wallos", "Subscription tracker", "Services", apiKey, nil},
	{"linkwarden", "linkwarden", "Linkwarden", "linkwarden", "Bookmark manager", "Services", apiKey, nil},
	{"freshrss", "freshrss", "FreshRSS", "freshrss", "Self-hosted RSS aggregator", "Services", userPass, nil},
	{"miniflux", "miniflux", "Miniflux", "miniflux", "Minimalist RSS reader", "Services", apiKey, nil},
	{"changedetection", "changedetectionio", "Changedetection.io", "changedetection", "Website change detection", "Services", apiKey, nil},
	{"paperless", "paperlessngx", "Paperless-ngx", "paperless-ngx", "Document management system", "Services", apiKey, nil},
	{"trilium", "trilium", "Trilium", "trilium", "Hierarchical note taking", "Services", noAuth, nil},
	{"romm", "romm", "Romm", "romm", "Self-hosted ROM manager", "Services", userPass, nil},
	{"karakeep", "karakeep", "Karakeep", "karakeep", "Bookmark and read-it-later app", "Services", apiKey, nil},
	{"komodo", "komodo", "Komodo", "komodo", "Server and deployment manager", "Services", apiKey, nil},
	{"mastodon", "mastodon", "Mastodon", "mastodon", "Federated social network", "Services", apiKey, nil},
	{"gotify", "gotify", "Gotify", "gotify", "Self-hosted push notifications", "Services", apiKey, nil},
	{"mailcow", "mailcow", "Mailcow", "mailcow", "Email server suite", "Services", apiKey, nil},
	{"stirlingpdf", "", "Stirling PDF", "stirlingpdf", "PDF manipulation tool", "Services", noAuth, nil},
	{"stirling-pdf", "", "Stirling PDF", "stirlingpdf", "PDF manipulation tool", "Services", noAuth, nil},
}
