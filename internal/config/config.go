package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Host     HostConfig             `toml:"host"`
	Docker   DockerConfig           `toml:"docker"`
	Homepage HomepageConfig         `toml:"homepage"`
	Layout   map[string]LayoutGroup `toml:"layout"`
	Groups   map[string]string      `toml:"groups"`
}

type HostConfig struct {
	IP       string `toml:"ip"`
	Protocol string `toml:"protocol"`
}

type DockerConfig struct {
	Socket string   `toml:"socket"`
	Skip   []string `toml:"skip"`
}

type HomepageConfig struct {
	Title           string `toml:"title"`
	Color           string `toml:"color"`
	Theme           string `toml:"theme"`
	HeaderStyle     string `toml:"header_style"`
	Target          string `toml:"target"`
	ShowStats       bool   `toml:"show_stats"`
	StatusStyle     string `toml:"status_style"`
	UseEqualHeights bool   `toml:"use_equal_heights"`
	MaxGroupColumns int    `toml:"max_group_columns"`
	Tabs            bool   `toml:"tabs"`
	SortAlpha       bool   `toml:"sort_alphabetically"`
}

type LayoutGroup struct {
	Style   string `toml:"style"`
	Columns int    `toml:"columns"`
	Header  *bool  `toml:"header"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Host: HostConfig{
			Protocol: "http",
		},
		Docker: DockerConfig{
			Socket: "/var/run/docker.sock",
		},
		Homepage: HomepageConfig{
			Title:           "My Homepage",
			Color:           "slate",
			Theme:           "dark",
			HeaderStyle:     "underlined",
			Target:          "_blank",
			StatusStyle:     "dot",
			UseEqualHeights: true,
		},
		Layout: make(map[string]LayoutGroup),
		Groups: make(map[string]string),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %q: run `homepagectl init > homepagectl.toml` to create one", path)
	}

	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %q: %w", path, err)
	}

	return cfg, nil
}

func Example() string {
	return `[host]
ip = "192.168.1.100"
protocol = "http"

[docker]
socket = "/var/run/docker.sock"
skip = ["watchtower", "homepage", "portainer-agent"]

[homepage]
title = "My Homepage"
color = "slate"
theme = "dark"
header_style = "underlined"
target = "_blank"
show_stats = false
status_style = "dot"
use_equal_heights = true
tabs = false
sort_alphabetically = false

[layout.Services]
style = "row"
columns = 3

[layout.Media]
style = "row"
columns = 3

[groups]
qBittorrent = "Services"
`
}
