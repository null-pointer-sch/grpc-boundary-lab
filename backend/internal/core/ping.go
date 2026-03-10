package core

import "fmt"

// GeneratePong returns the standard pong message.
func GeneratePong(message string) string {
	if message == "" {
		message = "ping from frontend"
	}
	return fmt.Sprintf("pong: %s", message)
}
