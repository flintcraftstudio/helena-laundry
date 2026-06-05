package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// defaultAllowedZips is Chanté's pickup service area — the Helena Valley.
// Bookings with a ZIP outside this set are routed to the out-of-area waitlist.
// Override with the ALLOWED_ZIPS env var (comma- or space-separated).
const defaultAllowedZips = "59601 59602 59634 59635"

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
	// AllowedZips is the set of ZIP codes Chanté picks up from.
	AllowedZips []string
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
		AllowedZips:        parseZips(envDefault("ALLOWED_ZIPS", defaultAllowedZips)),
	}, nil
}

// parseZips splits a comma/space/newline-separated list into trimmed,
// non-empty ZIP strings.
func parseZips(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == '\t' || r == ';'
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if f = strings.TrimSpace(f); f != "" {
			out = append(out, f)
		}
	}
	return out
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