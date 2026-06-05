/* Primitives.jsx — shared building blocks for the Helena Laundry site */
const { useState } = React;

function Icon({ name, size = 24, color, style }) {
  const ref = React.useRef(null);
  React.useEffect(() => {
    if (ref.current && window.lucide) {
      ref.current.innerHTML = "";
      const el = document.createElement("i");
      el.setAttribute("data-lucide", name);
      ref.current.appendChild(el);
      window.lucide.createIcons({ attrs: { width: size, height: size } });
    }
  }, [name, size]);
  return <span ref={ref} style={{ display: "inline-flex", color: color || "inherit", ...style }} />;
}

function Eyebrow({ children, color = "var(--teal)", style }) {
  return (
    <div style={{
      fontFamily: "var(--font-sans)", fontWeight: 700, fontSize: 13,
      textTransform: "uppercase", letterSpacing: "0.14em", color, ...style
    }}>{children}</div>
  );
}

function Button({ children, variant = "primary", size = "md", icon, onClick, style }) {
  const [h, setH] = useState(false);
  const [p, setP] = useState(false);
  const pads = size === "lg" ? "16px 32px" : size === "sm" ? "10px 18px" : "13px 26px";
  const fs = size === "lg" ? 17 : size === "sm" ? 14 : 15;
  const base = {
    fontFamily: "var(--font-sans)", fontWeight: 600, fontSize: fs,
    borderRadius: 999, padding: pads, border: "1.5px solid transparent",
    cursor: "pointer", display: "inline-flex", alignItems: "center", gap: 8,
    transition: "all 160ms cubic-bezier(0.22,1,0.36,1)",
    transform: p ? "scale(0.97)" : "scale(1)", whiteSpace: "nowrap",
  };
  const variants = {
    primary: {
      background: h ? "var(--rust-deep)" : "var(--rust)", color: "var(--cream)",
      boxShadow: h ? "var(--shadow-rust)" : "0 4px 14px rgba(192,73,43,0.18)",
    },
    secondary: {
      background: h ? "var(--teal-deep)" : "var(--teal)", color: "var(--cream)",
    },
    ghost: {
      background: h ? "color-mix(in oklab, var(--rust) 8%, transparent)" : "transparent",
      color: "var(--rust)", borderColor: "var(--rust)",
    },
    text: { background: "transparent", color: h ? "var(--rust-deep)" : "var(--rust)", padding: "13px 8px" },
  };
  return (
    <button onClick={onClick} style={{ ...base, ...variants[variant], ...style }}
      onMouseEnter={() => setH(true)} onMouseLeave={() => { setH(false); setP(false); }}
      onMouseDown={() => setP(true)} onMouseUp={() => setP(false)}>
      {icon && <Icon name={icon} size={fs + 3} />}
      {children}
    </button>
  );
}

function Tag({ children, tone = "ochre", style }) {
  const tones = {
    ochre: { background: "var(--ochre)", color: "var(--brown)" },
    teal: { background: "color-mix(in oklab, var(--teal) 16%, var(--white))", color: "var(--teal-deep)" },
    rust: { background: "color-mix(in oklab, var(--rust) 12%, var(--white))", color: "var(--rust)" },
    outline: { background: "transparent", color: "var(--fg2)", border: "1.5px solid var(--rule-strong)" },
  };
  return (
    <span style={{
      fontFamily: "var(--font-sans)", fontWeight: 700, fontSize: 12, letterSpacing: "0.04em",
      borderRadius: 999, padding: "6px 14px", display: "inline-flex", alignItems: "center", gap: 6,
      ...tones[tone], ...style
    }}>{children}</span>
  );
}

function Stars({ n = 5, size = 18 }) {
  return (
    <span style={{ color: "var(--ochre)", letterSpacing: 2, fontSize: size, lineHeight: 1 }}>
      {"★".repeat(n)}
    </span>
  );
}

function DoubleRule({ width = 64, color = "var(--rust)", style }) {
  return (
    <div style={{ width, ...style }}>
      <div style={{ borderTop: `2px solid ${color}` }} />
      <div style={{ borderTop: `1px solid ${color}`, marginTop: 4 }} />
    </div>
  );
}

Object.assign(window, { Icon, Eyebrow, Button, Tag, Stars, DoubleRule });
