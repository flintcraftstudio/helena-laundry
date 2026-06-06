// Command seeddata wipes and repopulates the operational tables
// (bookings, questions, waitlist) with a realistic spread of sample rows so the
// admin dashboard panels can be exercised. It leaves users, sessions, and
// settings untouched. Safe to re-run — it clears the three tables first.
//
//	go run ./cmd/seeddata          # ./data/app.db (or $DB_PATH)
//	DB_PATH=/tmp/x.db go run ./cmd/seeddata
package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const (
	nBookings = 60
	nQuestions = 60
	nWaitlist  = 60
)

// Helena-area flavor so the data reads true to brand.
var firstNames = []string{
	"Maggie", "Dale", "Bonnie", "Hank", "Shauna", "Cole", "Della", "Wyatt",
	"Rosa", "Travis", "Junie", "Marv", "Patty", "Gus", "Lena", "Roy",
	"Cara", "Dean", "Birdie", "Sal", "Nora", "Cody", "Faye", "Hollis",
}
var lastNames = []string{
	"Whitlock", "Behan", "Cobell", "Running Crane", "Yellowtail", "Doyle",
	"Kowalski", "Andersen", "Begay", "Petrović", "Hauser", "Stillwater",
	"O'Rourke", "Beaulieu", "Schmidt", "Larson", "Greer", "Tallman",
}
var inAreaZips = []string{"59601", "59602", "59634", "59635"}
var plans = []string{"standard", "weekly", "rush"}
var windows = []string{"Morning (8–11)", "Midday (11–2)", "Afternoon (2–5)"}
var weekdays = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

var bookingNotes = []string{
	"", "", "Gate code 1847.", "Leave on the porch bench, please.",
	"Two loads — one's mostly towels.", "Allergic to scented detergent.",
	"Dog's friendly, ignore the bark.", "Ring twice, baby's napping.",
	"Hang-dry the linen shirts if you can.", "Back door is easier.",
}
var questionMsgs = []string{
	"Do you do same-day if I call early?",
	"What's the turnaround on bedding?",
	"Can you handle a king comforter?",
	"Just wanted to say my towels came back perfect. Thank you!",
	"Do you take cash, or just card?",
	"Is there a minimum order for a one-off pickup?",
	"Could you do unscented detergent for my whole order?",
	"We run an Airbnb — do you do recurring linen service?",
	"Lost a sock last time, no worries, just FYI.",
	"How far out are you booking right now?",
}
var topics = []string{"question", "feedback", "schedule"}
var outOfArea = []string{
	"Townsend", "Boulder", "Lincoln", "Wolf Creek", "Augusta", "Three Forks",
	"Whitehall", "Clancy", "Montana City", "Basin", "Elliston", "Avon",
	"Canyon Creek", "Marysville", "Craig", "Cascade",
}
var waitlistNotes = []string{
	"", "", "We go through a lot of towels.", "Salon — color towels, daily.",
	"Two short-term rentals downtown.", "Gym, mats and towels.",
	"Just moved here, no machine yet.", "Spa linens, want a standing day.",
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/app.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fail("open database", err)
	}
	defer db.Close()
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		fail("pragma", err)
	}

	tx, err := db.Begin()
	if err != nil {
		fail("begin tx", err)
	}
	defer tx.Rollback()

	// Wipe the operational tables (leave users/sessions/settings alone) and reset
	// their autoincrement counters so IDs start clean each run.
	for _, t := range []string{"bookings", "questions", "waitlist"} {
		if _, err := tx.Exec("DELETE FROM " + t); err != nil {
			fail("wipe "+t, err)
		}
		if _, err := tx.Exec("DELETE FROM sqlite_sequence WHERE name = ?", t); err != nil {
			fail("reset seq "+t, err)
		}
	}

	now := time.Now()
	if err := seedBookings(tx, now); err != nil {
		fail("seed bookings", err)
	}
	if err := seedQuestions(tx, now); err != nil {
		fail("seed questions", err)
	}
	if err := seedWaitlist(tx, now); err != nil {
		fail("seed waitlist", err)
	}

	if err := tx.Commit(); err != nil {
		fail("commit", err)
	}

	fmt.Printf("Seeded %d bookings, %d inquiries, %d waitlist entries into %s\n",
		nBookings, nQuestions, nWaitlist, dbPath)
	fmt.Println("(users / sessions / settings left untouched)")
}

// seedBookings spreads rows across every status. "new" rows get no scheduled_date
// (they populate the calendar's "needs a date" strip); scheduled-and-later rows
// get a real date + window, biased past for delivered/picked_up and future for
// scheduled/in_progress/ready, so the month grid and prev/next nav all have data.
func seedBookings(tx *sql.Tx, now time.Time) error {
	stmt, err := tx.Prepare(`INSERT INTO bookings
		(name, contact, pickup_day, plan, address, zip, notes, status, scheduled_date, pickup_window, admin_notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < nBookings; i++ {
		name := person(i)
		status := bookingStatus(i)

		var schedDate, window, day string
		switch status {
		case "new":
			// Unscheduled — requested weekday only, no concrete date.
			day = weekdays[i%len(weekdays)]
		case "delivered", "picked_up":
			d := now.AddDate(0, 0, -(1 + i%18)) // past
			schedDate = d.Format("2006-01-02")
			window = windows[i%len(windows)]
			day = d.Weekday().String()
		case "canceled":
			d := now.AddDate(0, 0, (i%30)-12) // mixed past/future
			schedDate = d.Format("2006-01-02")
			window = windows[i%len(windows)]
			day = d.Weekday().String()
		default: // scheduled | in_progress | ready -> upcoming
			d := now.AddDate(0, 0, i%22) // today..~3 weeks out
			schedDate = d.Format("2006-01-02")
			window = windows[i%len(windows)]
			day = d.Weekday().String()
		}

		adminNote := ""
		if status != "new" && i%4 == 0 {
			adminNote = "Regular — knows the drill."
		}

		created := now.Add(-time.Duration(i*7) * time.Hour).Format("2006-01-02 15:04:05")

		if _, err := stmt.Exec(
			name, contact(i, name), day, plans[i%len(plans)],
			address(i), inAreaZips[i%len(inAreaZips)], bookingNotes[i%len(bookingNotes)],
			status, schedDate, window, adminNote, created,
		); err != nil {
			return err
		}
	}
	return nil
}

func seedQuestions(tx *sql.Tx, now time.Time) error {
	stmt, err := tx.Prepare(`INSERT INTO questions
		(name, contact, topic, message, status, admin_notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < nQuestions; i++ {
		name := person(i + 3)
		created := now.Add(-time.Duration(i*9) * time.Hour).Format("2006-01-02 15:04:05")
		if _, err := stmt.Exec(
			name, contact(i+3, name), topics[i%len(topics)], questionMsgs[i%len(questionMsgs)],
			triage(i, "handled", "archived"), "", created,
		); err != nil {
			return err
		}
	}
	return nil
}

func seedWaitlist(tx *sql.Tx, now time.Time) error {
	stmt, err := tx.Prepare(`INSERT INTO waitlist
		(name, location, contact, kind, notes, status, admin_notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < nWaitlist; i++ {
		name := person(i + 7)
		kind := "residential"
		if i%3 == 0 {
			kind = "business"
		}
		created := now.Add(-time.Duration(i*11) * time.Hour).Format("2006-01-02 15:04:05")
		if _, err := stmt.Exec(
			name, outOfArea[i%len(outOfArea)], contact(i+7, name), kind,
			waitlistNotes[i%len(waitlistNotes)], triage(i, "contacted", "archived"), "", created,
		); err != nil {
			return err
		}
	}
	return nil
}

// bookingStatus walks every status in the lifecycle, weighted so "new" (the
// unscheduled strip) and "delivered" appear often.
func bookingStatus(i int) string {
	switch i % 10 {
	case 0, 1:
		return "new"
	case 2, 3:
		return "scheduled"
	case 4:
		return "picked_up"
	case 5:
		return "in_progress"
	case 6:
		return "ready"
	case 7, 8:
		return "delivered"
	default:
		return "canceled"
	}
}

// triage maps an index to new|<mid>|<end>, keeping a majority "new" so the
// dashboard's new-count badges are populated.
func triage(i int, mid, end string) string {
	switch i % 5 {
	case 3:
		return mid
	case 4:
		return end
	default:
		return "new"
	}
}

func person(i int) string {
	return firstNames[i%len(firstNames)] + " " + lastNames[(i*3)%len(lastNames)]
}

// contact alternates phone and email so both render paths get exercised.
func contact(i int, name string) string {
	if i%2 == 0 {
		return fmt.Sprintf("406-555-%04d", 1000+i)
	}
	first := firstNames[i%len(firstNames)]
	return fmt.Sprintf("%s%d@example.com", first, i)
}

func address(i int) string {
	streets := []string{"Last Chance Gulch", "Benton Ave", "Euclid Ave", "Hauser Blvd", "Cedar St", "Henderson St", "Broadway"}
	return fmt.Sprintf("%d %s", 100+(i*37)%900, streets[i%len(streets)])
}

func fail(what string, err error) {
	fmt.Fprintf(os.Stderr, "seeddata: %s: %v\n", what, err)
	os.Exit(1)
}
