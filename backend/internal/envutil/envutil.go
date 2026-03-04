// Package envutil provides helpers for reading environment variables with fallbacks.
package envutil

import (
	"log"
	"os"
	"strconv"
)

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
