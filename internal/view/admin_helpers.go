package view

import (
	"time"

	"github.com/flintcraftstudio/flint-ui/components/badge"
)

// statusMeta pairs a human label with the flint-ui badge variant for a status.
type statusMeta struct {
	Label   string
	Variant badge.Variant
}

// Status vocabularies, in workflow order (drives select options + filter tabs).
var (
	bookingStatusOrder  = []string{"new", "scheduled", "picked_up", "in_progress", "ready", "delivered", "canceled"}
	questionStatusOrder = []string{"new", "handled", "archived"}
	waitlistStatusOrder = []string{"new", "contacted", "archived"}
)

var bookingStatusMeta = map[string]statusMeta{
	"new":         {"New", badge.VariantWarning},
	"scheduled":   {"Scheduled", badge.VariantPrimary},
	"picked_up":   {"Picked up", badge.VariantPrimary},
	"in_progress": {"In progress", badge.VariantPrimary},
	"ready":       {"Ready", badge.VariantAccent},
	"delivered":   {"Delivered", badge.VariantSuccess},
	"canceled":    {"Canceled", badge.VariantMuted},
}

var questionStatusMeta = map[string]statusMeta{
	"new":      {"New", badge.VariantWarning},
	"handled":  {"Handled", badge.VariantSuccess},
	"archived": {"Archived", badge.VariantMuted},
}

var waitlistStatusMeta = map[string]statusMeta{
	"new":       {"New", badge.VariantWarning},
	"contacted": {"Contacted", badge.VariantSuccess},
	"archived":  {"Archived", badge.VariantMuted},
}

func metaOr(m map[string]statusMeta, s string) statusMeta {
	if v, ok := m[s]; ok {
		return v
	}
	return statusMeta{Label: s, Variant: badge.VariantMuted}
}

func bookingStatusLabel(s string) string          { return metaOr(bookingStatusMeta, s).Label }
func bookingStatusVariant(s string) badge.Variant { return metaOr(bookingStatusMeta, s).Variant }
func questionStatusLabel(s string) string          { return metaOr(questionStatusMeta, s).Label }
func questionStatusVariant(s string) badge.Variant { return metaOr(questionStatusMeta, s).Variant }
func waitlistStatusLabel(s string) string          { return metaOr(waitlistStatusMeta, s).Label }
func waitlistStatusVariant(s string) badge.Variant { return metaOr(waitlistStatusMeta, s).Variant }

// humanDate renders a date like "Jun 5, 2026".
func humanDate(t time.Time) string { return t.Format("Jan 2, 2006") }

// humanDateTime renders a timestamp like "Jun 5 · 3:04 PM".
func humanDateTime(t time.Time) string { return t.Format("Jan 2 · 3:04 PM") }

// orDash renders "—" for an empty string.
func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

// truncate shortens s to n runes, adding an ellipsis when cut.
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

// scheduledLabel renders a booking's scheduled date (+ window) or a dash.
func scheduledLabel(date, window string) string {
	if date == "" {
		return "—"
	}
	if window != "" {
		return date + " · " + window
	}
	return date
}

// AdminCounts holds the "new item" tallies shown on the dashboard.
type AdminCounts struct {
	Bookings  int
	Inquiries int
	Waitlist  int
}
