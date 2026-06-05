/* Hero.jsx — the bold brand moment */
function Hero({ onSchedule, onPrices }) {
  return (
    <section id="top" style={{ background: "var(--cream)", overflow: "hidden" }}>
      <div style={{
        maxWidth: 1140, margin: "0 auto", padding: "72px 32px 88px",
        display: "grid", gridTemplateColumns: "1.05fr 0.95fr", gap: 48, alignItems: "center",
      }}>
        <div>
          <Eyebrow>Wash · Fold · Deliver</Eyebrow>
          <h1 style={{
            fontFamily: "var(--font-display)", fontSize: "clamp(3rem, 5.2vw, 4.5rem)",
            lineHeight: 1.02, letterSpacing: "-0.01em", margin: "16px 0 0", color: "var(--brown)",
          }}>
            Helena's laundry,<br /><span style={{ color: "var(--rust)" }}>lifted off</span> your shoulders.
          </h1>
          <p style={{ fontSize: 19, lineHeight: 1.6, color: "var(--fg2)", margin: "20px 0 0", maxWidth: "44ch" }}>
            Wash, dry, fold, and back to your door by Friday. Strong hands, soft towels —
            and your whole weekend handed back to you.
          </p>
          <div style={{ display: "flex", gap: 14, marginTop: 32, flexWrap: "wrap" }}>
            <Button size="lg" icon="calendar-check" onClick={onSchedule}>Schedule a pickup</Button>
            <Button size="lg" variant="ghost" onClick={onPrices}>See our prices</Button>
          </div>
          <div style={{ display: "flex", alignItems: "center", gap: 12, marginTop: 28 }}>
            <Stars n={5} size={18} />
            <span style={{ fontSize: 14, color: "var(--fg2)" }}>
              <b style={{ color: "var(--fg1)" }}>4.9</b> from 300+ Helena neighbors
            </span>
          </div>
        </div>

        <div style={{ position: "relative", display: "flex", justifyContent: "center" }}>
          <div style={{
            position: "absolute", width: 380, height: 380, borderRadius: "50%",
            background: "radial-gradient(circle at 50% 45%, color-mix(in oklab, var(--teal) 22%, var(--cream)), var(--cream) 72%)",
          }} />
          <div style={{
            position: "absolute", width: 360, height: 360, borderRadius: "50%",
            border: "2px solid var(--rust)", boxShadow: "inset 0 0 0 5px var(--cream), inset 0 0 0 7px var(--rust)",
          }} />
          <img src="../../assets/helena_laundry_badge.png" alt="Helena Laundry badge"
               style={{ position: "relative", width: 330, height: 330, filter: "drop-shadow(0 16px 30px rgba(58,42,30,0.18))" }} />
          <div style={{
            position: "absolute", bottom: 8, right: 18, background: "var(--white)",
            borderRadius: 999, padding: "9px 16px", boxShadow: "var(--shadow-md)",
            display: "flex", alignItems: "center", gap: 8, fontSize: 14, fontWeight: 600, color: "var(--fg1)",
          }}>
            <Icon name="truck" size={18} color="var(--teal)" /> Free pickup over $35
          </div>
        </div>
      </div>
    </section>
  );
}
Object.assign(window, { Hero });
