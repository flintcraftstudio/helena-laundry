/* Testimonials.jsx — neighbor quotes */
function Testimonials() {
  const quotes = [
    { q: "They picked up Tuesday, dropped it back Friday, folded better than I ever do. My weekend is mine again.", n: "Marie L.", r: "Helena · West Side" },
    { q: "Strong hands, soft towels — they weren't kidding. The whole house smells like a fresh start.", n: "Dale & Pat", r: "Helena · Downtown" },
    { q: "Friendly every single time. It feels like handing the chore to a capable neighbor, because it is.", n: "Sofia R.", r: "East Helena" },
  ];
  return (
    <section id="about" style={{ background: "var(--cream-deep)", padding: "84px 0" }}>
      <div style={{ maxWidth: 1140, margin: "0 auto", padding: "0 32px" }}>
        <div style={{ textAlign: "center", marginBottom: 48 }}>
          <Eyebrow>From the neighborhood</Eyebrow>
          <h2 style={{ fontFamily: "var(--font-display)", fontSize: 42, margin: "12px 0 0", color: "var(--brown)" }}>
            Helena keeps coming back.
          </h2>
        </div>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3,1fr)", gap: 24 }}>
          {quotes.map((c) => (
            <div key={c.n} style={{ background: "var(--white)", borderRadius: 16, padding: 28,
                                    boxShadow: "var(--shadow-sm)", display: "flex", flexDirection: "column", gap: 16 }}>
              <Stars n={5} />
              <p style={{ fontFamily: "var(--font-display)", fontSize: 21, lineHeight: 1.4,
                          color: "var(--brown)", margin: 0 }}>"{c.q}"</p>
              <div style={{ marginTop: "auto", display: "flex", alignItems: "center", gap: 12 }}>
                <div style={{ width: 40, height: 40, borderRadius: "50%", background: "var(--teal)",
                              color: "#fff", display: "flex", alignItems: "center", justifyContent: "center",
                              fontWeight: 700, fontFamily: "var(--font-sans)" }}>{c.n[0]}</div>
                <div>
                  <div style={{ fontWeight: 700, fontSize: 14, color: "var(--fg1)" }}>{c.n}</div>
                  <div style={{ fontSize: 13, color: "var(--fg3)" }}>{c.r}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
Object.assign(window, { Testimonials });
