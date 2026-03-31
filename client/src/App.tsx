import { useState, useEffect } from "react";
import { Component } from "@/components/ui/gradient-bars-background";
import { Copy, Check } from "lucide-react";

const FULL_TITLE = "Master your macOS environment";

export default function App() {
  const [copied, setCopied] = useState(false);
  const [displayedTitle, setDisplayedTitle] = useState("");
  const [showCursor, setShowCursor] = useState(true);
  const [typingDone, setTypingDone] = useState(false);

  const command = "curl -fsSL https://package-mate.com/install.sh | bash";

  useEffect(() => {
    let i = 0;
    const delay = 2000 / FULL_TITLE.length;
    const typer = setInterval(() => {
      i++;
      setDisplayedTitle(FULL_TITLE.slice(0, i));
      if (i === FULL_TITLE.length) {
        clearInterval(typer);
        setTypingDone(true);
      }
    }, delay);
    return () => clearInterval(typer);
  }, []);

  useEffect(() => {
    if (typingDone) return;
    const blink = setInterval(() => setShowCursor(p => !p), 500);
    return () => clearInterval(blink);
  }, [typingDone]);

  const handleCopy = () => {
    navigator.clipboard.writeText(command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  useEffect(() => {
    document.title = "Package Mate | Open-Source CLI for macOS Development";
  }, []);

  return (
    <>
      <link rel="icon" href="https://i.postimg.cc/MHgz7r8D/Screenshot-2026-03-31-at-21-03-35.png" />
      <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" />
      <style>{`
        .font-modern { font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; }
        @keyframes blink { 0%,100% { opacity: 1; } 50% { opacity: 0; } }
        .cursor { animation: blink 1s step-start infinite; }
      `}</style>

      <Component
        numBars={15}
        gradientFrom="rgb(58, 134, 255)"
        gradientTo="transparent"
        animationDuration={2}
        backgroundColor="rgb(10, 10, 10)"
      >
        <div className="text-center font-modern -mt-52">
          <h1 className="text-white text-5xl md:text-7xl font-bold mb-6 tracking-tight">
            {displayedTitle}
            {!typingDone && (
              <span className="cursor" style={{ opacity: showCursor ? 1 : 0 }}>|</span>
            )}
          </h1>

          <p className="text-gray-400 font-medium mb-8 mx-auto max-w-xl leading-relaxed" style={{ fontSize: '1.08rem' }}>
            Manage your macOS dev stack from the terminal. Package Mate handles Homebrew Casks, system binaries, and conflicts.
          </p>

          <div className="mx-auto w-full max-w-xl rounded-xl overflow-hidden border border-white/[0.08] text-left">
            <div className="flex items-center justify-between px-4 py-2 border-b border-white/[0.06]">
              <span className="text-white/30 text-sm font-medium">Bash</span>
              <button onClick={handleCopy} className="text-white/30 hover:text-white/70 transition-colors">
                {copied ? <Check className="w-4 h-4 text-green-400" /> : <Copy className="w-4 h-4" />}
              </button>
            </div>
            <div className="px-5 py-4">
              <span className="text-sm md:text-base font-mono">
                <span style={{ color: '#3A86FF' }}>curl </span>
                <span style={{ color: '#3A86FF' }}>-fsSL </span>
                <span className="text-white/70">https://package-mate.com/install.sh | bash</span>
              </span>
            </div>
          </div>
        </div>
      </Component>
    </>
  );
}
 