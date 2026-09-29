package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// LoadDotEnv loads optional dotenv files. Environment variables already set win.
func LoadDotEnv(files ...string) error {
	if len(files) == 0 {
		return godotenv.Load()
	}
	paths := make([]string, 0, len(files))
	for _, key := range files {
		if value := os.Getenv(key); value != "" {
			paths = append(paths, value)
		}
	}
	if len(paths) == 0 {
		return godotenv.Load()
	}
	return godotenv.Load(paths...)
}

// String returns the environment value or fallback when unset.
func String(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Prefixed returns the first configured value for key under the given prefixes.
func Prefixed(prefixes []string, key, fallback string) string {
	for _, prefix := range prefixes {
		if value := os.Getenv(prefix + key); value != "" {
			return value
		}
	}
	return fallback
}

// Int returns a parsed integer environment value or fallback on invalid input.
func Int(key string, fallback int) int {
	value, err := strconv.Atoi(String(key, strconv.Itoa(fallback)))
	if err != nil {
		return fallback
	}
	return value
}

// IntPrefixed returns a parsed prefixed integer or fallback on invalid input.
func IntPrefixed(prefixes []string, key string, fallback int) int {
	value, err := strconv.Atoi(Prefixed(prefixes, key, strconv.Itoa(fallback)))
	if err != nil {
		return fallback
	}
	return value
}

// Bool returns a parsed boolean environment value or fallback on invalid input.
func Bool(key string, fallback bool) bool {
	value, err := strconv.ParseBool(String(key, strconv.FormatBool(fallback)))
	if err != nil {
		return fallback
	}
	return value
}

// BoolPrefixed returns a parsed prefixed boolean or fallback on invalid input.
func BoolPrefixed(prefixes []string, key string, fallback bool) bool {
	value, err := strconv.ParseBool(Prefixed(prefixes, key, strconv.FormatBool(fallback)))
	if err != nil {
		return fallback
	}
	return value
}

// Duration returns a parsed duration environment value or fallback on invalid input.
func Duration(key string, fallback time.Duration) time.Duration {
	return ParseDuration(String(key, fallback.String()), fallback)
}

// DurationPrefixed returns a parsed prefixed duration or fallback on invalid input.
func DurationPrefixed(prefixes []string, key string, fallback time.Duration) time.Duration {
	return ParseDuration(Prefixed(prefixes, key, fallback.String()), fallback)
}

// Slice splits a comma-separated environment value into trimmed entries.
func Slice(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil
	}
	values := make([]string, 0)
	for _, value := range strings.Split(raw, ",") {
		if value = strings.TrimSpace(value); value != "" {
			values = append(values, value)
		}
	}
	return values
}

// ParseDuration parses Go or HH:MM:SS duration syntax, with a fallback on failure.
func ParseDuration(raw string, fallback time.Duration) time.Duration {
	if value, err := time.ParseDuration(raw); err == nil {
		return value
	}
	var h, m, s int
	if _, err := fmt.Sscanf(raw, "%d:%d:%d", &h, &m, &s); err == nil {
		return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second
	}
	return fallback
}
