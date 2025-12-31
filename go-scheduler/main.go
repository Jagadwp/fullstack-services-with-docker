package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
User represents the JSON structure returned by PHP API
*/
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

/*
Configuration loaded from environment variables
*/
type Config struct {
	PhpAPIBaseURL    string
	PythonAPIBaseURL string
	Interval         time.Duration
	DataDir          string
	IncomingDir      string
	ProcessedDir     string
}

/*
Load configuration from environment variables.
Fail fast if required config is missing.
*/
func loadConfig() Config {
	phpURL := os.Getenv("PHP_API_BASE_URL")
	if phpURL == "" {
		log.Fatal("PHP_API_BASE_URL is required")
	}

	pythonURL := os.Getenv("PYTHON_API_BASE_URL")
	if pythonURL == "" {
		log.Fatal("PYTHON_API_BASE_URL is required")
	}

	intervalStr := os.Getenv("SCHEDULER_INTERVAL_SECONDS")
	if intervalStr == "" {
		intervalStr = "10"
	}

	intervalSeconds, err := time.ParseDuration(intervalStr + "s")
	if err != nil {
		log.Fatalf("Invalid SCHEDULER_INTERVAL_SECONDS: %v", err)
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/data"
	}

	return Config{
		PhpAPIBaseURL:    phpURL,
		PythonAPIBaseURL: pythonURL,
		Interval:         intervalSeconds,
		DataDir:          dataDir,
		IncomingDir:      filepath.Join(dataDir, "incoming"),
		ProcessedDir:     filepath.Join(dataDir, "processed"),
	}
}

/*
Ensure required directories exist.
*/
func ensureDirectories(cfg Config) {
	dirs := []string{cfg.IncomingDir, cfg.ProcessedDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
}

/*
Generate sample user data.
This simulates incoming traffic to the PHP API.
*/
func generateUserData(counter int) map[string]string {
	names := []string{
		"David Smith",
		"Alice Brown",
		"David Lee",
		"John Doe",
		"David De Gea",
	}

	// Counter % len(names) ensures the index is always within bounds
	name := names[counter%len(names)]
	email := strings.ToLower(strings.ReplaceAll(name, " ", ".")) + "@example.com"

	return map[string]string{
		"name":  name,
		"email": email,
	}
}

/*
Call PHP API to create a user.
*/
func createUser(cfg Config, payload map[string]string) (*User, error) {
	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		cfg.PhpAPIBaseURL+"/users",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call PHP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("PHP API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode PHP API response: %w", err)
	}

	return &user, nil
}

/*
Write user JSON to file in incoming directory.
*/
func writeUserToFile(cfg Config, user *User) error {
	fileName := fmt.Sprintf("user_%d.json", user.ID)
	filePath := filepath.Join(cfg.IncomingDir, fileName)

	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

/*
Process incoming files:
- Read JSON
- If name starts with "David", forward to Python API
- Move file to processed directory
*/
func processFiles(cfg Config) {
	files, err := os.ReadDir(cfg.IncomingDir)
	if err != nil {
		log.Printf("Failed to read incoming directory: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		inPath := filepath.Join(cfg.IncomingDir, file.Name())
		outPath := filepath.Join(cfg.ProcessedDir, file.Name())

		data, err := os.ReadFile(inPath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", file.Name(), err)
			continue
		}

		var user User
		if err := json.Unmarshal(data, &user); err != nil {
			log.Printf("Invalid JSON in file %s: %v", file.Name(), err)
			_ = os.Rename(inPath, outPath)
			continue
		}

		if strings.HasPrefix(user.Name, "David") {
			log.Printf("Forwarding user %d (%s) to Python API", user.ID, user.Name)
			sendToPython(cfg, data)
		}

		if err := os.Rename(inPath, outPath); err != nil {
			log.Printf("Failed to move file %s: %v", file.Name(), err)
		}
	}
}

/*
Send user data to Python API.
This is best-effort delivery.
*/
func sendToPython(cfg Config, data []byte) {
	resp, err := http.Post(
		cfg.PythonAPIBaseURL+"/process",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Printf("Failed to send data to Python API: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Python API responded with status %d", resp.StatusCode)
}

func main() {
	log.Println("Starting Go Scheduler Service")

	cfg := loadConfig()
	ensureDirectories(cfg)

	counter := 0

	for {
		log.Println("Scheduler iteration started")

		payload := generateUserData(counter)
		user, err := createUser(cfg, payload)
		if err != nil {
			log.Printf("Error creating user: %v", err)
		} else {
			log.Printf("Created user ID=%d Name=%s", user.ID, user.Name)
			if err := writeUserToFile(cfg, user); err != nil {
				log.Printf("Failed to write user to file: %v", err)
			}
		}

		processFiles(cfg)

		counter++
		log.Printf("Sleeping for %s", cfg.Interval)
		time.Sleep(cfg.Interval)
	}
}
