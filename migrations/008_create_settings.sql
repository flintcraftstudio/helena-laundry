-- +goose Up
-- Single-row site settings: the operator edits these in the dashboard and they
-- drive the booking flow (service area, days, pricing, accepting toggle) and the
-- public site's identity. Seeded from the values previously hardcoded in
-- view/shared.go, config (ALLOWED_ZIPS), and the booking handler.
CREATE TABLE settings (
    id                 INTEGER PRIMARY KEY CHECK (id = 1),
    site_name          TEXT NOT NULL,
    tagline            TEXT NOT NULL,
    phone              TEXT NOT NULL,           -- display form, e.g. 406-471-8508
    email              TEXT NOT NULL DEFAULT '',
    hours              TEXT NOT NULL DEFAULT '',
    turnaround         TEXT NOT NULL,           -- e.g. "24–48 hours"
    accepting_bookings INTEGER NOT NULL DEFAULT 1,
    lb_rate            TEXT NOT NULL,           -- dollars, e.g. "2.00"
    rush_rate          TEXT NOT NULL,           -- dollars, e.g. "3.00"
    minimum_order      TEXT NOT NULL,           -- dollars, e.g. "35"
    allowed_zips       TEXT NOT NULL,           -- CSV of 5-digit ZIPs
    pickup_days        TEXT NOT NULL,           -- CSV of weekday names
    updated_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO settings (
    id, site_name, tagline, phone, email, hours, turnaround,
    accepting_bookings, lb_rate, rush_rate, minimum_order, allowed_zips, pickup_days
) VALUES (
    1,
    'Helena Laundry',
    'Helena, Montana · Serving the Helena Valley',
    '406-471-8508',
    '',
    '',
    '24–48 hours',
    1,
    '2.00',
    '3.00',
    '35',
    '59601,59602,59634,59635',
    'Monday,Tuesday,Wednesday,Thursday,Friday'
);

-- Flat-rate line items shown on the pricing page / managed in settings.
CREATE TABLE flat_rates (
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    label TEXT NOT NULL,
    price TEXT NOT NULL,                 -- dollars, e.g. "15"
    sort  INTEGER NOT NULL DEFAULT 0
);

INSERT INTO flat_rates (label, price, sort) VALUES
    ('Twin bedding', '15', 1),
    ('Full bedding', '20', 2),
    ('Queen bedding', '30', 3),
    ('King bedding', '35', 4),
    ('Sneaker deep clean', '20', 5),
    ('13-gallon bag of towels or rags', '35', 6);

-- +goose Down
DROP TABLE flat_rates;
DROP TABLE settings;
