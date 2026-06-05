/* Booking.jsx — interactive 3-step pickup scheduler modal */
const { useState: useStateB } = React;

function BookingModal({ open, onClose }) {
  const [step, setStep] = useStateB(0);
  const [name, setName] = useStateB("");
  const [day, setDay] = useStateB("Tuesday");
  const [plan, setPlan] = useStateB("Weekly");

  React.useEffect(() => { if (open) { setStep(0); setName(""); setDay("Tuesday"); setPlan("Weekly"); } }, [open]);
  if (!open) return null;

  const days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"];
  const plans = [
    { k: "By the pound", d: "$2.10 / lb · one-time" },
    { k: "Weekly", d: "$36 / week · standing pickup" },
    { k: "Specialty", d: "$9 / item · gentle care" },
  ];

  const Field = ({ label, children }) => (
    <label style={{ display: "flex", flexDirection: "column", gap: 7 }}>
      <span style={{ fontSize: 13, fontWeight: 600, color: "var(--fg1)" }}>{label}</span>
      {children}
    </label>
  );

  return (
    <div onClick={onClose} style={{
      position: "fixed", inset: 0, zIndex: 100, background: "rgba(58,42,30,0.42)",
      backdropFilter: "blur(3px)", display: "flex", alignItems: "center", justifyContent: "center", padding: 20,
      animation: "hl-fade 200ms ease",
    }}>
      <div onClick={(e) => e.stopPropagation()} style={{
        background: "var(--cream)", borderRadius: 24, width: "min(520px, 100%)",
        boxShadow: "var(--shadow-lg)", overflow: "hidden", animation: "hl-rise 280ms cubic-bezier(0.22,1,0.36,1)",
      }}>
        {/* header */}
        <div style={{ background: "var(--rust)", padding: "22px 28px", display: "flex", alignItems: "center", gap: 14 }}>
          <img src="../../assets/helena_laundry_badge.png" alt="" style={{ width: 44, height: 44 }} />
          <div style={{ flex: 1 }}>
            <div style={{ fontFamily: "var(--font-display)", fontSize: 22, color: "#fff", lineHeight: 1 }}>
              {step < 2 ? "Schedule a pickup" : "You're all set!"}
            </div>
            <div style={{ fontSize: 13, color: "color-mix(in oklab,#fff 80%,transparent)", marginTop: 3 }}>
              {step < 2 ? `Step ${step + 1} of 2` : "We'll text you a confirmation"}
            </div>
          </div>
          <button onClick={onClose} style={{ background: "color-mix(in oklab,#fff 16%,transparent)", border: "none",
            color: "#fff", width: 34, height: 34, borderRadius: "50%", cursor: "pointer", fontSize: 18, lineHeight: 1 }}>×</button>
        </div>

        <div style={{ padding: 28 }}>
          {step === 0 && (
            <div style={{ display: "flex", flexDirection: "column", gap: 18 }}>
              <Field label="Your name">
                <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Jane from Helena"
                  style={inputStyleB} />
              </Field>
              <Field label="Pickup day">
                <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
                  {days.map((d) => (
                    <button key={d} onClick={() => setDay(d)} style={{
                      fontFamily: "var(--font-sans)", fontSize: 14, fontWeight: 600, cursor: "pointer",
                      borderRadius: 999, padding: "9px 16px",
                      border: `1.5px solid ${day === d ? "var(--rust)" : "var(--rule-strong)"}`,
                      background: day === d ? "var(--rust)" : "transparent",
                      color: day === d ? "#fff" : "var(--fg2)", transition: "all 140ms",
                    }}>{d.slice(0, 3)}</button>
                  ))}
                </div>
              </Field>
              <Button onClick={() => setStep(1)} style={{ width: "100%", justifyContent: "center", marginTop: 4 }}>
                Continue
              </Button>
            </div>
          )}

          {step === 1 && (
            <div style={{ display: "flex", flexDirection: "column", gap: 14 }}>
              <span style={{ fontSize: 13, fontWeight: 600, color: "var(--fg1)" }}>Choose a plan</span>
              {plans.map((p) => (
                <button key={p.k} onClick={() => setPlan(p.k)} style={{
                  display: "flex", alignItems: "center", gap: 14, textAlign: "left", cursor: "pointer",
                  borderRadius: 14, padding: "14px 16px", background: "var(--white)",
                  border: `1.5px solid ${plan === p.k ? "var(--teal)" : "var(--rule)"}`,
                  boxShadow: plan === p.k ? "0 0 0 3px color-mix(in oklab,var(--teal) 20%,transparent)" : "none",
                  transition: "all 140ms",
                }}>
                  <span style={{ width: 22, height: 22, borderRadius: "50%", flex: "none",
                    border: `2px solid ${plan === p.k ? "var(--teal)" : "var(--rule-strong)"}`,
                    background: plan === p.k ? "var(--teal)" : "transparent",
                    display: "flex", alignItems: "center", justifyContent: "center" }}>
                    {plan === p.k && <Icon name="check" size={13} color="#fff" />}
                  </span>
                  <span style={{ flex: 1 }}>
                    <span style={{ display: "block", fontWeight: 700, fontSize: 15, color: "var(--fg1)" }}>{p.k}</span>
                    <span style={{ display: "block", fontSize: 13, color: "var(--fg3)" }}>{p.d}</span>
                  </span>
                </button>
              ))}
              <div style={{ display: "flex", gap: 10, marginTop: 6 }}>
                <Button variant="ghost" onClick={() => setStep(0)} style={{ flex: "none" }}>Back</Button>
                <Button onClick={() => setStep(2)} style={{ flex: 1, justifyContent: "center" }}>Confirm pickup</Button>
              </div>
            </div>
          )}

          {step === 2 && (
            <div style={{ display: "flex", flexDirection: "column", alignItems: "center", textAlign: "center", gap: 14, padding: "8px 0" }}>
              <div style={{ width: 64, height: 64, borderRadius: "50%", background: "color-mix(in oklab,var(--teal) 16%,var(--white))",
                            display: "flex", alignItems: "center", justifyContent: "center" }}>
                <Icon name="check" size={34} color="var(--teal)" />
              </div>
              <h3 style={{ fontFamily: "var(--font-display)", fontSize: 28, margin: 0, color: "var(--brown)" }}>
                See you {day}{name ? `, ${name.split(" ")[0]}` : ""}!
              </h3>
              <p style={{ fontSize: 15, color: "var(--fg2)", margin: 0, maxWidth: "34ch", lineHeight: 1.6 }}>
                Your <b style={{ color: "var(--rust)" }}>{plan}</b> pickup is booked. We'll text a reminder the morning of —
                just leave your bag by the door.
              </p>
              <Button onClick={onClose} style={{ marginTop: 6 }}>Done</Button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

const inputStyleB = {
  fontFamily: "var(--font-sans)", fontSize: 15, color: "var(--fg1)", background: "var(--white)",
  border: "1.5px solid var(--rule-strong)", borderRadius: 12, padding: "12px 14px", outline: "none", width: "100%",
};

Object.assign(window, { BookingModal });
