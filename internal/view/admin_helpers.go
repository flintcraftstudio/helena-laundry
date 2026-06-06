package view

import (
	"fmt"
	"net/url"
	"regexp"
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

// InquiriesView is everything the inquiries section needs to render: the active
// status tab, the search term, the current page, and the page of results plus the
// total match count for the pager. Mirrors BookingsView.
type InquiriesView struct {
	Status    string
	Search    string
	Page      int
	PageSize  int
	Total     int
	Inquiries []store.Question
}

// inquiriesTabExtra is the extra query suffix the filter tabs carry so switching
// status preserves the active search (already URL-escaped), or "" when no search.
func inquiriesTabExtra(search string) string {
	if search == "" {
		return ""
	}
	return "q=" + url.QueryEscape(search)
}

// inquiriesURL builds an inquiries URL carrying status, search, and page (page 1
// and empty filters are omitted for clean canonical URLs).
func inquiriesURL(status, search string, page int) string {
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
		return "/admin/inquiries"
	}
	return "/admin/inquiries?" + v.Encode()
}

// inquiriesPagerLabel renders the "1–25 of 60" range for the current page.
func inquiriesPagerLabel(iv InquiriesView) string {
	if iv.Total == 0 {
		return "No inquiries"
	}
	start := (iv.Page-1)*iv.PageSize + 1
	return fmt.Sprintf("%d–%d of %d", start, start+len(iv.Inquiries)-1, iv.Total)
}

// inquiriesHasNext reports whether another page of inquiries exists after this one.
func inquiriesHasNext(iv InquiriesView) bool {
	return iv.Page*iv.PageSize < iv.Total
}

// inquiriesEmptyMsg is the contextual empty-state line (search vs status vs none),
// in Chanté's warm, first-person voice.
func inquiriesEmptyMsg(iv InquiriesView) string {
	switch {
	case iv.Search != "":
		return fmt.Sprintf("Nothing matches “%s.”", iv.Search)
	case iv.Status != "":
		return "Nothing " + strings.ToLower(questionStatusLabel(iv.Status)) + " right now."
	default:
		return "No inquiries yet — questions and feedback from the site land right here."
	}
}

// inquiryRowClass tints a new (unhandled) inquiry row so the work to do stands
// out from already-handled rows during triage.
func inquiryRowClass(status string) string {
	if status == "new" {
		return rowClickClass + " bg-accent/[0.04]"
	}
	return rowClickClass
}

// inquiryAccentCell paints a thin rust bar down the leading cell of a new
// inquiry (table cells are position:relative, so the before-element anchors to
// the row's left edge regardless of border-collapse).
func inquiryAccentCell(status string) string {
	if status == "new" {
		return "before:absolute before:inset-y-0 before:left-0 before:w-0.5 before:bg-accent"
	}
	return ""
}

// Contact is one free-text "Best way to reach you" field, so an inquirer may
// leave a phone, an email, or both ("406-555-1234 or jane@example.com"). These
// patterns pull each out independently.
var (
	contactEmailRe = regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)
	contactPhoneRe = regexp.MustCompile(`\+?\d[\d\s().\-]{5,}\d`)
)

// ContactLink is one tappable way to reach an inquirer, parsed from the contact
// field. Kind is "email" or "phone" (drives the icon/label in the reply button).
type ContactLink struct {
	Kind  string
	Label string
	Href  string
}

// contactLinks pulls an email and/or a phone number out of the free-text contact
// field and returns a tappable link for each it finds (email first), so the
// operator can reply in one tap whether someone left a phone, an email, or both.
// Returns nil when nothing linkable is found (the caller falls back to plain text).
func contactLinks(contact string) []ContactLink {
	var out []ContactLink
	if m := contactEmailRe.FindString(contact); m != "" {
		out = append(out, ContactLink{Kind: "email", Label: m, Href: "mailto:" + m})
	}
	// Strip any email before scanning for a phone so its digits aren't mistaken
	// for a number.
	rest := contactEmailRe.ReplaceAllString(contact, " ")
	if m := contactPhoneRe.FindString(rest); m != "" {
		if tel := telDigits(m); len(strings.TrimPrefix(tel, "+")) >= 7 {
			out = append(out, ContactLink{Kind: "phone", Label: strings.TrimSpace(m), Href: "tel:" + tel})
		}
	}
	return out
}

// telDigits reduces a matched phone string to a dialable tel: value — digits
// only, preserving a single leading "+".
func telDigits(s string) string {
	var b strings.Builder
	for _, r := range strings.TrimSpace(s) {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '+' && b.Len() == 0:
			b.WriteRune(r)
		}
	}
	return b.String()
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

// questionTopicMeta maps the contact form's topic codes to human labels. A
// "schedule" inquiry is someone who came to ask but actually wants to book — a
// hot lead — so it gets the accent variant to stand out.
var questionTopicMeta = map[string]statusMeta{
	"question": {"Question", badge.VariantMuted},
	"feedback": {"Feedback", badge.VariantMuted},
	"schedule": {"Wants to schedule", badge.VariantAccent},
}

func questionTopicMetaOr(s string) statusMeta {
	if v, ok := questionTopicMeta[s]; ok {
		return v
	}
	if s == "" {
		return statusMeta{Label: "—", Variant: badge.VariantMuted}
	}
	return statusMeta{Label: s, Variant: badge.VariantMuted}
}

func questionTopicLabel(s string) string          { return questionTopicMetaOr(s).Label }
func questionTopicVariant(s string) badge.Variant { return questionTopicMetaOr(s).Variant }

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

func bookingStatusLabel(s string) string           { return metaOr(bookingStatusMeta, s).Label }
func bookingStatusVariant(s string) badge.Variant  { return metaOr(bookingStatusMeta, s).Variant }
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
