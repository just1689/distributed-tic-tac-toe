package config

import "os"

func GetVar(key, def string) string {
	result := def
	if result == "" {
		result = os.Getenv(key)
	}
	return result
}
