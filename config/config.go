package config

import "os"

func GetVar(key, def string) string {
	result := os.Getenv(key)
	if result == "" {
		result = def
	}
	return result
}
