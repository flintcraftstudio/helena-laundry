/* Footer.jsx */
function Footer({ onSchedule }) {
  return (
    <footer style={{ background: "var(--brown)", color: "var(--cream)", padding: "64px 0 36px" }}>
      <div style={{ maxWidth: 1140, margin: "0 auto", padding: "0 32px" }}>
        <div style={{ display: "grid", gridTemplateColumns: "1.4fr 1fr 1fr 1fr", gap: 32, alignItems: "start" }}>
          <div>
            <div style={{ display: "flex", alignItems: "center", gap: 12, marginBottom: 14 }}>
              <img src="../../assets/helena_laundry_badge.png" alt="" style={{ width: 48, height: 48 }} />
              <span style={{ fontFamily: "var(--font-display)", fontSize: 24, color: "var(--cream)" }}>
                Helena <span style={{ color: "var(--ochre)" }}>Laundry</span>
              </span>
            </div>
            <p style={{ fontSize: 14, lineHeight: 1.6, color: "color-mix(in oklab,var(--cream) 75%,transparent)", maxWidth: "32ch" }}>
              Wash & fold, pickup & delivery for the whole Helena valley. Strong hands, soft towels.
            </p>
          </div>
          <FootCol title="Services" links={["Wash & fold", "Pickup & delivery", "Specialty care", "Pricing"]} />
          <FootCol title="Company" links={["About us", "How it works", "Careers", "Contact"]} />
          <div>
            <h4 style={{ fontFamily: "var(--font-sans)", fontWeight: 700, fontSize: 14, textTransform: "uppercase",
                         letterSpacing: "0.1em", color: "var(--ochre)", margin: "0 0 14px" }}>Visit</h4>
            <div style={{ fontSize: 14, lineHeight: 1.7, color: "color-mix(in oklab,var(--cream) 80%,transparent)" }}>
              412 Last Chance Gulch<br />Helena, MT 59601<br />Mon–Sat · 7am–7pm<br />(406) 555-0142
            </div>
          </div>
        </div>
        <div style={{ margin: "40px 0 24px" }}>
          <div style={{ borderTop: "2px solid var(--ochre)" }} />
          <div style={{ borderTop: "1px solid var(--ochre)", marginTop: 4, opacity: 0.6 }} />
        </div>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", flexWrap: "wrap", gap: 12 }}>
          <span style={{ fontSize: 13, color: "color-mix(in oklab,var(--cream) 65%,transparent)" }}>
            Made with strong arms in Helena, Montana. © 2026 Helena Laundry.
          </span>
          <Button size="sm" variant="secondary" onClick={onSchedule}
            style={{ background: "var(--ochre)", color: "var(--brown)" }}>Schedule a pickup</Button>
        </div>
      </div>
    </footer>
  );
}

function FootCol({ title, links }) {
  return (
    <div>
      <h4 style={{ fontFamily: "var(--font-sans)", fontWeight: 700, fontSize: 14, textTransform: "uppercase",
                   letterSpacing: "0.1em", color: "var(--ochre)", margin: "0 0 14px" }}>{title}</h4>
      <div style={{ display: "flex", flexDirection: "column", gap: 9 }}>
        {links.map((l) => (
          <a key={l} href="#" onClick={(e) => e.preventDefault()}
             style={{ fontSize: 14, color: "color-mix(in oklab,var(--cream) 80%,transparent)", textDecoration: "none" }}
             onMouseEnter={(e) => (e.target.style.color = "var(--cream)")}
             onMouseLeave={(e) => (e.target.style.color = "color-mix(in oklab,var(--cream) 80%,transparent)")}>{l}</a>
        ))}
      </div>
    </div>
  );
}
Object.assign(window, { Footer });
