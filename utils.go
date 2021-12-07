package main

import "os"

func GetEnvOrDefault(env, defaultVal string) string {
	e := os.Getenv(env)
	if env == "" {
		return defaultVal
	} else {
		return e
	}
}
