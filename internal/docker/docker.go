package docker

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
)

type Container struct {
	ID     string
	Name   string
	Image  string
	Ports  []uint16
	Labels map[string]string
}

type apiContainer struct {
	ID     string            `json:"Id"`
	Names  []string          `json:"Names"`
	Image  string            `json:"Image"`
	Ports  []apiPort         `json:"Ports"`
	Labels map[string]string `json:"Labels"`
}

type apiPort struct {
	PublicPort uint16 `json:"PublicPort"`
	Type       string `json:"Type"`
}

func ListContainers(socketPath string) ([]Container, error) {
	transport := &http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get("http://localhost/containers/json")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker socket at %q: %w", socketPath, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw []apiContainer
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse Docker API response: %w", err)
	}

	var containers []Container
	for _, c := range raw {
		name := c.ID[:12]
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		portSet := make(map[uint16]struct{})
		for _, p := range c.Ports {
			if p.PublicPort > 0 && p.Type == "tcp" {
				portSet[p.PublicPort] = struct{}{}
			}
		}
		ports := make([]uint16, 0, len(portSet))
		for p := range portSet {
			ports = append(ports, p)
		}
		sort.Slice(ports, func(i, j int) bool { return ports[i] < ports[j] })

		containers = append(containers, Container{
			ID:     c.ID[:12],
			Name:   name,
			Image:  c.Image,
			Ports:  ports,
			Labels: c.Labels,
		})
	}

	return containers, nil
}
