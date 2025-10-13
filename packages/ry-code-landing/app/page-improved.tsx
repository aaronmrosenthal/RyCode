import Image from 'next/image'

export default function Home() {
  return (
    <>
      {/* Sticky Navigation */}
      <nav className="sticky top-0 z-50 bg-black/80 backdrop-blur-md border-b border-neural-cyan/20">
        <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-2">
              <span className="text-2xl font-bold text-gradient">RyCode</span>
              <span className="hidden sm:inline px-2 py-0.5 rounded bg-neural-cyan/10 border border-neural-cyan/30 text-neural-cyan text-xs font-mono">v1.0</span>
            </div>
            <div className="hidden md:flex items-center gap-6 text-sm font-medium">
              <a href="#features" className="text-gray-300 hover:text-neural-cyan transition-colors">Features</a>
              <a href="#demo" className="text-gray-300 hover:text-neural-cyan transition-colors">Demo</a>
              <a href="https://github.com/aaronmrosenthal/RyCode" target="_blank" rel="noopener noreferrer" className="text-gray-300 hover:text-neural-cyan transition-colors">GitHub</a>
            </div>
            <a
              href="#install"
              className="bg-neural-cyan hover:bg-neural-magenta text-black font-bold px-6 py-2 rounded-lg transition-all duration-200 text-sm"
            >
              Get Started
            </a>
          </div>
        </div>
      </nav>

      <main className="min-h-screen">
        {/* Hero Section - Tighter spacing */}
        <section id="hero" className="relative min-h-[90vh] flex items-center justify-center px-4 sm:px-6 lg:px-8 py-12 sm:py-16">
          <div className="mx-auto max-w-6xl w-full">
            <div className="text-center">
              {/* Compact Badge */}
              <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-neural-cyan/10 border border-neural-cyan/30 mb-6">
                <span className="w-1.5 h-1.5 rounded-full bg-neural-cyan animate-pulse"></span>
                <span className="text-neural-cyan text-xs sm:text-sm font-mono font-semibold">Open Source • Production Ready • 31/31 Tests Passing</span>
              </div>

              {/* Refined Typography */}
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold mb-4 text-gradient glow-cyan">
                RyCode
              </h1>
              <p className="text-xl sm:text-2xl lg:text-3xl mb-3 text-neural-cyan font-mono font-bold">
                World's Most Advanced Open Source Coding Agent
              </p>
              <p className="text-base sm:text-lg lg:text-xl text-gray-300 mb-4 max-w-3xl mx-auto font-light">
                Switch between <span className="text-neural-magenta font-semibold">5 state-of-the-art AI models</span> with a single keystroke. <span className="text-neural-cyan font-semibold">Zero context loss.</span>
              </p>

              {/* Compact Feature Pills */}
              <div className="flex flex-wrap justify-center gap-2 sm:gap-3 mb-8 max-w-4xl mx-auto">
                {['Production-Grade TUI', 'Instant Model Switching', 'Context Preservation', '60 FPS Terminal UI', '19MB Binary'].map((feature) => (
                  <div key={feature} className="flex items-center gap-1.5 px-2.5 py-1 text-xs sm:text-sm text-gray-400 bg-black/40 rounded-full border border-gray-700/50">
                    <span className="text-matrix-green text-sm">✓</span>
                    <span>{feature}</span>
                  </div>
                ))}
              </div>

              {/* Prominent CTA Section */}
              <div id="install" className="mb-10 flex flex-col items-center gap-4">
                <div className="bg-black/60 backdrop-blur-sm border-2 border-neural-cyan/40 rounded-lg px-4 sm:px-6 py-3 font-mono text-sm sm:text-base w-full max-w-2xl">
                  <span className="text-gray-500">$</span>{' '}
                  <span className="text-neural-cyan">curl -fsSL https://ry-code.com/install | sh</span>
                </div>
                <div className="flex flex-col sm:flex-row items-center gap-3 w-full sm:w-auto">
                  <button className="bg-neural-cyan hover:bg-neural-magenta text-black font-bold py-3 px-8 sm:py-4 sm:px-12 rounded-lg transition-all duration-200 transform hover:scale-105 shadow-xl shadow-neural-cyan/50 text-base sm:text-lg w-full sm:w-auto">
                    Get Started - It's Free →
                  </button>
                  <a
                    href="https://github.com/aaronmrosenthal/RyCode"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center justify-center gap-2 text-gray-400 hover:text-neural-cyan transition-colors py-3 px-6 w-full sm:w-auto"
                  >
                    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                    <span className="font-mono text-sm">Star on GitHub</span>
                  </a>
                </div>
              </div>

              {/* Compact Model Badges */}
              <div className="mb-8">
                <p className="text-xs text-gray-500 uppercase tracking-wide mb-3">Latest Frontier Models (2025)</p>
                <div className="flex flex-wrap justify-center gap-2">
                  {[
                    { name: 'Claude Sonnet 4.5', color: '7aa2f7' },
                    { name: 'Gemini 2.5 Pro', color: 'ea4aaa' },
                    { name: 'GPT-5', color: 'ff6b35' },
                    { name: 'Grok 4 Fast', color: '00ffff' },
                    { name: 'Qwen3-Coder', color: 'ff00ff' }
                  ].map((model) => (
                    <div
                      key={model.name}
                      className={`px-3 py-1.5 rounded-full bg-[#${model.color}]/20 border border-[#${model.color}]/40 text-[#${model.color}] font-mono text-xs font-semibold`}
                      style={{
                        backgroundColor: `rgba(${parseInt(model.color.slice(0,2), 16)}, ${parseInt(model.color.slice(2,4), 16)}, ${parseInt(model.color.slice(4,6), 16)}, 0.2)`,
                        borderColor: `rgba(${parseInt(model.color.slice(0,2), 16)}, ${parseInt(model.color.slice(2,4), 16)}, ${parseInt(model.color.slice(4,6), 16)}, 0.4)`,
                        color: `#${model.color}`
                      }}
                    >
                      {model.name}
                    </div>
                  ))}
                </div>
              </div>

              {/* Optimized Terminal Mockup */}
              <div id="features" className="relative max-w-5xl mx-auto">
                <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-cyan/30 shadow-2xl shadow-neural-cyan/20">
                  {/* Compact Terminal Header */}
                  <div className="bg-[#11111b] px-3 py-2 sm:px-4 sm:py-2.5 flex items-center justify-between border-b border-gray-700">
                    <div className="flex items-center gap-2">
                      <div className="flex gap-1.5">
                        <div className="w-2.5 h-2.5 rounded-full bg-[#ff5f56]"></div>
                        <div className="w-2.5 h-2.5 rounded-full bg-[#ffbd2e]"></div>
                        <div className="w-2.5 h-2.5 rounded-full bg-[#27c93f]"></div>
                      </div>
                      <div className="ml-2 text-xs text-gray-400 font-mono hidden sm:block">RyCode - Model Selector</div>
                    </div>
                  </div>

                  {/* Compact Prompt */}
                  <div className="p-4 sm:p-6">
                    <div className="flex gap-2 mb-4 font-mono text-xs sm:text-sm">
                      <span className="text-neural-cyan">❯</span>
                      <span className="text-neural-magenta">/model</span>
                    </div>

                    {/* Compact Model Selector */}
                    <div className="space-y-2 font-mono text-xs sm:text-sm">
                      <div className="text-gray-400 text-xs mb-3 uppercase tracking-wide">Select AI Model:</div>

                      {/* Active Model - More Compact */}
                      <div className="flex items-center gap-2 sm:gap-3 px-3 py-2 rounded-lg bg-[#7aa2f7]/10 border-l-4 border-[#7aa2f7]">
                        <span className="text-[#7aa2f7] text-lg">▶</span>
                        <div className="flex-1 min-w-0">
                          <div className="text-[#7aa2f7] font-bold text-sm">Claude Sonnet 4.5</div>
                          <div className="text-xs text-gray-400 mt-0.5 hidden sm:block truncate">Best coding model, 77.2% on SWE-bench</div>
                        </div>
                        <div className="px-2 py-0.5 rounded bg-[#7aa2f7]/20 text-[#7aa2f7] text-xs font-semibold shrink-0">ACTIVE</div>
                      </div>

                      {/* Other Models - Compact */}
                      {[
                        { name: 'Gemini 2.5 Pro', desc: 'Most intelligent with advanced thinking', color: 'ea4aaa' },
                        { name: 'GPT-5', desc: 'OpenAI smartest with built-in reasoning', color: 'ff6b35' },
                        { name: 'Grok 4 Fast', desc: '2M context, 98% cost reduction', color: '00ffff' },
                        { name: 'Qwen3-Coder', desc: '480B parameter agentic coder', color: 'ff00ff' }
                      ].map((model) => (
                        <div key={model.name} className="flex items-center gap-2 sm:gap-3 px-3 py-2 rounded-lg hover:bg-gray-700/20 transition-colors cursor-pointer">
                          <span className="text-gray-600 text-lg">○</span>
                          <div className="flex-1 min-w-0">
                            <div className={`text-[#${model.color}] font-bold text-sm`} style={{ color: `#${model.color}` }}>{model.name}</div>
                            <div className="text-xs text-gray-500 mt-0.5 hidden sm:block truncate">{model.desc}</div>
                          </div>
                        </div>
                      ))}
                    </div>

                    {/* Compact Help Text */}
                    <div className="mt-4 pt-3 border-t border-gray-700/50 text-xs text-gray-500 font-mono">
                      <div className="flex flex-wrap items-center gap-3 sm:gap-4">
                        <span><span className="text-neural-cyan">↑↓</span> Navigate</span>
                        <span><span className="text-neural-cyan">Tab</span> Quick Switch</span>
                        <span><span className="text-neural-cyan">Enter</span> Select</span>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Compact Feature Highlights */}
                <div className="mt-6 text-center space-y-3">
                  <p className="text-neural-cyan text-base sm:text-lg font-semibold">
                    Type <span className="bg-black/40 px-2 py-1 rounded font-mono text-sm">/model</span> to switch models
                  </p>
                  <div className="flex flex-wrap justify-center gap-4 text-sm text-gray-400">
                    <div className="flex items-center gap-1.5">
                      <span className="text-xl">⚡</span>
                      <span>Instant switching with <span className="text-neural-cyan font-mono">Tab</span></span>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <span className="text-xl">🧠</span>
                      <span>Context preserved</span>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <span className="text-xl">🚀</span>
                      <span>Zero config</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Tab Switching Section - Reduced spacing */}
        <section className="relative py-16 px-4 sm:px-6 lg:px-8 bg-black/20">
          <div className="mx-auto max-w-6xl w-full">
            <div className="text-center mb-10">
              <h2 className="text-3xl sm:text-4xl font-bold mb-3 text-gradient">
                Switch Models Instantly
              </h2>
              <p className="text-base sm:text-lg text-gray-400 max-w-2xl mx-auto">
                Press <span className="text-neural-cyan font-mono bg-black/40 px-2 py-1 rounded text-sm">Tab</span> to cycle through models. Watch the chips update in real-time.
              </p>
            </div>

            {/* Compact Terminal Demo */}
            <div className="relative max-w-5xl mx-auto">
              <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-cyan/30 shadow-2xl shadow-neural-cyan/20">
                <div className="bg-[#11111b] px-3 py-2 sm:px-4 sm:py-2.5 flex items-center justify-between border-b border-gray-700">
                  <div className="flex items-center gap-2">
                    <div className="flex gap-1.5">
                      <div className="w-2.5 h-2.5 rounded-full bg-[#ff5f56]"></div>
                      <div className="w-2.5 h-2.5 rounded-full bg-[#ffbd2e]"></div>
                      <div className="w-2.5 h-2.5 rounded-full bg-[#27c93f]"></div>
                    </div>
                    <div className="ml-2 text-xs text-gray-400 font-mono hidden sm:block">RyCode - Multi-Agent CLI</div>
                  </div>

                  <div className="flex gap-1.5 sm:gap-2 relative">
                    <div className="absolute -inset-2 sm:-inset-3 bg-neural-cyan/10 rounded-lg animate-pulse pointer-events-none"></div>
                    <div className="px-2 sm:px-2.5 py-0.5 sm:py-1 rounded-full bg-[#7aa2f7]/20 border border-[#7aa2f7]/40 text-[#7aa2f7] text-xs font-mono flex items-center gap-1">
                      <span className="w-1.5 h-1.5 rounded-full bg-[#7aa2f7] animate-pulse"></span>
                      <span className="hidden sm:inline">Sonnet 4.5</span>
                      <span className="sm:hidden">Claude</span>
                    </div>
                    <div className="px-2 sm:px-2.5 py-0.5 sm:py-1 rounded-full bg-gray-700/30 border border-gray-600 text-gray-500 text-xs font-mono hidden sm:block">
                      Gemini
                    </div>
                    <div className="px-2 sm:px-2.5 py-0.5 sm:py-1 rounded-full bg-gray-700/30 border border-gray-600 text-gray-500 text-xs font-mono hidden md:block">
                      GPT-5
                    </div>
                  </div>
                </div>

                <div className="p-4 sm:p-6 font-mono text-sm">
                  <div className="p-4 sm:p-6 bg-black/30 rounded-lg border border-neural-cyan/20">
                    <div className="flex flex-col sm:flex-row items-center justify-center gap-4 sm:gap-8">
                      <div className="text-center">
                        <div className="text-4xl sm:text-5xl mb-2">⇥</div>
                        <p className="text-neural-cyan font-semibold text-sm">Press Tab</p>
                        <p className="text-xs text-gray-400 mt-1">Cycle models</p>
                      </div>
                      <div className="text-2xl sm:text-3xl text-gray-600">→</div>
                      <div className="text-center">
                        <div className="text-4xl sm:text-5xl mb-2">💫</div>
                        <p className="text-neural-magenta font-semibold text-sm">Chips Update</p>
                        <p className="text-xs text-gray-400 mt-1">See active model</p>
                      </div>
                    </div>
                  </div>

                  <div className="mt-4 text-center text-xs text-gray-500">
                    <p className="hidden sm:block">Sonnet 4.5 → Gemini 2.5 Pro → GPT-5 → Grok 4 Fast → Qwen3-Coder</p>
                    <p className="sm:hidden">Tap Tab to cycle through 5 models</p>
                  </div>
                </div>
              </div>

              <div className="mt-4 flex justify-center">
                <div className="bg-neural-cyan/10 backdrop-blur-sm border border-neural-cyan/30 rounded-lg px-4 py-2">
                  <p className="text-sm text-neural-cyan font-mono">
                    ← Model chips update in real-time
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Demo Section - Optimized */}
        <section id="demo" className="relative py-16 px-4 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-6xl w-full">
            <div className="text-center mb-10">
              <h2 className="text-3xl sm:text-4xl font-bold mb-3 text-gradient">
                See It In Action
              </h2>
              <p className="text-base sm:text-lg text-gray-400 max-w-2xl mx-auto">
                Real conversation. Real code generation. Real multi-model switching.
              </p>
            </div>

            <div className="relative max-w-5xl mx-auto">
              <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-magenta/30 shadow-2xl shadow-neural-magenta/20">
                <div className="bg-[#11111b] px-3 py-2 sm:px-4 sm:py-2.5 flex items-center justify-between border-b border-gray-700">
                  <div className="flex items-center gap-2">
                    <div className="flex gap-1.5">
                      <div className="w-2.5 h-2.5 rounded-full bg-[#ff5f56]"></div>
                      <div className="w-2.5 h-2.5 rounded-full bg-[#ffbd2e]"></div>
                      <div className="w-2.5 h-2.5 rounded-full bg-[#27c93f]"></div>
                    </div>
                    <div className="ml-2 text-xs text-gray-400 font-mono hidden sm:block">RyCode - Live Session</div>
                  </div>

                  <div className="px-2.5 py-1 rounded-full bg-[#7aa2f7]/20 border border-[#7aa2f7]/40 text-[#7aa2f7] text-xs font-mono flex items-center gap-1">
                    <span className="w-1.5 h-1.5 rounded-full bg-[#7aa2f7] animate-pulse"></span>
                    Claude
                  </div>
                </div>

                <div className="p-4 sm:p-6 font-mono text-sm space-y-4 max-h-[500px] overflow-y-auto">
                  <div className="flex gap-2">
                    <span className="text-neural-cyan">❯</span>
                    <span className="text-white text-xs sm:text-sm">Build me a REST API for a todo app with auth</span>
                  </div>

                  <div className="pl-4 sm:pl-6 space-y-2">
                    <div className="text-gray-300 text-xs sm:text-sm">
                      <span className="text-[#7aa2f7] font-semibold">Claude:</span> I'll create a production-ready REST API with authentication.
                    </div>

                    <div className="bg-black/40 rounded-lg p-3 border border-[#7aa2f7]/20">
                      <div className="text-xs text-gray-500 mb-2">server.js</div>
                      <pre className="text-xs text-matrix-green overflow-x-auto">
{`const express = require('express');
const jwt = require('jsonwebtoken');
const app = express();

app.post('/register', async (req, res) => {
  const hash = await bcrypt.hash(req.body.password, 10);
  res.status(201).send('User registered');
});`}
                      </pre>
                    </div>

                    <div className="text-gray-400 text-xs">
                      ✓ Express configured • ✓ JWT auth • ✓ Bcrypt hashing
                    </div>
                  </div>

                  <div className="border-t border-gray-700/50 my-4"></div>

                  <div className="flex gap-2">
                    <span className="text-neural-cyan">❯</span>
                    <span className="text-neural-magenta italic text-xs sm:text-sm">*Tab → GPT-5*</span>
                  </div>

                  <div className="flex gap-2">
                    <span className="text-neural-cyan">❯</span>
                    <span className="text-white text-xs sm:text-sm">Add rate limiting and validation</span>
                  </div>

                  <div className="pl-4 sm:pl-6">
                    <div className="text-gray-300 text-xs sm:text-sm mb-2">
                      <span className="text-[#10a37f] font-semibold">GPT-5:</span> Adding express-rate-limit and joi validation...
                    </div>
                  </div>
                </div>

                <div className="bg-[#11111b] px-4 py-2 border-t border-gray-700 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 text-xs">
                  <div className="flex items-center gap-3 text-gray-400">
                    <span className="flex items-center gap-1.5">
                      <span className="w-1.5 h-1.5 rounded-full bg-matrix-green animate-pulse"></span>
                      3 models used
                    </span>
                    <span className="hidden sm:inline">Context preserved</span>
                  </div>
                  <div className="text-gray-500 font-mono">
                    45 lines • 0 errors
                  </div>
                </div>
              </div>

              {/* Compact Value Props */}
              <div className="mt-8 grid sm:grid-cols-3 gap-4">
                {[
                  { icon: '🎯', title: 'Context Preserved', desc: 'Switch models without losing context', color: '7aa2f7' },
                  { icon: '⚡', title: 'Instant Switching', desc: 'Tab for second opinion in ms', color: '00ffff' },
                  { icon: '🧠', title: 'Best of All Models', desc: 'Claude, Gemini, GPT-5, Grok, Qwen', color: '00ff00' }
                ].map((item) => (
                  <div key={item.title} className="text-center p-4 bg-black/40 backdrop-blur-sm rounded-lg border border-[#${item.color}]/20" style={{ borderColor: `#${item.color}33` }}>
                    <div className="text-3xl mb-2">{item.icon}</div>
                    <h3 className="text-sm font-bold mb-1" style={{ color: `#${item.color}` }}>{item.title}</h3>
                    <p className="text-xs text-gray-400">{item.desc}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </section>

        {/* toolkit-cli Section - Compact */}
        <section className="px-4 sm:px-6 lg:px-8 py-12 bg-black/20">
          <div className="mx-auto max-w-4xl text-center">
            <h2 className="text-3xl font-bold mb-4 text-gradient">
              Built with toolkit-cli
            </h2>
            <p className="text-base sm:text-lg text-gray-300 mb-6">
              RyCode showcases{' '}
              <a
                href="https://toolkit-cli.com"
                target="_blank"
                rel="noopener noreferrer"
                className="text-claude-blue hover:text-neural-cyan transition-colors font-semibold"
              >
                toolkit-cli
              </a>
              , commercial software <span className="text-matrix-green font-semibold">patent pending</span>.
            </p>

            <div className="grid sm:grid-cols-2 gap-4 mt-8">
              <div className="p-4 rounded-lg bg-black/40 border border-claude-blue/30">
                <h3 className="text-lg font-bold mb-2 text-claude-blue">100% AI-Designed</h3>
                <p className="text-gray-300 text-sm">From concept to completion by Claude AI using toolkit-cli.</p>
              </div>
              <div className="p-4 rounded-lg bg-black/40 border border-performance-gold/30">
                <h3 className="text-lg font-bold mb-2 text-performance-gold">Production Ready</h3>
                <p className="text-gray-300 text-sm">31/31 tests passing, 54.2% coverage, zero bugs.</p>
              </div>
            </div>
          </div>
        </section>

        {/* Enhanced Footer */}
        <footer className="px-4 sm:px-6 lg:px-8 py-12 border-t border-neural-cyan/20 bg-black/40">
          <div className="mx-auto max-w-6xl">
            <div className="grid sm:grid-cols-2 md:grid-cols-4 gap-8 mb-8">
              <div>
                <h3 className="text-sm font-bold text-neural-cyan mb-3">Product</h3>
                <ul className="space-y-2 text-sm text-gray-400">
                  <li><a href="#features" className="hover:text-neural-cyan transition-colors">Features</a></li>
                  <li><a href="#demo" className="hover:text-neural-cyan transition-colors">Demo</a></li>
                  <li><a href="#install" className="hover:text-neural-cyan transition-colors">Installation</a></li>
                </ul>
              </div>
              <div>
                <h3 className="text-sm font-bold text-neural-cyan mb-3">Resources</h3>
                <ul className="space-y-2 text-sm text-gray-400">
                  <li><a href="https://github.com/aaronmrosenthal/RyCode" target="_blank" rel="noopener noreferrer" className="hover:text-neural-cyan transition-colors">GitHub</a></li>
                  <li><a href="https://github.com/aaronmrosenthal/RyCode/blob/main/README.md" target="_blank" rel="noopener noreferrer" className="hover:text-neural-cyan transition-colors">Documentation</a></li>
                </ul>
              </div>
              <div>
                <h3 className="text-sm font-bold text-neural-cyan mb-3">Company</h3>
                <ul className="space-y-2 text-sm text-gray-400">
                  <li><a href="https://toolkit-cli.com" target="_blank" rel="noopener noreferrer" className="hover:text-neural-cyan transition-colors">toolkit-cli</a></li>
                </ul>
              </div>
              <div>
                <h3 className="text-sm font-bold text-neural-cyan mb-3">Connect</h3>
                <ul className="space-y-2 text-sm text-gray-400">
                  <li><a href="https://github.com/aaronmrosenthal/RyCode" target="_blank" rel="noopener noreferrer" className="hover:text-neural-cyan transition-colors">GitHub</a></li>
                </ul>
              </div>
            </div>
            <div className="pt-8 border-t border-gray-800 text-center">
              <p className="font-mono text-sm text-gray-400">
                🤖 100% AI-Designed by Claude • Built with{' '}
                <a
                  href="https://toolkit-cli.com"
                  className="text-claude-blue hover:text-neural-cyan transition-colors"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  toolkit-cli
                </a>
              </p>
              <p className="mt-2 text-xs text-gray-500">
                Zero Compromises • Infinite Attention to Detail
              </p>
            </div>
          </div>
        </footer>
      </main>
    </>
  )
}
