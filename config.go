package main

import "os"

const (
	TRAFFIC_PATH     = "/traffic"
	CONNECTIONS_PATH = "/connections"
	TRACING_PATH     = "/profile/tracing"
)

var (
	CLASH_HOST  = getEnvOrDefault("CLASH_HOST", "127.0.0.1:9090")
	CLASH_TOKEN = getEnvOrDefault("CLASH_TOKEN", "clash")
)

func getEnvOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
