# Helena Laundry — Website UI Kit

A high-fidelity, interactive recreation of the Helena Laundry marketing website.
Built with React (via Babel standalone) in modular JSX files, all consuming the
shared tokens from the root `colors_and_type.css`.

## Run it
Open `index.html`. It renders the full one-page marketing site with a working,
click-through **"Schedule a pickup"** flow (a 3-step modal: details → plan →
confirmation).

## Components
| File | Exports | Purpose |
|---|---|---|
| `Primitives.jsx` | `Button`, `Eyebrow`, `Tag`, `Stars`, `Icon`, `DoubleRule`, `Badge` | Shared building blocks. |
| `Header.jsx` | `Header` | Sticky top nav with badge lockup + CTA. |
| `Hero.jsx` | `Hero` | Hero with the badge as the bold moment + dual CTAs. |
| `Services.jsx` | `Services` | Three feature cards (wash & fold, pickup, specialty). |
| `HowItWorks.jsx` | `HowItWorks` | Teal band, numbered 3-step process. |
| `Pricing.jsx` | `Pricing` | Pricing cards with one highlighted plan. |
| `Testimonials.jsx` | `Testimonials` | Neighbor quotes + star ratings. |
| `Booking.jsx` | `BookingModal` | Interactive 3-step pickup scheduler. |
| `Footer.jsx` | `Footer` | Sign-off, hours, contact, double-ring rule. |
| `App.jsx` | mounts everything | Page assembly + modal/nav state. |

## Notes
- Icons: **Lucide** via CDN (substitution — see root README ICONOGRAPHY).
- This is a cosmetic recreation: the booking flow is faked (no network), inputs are
  illustrative. It demonstrates layout, components, states, and interactions — not
  production logic.
- All color/type/spacing comes from `../../colors_and_type.css`. Change tokens there
  and the whole kit updates.
