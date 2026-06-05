# Helena Laundry — Design System

A complete brand + UI design system for **Helena Laundry**, a local wash-and-fold
and laundry service in Helena, Montana. Warm, friendly, vintage-Americana with a
clean modern execution — bright, fun, and human, never corporate or sterile.

> **The hero is the badge.** A circular Art Nouveau medallion of an auburn-haired
> woman in a floral headscarf, flexing (a neighborly Rosie-the-Riveter). She does
> the heavy lifting so the layout doesn't have to. Give her room to breathe;
> keep everything around her simple.

---

## Sources & provenance

- **Brand mark:** `uploads/helena_laundry_logo.png` (provided by the user, 2048×2048).
  Cleaned circular cut-out lives at `assets/helena_laundry_badge.png`; favicon at
  `assets/favicon.png`.
- **Brand brief:** supplied as text by the user (palette, type direction, tone,
  motifs, things to avoid). This README captures and extends it.
- No codebase or Figma file was provided — this system is built from the brand brief
  and logo. Where a choice was inferred (e.g. exact font pairing), it is flagged below.

---

## The product

A single core product: the **marketing website** for the laundry service —
booking wash-and-fold pickup/delivery, browsing services & pricing, and learning
about the neighborhood shop. The UI kit in `ui_kits/website/` recreates it.

There is no app or dashboard in scope. If one is added later, it should inherit
the same tokens from `colors_and_type.css`.

---

## CONTENT FUNDAMENTALS — how Helena Laundry talks

**Voice:** Cheerful, neighborly, plainspoken Montana-friendly. Confident and
capable (she's flexing!) but never shouty or salesy. Short, warm sentences.

- **Person:** Speaks as *we* ("We pick up Tuesdays"); addresses the customer as
  *you* ("Drop the chores, keep the day"). Warm and direct.
- **Casing:** Sentence case everywhere — headlines, buttons, labels. Reserve
  ALL-CAPS only for tiny tracked eyebrow labels (e.g. `WASH · FOLD · DELIVER`).
  Never shout in body copy.
- **Length:** Short. Headlines are a phrase, not a paragraph. Body in 1–2 sentence
  beats. Generous whitespace does the rest.
- **Punctuation:** Friendly, light. The occasional em-dash and ampersand. No
  exclamation-point spam — at most one, and only when genuinely cheerful.
- **Emoji:** None. The warmth comes from the badge, the palette, and the words —
  not emoji. (Decorative ✦ / · separators are fine in eyebrows.)
- **Numbers & claims:** Plain and honest. "Free pickup on orders over $35," not
  "UNBEATABLE SAVINGS." Montana-modest.

**Vibe words:** neighborly · capable · sunny · tactile · honest · unhurried.

**Sample copy (use as reference for tone):**
- Hero: *"Helena's laundry, lifted off your shoulders."*
- Sub: *"Wash, dry, fold, and back to your door by Friday. Strong hands, soft towels."*
- Button: *"Schedule a pickup"* / *"See our prices"*
- Eyebrow: *"WASH · FOLD · DELIVER"*
- Service blurb: *"We sort lights from darks like it's a sport. You just enjoy the
  fresh-folded stack."*
- Footer sign-off: *"Made with strong arms in Helena, Montana."*

**Avoid:** corporate jargon ("solutions," "leverage"), urgency/FOMO, all-caps
sentences, dense paragraphs, cold or clinical phrasing.

---

## VISUAL FOUNDATIONS

### Color
Warm/cool complementary pairing of **rust** (#C0492B) and **teal** (#3E7E8C),
grounded on **cream** (#F2EBD8), with **golden ochre** (#D99A3D) as a sparing pop
and **dark brown** (#3A2A1E) for all text and linework. Roughly **60% cream/white,
30% one accent, 10% the other + ochre.**

- **Cream is the page** — never stark white as the base. White is reserved for
  cards/panels that need to lift off the cream.
- **Concentrate color into accents** — buttons, links, key headings, dividers,
  icons. Keep large surfaces calm.
- **One bold hero moment** per page (usually a rust or teal band carrying the badge).
- **Never** pure black (#000), neon, cold grays, or corporate blue.
- Alternating sections may swap between cream and a deep-cream or a full teal/rust band.

### Type
- **Display / headings:** **DM Serif Display** — a friendly, high-contrast serif
  with Art-Nouveau-era flavor, clean modern cut. Used at large sizes; can be
  expressive and characterful. *(Substitution note: brief suggested Fraunces /
  Recoleta / Cormorant. Fraunces was on the avoid list; Recoleta is not free.
  DM Serif Display is a free, warm, high-contrast match. Swap if you license
  Recoleta — change `--font-display` in `colors_and_type.css`.)*
- **Body / UI:** **DM Sans** — clean, warm humanist sans, highly legible.
- **Two display fonts max.** Generous line spacing (body line-height 1.6).
- Eyebrow labels: DM Sans 700, uppercase, wide tracking (0.14em), in teal.

### Spacing & layout
- Generous whitespace, comfortable padding, light + airy. The warm palette stays
  clean because the spacing is roomy.
- 8px-based spacing scale (4 → 128). Section padding is large (64–128px vertical).
- Content max-width ~1140px; comfortable line lengths (~62ch for body).

### Corners & borders
- Soft rounded corners: **cards 16px, buttons pill (999px), insets 8–12px.**
- **Thin rust or brown rules** echo the badge's border ring. A signature
  **double-ring** device (`.rule-double`) repeats the medallion's twin border.
- Cards: white or warm near-white on cream, 16px radius, soft shadow, sometimes a
  hairline rust/brown border.

### Shadows & elevation
- Soft, warm, **brown-tinted** shadows — never harsh black. Three steps (sm/md/lg)
  plus a rust-tinted glow for primary buttons. No hard drop shadows, no neon glow.

### Backgrounds
- Primarily flat cream. No busy patterns, no heavy gradients. At most a very subtle
  paper-warmth. Hero band may be a solid teal or rust panel. Floral accents from the
  headscarf appear **very sparingly** (a corner flourish), never as a dense pattern.

### Imagery
- The badge is the central illustration — warm, slightly vintage line-art palette
  (auburn, teal, ochre, cream). Photography, if used, should be warm-toned, natural
  light, real and human — folded laundry, hands, the shop — never cold stock.

### Motion
- Gentle and tactile. Soft fades and short rises (8–12px) on scroll-in.
  Easing `cubic-bezier(0.22,1,0.36,1)`. Durations 140–420ms.
- Hover: buttons darken (rust→rust-deep) and lift slightly with a soft shadow.
- Press: small scale-down (~0.97) and a touch darker — tactile, hand-friendly.
- No bouncy/springy excess, no infinite looping decoration on content.

### Transparency & blur
- Used minimally. A translucent cream may sit over a teal band; otherwise solid
  surfaces. No heavy glassmorphism.

---

## ICONOGRAPHY

No bespoke icon set was provided. The system uses **Lucide** (lucide.dev) via CDN —
a friendly, rounded, consistent-stroke open-source set whose ~1.75px hand-drawn
warmth matches the brand's tactile feel. *(Substitution flagged: swap for a custom
set if Helena Laundry commissions one.)*

- **Stroke style:** 1.75–2px, round caps + round joins. Brown (`--fg1`) by default;
  teal for interactive/feature icons; rust for emphasis.
- **Sizing:** 20px inline, 24px UI, 28–32px feature blocks. Generous padding.
- **Emoji:** never used as iconography.
- **Brand motifs as icons:** a small flexed-arm and a laundry-basket mark recur as
  brand glyphs; the circular badge frame and double-ring border are decorative
  devices, not icons. Floral accents (from the headscarf) appear sparingly.
- Load: `<script src="https://unpkg.com/lucide@latest"></script>` then
  `lucide.createIcons()`, or use `<i data-lucide="washing-machine"></i>`.

---

## Files in this system (index)

| Path | What it is |
|---|---|
| `README.md` | This file — context, voice, visual foundations, iconography, index. |
| `SKILL.md` | Agent-Skills-compatible entry point. |
| `colors_and_type.css` | Source of truth: CSS vars for color, type, spacing, radius, shadow, motion + semantic base styles. **Import this everywhere.** |
| `assets/` | Brand assets — full logo, circular badge cut-out, favicon. |
| `preview/` | Small HTML cards that populate the Design System tab (colors, type, components, etc). |
| `ui_kits/website/` | High-fidelity recreation of the marketing website (`index.html` + JSX components). |

### assets/
- `helena_laundry_logo.png` — full square logo on cream (original).
- `helena_laundry_badge.png` — circular cut-out, transparent outside the ring (use on any surface).
- `favicon.png` — 128px circular favicon.
