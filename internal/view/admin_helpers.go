package view

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/firefly-software-mt/advanced-template/internal/store"
	"github.com/flintcraftstudio/flint-ui/components/badge"
)

// BookingsView is everything the bookings section needs to render: the active
// status tab, the search term, the current page, and the page of results plus the
// total match count for the pager.
type BookingsView struct {
	Status   string
	Search   string
	Page     int
	PageSize int
	Total    int
	Bookings []store.Booking
}

// statusQuery is the "?status=x" suffix for a base path, or "" for the all tab.
func statusQuery(status string) string {
	if status == "" {
		return ""
	}
	return "?status=" + status
}

// bookingsTabExtra is the extra query suffix the filter tabs carry so switching
// status preserves the active search (already URL-escaped), or "" when no search.
func bookingsTabExtra(search string) string {
	if search == "" {
		return ""
	}
	return "q=" + url.QueryEscape(search)
}

// bookingsURL builds a bookings URL carrying status, search, and page (page 1 and
// empty filters are omitted for clean canonical URLs).
func bookingsURL(status, search string, page int) string {
	v := url.Values{}
	if status != "" {
		v.Set("status", status)
	}
	if search != "" {
		v.Set("q", search)
	}
	if page > 1 {
		v.Set("page", strconv.Itoa(page))
	}
	if len(v) == 0 {
		return "/admin/bookings"
	}
	return "/admin/bookings?" + v.Encode()
}

// pagerLabel renders the "1–25 of 60" range for the current page.
func pagerLabel(bv BookingsView) string {
	if bv.Total == 0 {
		return "No pickups"
	}
	start := (bv.Page-1)*bv.PageSize + 1
	return fmt.Sprintf("%d–%d of %d", start, start+len(bv.Bookings)-1, bv.Total)
}

// hasNextPage reports whether another page of bookings exists after this one.
func hasNextPage(bv BookingsView) bool {
	return bv.Page*bv.PageSize < bv.Total
}

// bookingsEmptyMsg is the contextual empty-state line (search vs status vs none).
func bookingsEmptyMsg(bv BookingsView) string {
	switch {
	case bv.Search != "":
		return fmt.Sprintf("No pickups match “%s.”", bv.Search)
	case bv.Status != "":
		return "No " + strings.ToLower(bookingStatusLabel(bv.Status)) + " pickups right now."
	default:
		return "No pickups yet — they'll land here the moment someone books."
	}
}

// bookingTimeFlag marks an active booking as "today" or "overdue" by its
// scheduled date (ISO strings compare lexicographically). Done/empty rows are
// unflagged. Drives the row emphasis in the pickups table.
func bookingTimeFlag(b store.Booking) string {
	if b.ScheduledDate == "" || b.Status == "delivered" || b.Status == "canceled" {
		return ""
	}
	today := time.Now().Format("2006-01-02")
	switch {
	case b.ScheduledDate < today:
		return "overdue"
	case b.ScheduledDate == today:
		return "today"
	default:
		return ""
	}
}

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
