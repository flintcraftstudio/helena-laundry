# CLAUDE.md — Helena Laundry

Project-scoped guidance for `helena-laundry/`. Inherits the FlintCraft repo-root
`CLAUDE.md` (stack, commands, conventions) but **overrides its brand**: Helena
Laundry is a separate entity with its own voice, palette, and tokens.

- **Voice source of truth:** the `helena-laundry-brand-voice` skill +
  `helena-laundry-voice.md`. First-person "I" (Chanté, named on About only).
  Do NOT use the Firefly "we" voice or the design-system README's "we" samples.
  Use `helena-laundry-site-copy.md` copy verbatim.
- **Design tokens source of truth:** `Helena Laundry Design System/colors_and_type.css`
  and `tailwind/tailwind.config.js` (`hl-*` prefix). Warm vintage-Americana,
  light theme — NOT the FlintCraft dark `fc-*`/`ff-*` palette.
- **Full design guidelines:** `.impeccable.md` (read by the impeccable skills).

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
