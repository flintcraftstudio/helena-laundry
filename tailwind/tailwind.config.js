/** @type {import('tailwindcss').Config} */
const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  content: [
    "./internal/view/**/*.templ",
  ],
  theme: {
    extend: {
      colors: {
        // Helena Laundry — warm vintage-Americana palette
        hl: {
          cream:        "#F2EBD8", // primary page background
          "cream-deep": "#EAE0C7", // alt bands / insets
          tint:         "#FBF7EC", // near-white warm card
          white:        "#FFFFFF",
          rust:         "#C0492B", // primary accent: buttons, links, key headings
          "rust-deep":  "#A23B22", // rust hover / pressed
          teal:         "#3E7E8C", // secondary accent: eyebrows, icons, alt sections
          "teal-deep":  "#336774",
          ochre:        "#D99A3D", // tertiary pop: tags, stars
          "ochre-deep": "#C2862C",
          brown:        "#3A2A1E", // all body text + linework (never pure black)
          "brown-soft": "#6B5848", // secondary text
          "brown-faint":"#9A8A78", // tertiary / captions / placeholders
        },
      },
      fontFamily: {
        display: ['"DM Serif Display"', "Georgia", ...defaultTheme.fontFamily.serif],
        body:    ['"DM Sans"', ...defaultTheme.fontFamily.sans],
      },
      fontSize: {
        // Brand type scale (tokens from the design system). Overrides Tailwind
        // defaults at xl/2xl/3xl so headings align to the system, and adds three
        // fluid display sizes so hero/section heads stop using one-off clamps.
        sm:   ["0.875rem", { lineHeight: "1.5" }],
        base: ["1rem",     { lineHeight: "1.6" }],
        lg:   ["1.125rem", { lineHeight: "1.65" }],
        xl:   ["1.375rem", { lineHeight: "1.4" }],
        "2xl":["1.75rem",  { lineHeight: "1.2" }],
        "3xl":["2.25rem",  { lineHeight: "1.1" }],
        "display-lg": ["clamp(2rem, 4vw, 3rem)",       { lineHeight: "1.1",  letterSpacing: "-0.01em" }],
        "display-xl": ["clamp(2.1rem, 4.6vw, 3.5rem)", { lineHeight: "1.08", letterSpacing: "-0.01em" }],
        hero:         ["clamp(2.6rem, 6vw, 4.4rem)",   { lineHeight: "1.04", letterSpacing: "-0.015em" }],
      },
      borderRadius: {
        // soft, rounded — cards 16px, pill buttons handled by rounded-full
        pill: "999px",
      },
      boxShadow: {
        // soft, warm, brown-tinted — never harsh black
        "hl-sm":   "0 2px 6px rgba(58,42,30,0.08)",
        "hl-md":   "0 6px 18px rgba(58,42,30,0.10)",
        "hl-lg":   "0 16px 40px rgba(58,42,30,0.12)",
        "hl-rust": "0 8px 22px rgba(192,73,43,0.22)",
      },
      letterSpacing: {
        caps: "0.14em",
      },
      maxWidth: {
        content: "1140px",
      },
      transitionTimingFunction: {
        "hl-out": "cubic-bezier(0.22, 1, 0.36, 1)",
      },
    },
  },
  plugins: [],
}
