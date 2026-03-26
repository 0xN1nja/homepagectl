package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/0xN1nja/homepagectl/internal/config"
	"github.com/0xN1nja/homepagectl/internal/docker"
	"github.com/0xN1nja/homepagectl/internal/generator"
	"github.com/0xN1nja/homepagectl/internal/widgets"
	"github.com/spf13/cobra"
)

var cfgPath string

var root = &cobra.Command{
	Use:   "homepagectl",
	Short: "Generate gethomepage.dev configs from running Docker containers",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Print an example homepagectl.toml to stdout",
	Long:  "Print an example config to stdout.\n\nWrite to file:\n  homepagectl init | tee homepagectl.toml\n  homepagectl init | sudo tee /etc/homepagectl.toml",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(config.Example())
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List running containers with detected widget and group",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}

		containers, err := docker.ListContainers(cfg.Docker.Socket)
		if err != nil {
			return err
		}

		if len(containers) == 0 {
			fmt.Println("No running containers found.")
			return nil
		}

		fmt.Printf("%-30s %-25s %-10s %-10s\n", "CONTAINER", "WIDGET", "PORT", "GROUP")
		fmt.Println(repeat("-", 80))

		for _, c := range containers {
			info := widgets.Lookup(c.Name)
			port := "-"
			if len(c.Ports) > 0 {
				port = fmt.Sprintf("%d", c.Ports[0])
			}
			widget := info.WidgetType
			if widget == "" {
				widget = "(unknown)"
			}
			group := generator.GuessGroup(c.Name, cfg)
			fmt.Printf("%-30s %-25s %-10s %-10s\n", c.Name, widget, port, group)
		}

		return nil
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate services.yaml, settings.yaml and .env",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}

		containers, err := docker.ListContainers(cfg.Docker.Socket)
		if err != nil {
			return err
		}

		if len(containers) == 0 {
			fmt.Println("No running containers found.")
			return nil
		}

		fmt.Printf("Found %d running container(s).\n", len(containers))

		servicesYAML := generator.Services(cfg, containers)
		settingsYAML := generator.Settings(cfg)
		envFile := generator.Env(cfg, containers)

		if dryRun {
			fmt.Println("\n========== services.yaml ==========\n")
			fmt.Println(servicesYAML)
			fmt.Println("\n========== settings.yaml ==========\n")
			fmt.Println(settingsYAML)
			fmt.Println("\n========== .env ==========\n")
			fmt.Println(envFile)
			return nil
		}

		files := map[string]string{
			"services.yaml": servicesYAML,
			"settings.yaml": settingsYAML,
			".env":          envFile,
		}

		for name, content := range files {
			path := filepath.Join(output, name)
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return err
			}
			fmt.Println("Written:", path)
		}

		return nil
	},
}

func repeat(s string, n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = s[0]
	}
	return string(out)
}

func init() {
	root.PersistentFlags().StringVarP(&cfgPath, "config", "c", "homepagectl.toml", "path to homepagectl.toml")
	generateCmd.Flags().StringP("output", "o", ".", "output directory for generated files")
	generateCmd.Flags().BoolP("dry-run", "d", false, "print to stdout instead of writing files")
	root.AddCommand(initCmd, listCmd, generateCmd)
}

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
