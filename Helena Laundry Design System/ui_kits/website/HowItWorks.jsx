/* HowItWorks.jsx — teal band, numbered 3-step process */
function HowItWorks({ onSchedule }) {
  const steps = [
    { n: "01", icon: "calendar-check", t: "Schedule a pickup", d: "Pick a day that works. We'll text when we're on the way." },
    { n: "02", icon: "washing-machine", t: "We wash & fold", d: "Sorted, washed in eco detergent, dried gentle, folded neat." },
    { n: "03", icon: "package-check", t: "Back to your door", d: "Fresh stacks delivered by Friday. Drop the chores, keep the day." },
  ];
  return (
    <section id="how-it-works" style={{ background: "var(--teal-deep)", padding: "84px 0", color: "#fff" }}>
      <div style={{ maxWidth: 1140, margin: "0 auto", padding: "0 32px" }}>
        <div style={{ textAlign: "center", marginBottom: 52 }}>
          <Eyebrow color="var(--ochre)">How it works</Eyebrow>
          <h2 style={{ fontFamily: "var(--font-display)", fontSize: 42, margin: "12px 0 0", color: "#fff" }}>
            Three steps. That's the whole thing.
          </h2>
        </div>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3,1fr)", gap: 28 }}>
          {steps.map((s) => (
            <div key={s.n} style={{ display: "flex", flexDirection: "column", gap: 14 }}>
              <div style={{ display: "flex", alignItems: "center", gap: 14 }}>
                <span style={{ fontFamily: "var(--font-display)", fontSize: 40, color: "var(--ochre)", lineHeight: 1 }}>{s.n}</span>
                <div style={{ width: 48, height: 48, borderRadius: "50%",
                              background: "color-mix(in oklab, #fff 14%, transparent)",
                              display: "flex", alignItems: "center", justifyContent: "center" }}>
                  <Icon name={s.icon} size={24} color="#fff" />
                </div>
              </div>
              <h3 style={{ fontFamily: "var(--font-display)", fontSize: 24, margin: 0, color: "#fff" }}>{s.t}</h3>
              <p style={{ fontSize: 15, lineHeight: 1.6, margin: 0, color: "color-mix(in oklab, #fff 86%, transparent)" }}>{s.d}</p>
            </div>
          ))}
        </div>
        <div style={{ display: "flex", justifyContent: "center", marginTop: 48 }}>
          <Button size="lg" icon="calendar-check" onClick={onSchedule}>Schedule your first pickup</Button>
        </div>
      </div>
    </section>
  );
}
Object.assign(window, { HowItWorks });
