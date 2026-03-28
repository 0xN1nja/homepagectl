# homepagectl

Tired of manually writing homepage configs every time you spin up a new container? `homepagectl` fixes that.

Queries your Docker socket, matches containers against homepage known services, auto-detects ports, and writes `services.yaml`, `settings.yaml`, and `.env` — no manual YAML editing required.

## Installation

```bash
git clone https://github.com/0xN1nja/homepagectl.git
cd homepagectl
make build
sudo mv homepagectl /usr/local/bin/
```

## Usage

```bash
homepagectl init | tee <path-to-homepage-config>/homepagectl.toml
homepagectl list
homepagectl generate --output <path-to-homepage-config>
homepagectl generate --output <path-to-homepage-config> --dry-run
```

| Command     | Description                                       |
| ----------- | ------------------------------------------------- |
| `init`      | Generate a `homepagectl.toml` config file         |
| `list`      | Preview detected containers, widgets and groups   |
| `generate`  | Write `services.yaml`, `settings.yaml` and `.env` |
| `--dry-run` | Print output to stdout without writing files      |

> Container labels (`homepage.name`, `homepage.icon`, `homepage.group`, `homepage.href`, etc.) are read at generation time and take priority over auto-detection. See [homepage Docker labels](https://gethomepage.dev/configs/docker/#automatic-service-discovery).

## Config

Everything is controlled through `homepagectl.toml`. Running `homepagectl generate` produces `services.yaml`, `settings.yaml` and `.env` on what you set here.

```toml
[host]
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
jellyfin = "Media"
portainer = "Services"
```
