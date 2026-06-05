-- +goose Up
-- Pickup bookings from the multi-step "Book a pickup" modal. Captures enough
-- for Chanté to reach the person and run the pickup; plan/day stay as plain
-- text so the modal's options can change without a schema migration.
CREATE TABLE bookings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    contact TEXT NOT NULL,        -- phone or email, as typed
    pickup_day TEXT NOT NULL,     -- Monday..Friday
    plan TEXT NOT NULL,           -- standard | weekly | rush
    address TEXT NOT NULL DEFAULT '', -- street / neighborhood (optional)
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE bookings;
