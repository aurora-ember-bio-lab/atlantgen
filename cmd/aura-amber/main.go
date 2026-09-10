package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	version     = "0.2.0"
	appName     = "Aura Amber"
	defaultAPI  = "http://localhost:8080"
	defaultUI   = "http://localhost:3000"
	defaultWork = "http://localhost:8090"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "start":
		cmdStart()
	case "stop":
		cmdStop()
	case "status":
		cmdStatus()
	case "open":
		cmdOpen()
	case "migrate":
		cmdMigrate()
	case "version":
		fmt.Printf("aura-amber %s\n", version)
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stdout, `%s v%s — Migration SaaS CLI

Usage:
  aura-amber <command> [options]

Commands:
  start           Start all services (Docker required)
  stop            Stop all services
  status          Show status of all services
  open            Open dashboard in default browser
  migrate <type>  Start a migration (wordpress|shopify|woocommerce|magento)
  version         Show version
  help            Show this help

Examples:
  aura-amber start
  aura-amber status
  aura-amber open
  aura-amber migrate wordpress

`, appName, version)
}

func cmdStart() {
	fmt.Println("Starting Aura Amber services...")
	dockerCompose("up", "-d", "--build")
	time.Sleep(3 * time.Second)
	cmdStatus()
	fmt.Println("\nDashboard:  http://localhost:3000")
	fmt.Println("API:        http://localhost:8080")
	fmt.Println("Worker:     http://localhost:8090")
}

func cmdStop() {
	fmt.Println("Stopping Aura Amber services...")
	dockerCompose("down")
	fmt.Println("All services stopped.")
}

func cmdStatus() {
	fmt.Println("Service Status:")
	fmt.Println(strings.Repeat("─", 50))

	services := []struct {
		name string
		url  string
		port string
	}{
		{"API", defaultAPI + "/health", "8080"},
		{"Worker", defaultWork + "/health", "8090"},
		{"PostgreSQL", "", "5432"},
		{"Redis", "", "6379"},
		{"Dashboard", defaultUI, "3000"},
	}

	for _, s := range services {
		status := "● DOWN"
		if s.url != "" {
			if checkHealth(s.url) {
				status = "● UP"
			}
		} else {
			if portOpen(s.port) {
				status = "● UP"
			}
		}
		fmt.Printf("  %-12s %s  :%s\n", s.name, status, s.port)
	}
	fmt.Println(strings.Repeat("─", 50))
}

func cmdOpen() {
	url := defaultUI
	fmt.Printf("Opening %s in browser...\n", url)
	openBrowser(url)
}

func cmdMigrate() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: aura-amber migrate <wordpress|shopify|woocommerce|magento>")
		os.Exit(1)
	}
	source := os.Args[2]
	valid := map[string]bool{"wordpress": true, "shopify": true, "woocommerce": true, "magento": true}
	if !valid[source] {
		fmt.Fprintf(os.Stderr, "Invalid source: %s. Use: wordpress, shopify, woocommerce, magento\n", source)
		os.Exit(1)
	}

	fmt.Printf("Starting %s migration...\n", source)
	payload := map[string]interface{}{
		"name":                 fmt.Sprintf("%s-migration-%d", source, time.Now().UnixMilli()),
		"sourceType":           source,
		"sourceConnectionInfo": map[string]string{},
		"targetDbUrl":          "postgres://aura:aura@localhost:5432/aura_amber?sslmode=disable",
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(defaultAPI+"/api/migrations", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to API at %s: %v\n", defaultAPI, err)
		fmt.Println("Make sure services are running: aura-amber start")
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == 201 {
		fmt.Printf("Migration created: %s (status: %v)\n", result["id"], result["status"])
		fmt.Printf("Worker job: %v\n", result["workerJobId"])
		fmt.Println("Track progress: http://localhost:3000/dashboard/migrations")
	} else {
		fmt.Fprintf(os.Stderr, "Migration failed: %v\n", result["error"])
		os.Exit(1)
	}
}

func dockerCompose(args ...string) {
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func checkHealth(url string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func portOpen(port string) bool {
	cmd := exec.Command("netstat", "-an")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":"+port) && strings.Contains(line, "LISTENING") {
			return true
		}
	}
	return false
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}

// init project directory detection
func init() {
	// Try to find docker-compose.yml in current or parent dirs
	if _, err := os.Stat("docker-compose.yml"); err != nil {
		if _, err := os.Stat("../docker-compose.yml"); err == nil {
			os.Chdir("..")
		}
	}
}
