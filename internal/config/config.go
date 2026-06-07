package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// The pickup service area (allowed ZIPs) now lives in the settings table, edited
// from the admin dashboard — it is no longer an env var.

type Config struct {
	Port           int
	PostmarkToken  string
	PostmarkFrom   string
	PostmarkTo     string
	PixelID            string
	GtagID             string
	TurnstileSiteKey   string
	TurnstileSecretKey string
	DBPath             string
	SessionSecret      string
	// NoIndex, when true, asks search engines not to crawl or index the site —
	// for a testing/staging deploy that shouldn't surface in results yet. Drives
	// robots.txt, a robots meta tag, and an X-Robots-Tag header. Default false.
	NoIndex bool
}

// Load reads configuration from environment variables, applying defaults where not set.
func Load() (*Config, error) {
	port, err := parseInt("PORT", 8080)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:          port,
		PostmarkToken: os.Getenv("POSTMARK_SERVER_TOKEN"),
		PostmarkFrom:  os.Getenv("POSTMARK_FROM"),
		PostmarkTo:    os.Getenv("POSTMARK_TO"),
		PixelID:            os.Getenv("PIXEL_ID"),
		GtagID:             os.Getenv("GTAG_ID"),
		TurnstileSiteKey:   os.Getenv("TURNSTILE_SITE_KEY"),
		TurnstileSecretKey: os.Getenv("TURNSTILE_SECRET_KEY"),
		DBPath:             envDefault("DB_PATH", "./data/app.db"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		NoIndex:            parseBool("SITE_NOINDEX"),
	}, nil
}

// parseBool reads an environment variable as a boolean flag. Accepts the usual
// truthy spellings (1, true, yes, on); anything else — including unset — is false.
func parseBool(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// Addr returns the server address string in the format expected by http.ListenAndServe.
func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

// envDefault reads an environment variable, returning the fallback if unset.
func envDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// parseInt reads an environment variable as an integer, returning the fallback if unset.
func parseInt(key string, fallback int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %q", key, val)
	}
	return n, nil
}