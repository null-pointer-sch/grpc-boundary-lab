// Package envutil provides helpers for reading environment variables with fallbacks.
package envutil

import (
	"log"
	"os"
	"strconv"
)

// Config represents the application configuration.
type Config struct {
	BackendPort        string
	BackendPortTLS     string
	BackendRESTPort    string
	BackendRESTPortTLS string
	CertDir            string
	GatewayPort        string
	GatewayRESTPort    string
	BackendHost        string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		BackendPort:        GetOrDefault("BACKEND_PORT", "50051"),
		BackendPortTLS:     GetOrDefault("BACKEND_PORT_TLS", "50151"),
		BackendRESTPort:    GetOrDefault("BACKEND_REST_PORT", "8081"),
		BackendRESTPortTLS: GetOrDefault("BACKEND_REST_PORT_TLS", "8181"),
		CertDir:            GetOrDefault("CERT_DIR", "/certs"),
		GatewayPort:        GetOrDefault("GATEWAY_PORT", "50052"),
		GatewayRESTPort:    GetOrDefault("REST_PORT", "8080"),
		BackendHost:        GetOrDefault("BACKEND_HOST", "127.0.0.1"),
	}
}

// GetOrDefault returns the value of the environment variable named by key,
// or fallback if the variable is not set or empty.
func GetOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// GetInt returns the integer value of the environment variable named by key,
// or fallback if the variable is not set. It terminates the process on parse errors.
func GetInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("envutil: invalid integer for %s: %v", key, err)
	}
	return n
}

// GetInt64 returns the int64 value of the environment variable named by key,
// or fallback if the variable is not set. It terminates the process on parse errors.
func GetInt64(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatalf("envutil: invalid integer for %s: %v", key, err)
	}
	return n
}
