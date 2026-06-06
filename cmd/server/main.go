package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/firefly-software-mt/advanced-template/internal/config"
	"github.com/firefly-software-mt/advanced-template/internal/handler"
	"github.com/firefly-software-mt/advanced-template/internal/mail"
	"github.com/firefly-software-mt/advanced-template/internal/middleware"
	"github.com/firefly-software-mt/advanced-template/internal/session"
	"github.com/firefly-software-mt/advanced-template/internal/settings"
	"github.com/firefly-software-mt/advanced-template/internal/store"
	"github.com/firefly-software-mt/advanced-template/internal/view"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := loadEnv(".env"); err != nil {
		slog.Error("env error", "err", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	// Tracking pixels
	view.GtagID = cfg.GtagID
	view.PixelID = cfg.PixelID
	if cfg.GtagID == "" {
		slog.Warn("GTAG_ID not set, Google Analytics disabled")
	}
	if cfg.PixelID == "" {
		slog.Warn("PIXEL_ID not set, Facebook Pixel disabled")
	}

	// Turnstile
	view.TurnstileSiteKey = cfg.TurnstileSiteKey
	if cfg.TurnstileSiteKey == "" || cfg.TurnstileSecretKey == "" {
		slog.Warn("TURNSTILE_SITE_KEY or TURNSTILE_SECRET_KEY not set, Turnstile disabled")
	}

	// Database
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0755); err != nil {
		slog.Error("failed to create database directory", "err", err)
		os.Exit(1)
	}
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		slog.Error("database open error", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// Enable WAL mode and foreign keys for SQLite
	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;"); err != nil {
		slog.Error("database pragma error", "err", err)
		os.Exit(1)
	}

	// Run migrations
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("goose dialect error", "err", err)
		os.Exit(1)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("migration error", "err", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	// Store
	st := store.New(db)

	// Load site settings into the in-memory cache + the view package vars. These
	// drive the booking flow (service area, days, pricing, accepting toggle) and
	// the public site's identity, and are refreshed when the operator saves.
	initialSettings, err := st.GetSettings(context.Background())
	if err != nil {
		slog.Error("settings load error", "err", err)
		os.Exit(1)
	}
	settings.Set(initialSettings)
	view.ApplySettings(initialSettings)
	slog.Info("settings loaded", "accepting", initialSettings.AcceptingBookings, "zips", initialSettings.AllowedZips)

	// Mail client (nil if Postmark is not configured)
	var mailer *mail.Client
	if cfg.PostmarkToken != "" {
		mailer = mail.NewClient(cfg.PostmarkToken, cfg.PostmarkFrom, cfg.PostmarkTo)
		slog.Info("postmark configured")
	} else {
		slog.Info("postmark not configured, contact form emails disabled")
	}

	mux := http.NewServeMux()

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Pages
	mux.Handle("GET /{$}", handler.Home())
	mux.Handle("GET /services", handler.Services())
	mux.Handle("GET /about", handler.About())
	mux.Handle("GET /businesses", handler.Businesses())
	mux.Handle("GET /contact", handler.Contact())

	// Contact forms (deliberately distinct capture paths)
	mux.Handle("POST /contact/question", handler.QuestionSubmit(st, mailer, cfg.TurnstileSecretKey))
	mux.Handle("POST /contact/waitlist", handler.WaitlistSubmit(st, mailer, cfg.TurnstileSecretKey))
	mux.Handle("POST /contact/booking", handler.BookingSubmit(st, mailer, cfg.TurnstileSecretKey))

	// Auth
	mux.Handle("GET /login", handler.LoginPage())
	mux.Handle("POST /login", handler.LoginSubmit(st))
	mux.Handle("POST /logout", handler.Logout(st))

	// Admin dashboard (all routes gated by RequireAuth -> /login).
	protected := func(h http.Handler) http.Handler { return session.RequireAuth(h) }
	mux.Handle("GET /admin", protected(handler.AdminDashboard(st)))
	mux.Handle("GET /admin/calendar", protected(handler.AdminCalendar(st)))
	mux.Handle("GET /admin/calendar/pickup/{id}", protected(handler.AdminCalendarPickup(st)))
	mux.Handle("POST /admin/calendar/pickup/{id}", protected(handler.AdminCalendarPickupUpdate(st)))
	mux.Handle("GET /admin/bookings", protected(handler.AdminBookings(st)))
	mux.Handle("GET /admin/bookings/{id}", protected(handler.AdminBookingRow(st)))
	mux.Handle("GET /admin/bookings/{id}/edit", protected(handler.AdminBookingEdit(st)))
	mux.Handle("POST /admin/bookings/{id}", protected(handler.AdminBookingUpdate(st)))
	mux.Handle("GET /admin/inquiries", protected(handler.AdminInquiries(st)))
	mux.Handle("GET /admin/inquiries/{id}", protected(handler.AdminInquiryRow(st)))
	mux.Handle("GET /admin/inquiries/{id}/edit", protected(handler.AdminInquiryEdit(st)))
	mux.Handle("POST /admin/inquiries/{id}", protected(handler.AdminInquiryUpdate(st)))
	mux.Handle("GET /admin/settings", protected(handler.AdminSettings(st)))
	mux.Handle("POST /admin/settings", protected(handler.AdminSettingsUpdate(st)))
	mux.Handle("GET /admin/waitlist", protected(handler.AdminWaitlist(st)))
	mux.Handle("GET /admin/waitlist/{id}", protected(handler.AdminWaitlistRow(st)))
	mux.Handle("GET /admin/waitlist/{id}/edit", protected(handler.AdminWaitlistEdit(st)))
	mux.Handle("POST /admin/waitlist/{id}", protected(handler.AdminWaitlistUpdate(st)))

	// Branded 404 for any unmatched GET path
	mux.Handle("GET /", handler.NotFound())

	// Session + logging middleware
	srv := session.Middleware(st)(mux)
	srv = middleware.Logging(logger)(srv)

	// --- Graceful shutdown sequence ---

	// 1. Configure the HTTP server with timeouts to bound slow clients.
	//    ReadTimeout:  max time to read the entire request (headers + body).
	//    WriteTimeout: max time to write the response.
	//    IdleTimeout:  max time a keep-alive connection sits idle.
	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      srv,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 2. Start serving in a background goroutine. Any fatal listen error
	//    (e.g. port already in use) is sent to errCh so we can react.
	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", cfg.Addr())
		fmt.Printf("listening on %s\n", cfg.Addr())
		errCh <- server.ListenAndServe()
	}()

	// 3. Register for SIGINT (Ctrl-C) and SIGTERM (Docker/systemd stop).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 4. Block until we receive a shutdown signal or a server error.
	select {
	case sig := <-quit:
		slog.Info("shutdown signal received", "signal", sig)
	case err := <-errCh:
		// ErrServerClosed is expected after Shutdown(); anything else is fatal.
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}

	// 5. Begin graceful shutdown: stop accepting new connections and give
	//    in-flight requests up to 10 seconds to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		// 6. Deadline exceeded — force-close remaining connections.
		slog.Error("shutdown deadline exceeded, forcing close", "err", err)
		server.Close()
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}

// loadEnv reads a .env file and sets environment variables if not already set.
func loadEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, line := range splitLines(string(data)) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		key, val, ok := splitOnce(line, '=')
		if !ok {
			continue
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	return nil
}

// splitLines splits a string into non-empty lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// splitOnce splits a string on the first occurrence of sep.
func splitOnce(s string, sep byte) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}
