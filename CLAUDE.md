# CLAUDE.md — Helena Laundry

Project-scoped guidance for `helena-laundry/`. Borrows conventions from the
FlintCraft repo-root `CLAUDE.md` but **overrides its brand AND its stack**:
Helena Laundry is a separate entity with its own voice, palette, and tokens, and
this is the Firefly **advanced tier** template — a superset of the root's
brochure stack that adds a database, migrations, sessions, auth, and a full admin
dashboard. Where the root says "stdlib only, no DB, no external frameworks,"
trust the **Architecture** section below instead (module path is
`github.com/firefly-software-mt/advanced-template`, not `standard-template`).

- **Voice source of truth:** the `helena-laundry-brand-voice` skill +
  `helena-laundry-voice.md`. First-person "I" (Chanté, named on About only).
  Do NOT use the Firefly "we" voice or the design-system README's "we" samples.
  Use `helena-laundry-site-copy.md` copy verbatim.
- **Design tokens source of truth:** `Helena Laundry Design System/colors_and_type.css`
  and `tailwind/tailwind.config.js` (`hl-*` prefix). Warm vintage-Americana,
  light theme — NOT the FlintCraft dark `fc-*`/`ff-*` palette.
- **Full design guidelines:** `.impeccable.md` (read by the impeccable skills).

## Architecture

Go `net/http` (no web framework) + templ + htmx + Alpine.js + Tailwind **v4** +
[`flint-ui`](https://github.com/flintcraftstudio/flint-ui) components, backed by
SQLite. The README covers setup; this section is the parts you can't see from one
file.

**Build/dev (Mage, not Make).** `mage Dev` = `Build` (`BuildCSS` + `Generate` +
`go build`) then `Run`. Two codegen passes feed the build — **always regenerate
after editing their sources, or you'll build stale code:**
- `mage GenerateTempl` — `.templ` → `*_templ.go` (run after editing any `.templ`).
- `mage GenerateSqlc` — `queries/*.sql` → `internal/db/` (run after editing
  queries). `mage Generate` runs both.
- `mage BuildCSS` — Tailwind v4 is **CSS-first**: no `-c` config flag; theme and
  `@source` globs live in `tailwind/input.css`. `BuildCSS` first writes the
  gitignored `tailwind/flint-source.css` pointing Tailwind at the resolved
  flint-ui module dir so its utility classes (kept in flint-ui's `*.go`, not
  `.templ`) get compiled. Output → `web/static/css/site.css` (gitignored).

**Database & migrations.** SQLite via `modernc.org/sqlite` (pure Go, no CGo),
WAL + foreign keys on. Migrations are **goose** SQL files in `migrations/`,
applied two ways: `main.go` runs `goose.Up` on every startup, and `mage MigrateUp`
/ `MigrateDown` / `MigrateStatus` / `CreateMigration <name>` for manual control.
`mage Seed <email> <password>` (→ `cmd/seed`) creates an admin user.

**Data access — read this before touching the DB.** There are *two* layers and
they don't overlap. `internal/db/` is **sqlc-generated** but currently only
covers `users` + `sessions` and **is not imported anywhere** — effectively dead.
All live queries are **hand-written `database/sql`** on `*store.Store`
(`internal/store/`: `store.go` writes/lists, `admin.go` get/update + dashboard
counts, `settings.go` the settings singleton). Follow the hand-written pattern
(parameterized SQL, a `scanX` helper, domain structs) — do **not** assume sqlc.
The `bookings` "needs a date" → calendar flow is pure SQL
(`ListUnscheduledBookings` / `ListBookingsScheduledBetween`), no schema for it.

**Settings cache.** `internal/settings/` holds an in-memory `sync.RWMutex`
snapshot of the single `settings` row (service-area ZIPs, pickup days, pricing,
`accepting_bookings`, business identity), loaded once at startup. The booking
flow reads it via `settings.Get()` / `IsZipAllowed` / `IsPickupDay` (no DB hit
per request). On an admin save you must refresh **both** sinks: `settings.Set(s)`
*and* `view.ApplySettings(s)` (the latter pushes identity strings into `view`
package vars) — same two calls `main.go` makes at boot. Service area and pricing
are **DB-backed now, not env vars** — `.env.example` notes this.

**Auth & sessions.** Cookie-based, in `internal/session/`. `session.Middleware`
wraps the whole mux, loads the `session_token` cookie, and attaches a `*User` to
the request context (`session.FromContext`). `session.RequireAuth` gates every
`/admin/*` route and is htmx-aware: unauthenticated htmx requests get an
`HX-Redirect: /login` on a 200 (a 303 would be swallowed by the XHR and swap the
login page into the target). Passwords are bcrypt; tokens are 32 random bytes.

**Routing & handlers.** Go 1.22 method+pattern mux in `cmd/server/main.go` —
that file is the route map. Handlers (`internal/handler/`) are closures that take
their deps (`*store.Store`, mailer, Turnstile secret) as args. Public pages, three
distinct contact-capture POSTs (`/contact/question|waitlist|booking`), auth, and
the `RequireAuth`-gated admin dashboard (bookings, calendar, inquiries, waitlist,
settings). Most admin interactions are htmx partial swaps — handlers render a row
/ edit panel / region fragment, not a full page.

**Graceful degradation.** Postmark, Turnstile, GA, and Pixel are all optional;
missing env vars log a warning and disable the feature — form submissions are
still persisted. Server does timed graceful shutdown on SIGINT/SIGTERM.

## Design Context

### Users

Two audiences, one action — **book a pickup**. (1) Helena Valley **residents**,
mostly on phones mid-errand, deciding whether to trust one local person with
their clothes. (2) Local **businesses** (salons, spas, short-term rentals, gyms)
wanting standing pickups, bulk rates, and reliable turnaround. Local, low-pressure,
warm. Helena has many transplants and a large Air Force community — celebrate
*local human work*, never resent newcomers (the owner was a military spouse).

### Brand Personality

**Spunky, no-nonsense Montana laundress — Beth Dutton business sense fused with
Mrs. Weasley warmth.** Capable · neighborly · honest. First-person singular "I":
she's one real person, the moat against faceless apps. Lead with the answer (say
the price out loud, no apology), then the warmth. Specific over vague. Warmth =
competence and care, never gushing. Dry humor once per piece; rationed
exclamation points; light Montana/Irish texture, not costume.

### Aesthetic Direction

**Warm vintage-Americana, light theme, clean modern execution** — bright, human,
tactile; never corporate or sterile. Light mode only.

- **Hero is the badge** — circular Art-Nouveau medallion (auburn woman in floral
  headscarf, flexing) with a CSS double-ring. Keep everything around it simple.
- **Palette (`hl-*`):** cream `#F2EBD8` is the page (white only for cards that
  lift off it); rust `#C0492B` primary accent; teal `#3E7E8C` secondary; ochre
  `#D99A3D` sparing pop; brown `#3A2A1E` for ALL text/linework. ~60% cream/white
  · 30% one accent · 10% the other + ochre. Never pure black, neon, cold gray,
  or corporate blue.
- **Type:** DM Serif Display (headings) + DM Sans (body). Sentence case; ALL-CAPS
  only for tiny tracked teal eyebrow labels. Generous line-height (1.6 body).
- **Shape:** soft corners (cards 16px, pill buttons, insets 8–12px), thin
  rust/brown rules + signature double-ring, brown-tinted soft shadows, flat cream
  backgrounds. Floral accents only as rare corner flourishes.
- **Icons:** inline SVG templ components (no Lucide CDN/JS). No emoji as icons.
- **Anti-references:** SaaS landing pages, greeting cards, faceless laundry apps,
  cold stock photography, jargon, urgency/FOMO, exclamation spam.

### Design Principles

1. **One real person, said out loud.** Every screen feels like it comes from
   Chanté — first-person, direct, price stated plainly. Never flatten her into a
   faceless company.
2. **Cream breathes; color concentrates.** Calm, roomy surfaces (64–128px section
   padding, ~62ch lines); pour color into accents and one hero band per page.
3. **Competence first, warmth always.** In a trade-off, lean toward looking
   capable and confident (it earns the price and trust); let warmth show through
   care and specificity, not decoration.
4. **Tactile, not flashy.** Soft corners, brown-tinted shadows, gentle motion.
   Restrained by spec, with ONE signature beat — the badge double-ring drawing in
   on load. No parallax, no looping decoration. Always honor
   `prefers-reduced-motion`.
5. **Honest and accessible by default.** Plain numbers, no FOMO, **WCAG 2.1 AA**
   contrast (ochre and light teal are decorative/large-text only — never small
   body text on cream), reduced-motion respected. Montana-modest in copy and craft.
