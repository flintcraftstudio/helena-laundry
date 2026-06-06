package view

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
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
	// A new out-of-area signup is an opportunity (where to grow next), not a
	// problem — accent rust, not warning amber.
	"new":       {"New", badge.VariantAccent},
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

// ---------------------------------------------------------------------------
// Waitlist demand view — turns the flat signup list into "where to grow next".
// ---------------------------------------------------------------------------

// waitlistTopAreas caps how many areas the demand band shows before collapsing
// the rest into a "+N more" note (kept scannable; the full set is in the table).
const waitlistTopAreas = 6

// AreaDemand is one location's pull: how many people are waiting there and how
// many of them are businesses (the higher-value, standing-pickup leads).
type AreaDemand struct {
	Label    string // location as the signer typed it (first seen)
	Key      string // case-folded location, used for matching + the ?loc= filter
	Count    int
	Business int
}

// WaitlistView is everything the waitlist section needs: the active status tab,
// the active location filter ("" = all areas), the table rows (status- AND
// location-filtered), and the demand band (aggregated over the status-filtered
// set only, so picking an area never collapses the band to one chip).
type WaitlistView struct {
	Status    string
	Loc       string
	Entries   []store.WaitlistEntry
	Areas     []AreaDemand
	MoreAreas int // distinct areas beyond the ones shown in the band
	Waiting   int // total entries in the status-filtered set
}

// NewWaitlistView aggregates the demand band from the full status-filtered set,
// then narrows the table rows to the active location. items is already limited +
// status-filtered by the store; this adds no DB hit.
func NewWaitlistView(status, loc string, items []store.WaitlistEntry) WaitlistView {
	loc = normalizeLoc(loc)
	areas, totalAreas := topAreas(items, waitlistTopAreas)

	entries := items
	if loc != "" {
		filtered := make([]store.WaitlistEntry, 0, len(items))
		for _, e := range items {
			if normalizeLoc(e.Location) == loc {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	return WaitlistView{
		Status:    status,
		Loc:       loc,
		Entries:   entries,
		Areas:     areas,
		MoreAreas: totalAreas - len(areas),
		Waiting:   len(items),
	}
}

// topAreas groups entries by location (case-folded), busiest first, and returns
// the top n plus the total distinct-area count so the view can say how many were
// left off. Entries with a blank location are skipped — they tell us nothing
// about where to head next.
func topAreas(entries []store.WaitlistEntry, n int) (top []AreaDemand, total int) {
	idx := map[string]int{}
	var all []AreaDemand
	for _, e := range entries {
		key := normalizeLoc(e.Location)
		if key == "" {
			continue
		}
		i, ok := idx[key]
		if !ok {
			i = len(all)
			idx[key] = i
			all = append(all, AreaDemand{Label: strings.TrimSpace(e.Location), Key: key})
		}
		all[i].Count++
		if e.Kind == "business" {
			all[i].Business++
		}
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].Count > all[j].Count })
	if n < len(all) {
		return all[:n], len(all)
	}
	return all, len(all)
}

// waitlistKindLabel turns the stored kind code into a display label. Anything
// that isn't an explicit "business" is treated as residential (matches the
// submit handler, which defaults non-business to residential).
func waitlistKindLabel(kind string) string {
	if kind == "business" {
		return "Business"
	}
	return "Residential"
}

// normalizeLoc folds a location to its grouping/URL key. Free-typed locations
// won't always agree ("E. Helena" vs "East Helena") — this only collapses the
// easy cases (whitespace + casing); it deliberately doesn't guess at the rest.
func normalizeLoc(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// areaSplit phrases an area's residential/business mix in a few words.
func areaSplit(a AreaDemand) string {
	switch {
	case a.Business == 0:
		return "all residential"
	case a.Business == a.Count:
		return "all business"
	default:
		return fmt.Sprintf("%d business", a.Business)
	}
}

// waitlistBandSummary is the band's one-line gloss, e.g. "12 people across 4 areas".
func waitlistBandSummary(wv WaitlistView) string {
	areas := len(wv.Areas) + wv.MoreAreas
	return fmt.Sprintf("%d %s across %d %s",
		wv.Waiting, plural(wv.Waiting, "person", "people"),
		areas, plural(areas, "area", "areas"))
}

// waitlistTabExtra carries the active location filter across a status-tab switch.
func waitlistTabExtra(loc string) string {
	if loc == "" {
		return ""
	}
	return "loc=" + url.QueryEscape(loc)
}

// waitlistAreaURL builds a demand-chip's hx-get: the current status plus the
// area's location key (empty key = the "All areas" reset).
func waitlistAreaURL(status, key string) string {
	q := url.Values{}
	if status != "" {
		q.Set("status", status)
	}
	if key != "" {
		q.Set("loc", key)
	}
	if enc := q.Encode(); enc != "" {
		return "/admin/waitlist?" + enc
	}
	return "/admin/waitlist"
}

// waitlistEmptyMsg is the contextual empty state — it names why the view is empty
// (a location filter or a status tab) rather than a flat "nothing here".
func waitlistEmptyMsg(wv WaitlistView) string {
	switch {
	case wv.Loc != "":
		return "Nobody from that area in this view."
	case wv.Status == "contacted":
		return "Nobody's marked contacted yet."
	case wv.Status == "archived":
		return "Nothing archived yet."
	case wv.Status == "new":
		return "No new signups right now — this fills as out-of-area folks raise a hand."
	default:
		return "No waitlist signups yet. When someone outside the route signs up, this is where you'll see where to head next."
	}
}

// plural picks the singular or plural word for n.
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

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
