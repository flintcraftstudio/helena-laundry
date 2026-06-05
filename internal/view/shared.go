package view

import "time"

// Brand constants for Helena Laundry.
const (
	// SiteName is the display name used in titles and the footer.
	SiteName = "Helena Laundry"
	// Tagline appears in the footer sign-off line.
	Tagline = "Helena, Montana · Serving the Helena Valley"
	// Phone is the single call/text number — the primary CTA everywhere.
	Phone = "406-471-8508"
	// PhoneTel is the tel: href form of Phone (E.164, US).
	PhoneTel = "tel:+14064718508"
	// PhoneSMS is the sms: href form of Phone.
	PhoneSMS = "sms:+14064718508"
)

// Tracking IDs and Turnstile site key, set once at startup from config.
var (
	PixelID          string
	GtagID           string
	TurnstileSiteKey string
)

// Year returns the current year for copyright notices.
func Year() int {
	return time.Now().Year()
}
