/* Header.jsx — sticky top navigation */
function Header({ onSchedule, onNav }) {
  const [scrolled, setScrolled] = React.useState(false);
  React.useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 12);
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, []);
  const links = ["Services", "Pricing", "How it works", "About"];
  return (
    <header style={{
      position: "sticky", top: 0, zIndex: 50,
      background: scrolled ? "color-mix(in oklab, var(--cream) 88%, transparent)" : "var(--cream)",
      backdropFilter: scrolled ? "blur(8px)" : "none",
      borderBottom: scrolled ? "1px solid var(--rule)" : "1px solid transparent",
      transition: "all 240ms cubic-bezier(0.22,1,0.36,1)",
    }}>
      <div style={{
        maxWidth: 1140, margin: "0 auto", padding: "14px 32px",
        display: "flex", alignItems: "center", gap: 28,
      }}>
        <a href="#top" onClick={(e) => { e.preventDefault(); onNav && onNav("top"); }}
           style={{ display: "flex", alignItems: "center", gap: 12, textDecoration: "none" }}>
          <img src="../../assets/helena_laundry_badge.png" alt="" style={{ width: 44, height: 44 }} />
          <span style={{ fontFamily: "var(--font-display)", fontSize: 23, color: "var(--rust)", lineHeight: 1 }}>
            Helena <span style={{ color: "var(--brown)" }}>Laundry</span>
          </span>
        </a>
        <nav style={{ display: "flex", gap: 26, marginLeft: 8 }}>
          {links.map((l) => (
            <a key={l} href="#" onClick={(e) => { e.preventDefault(); onNav && onNav(l); }}
               style={{ fontFamily: "var(--font-sans)", fontSize: 15, fontWeight: 500,
                        color: "var(--fg1)", textDecoration: "none", transition: "color 140ms" }}
               onMouseEnter={(e) => (e.target.style.color = "var(--rust)")}
               onMouseLeave={(e) => (e.target.style.color = "var(--fg1)")}>{l}</a>
          ))}
        </nav>
        <div style={{ flex: 1 }} />
        <Button size="sm" onClick={onSchedule}>Schedule a pickup</Button>
      </div>
    </header>
  );
}
Object.assign(window, { Header });
