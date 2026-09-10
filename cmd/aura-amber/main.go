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
	version     = "0.3.0"
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
	case "test":
		cmdTest()
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
  start                          Start all services (Docker required)
  stop                           Stop all services
  status                         Show status of all services
  open                           Open dashboard in default browser
  migrate <type> [flags]         Start a migration
  test                           Run self-test (no real credentials needed)
  version                        Show version
  help                           Show this help

Migration flags:
  --url <url>                    Source base URL
  --user <username>              Username / consumer key
  --pass <password>              Password / consumer secret / access token
  --target <db-url>              Target PostgreSQL URL (default: local)

Sources:
  wordpress      --url https://my-site.com --user admin --pass xxxx
  shopify        --url my-shop.myshopify.com --pass shpat_xxxx
  woocommerce    --url https://my-site.com --user ck_xxxx --pass cs_xxxx
  magento        --url https://my-site.com --pass bpat_xxxx

Examples:
  aura-amber start
  aura-amber status
  aura-amber open
  aura-amber migrate wordpress --url https://my-site.com --user admin --pass xxxx
  aura-amber test

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
		fmt.Fprintln(os.Stderr, "Usage: aura-amber migrate <wordpress|shopify|woocommerce|magento> [flags]")
		fmt.Fprintln(os.Stderr, "Run 'aura-amber help' for details.")
		os.Exit(1)
	}

	source := os.Args[2]
	valid := map[string]bool{"wordpress": true, "shopify": true, "woocommerce": true, "magento": true}
	if !valid[source] {
		fmt.Fprintf(os.Stderr, "Invalid source: %s\nValid sources: wordpress, shopify, woocommerce, magento\n", source)
		os.Exit(1)
	}

	// Parse flags
	flags := parseFlags()
	connInfo := buildConnectionInfo(source, flags)

	// Validate required fields
	if err := validateConnection(source, connInfo); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Required fields for %s:\n", source)
		printRequiredFields(source)
		os.Exit(1)
	}

	targetDb := flags["target"]
	if targetDb == "" {
		targetDb = "postgres://aura:aura@localhost:5432/aura_amber?sslmode=disable"
	}

	fmt.Printf("Starting %s migration...\n", source)
	payload := map[string]interface{}{
		"name":                 fmt.Sprintf("%s-migration-%d", source, time.Now().UnixMilli()),
		"sourceType":           source,
		"sourceConnectionInfo": connInfo,
		"targetDbUrl":          targetDb,
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

func cmdTest() {
	fmt.Println("Running self-test...")
	fmt.Println()

	// Test 1: API health
	fmt.Print("  API health... ")
	if checkHealth(defaultAPI + "/health") {
		fmt.Println("✓ OK")
	} else {
		fmt.Println("✗ DOWN")
	}

	// Test 2: Worker health
	fmt.Print("  Worker health... ")
	if checkHealth(defaultWork + "/health") {
		fmt.Println("✓ OK")
	} else {
		fmt.Println("✗ DOWN")
	}

	// Test 3: Database
	fmt.Print("  PostgreSQL... ")
	if portOpen("5432") {
		fmt.Println("✓ OK")
	} else {
		fmt.Println("✗ DOWN")
	}

	// Test 4: Redis
	fmt.Print("  Redis... ")
	if portOpen("6379") {
		fmt.Println("✓ OK")
	} else {
		fmt.Println("✗ DOWN")
	}

	// Test 5: List migrations
	fmt.Print("  GET /api/migrations... ")
	resp, err := http.Get(defaultAPI + "/api/migrations")
	if err == nil && resp.StatusCode == 200 {
		fmt.Println("✓ OK")
		resp.Body.Close()
	} else {
		fmt.Println("✗ FAIL")
	}

	// Test 6: List connectors
	fmt.Print("  GET /api/connectors... ")
	resp2, err := http.Get(defaultAPI + "/api/connectors")
	if err == nil && resp2.StatusCode == 200 {
		var connectors []map[string]interface{}
		json.NewDecoder(resp2.Body).Decode(&connectors)
		resp2.Body.Close()
		fmt.Printf("✓ OK (%d connectors)\n", len(connectors))
	} else {
		fmt.Println("✗ FAIL")
	}

	// Test 7: Create + immediately check a test migration
	fmt.Print("  POST /api/migrations (dry)... ")
	testPayload := map[string]interface{}{
		"name":                 fmt.Sprintf("test-%d", time.Now().UnixMilli()),
		"sourceType":           "wordpress",
		"sourceConnectionInfo": map[string]string{"baseUrl": "http://test.example.com", "username": "test", "appPassword": "test"},
		"targetDbUrl":          "postgres://aura:aura@localhost:5432/aura_amber?sslmode=disable",
	}
	testBody, _ := json.Marshal(testPayload)
	resp3, err := http.Post(defaultAPI+"/api/migrations", "application/json", bytes.NewReader(testBody))
	if err == nil && resp3.StatusCode == 201 {
		fmt.Println("✓ OK (migration created, will fail at worker — expected)")
		resp3.Body.Close()
	} else if err == nil {
		fmt.Printf("✗ FAIL (status %d)\n", resp3.StatusCode)
		resp3.Body.Close()
	} else {
		fmt.Println("✗ FAIL")
	}

	fmt.Println()
	fmt.Println("Self-test complete.")
}

func parseFlags() map[string]string {
	flags := map[string]string{}
	for i := 3; i < len(os.Args)-1; i++ {
		if strings.HasPrefix(os.Args[i], "--") {
			key := strings.TrimPrefix(os.Args[i], "--")
			flags[key] = os.Args[i+1]
			i++
		}
	}
	return flags
}

func buildConnectionInfo(source string, flags map[string]string) map[string]string {
	info := map[string]string{}
	url := flags["url"]
	user := flags["user"]
	pass := flags["pass"]

	switch source {
	case "wordpress":
		info["baseUrl"] = url
		info["username"] = user
		info["appPassword"] = pass
	case "shopify":
		info["shopDomain"] = url
		info["accessToken"] = pass
	case "woocommerce":
		info["baseUrl"] = url
		info["consumerKey"] = user
		info["consumerSecret"] = pass
	case "magento":
		info["baseUrl"] = url
		info["accessToken"] = pass
	}
	return info
}

func validateConnection(source string, info map[string]string) error {
	switch source {
	case "wordpress":
		if info["baseUrl"] == "" || info["username"] == "" || info["appPassword"] == "" {
			return fmt.Errorf("wordpress requires --url, --user, --pass")
		}
	case "shopify":
		if info["shopDomain"] == "" || info["accessToken"] == "" {
			return fmt.Errorf("shopify requires --url (shop domain) and --pass (access token)")
		}
	case "woocommerce":
		if info["baseUrl"] == "" || info["consumerKey"] == "" || info["consumerSecret"] == "" {
			return fmt.Errorf("woocommerce requires --url, --user (consumer key), --pass (consumer secret)")
		}
	case "magento":
		if info["baseUrl"] == "" || info["accessToken"] == "" {
			return fmt.Errorf("magento requires --url and --pass (access token)")
		}
	}
	return nil
}

func printRequiredFields(source string) {
	switch source {
	case "wordpress":
		fmt.Println("  --url      WordPress site URL (e.g. https://my-site.com)")
		fmt.Println("  --user     WordPress username")
		fmt.Println("  --pass     Application password")
	case "shopify":
		fmt.Println("  --url      Shop domain (e.g. my-shop.myshopify.com)")
		fmt.Println("  --pass     Admin API access token")
	case "woocommerce":
		fmt.Println("  --url      WooCommerce site URL (e.g. https://my-site.com)")
		fmt.Println("  --user     Consumer key")
		fmt.Println("  --pass     Consumer secret")
	case "magento":
		fmt.Println("  --url      Magento site URL (e.g. https://my-site.com)")
		fmt.Println("  --pass     Integration access token")
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

func init() {
	if _, err := os.Stat("docker-compose.yml"); err != nil {
		if _, err := os.Stat("../docker-compose.yml"); err == nil {
			os.Chdir("..")
		}
	}
}
