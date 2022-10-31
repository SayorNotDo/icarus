package utils

import "os"

// GetEnv Run if docker then we have to get from env
func GetEnv(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
