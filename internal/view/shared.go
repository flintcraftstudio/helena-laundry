package view

import (
	"strings"
	"time"

	"github.com/firefly-software-mt/advanced-template/internal/store"
)

// Site/brand values used throughout the templates. Seeded with the original
// defaults and overwritten from the settings row at startup and on every save
// (ApplySettings), so editing settings in the dashboard updates the public site.
var (
	// SiteName is the display name used in titles and the footer.
	SiteName = "Helena Laundry"
	// Tagline appears in the footer sign-off line.
	Tagline = "Helena, Montana · Serving the Helena Valley"
	// Phone is the single call/text number — the primary CTA everywhere.
	Phone = "406-471-8508"
	// PhoneTel / PhoneSMS are the tel:/sms: hrefs derived from Phone.
	PhoneTel = "tel:+14064718508"
	PhoneSMS = "sms:+14064718508"
	// Email and Hours are optional identity fields (blank by default).
	Email string
	Hours string
	// Turnaround is the standard turnaround phrase, e.g. "24–48 hours".
	Turnaround = "24–48 hours"
	// Pricing display values (dollars, no $); the $ is added at the use site.
	LbRate       = "2.00"
	RushRate     = "3.00"
	MinimumOrder = "35"
	// AcceptingBookings gates the public "Book a pickup" path.
	AcceptingBookings = true
)

// Tracking IDs and Turnstile site key, set once at startup from config.
var (
	PixelID          string
	GtagID           string
	TurnstileSiteKey string
)

// NoIndex, when true, renders a robots noindex/nofollow meta tag on every public
// page — set once at startup from config for testing/staging deploys.
var NoIndex bool

// ApplySettings copies the settings row into the package-level view vars the
// templates read. Safe to call at startup and after each settings save.
func ApplySettings(s store.Settings) {
	SiteName = s.SiteName
	Tagline = s.Tagline
	Phone = s.Phone
	PhoneTel = "tel:" + telTarget(s.Phone)
	PhoneSMS = "sms:" + telTarget(s.Phone)
	Email = s.Email
	Hours = s.Hours
	Turnaround = s.Turnaround
	LbRate = s.LbRate
	RushRate = s.RushRate
	MinimumOrder = s.MinimumOrder
	AcceptingBookings = s.AcceptingBookings
}

// telTarget builds a US tel/sms target ("+1" + 10 digits) from a display phone.
func telTarget(phone string) string {
	var digits strings.Builder
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	d := digits.String()
	d = strings.TrimPrefix(d, "1") // drop a leading country-code 1 if present
	return "+1" + d
}

// Year returns the current year for copyright notices.
func Year() int {
	return time.Now().Year()
}
