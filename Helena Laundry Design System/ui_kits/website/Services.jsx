/* Services.jsx — three feature cards */
function ServiceCard({ icon, title, body, tone }) {
  const [h, setH] = React.useState(false);
  const tint = tone === "rust"
    ? "color-mix(in oklab, var(--rust) 12%, var(--white))"
    : tone === "ochre"
    ? "color-mix(in oklab, var(--ochre) 20%, var(--white))"
    : "color-mix(in oklab, var(--teal) 14%, var(--white))";
  const ic = tone === "rust" ? "var(--rust)" : tone === "ochre" ? "var(--ochre-deep)" : "var(--teal)";
  return (
    <div onMouseEnter={() => setH(true)} onMouseLeave={() => setH(false)}
      style={{
        background: "var(--white)", borderRadius: 16, padding: 28,
        boxShadow: h ? "var(--shadow-lg)" : "var(--shadow-sm)",
        transform: h ? "translateY(-4px)" : "none",
        transition: "all 240ms cubic-bezier(0.22,1,0.36,1)",
        display: "flex", flexDirection: "column", gap: 14,
      }}>
      <div style={{ width: 52, height: 52, borderRadius: 14, background: tint,
                    display: "flex", alignItems: "center", justifyContent: "center" }}>
        <Icon name={icon} size={26} color={ic} />
      </div>
      <h3 style={{ fontFamily: "var(--font-display)", fontSize: 25, margin: 0, color: "var(--brown)" }}>{title}</h3>
      <p style={{ fontSize: 15, lineHeight: 1.6, color: "var(--fg2)", margin: 0 }}>{body}</p>
    </div>
  );
}

function Services() {
  return (
    <section id="services" style={{ background: "var(--cream)", padding: "88px 0" }}>
      <div style={{ maxWidth: 1140, margin: "0 auto", padding: "0 32px" }}>
        <div style={{ display: "flex", flexDirection: "column", alignItems: "center", textAlign: "center", marginBottom: 48 }}>
          <Eyebrow>What we do</Eyebrow>
          <h2 style={{ fontFamily: "var(--font-display)", fontSize: 42, margin: "12px 0 0", color: "var(--brown)" }}>
            Laundry day, handled.
          </h2>
          <DoubleRule width={72} style={{ marginTop: 18 }} />
        </div>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3,1fr)", gap: 24 }}>
          <ServiceCard tone="rust" icon="washing-machine" title="Wash & fold"
            body="We sort lights from darks like it's a sport. By the pound, washed in eco detergent, and folded into a fresh stack." />
          <ServiceCard tone="teal" icon="truck" title="Pickup & delivery"
            body="We pick up Tuesdays and have it back to your door by Friday. Free on orders over $35 — no need to be home." />
          <ServiceCard tone="ochre" icon="shirt" title="Specialty care"
            body="Comforters, delicates, and game-day jerseys. The things that don't fit a regular load, treated with extra care." />
        </div>
      </div>
    </section>
  );
}
Object.assign(window, { Services, ServiceCard });
