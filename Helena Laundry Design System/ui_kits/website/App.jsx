/* App.jsx — assembles the full Helena Laundry marketing site */
function App() {
  const [booking, setBooking] = React.useState(false);
  const openBooking = () => setBooking(true);

  const scrollTo = (key) => {
    const map = { "Services": "services", "Pricing": "pricing", "How it works": "how-it-works", "About": "about", "top": "top" };
    const id = map[key];
    const el = id && document.getElementById(id);
    if (el) {
      const y = el.getBoundingClientRect().top + window.scrollY - 70;
      window.scrollTo({ top: Math.max(0, y), behavior: "smooth" });
    }
  };

  return (
    <React.Fragment>
      <Header onSchedule={openBooking} onNav={scrollTo} />
      <Hero onSchedule={openBooking} onPrices={() => scrollTo("Pricing")} />
      <Services />
      <HowItWorks onSchedule={openBooking} />
      <Pricing onSchedule={openBooking} />
      <Testimonials />
      <Footer onSchedule={openBooking} />
      <BookingModal open={booking} onClose={() => setBooking(false)} />
    </React.Fragment>
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(<App />);
