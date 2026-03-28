package generator

import (
	"bufio"
	"fmt"
	"sort"
	"strings"

	"github.com/0xN1nja/homepagectl/internal/config"
	"github.com/0xN1nja/homepagectl/internal/docker"
	"github.com/0xN1nja/homepagectl/internal/widgets"
)

const sentinel = "# managed by homepagectl"

func IsManagedFile(content string) bool {
	for _, line := range strings.SplitN(content, "\n", 5) {
		if strings.TrimSpace(line) == sentinel {
			return true
		}
	}
	return false
}

func Services(cfg *config.Config, containers []docker.Container, existing string) string {
	groups := make(map[string][]docker.Container)

	for _, c := range containers {
		if shouldSkip(c.Name, cfg.Docker.Skip) {
			continue
		}
		group := resolveGroup(c, cfg)
		groups[group] = append(groups[group], c)
	}

	if cfg.Homepage.SortAlpha {
		for k := range groups {
			sort.Slice(groups[k], func(i, j int) bool {
				return strings.ToLower(groups[k][i].Name) < strings.ToLower(groups[k][j].Name)
			})
		}
	}

	if existing != "" && IsManagedFile(existing) {
		return mergeServices(cfg, groups, existing)
	}

	return freshServices(cfg, groups)
}

func freshServices(cfg *config.Config, groups map[string][]docker.Container) string {
	var b strings.Builder
	b.WriteString(sentinel + "\n\n")

	for _, groupName := range []string{"Services", "Media"} {
		members, ok := groups[groupName]
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "- %s:\n", groupName)
		for _, c := range members {
			b.WriteString(renderContainer(cfg, c))
		}
	}

	return b.String()
}

func mergeServices(cfg *config.Config, groups map[string][]docker.Container, existing string) string {
	present := parseExistingContainers(existing)

	var appended strings.Builder
	added := 0

	for _, groupName := range []string{"Services", "Media"} {
		members, ok := groups[groupName]
		if !ok {
			continue
		}

		var newMembers []docker.Container
		for _, c := range members {
			if !present[strings.ToLower(c.Name)] {
				newMembers = append(newMembers, c)
			}
		}

		if len(newMembers) == 0 {
			continue
		}

		fmt.Fprintf(&appended, "- %s:\n", groupName)
		for _, c := range newMembers {
			appended.WriteString(renderContainer(cfg, c))
			added++
		}
	}

	if added == 0 {
		return existing
	}

	result := strings.TrimRight(existing, "\n") + "\n\n" + appended.String()
	return result
}

func parseExistingContainers(content string) map[string]bool {
	present := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(line, "    - ") && strings.HasSuffix(trimmed, ":") {
			name := strings.TrimSuffix(strings.TrimPrefix(line, "    - "), ":")
			present[strings.ToLower(strings.TrimSpace(name))] = true
		}
	}
	return present
}

func renderContainer(cfg *config.Config, c docker.Container) string {
	var b strings.Builder

	info := widgets.Lookup(c.Name)

	if label, ok := c.Labels["homepage.name"]; ok {
		info.DisplayName = label
	}
	if label, ok := c.Labels["homepage.icon"]; ok {
		info.Icon = label
	}
	if label, ok := c.Labels["homepage.description"]; ok {
		info.Description = label
	}

	webPort := uint16(0)
	if len(c.Ports) > 0 {
		webPort = c.Ports[0]
	}

	href := labelOr(c.Labels, "homepage.href", buildURL(cfg, webPort))
	pingURL := buildURL(cfg, webPort)

	fmt.Fprintf(&b, "    - %s:\n", info.DisplayName)
	fmt.Fprintf(&b, "        icon: %s\n", info.Icon)
	fmt.Fprintf(&b, "        href: %s\n", href)
	fmt.Fprintf(&b, "        ping: %s\n", pingURL)
	fmt.Fprintf(&b, "        statusStyle: \"%s\"\n", cfg.Homepage.StatusStyle)
	fmt.Fprintf(&b, "        description: %s\n", info.Description)

	if cfg.Homepage.ShowStats {
		fmt.Fprintf(&b, "        showStats: true\n")
	}

	widgetType := labelOr(c.Labels, "homepage.widget.type", info.WidgetType)
	if widgetType != "" {
		widgetURL := labelOr(c.Labels, "homepage.widget.url", pingURL)
		fmt.Fprintf(&b, "        widget:\n")
		fmt.Fprintf(&b, "            type: %s\n", widgetType)
		fmt.Fprintf(&b, "            url: %s\n", widgetURL)

		varName := strings.ToUpper(strings.ReplaceAll(c.Name, "-", "_"))
		keyPlaceholder := "{{HOMEPAGE_VAR_" + varName + "_KEY}}"
		userPlaceholder := "{{HOMEPAGE_VAR_" + varName + "_USERNAME}}"
		passPlaceholder := "{{HOMEPAGE_VAR_" + varName + "_PASSWORD}}"
		tokenPlaceholder := "{{HOMEPAGE_VAR_" + varName + "_TOKEN}}"

		switch info.Auth.Type {
		case widgets.AuthAPIKey:
			fmt.Fprintf(&b, "            key: %s\n", labelOr(c.Labels, "homepage.widget.key", keyPlaceholder))
		case widgets.AuthToken:
			fmt.Fprintf(&b, "            token: %s\n", labelOr(c.Labels, "homepage.widget.key", tokenPlaceholder))
		case widgets.AuthUserPass:
			fmt.Fprintf(&b, "            username: %s\n", userPlaceholder)
			fmt.Fprintf(&b, "            password: %s\n", passPlaceholder)
		}

		for k, v := range info.ExtraFields {
			fmt.Fprintf(&b, "            %s: %s\n", k, v)
		}
	}

	if len(c.Ports) > 1 {
		extras := make([]string, 0, len(c.Ports)-1)
		for _, p := range c.Ports[1:] {
			extras = append(extras, fmt.Sprintf("%d", p))
		}
		fmt.Fprintf(&b, "        # extra ports: %s\n", strings.Join(extras, ", "))
	}

	b.WriteString("\n")
	return b.String()
}

func Settings(cfg *config.Config) string {
	var b strings.Builder

	fmt.Fprintf(&b, "title: %s\n", cfg.Homepage.Title)
	fmt.Fprintf(&b, "theme: %s\n", cfg.Homepage.Theme)
	fmt.Fprintf(&b, "color: %s\n", cfg.Homepage.Color)
	fmt.Fprintf(&b, "headerStyle: %s\n", cfg.Homepage.HeaderStyle)
	fmt.Fprintf(&b, "target: %s\n", cfg.Homepage.Target)
	fmt.Fprintf(&b, "useEqualHeights: %v\n", cfg.Homepage.UseEqualHeights)

	if cfg.Homepage.MaxGroupColumns > 0 {
		fmt.Fprintf(&b, "maxGroupColumns: %d\n", cfg.Homepage.MaxGroupColumns)
	}
	if cfg.Homepage.ShowStats {
		b.WriteString("showStats: true\n")
	}

	b.WriteString("\nlayout:\n")

	for _, groupName := range []string{"Services", "Media"} {
		fmt.Fprintf(&b, "  %s:\n", groupName)
		if cfg.Homepage.Tabs {
			fmt.Fprintf(&b, "    tab: %s\n", groupName)
		}
		if layout, ok := cfg.Layout[groupName]; ok {
			if layout.Style != "" {
				fmt.Fprintf(&b, "    style: %s\n", layout.Style)
			}
			if layout.Columns > 0 {
				fmt.Fprintf(&b, "    columns: %d\n", layout.Columns)
			}
			if layout.Header != nil {
				fmt.Fprintf(&b, "    header: %v\n", *layout.Header)
			}
		}
	}

	return b.String()
}

func Env(cfg *config.Config, containers []docker.Container, existing string) string {
	existingVals := parseEnv(existing)

	var b strings.Builder

	writeKey := func(key, fallback string) {
		if val, ok := existingVals[key]; ok && val != "" {
			fmt.Fprintf(&b, "%s=%s\n", key, val)
		} else {
			fmt.Fprintf(&b, "%s=%s\n", key, fallback)
		}
	}

	hostKey := "HOMEPAGE_VAR_HOST_IP"
	writeKey(hostKey, cfg.Host.IP)
	b.WriteString("\n")

	for _, c := range containers {
		if shouldSkip(c.Name, cfg.Docker.Skip) {
			continue
		}

		info := widgets.Lookup(c.Name)
		if info.WidgetType == "" {
			continue
		}

		prefix := "HOMEPAGE_VAR_" + strings.ToUpper(strings.ReplaceAll(c.Name, "-", "_")) + "_"

		switch info.Auth.Type {
		case widgets.AuthAPIKey:
			writeKey(prefix+"KEY", "")
		case widgets.AuthToken:
			writeKey(prefix+"TOKEN", "")
		case widgets.AuthUserPass:
			writeKey(prefix+"USERNAME", "")
			writeKey(prefix+"PASSWORD", "")
		}
	}

	return b.String()
}

func parseEnv(content string) map[string]string {
	result := make(map[string]string)
	if content == "" {
		return result
	}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func GuessGroup(containerName string, cfg *config.Config) string {
	lower := strings.ToLower(containerName)
	for pattern, group := range cfg.Groups {
		if strings.Contains(lower, strings.ToLower(pattern)) {
			return group
		}
	}
	return widgets.Lookup(containerName).Group
}

func resolveGroup(c docker.Container, cfg *config.Config) string {
	if label, ok := c.Labels["homepage.group"]; ok {
		return label
	}
	return GuessGroup(c.Name, cfg)
}

func shouldSkip(name string, skipList []string) bool {
	lower := strings.ToLower(name)
	for _, s := range skipList {
		if strings.HasPrefix(lower, strings.ToLower(s)) || lower == strings.ToLower(s) {
			return true
		}
	}
	return false
}

func buildURL(cfg *config.Config, port uint16) string {
	if port == 0 {
		return fmt.Sprintf("%s://%s", cfg.Host.Protocol, cfg.Host.IP)
	}
	return fmt.Sprintf("%s://%s:%d", cfg.Host.Protocol, cfg.Host.IP, port)
}

func labelOr(labels map[string]string, key, fallback string) string {
	if v, ok := labels[key]; ok && v != "" {
		return v
	}
	return fallback
}
