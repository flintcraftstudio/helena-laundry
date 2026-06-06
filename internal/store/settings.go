package store

import (
	"context"
	"strings"
)

// Settings is the single-row site configuration the operator edits in the
// dashboard. AllowedZips/PickupDays are stored as CSV but exposed as slices.
type Settings struct {
	SiteName          string
	Tagline           string
	Phone             string
	Email             string
	Hours             string
	Turnaround        string
	AcceptingBookings bool
	LbRate            string
	RushRate          string
	MinimumOrder      string
	AllowedZips       []string
	PickupDays        []string
	FlatRates         []FlatRate
}

// FlatRate is one flat-price line item (e.g. "Queen bedding" — "$30").
type FlatRate struct {
	ID    int64
	Label string
	Price string
	Sort  int
}

// GetSettings loads the singleton settings row plus its flat rates.
func (s *Store) GetSettings(ctx context.Context) (Settings, error) {
	var st Settings
	var accepting int
	var zips, days string
	err := s.db.QueryRowContext(ctx,
		`SELECT site_name, tagline, phone, email, hours, turnaround, accepting_bookings,
		        lb_rate, rush_rate, minimum_order, allowed_zips, pickup_days
		 FROM settings WHERE id = 1`,
	).Scan(&st.SiteName, &st.Tagline, &st.Phone, &st.Email, &st.Hours, &st.Turnaround,
		&accepting, &st.LbRate, &st.RushRate, &st.MinimumOrder, &zips, &days)
	if err != nil {
		return Settings{}, err
	}
	st.AcceptingBookings = accepting != 0
	st.AllowedZips = splitCSV(zips)
	st.PickupDays = splitCSV(days)

	rates, err := s.ListFlatRates(ctx)
	if err != nil {
		return Settings{}, err
	}
	st.FlatRates = rates
	return st, nil
}

// UpdateSettings writes the scalar settings (not flat rates — see ReplaceFlatRates).
func (s *Store) UpdateSettings(ctx context.Context, st Settings) error {
	accepting := 0
	if st.AcceptingBookings {
		accepting = 1
	}
	_, err := s.db.ExecContext(ctx,
		`UPDATE settings SET
		   site_name = ?, tagline = ?, phone = ?, email = ?, hours = ?, turnaround = ?,
		   accepting_bookings = ?, lb_rate = ?, rush_rate = ?, minimum_order = ?,
		   allowed_zips = ?, pickup_days = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE id = 1`,
		st.SiteName, st.Tagline, st.Phone, st.Email, st.Hours, st.Turnaround,
		accepting, st.LbRate, st.RushRate, st.MinimumOrder,
		joinCSV(st.AllowedZips), joinCSV(st.PickupDays),
	)
	return err
}

// ListFlatRates returns flat-rate line items in display order.
func (s *Store) ListFlatRates(ctx context.Context) ([]FlatRate, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, label, price, sort FROM flat_rates ORDER BY sort, id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []FlatRate
	for rows.Next() {
		var f FlatRate
		if err := rows.Scan(&f.ID, &f.Label, &f.Price, &f.Sort); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// ReplaceFlatRates swaps the entire flat-rate list in one transaction (the admin
// form posts the full set each save).
func (s *Store) ReplaceFlatRates(ctx context.Context, rates []FlatRate) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM flat_rates"); err != nil {
		return err
	}
	for i, r := range rates {
		if strings.TrimSpace(r.Label) == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO flat_rates (label, price, sort) VALUES (?, ?, ?)",
			strings.TrimSpace(r.Label), strings.TrimSpace(r.Price), i,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// splitCSV splits a comma-separated string into trimmed, non-empty values.
func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// joinCSV joins trimmed, non-empty values with commas.
func joinCSV(vals []string) string {
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
	}
	return strings.Join(out, ",")
}
