package view

import (
	"fmt"
	"time"

	"github.com/firefly-software-mt/advanced-template/internal/store"
)

// CalDay is one day in any calendar view — a month-grid cell, a week column, or
// the day run sheet. InMonth means "belongs to the active scope": in month view
// the real month days (vs leading/trailing spill); in week/day view it's always
// true (every shown day is in scope).
type CalDay struct {
	ISO      string // YYYY-MM-DD
	Day      int    // day of month
	InMonth  bool   // belongs to the displayed scope (vs month-grid spill)
	IsToday  bool
	Bookings []store.Booking // scheduled for this day
}

// Cal is the calendar view model for one displayed period, in one of three views.
// Day/week views fill Days (1 or 7 entries); month view fills Weeks (the 6×7
// grid) and Days (the flat 42 for the phone agenda).
type Cal struct {
	View        string     // "day" | "week" | "month"
	Anchor      string     // YYYY-MM-DD the view is centered on
	Label       string     // header: "June 2026" | "Jun 8–14, 2026" | "Saturday, Jun 13"
	Prev        string     // anchor date for the previous step
	Next        string     // anchor date for the next step
	IsCurrent   bool       // the shown period contains today (hides the "Today" button)
	Weeks       [][]CalDay // month grid only (nil for day/week)
	Days        []CalDay   // all in-scope days: 1 (day), 7 (week), or 42 (month, for agenda)
	Unscheduled []store.Booking
}

// CalRef is the minimal scope (view + anchor) round-tripped through every chip,
// day-peek, and edit panel so a save re-renders the same view on the same date.
type CalRef struct {
	View   string
	Anchor string
}

// Ref is the current scope of a built calendar.
func (c Cal) Ref() CalRef { return CalRef{View: c.View, Anchor: c.Anchor} }

// Query is the "v=…&d=…" suffix carrying this scope on a URL.
func (r CalRef) Query() string { return "v=" + r.View + "&d=" + r.Anchor }

// CalViews are the switcher options, in altitude order (broad → narrow shown
// right). Label is the button text.
var CalViews = []struct{ Key, Label string }{
	{"day", "Day"},
	{"week", "Week"},
	{"month", "Month"},
}

// WeekdayLabels are the month-grid column headers, Sunday-first (US convention).
var WeekdayLabels = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// gridDays is the number of cells in the month grid (6 weeks, always — keeps the
// grid height stable across months).
const gridDays = 42

// defaultCalView is the view used when none (or an unknown one) is requested.
// Week is the operating default: it matches Chanté's weekly pickup rhythm and
// has room to show each pickup's window without crushing the layout.
const defaultCalView = "week"

// normalizeCalView clamps an arbitrary "v" param to a known view.
func normalizeCalView(v string) string {
	switch v {
	case "day", "week", "month":
		return v
	default:
		return defaultCalView
	}
}

// parseAnchor parses a "2006-01-02" anchor (UTC, fine for date-only math),
// falling back to today's date on a bad/empty value.
func parseAnchor(anchor string, now time.Time) time.Time {
	if t, err := time.Parse("2006-01-02", anchor); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// firstOfMonth returns the first day of t's month.
func firstOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// gridStart is the Sunday on/just before the first of the month.
func gridStart(first time.Time) time.Time {
	return first.AddDate(0, 0, -int(first.Weekday()))
}

// weekStart is the Sunday on/just before the given date.
func weekStart(t time.Time) time.Time {
	return t.AddDate(0, 0, -int(t.Weekday()))
}

const isoFmt = "2006-01-02"

// CalRange returns the [start, end) ISO window the given view+anchor covers, so
// the handler fetches exactly the scheduled bookings the view will display.
func CalRange(vmode, anchor string, now time.Time) (start, end string) {
	a := parseAnchor(anchor, now)
	switch normalizeCalView(vmode) {
	case "day":
		return a.Format(isoFmt), a.AddDate(0, 0, 1).Format(isoFmt)
	case "week":
		ws := weekStart(a)
		return ws.Format(isoFmt), ws.AddDate(0, 0, 7).Format(isoFmt)
	default: // month
		gs := gridStart(firstOfMonth(a))
		return gs.Format(isoFmt), gs.AddDate(0, 0, gridDays).Format(isoFmt)
	}
}

// makeDay builds one CalDay, attaching any bookings scheduled on it.
func makeDay(d time.Time, byDate map[string][]store.Booking, todayISO string, inScope bool) CalDay {
	iso := d.Format(isoFmt)
	return CalDay{
		ISO:      iso,
		Day:      d.Day(),
		InMonth:  inScope,
		IsToday:  iso == todayISO,
		Bookings: byDate[iso],
	}
}

// BuildCal assembles the view model for one view+anchor from the scheduled +
// unscheduled bookings the handler fetched. now drives the "today" highlight and
// the empty-anchor fallback.
func BuildCal(vmode, anchor string, scheduled, unscheduled []store.Booking, now time.Time) Cal {
	vmode = normalizeCalView(vmode)
	a := parseAnchor(anchor, now)
	todayISO := now.Format(isoFmt)

	byDate := make(map[string][]store.Booking, len(scheduled))
	for _, b := range scheduled {
		byDate[b.ScheduledDate] = append(byDate[b.ScheduledDate], b)
	}

	c := Cal{View: vmode, Anchor: a.Format(isoFmt), Unscheduled: unscheduled}

	switch vmode {
	case "day":
		c.Days = []CalDay{makeDay(a, byDate, todayISO, true)}
		c.Label = a.Format("Monday, Jan 2")
		c.Prev = a.AddDate(0, 0, -1).Format(isoFmt)
		c.Next = a.AddDate(0, 0, 1).Format(isoFmt)
		c.IsCurrent = c.Anchor == todayISO

	case "week":
		ws := weekStart(a)
		days := make([]CalDay, 7)
		for i := 0; i < 7; i++ {
			days[i] = makeDay(ws.AddDate(0, 0, i), byDate, todayISO, true)
		}
		c.Days = days
		c.Label = weekLabel(ws)
		c.Prev = ws.AddDate(0, 0, -7).Format(isoFmt)
		c.Next = ws.AddDate(0, 0, 7).Format(isoFmt)
		c.IsCurrent = todayISO >= ws.Format(isoFmt) && todayISO < ws.AddDate(0, 0, 7).Format(isoFmt)

	default: // month
		first := firstOfMonth(a)
		gs := gridStart(first)
		weeks := make([][]CalDay, 0, gridDays/7)
		flat := make([]CalDay, 0, gridDays)
		week := make([]CalDay, 0, 7)
		for i := 0; i < gridDays; i++ {
			d := gs.AddDate(0, 0, i)
			inMonth := d.Month() == first.Month() && d.Year() == first.Year()
			cd := makeDay(d, byDate, todayISO, inMonth)
			week = append(week, cd)
			flat = append(flat, cd)
			if len(week) == 7 {
				weeks = append(weeks, week)
				week = make([]CalDay, 0, 7)
			}
		}
		c.Weeks = weeks
		c.Days = flat
		c.Label = first.Format("January 2006")
		c.Prev = first.AddDate(0, -1, 0).Format(isoFmt)
		c.Next = first.AddDate(0, 1, 0).Format(isoFmt)
		c.IsCurrent = first.Format("2006-01") == now.Format("2006-01")
	}

	return c
}

// weekLabel renders a week's date range, collapsing the shared month/year:
// "Jun 8–14, 2026", "Jun 29 – Jul 5, 2026", or spanning a year boundary.
func weekLabel(ws time.Time) string {
	we := ws.AddDate(0, 0, 6)
	switch {
	case ws.Month() == we.Month() && ws.Year() == we.Year():
		return fmt.Sprintf("%s %d–%d, %d", ws.Format("Jan"), ws.Day(), we.Day(), ws.Year())
	case ws.Year() == we.Year():
		return fmt.Sprintf("%s %d – %s %d, %d", ws.Format("Jan"), ws.Day(), we.Format("Jan"), we.Day(), ws.Year())
	default:
		return fmt.Sprintf("%s %d, %d – %s %d, %d", ws.Format("Jan"), ws.Day(), ws.Year(), we.Format("Jan"), we.Day(), we.Year())
	}
}

// maxChipsPerCell caps how many booking chips a month-grid day draws before the
// rest collapse into a "+N more" button. Kept small so the grid stays a stable
// height no matter how busy one day gets.
const maxChipsPerCell = 3

// cellChips returns the chips to draw in a month-grid day cell. When the day
// holds exactly one over the cap, all are shown anyway — a "+1 more" button to
// hide a single pickup reads worse than just showing it.
func cellChips(bs []store.Booking) []store.Booking {
	if len(bs) <= maxChipsPerCell+1 {
		return bs
	}
	return bs[:maxChipsPerCell]
}

// cellOverflow is the count hidden behind the "+N more" button (0 when none).
func cellOverflow(bs []store.Booking) int {
	if len(bs) <= maxChipsPerCell+1 {
		return 0
	}
	return len(bs) - maxChipsPerCell
}

// chipTitle is the hover tooltip for a grid chip: name, status, and window if set.
func chipTitle(b store.Booking) string {
	t := b.Name + " — " + bookingStatusLabel(b.Status)
	if b.PickupWindow != "" {
		t += " · " + b.PickupWindow
	}
	return t
}

// calViewURL switches to view v while staying on the current anchor date.
func calViewURL(v, anchor string) string {
	return "/admin/calendar?v=" + v + "&d=" + anchor
}

// calStepURL navigates the same view to a new anchor (prev/next step).
func calStepURL(view, anchor string) string {
	return "/admin/calendar?v=" + view + "&d=" + anchor
}

// calTodayURL jumps the current view back to today (no anchor → handler defaults
// to now).
func calTodayURL(view string) string {
	return "/admin/calendar?v=" + view
}

// dayURL is the GET that loads one day's full pickup list into the drawer. ref is
// the current calendar scope, round-tripped so a save from the day list returns
// to the view the operator was on.
func dayURL(iso string, ref CalRef) string {
	return fmt.Sprintf("/admin/calendar/day/%s?%s", iso, ref.Query())
}

// pickupURL is the GET that loads a pickup's edit panel, carrying the scope.
func pickupURL(id int64, ref CalRef) string {
	return fmt.Sprintf("/admin/calendar/pickup/%d?%s", id, ref.Query())
}

// DayPanelLabel renders a drawer day heading like "Saturday, Jun 13" from an ISO
// date, falling back to the raw value on a parse error.
func DayPanelLabel(iso string) string {
	if t, err := time.Parse(isoFmt, iso); err == nil {
		return t.Format("Monday, Jan 2")
	}
	return iso
}

// calHasScheduled reports whether any in-scope day in the slice carries a booking
// — drives the agenda / run-sheet empty state.
func calHasScheduled(days []CalDay) bool {
	for _, d := range days {
		if d.InMonth && len(d.Bookings) > 0 {
			return true
		}
	}
	return false
}

// agendaDayLabel renders an agenda/run-sheet day heading like "Mon, Jun 8".
func agendaDayLabel(d CalDay) string {
	if t, err := time.Parse(isoFmt, d.ISO); err == nil {
		return t.Format("Mon, Jan 2")
	}
	return d.ISO
}

// weekdayShort renders a week-column header weekday like "Mon".
func weekdayShort(d CalDay) string {
	if t, err := time.Parse(isoFmt, d.ISO); err == nil {
		return t.Format("Mon")
	}
	return ""
}

// agendaMeta is the secondary line on an agenda row: plan, plus the pickup window
// when one is set.
func agendaMeta(b store.Booking) string {
	m := planShort(b.Plan)
	if b.PickupWindow != "" {
		m += " · " + b.PickupWindow
	}
	return m
}

// ariaPressed returns "true"/"false" for a toggle button's aria-pressed state.
func ariaPressed(on bool) string {
	if on {
		return "true"
	}
	return "false"
}

// bookingDotClass returns the full class string for a status dot — a small
// circle paired on two channels so the seven statuses stay distinguishable
// without relying on hue alone: HUE groups the lifecycle (ochre = needs
// scheduling, teal = booked & waiting, rust = in Chanté's hands, brown = done,
// faint = off), and FILL splits each hue's two steps (a ring is the earlier /
// in-flight step, a solid is the later / action-needed step). Used by the grid
// chips, week events, run-sheet rows, and the legend, so a dot always means the
// same thing.
func bookingDotClass(status string) string {
	const base = "h-2.5 w-2.5 shrink-0 rounded-full border-[1.5px] "
	switch status {
	case "new":
		return base + "border-warning bg-warning"
	case "scheduled":
		return base + "border-success bg-transparent"
	case "picked_up":
		return base + "border-success bg-success"
	case "in_progress":
		return base + "border-accent bg-transparent"
	case "ready":
		return base + "border-accent bg-accent"
	case "delivered":
		return base + "border-foreground bg-foreground"
	default: // canceled / unknown
		return base + "border-muted-foreground bg-transparent"
	}
}
