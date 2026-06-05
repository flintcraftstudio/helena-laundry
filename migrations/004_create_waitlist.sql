-- +goose Up
-- Out-of-area waitlist signups. Location + kind are captured cleanly so
-- demand is sortable later for expansion decisions.
CREATE TABLE waitlist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    location TEXT NOT NULL,       -- town or neighborhood
    contact TEXT NOT NULL,        -- phone or email
    kind TEXT NOT NULL DEFAULT 'residential', -- residential | business
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_waitlist_location ON waitlist (location);

-- +goose Down
DROP TABLE waitlist;
