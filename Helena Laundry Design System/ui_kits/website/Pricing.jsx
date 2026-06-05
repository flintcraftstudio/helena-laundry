/* Pricing.jsx — pricing cards, one highlighted */
function PriceCard({ name, price, unit, desc, features, featured, onSchedule }) {
  const [h, setH] = React.useState(false);
  return (
    <div onMouseEnter={() => setH(true)} onMouseLeave={() => setH(false)}
      style={{
        background: featured ? "var(--rust)" : "var(--white)",
        color: featured ? "#fff" : "var(--fg1)",
        borderRadius: 20, padding: 30, position: "relative",
        boxShadow: featured ? "var(--shadow-rust)" : (h ? "var(--shadow-md)" : "var(--shadow-sm)"),
        transform: h && !featured ? "translateY(-3px)" : "none",
        transition: "all 240ms cubic-bezier(0.22,1,0.36,1)",
        border: featured ? "none" : "1px solid var(--rule)",
      }}>
      {featured && (
        <div style={{ position: "absolute", top: -13, left: 30 }}>
          <Tag tone="ochre">★ Most popular</Tag>
        </div>
      )}
      <div style={{ fontFamily: "var(--font-sans)", fontWeight: 700, fontSize: 15,
                    textTransform: "uppercase", letterSpacing: "0.08em",
                    color: featured ? "var(--ochre)" : "var(--teal)" }}>{name}</div>
      <div style={{ display: "flex", alignItems: "baseline", gap: 6, margin: "12px 0 6px" }}>
        <span style={{ fontFamily: "var(--font-display)", fontSize: 48, lineHeight: 1,
                       color: featured ? "#fff" : "var(--brown)" }}>{price}</span>
        <span style={{ fontSize: 15, color: featured ? "color-mix(in oklab,#fff 80%,transparent)" : "var(--fg3)" }}>{unit}</span>
      </div>
      <p style={{ fontSize: 14, lineHeight: 1.5, margin: "0 0 20px",
                  color: featured ? "color-mix(in oklab,#fff 86%,transparent)" : "var(--fg2)" }}>{desc}</p>
      <div style={{ display: "flex", flexDirection: "column", gap: 11, marginBottom: 24 }}>
        {features.map((f) => (
          <div key={f} style={{ display: "flex", alignItems: "center", gap: 10, fontSize: 14 }}>
            <Icon name="check" size={18} color={featured ? "var(--ochre)" : "var(--teal)"} />
            <span style={{ color: featured ? "color-mix(in oklab,#fff 92%,transparent)" : "var(--fg1)" }}>{f}</span>
          </div>
        ))}
      </div>
      <Button variant={featured ? "secondary" : "ghost"} onClick={onSchedule}
        style={featured ? { width: "100%", justifyContent: "center", background: "#fff", color: "var(--rust)" }
                        : { width: "100%", justifyContent: "center" }}>
        Choose {name}
      </Button>
    </div>
  );
}

function Pricing({ onSchedule }) {
  return (
    <section id="pricing" style={{ background: "var(--cream)", padding: "88px 0" }}>
      <div style={{ maxWidth: 1140, margin: "0 auto", padding: "0 32px" }}>
        <div style={{ textAlign: "center", marginBottom: 52 }}>
          <Eyebrow>Simple pricing</Eyebrow>
          <h2 style={{ fontFamily: "var(--font-display)", fontSize: 42, margin: "12px 0 0", color: "var(--brown)" }}>
            Honest rates, no surprises.
          </h2>
          <p style={{ fontSize: 17, color: "var(--fg2)", margin: "12px auto 0", maxWidth: "48ch" }}>
            Pay by the pound or settle into a weekly rhythm. Free pickup on every order over $35.
          </p>
        </div>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3,1fr)", gap: 24, alignItems: "start" }}>
          <PriceCard name="By the pound" price="$2.10" unit="/ lb" onSchedule={onSchedule}
            desc="Perfect for the occasional pile. Drop it off or schedule a one-time pickup."
            features={["48-hour turnaround", "Eco detergent included", "Folded & bagged"]} />
          <PriceCard name="Weekly" price="$36" unit="/ week" featured onSchedule={onSchedule}
            desc="Our neighborly favorite. A standing pickup so laundry never piles up again."
            features={["Up to 18 lbs each week", "Free pickup & delivery", "Priority Friday return", "Pause anytime"]} />
          <PriceCard name="Specialty" price="$9" unit="/ item" onSchedule={onSchedule}
            desc="Comforters, delicates, and the big stuff that needs a gentle, careful hand."
            features={["Hand-finished care", "Air-dry options", "Stain pre-treatment"]} />
        </div>
      </div>
    </section>
  );
}
Object.assign(window, { Pricing, PriceCard });
