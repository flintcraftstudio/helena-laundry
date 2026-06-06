// Package settings holds an in-memory snapshot of the site settings row so the
// booking flow and templates can read current config without a DB hit per
// request. Loaded once at startup and refreshed whenever the operator saves.
package settings

import (
	"sync"

	"github.com/firefly-software-mt/advanced-template/internal/store"
)

var (
	mu      sync.RWMutex
	current store.Settings
)

// Set replaces the cached settings (called at startup and after a save).
func Set(s store.Settings) {
	mu.Lock()
	current = s
	mu.Unlock()
}

// Get returns the current settings snapshot.
func Get() store.Settings {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

// IsZipAllowed reports whether zip is inside the current service area.
func IsZipAllowed(zip string) bool {
	for _, z := range Get().AllowedZips {
		if z == zip {
			return true
		}
	}
	return false
}

// IsPickupDay reports whether day is an offered pickup weekday.
func IsPickupDay(day string) bool {
	for _, d := range Get().PickupDays {
		if d == day {
			return true
		}
	}
	return false
}

// DefaultPickupDay is the fallback day when none/an invalid one is submitted.
func DefaultPickupDay() string {
	if days := Get().PickupDays; len(days) > 0 {
		return days[0]
	}
	return "Tuesday"
}
