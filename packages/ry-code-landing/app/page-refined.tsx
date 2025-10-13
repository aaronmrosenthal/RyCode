import Image from 'next/image'

export default function Home() {
  return (
    <>
      {/* Sticky Navigation */}
      <nav className="sticky top-0 z-50 bg-black/90 backdrop-blur-md border-b border-neural-cyan/20">
        <div className="mx-auto max-w-6xl px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-3">
              <span className="text-2xl font-bold text-gradient">RyCode</span>
              <span className="px-2 py-0.5 rounded bg-neural-cyan/10 border border-neural-cyan/30 text-neural-cyan text-xs font-mono">v1.0</span>
            </div>
            <div className="hidden md:flex items-center gap-8 text-sm font-medium">
              <a href="#features" className="text-gray-300 hover:text-neural-cyan transition-colors">Features</a>
              <a href="#demo" className="text-gray-300 hover:text-neural-cyan transition-colors">Demo</a>
              <a href="https://github.com/aaronmrosenthal/RyCode" target="_blank" rel="noopener noreferrer" className="flex items-center gap-2 text-gray-300 hover:text-neural-cyan transition-colors">
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                </svg>
                <span>Star</span>
              </a>
            </div>
          </div>
        </div>
      </nav>

      <main className="min-h-screen">
        {/* FOLD 1: Hero - Refined & Balanced */}
        <section id="hero" className="relative flex items-center justify-center px-6 lg:px-8 py-24 min-h-[95vh]">
          <div className="mx-auto max-w-5xl w-full">
            <div className="text-center">
              {/* Status Badge - Simplified */}
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-neural-cyan/10 border border-neural-cyan/30 mb-12">
                <span className="w-2 h-2 rounded-full bg-neural-cyan animate-pulse"></span>
                <span className="text-neural-cyan text-sm font-mono font-semibold">Open Source</span>
              </div>

              {/* Hero Lockup - Clear Hierarchy */}
              <h1 className="text-5xl lg:text-6xl font-bold mb-6 text-gradient glow-cyan">
                RyCode
              </h1>
              <p className="text-2xl lg:text-3xl mb-16 text-white font-light max-w-3xl mx-auto leading-relaxed">
                Switch Between 5 AI Models<br className="hidden sm:block" />
                With One Keystroke
              </p>

              {/* Primary CTA - Installation Command */}
              <div id="install" className="mb-16 flex flex-col items-center gap-4">
                <div className="bg-black/60 backdrop-blur-sm border-2 border-neural-cyan/40 rounded-xl px-8 py-5 font-mono text-lg w-full max-w-2xl shadow-2xl shadow-neural-cyan/20">
                  <span className="text-gray-500 select-none">$ </span>
                  <span className="text-neural-cyan select-all">curl -fsSL https://ry-code.com/install | sh</span>
                </div>
                <button className="bg-neural-cyan hover:bg-neural-magenta text-black font-bold py-4 px-12 rounded-lg transition-all duration-200 transform hover:scale-105 shadow-xl shadow-neural-cyan/50 text-lg">
                  Copy & Install →
                </button>
              </div>

              {/* Terminal Mockup - Hero Visual */}
              <div className="relative max-w-4xl mx-auto mb-12">
                <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-cyan/30 shadow-2xl shadow-neural-cyan/20">
                  {/* Terminal Header with Model Chips */}
                  <div className="bg-[#11111b] px-6 py-4 flex items-center justify-between border-b border-gray-700">
                    <div className="flex items-center gap-3">
                      <div className="flex gap-2">
                        <div className="w-3 h-3 rounded-full bg-[#ff5f56]"></div>
                        <div className="w-3 h-3 rounded-full bg-[#ffbd2e]"></div>
                        <div className="w-3 h-3 rounded-full bg-[#27c93f]"></div>
                      </div>
                      <div className="ml-4 text-sm text-gray-400 font-mono">rycode</div>
                    </div>

                    {/* Model Chips - Integrated */}
                    <div className="hidden md:flex gap-2">
                      <div className="px-3 py-1.5 rounded-full bg-[#7aa2f7]/20 border border-[#7aa2f7]/40 text-[#7aa2f7] text-xs font-mono font-semibold">
                        Claude
                      </div>
                      <div className="px-3 py-1.5 rounded-full bg-[#ea4aaa]/20 border border-[#ea4aaa]/40 text-[#ea4aaa] text-xs font-mono font-semibold">
                        Gemini
                      </div>
                      <div className="px-3 py-1.5 rounded-full bg-[#ff6b35]/20 border border-[#ff6b35]/40 text-[#ff6b35] text-xs font-mono font-semibold">
                        GPT-5
                      </div>
                      <div className="px-3 py-1.5 rounded-full bg-[#00ffff]/20 border border-[#00ffff]/40 text-[#00ffff] text-xs font-mono font-semibold">
                        Grok
                      </div>
                      <div className="px-3 py-1.5 rounded-full bg-[#ff00ff]/20 border border-[#ff00ff]/40 text-[#ff00ff] text-xs font-mono font-semibold">
                        Qwen
                      </div>
                    </div>
                  </div>

                  {/* Terminal Content */}
                  <div className="p-8 font-mono">
                    <div className="flex gap-3 mb-6">
                      <span className="text-neural-cyan text-lg">❯</span>
                      <span className="text-neural-magenta text-lg">/model</span>
                    </div>

                    {/* Model Selector - Simplified */}
                    <div className="space-y-3">
                      {/* Active Model */}
                      <div className="flex items-center gap-4 px-5 py-3 rounded-lg bg-[#7aa2f7]/10 border-l-4 border-[#7aa2f7]">
                        <span className="text-[#7aa2f7] text-xl">▶</span>
                        <div className="flex-1">
                          <div className="text-[#7aa2f7] font-bold text-base">Claude Sonnet 4.5</div>
                          <div className="text-xs text-gray-400 mt-1">Best for coding tasks</div>
                        </div>
                        <div className="px-2 py-1 rounded bg-[#7aa2f7]/20 text-[#7aa2f7] text-xs font-semibold">ACTIVE</div>
                      </div>

                      {/* Other Models */}
                      <div className="flex items-center gap-4 px-5 py-3 rounded-lg hover:bg-gray-700/20 transition-colors">
                        <span className="text-gray-600 text-xl">○</span>
                        <div className="flex-1">
                          <div className="text-[#ea4aaa] font-bold text-base">Gemini 2.5 Pro</div>
                          <div className="text-xs text-gray-500 mt-1">Advanced reasoning</div>
                        </div>
                      </div>

                      <div className="flex items-center gap-4 px-5 py-3 rounded-lg hover:bg-gray-700/20 transition-colors">
                        <span className="text-gray-600 text-xl">○</span>
                        <div className="flex-1">
                          <div className="text-[#ff6b35] font-bold text-base">GPT-5</div>
                          <div className="text-xs text-gray-500 mt-1">OpenAI's latest</div>
                        </div>
                      </div>

                      <div className="flex items-center gap-4 px-5 py-3 rounded-lg hover:bg-gray-700/20 transition-colors">
                        <span className="text-gray-600 text-xl">○</span>
                        <div className="flex-1">
                          <div className="text-[#00ffff] font-bold text-base">Grok 4 Fast</div>
                          <div className="text-xs text-gray-500 mt-1">2M context window</div>
                        </div>
                      </div>

                      <div className="flex items-center gap-4 px-5 py-3 rounded-lg hover:bg-gray-700/20 transition-colors">
                        <span className="text-gray-600 text-xl">○</span>
                        <div className="flex-1">
                          <div className="text-[#ff00ff] font-bold text-base">Qwen3-Coder</div>
                          <div className="text-xs text-gray-500 mt-1">480B parameters</div>
                        </div>
                      </div>
                    </div>

                    {/* Help Text */}
                    <div className="mt-6 pt-4 border-t border-gray-700/50 text-sm text-gray-500 font-mono flex items-center gap-6">
                      <span><span className="text-neural-cyan">Tab</span> Switch</span>
                      <span><span className="text-neural-cyan">↑↓</span> Navigate</span>
                      <span><span className="text-neural-cyan">Enter</span> Select</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Feature Highlights - Clean & Minimal */}
              <div className="flex flex-wrap justify-center gap-12 text-gray-400">
                <div className="flex items-center gap-3">
                  <span className="text-3xl">⚡</span>
                  <span className="text-base">Instant Switching</span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-3xl">🧠</span>
                  <span className="text-base">Context Preserved</span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-3xl">🚀</span>
                  <span className="text-base">Zero Configuration</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* FOLD 2: Tab Switching Demo - Refined */}
        <section id="features" className="relative py-24 px-6 lg:px-8 bg-gradient-to-b from-black/20 to-black/40">
          <div className="mx-auto max-w-5xl w-full">
            <div className="text-center mb-16">
              <h2 className="text-4xl lg:text-5xl font-bold mb-6 text-gradient">
                Switch Models Instantly
              </h2>
              <p className="text-xl text-gray-300 max-w-2xl mx-auto">
                Press <span className="text-neural-cyan font-mono bg-black/40 px-3 py-1.5 rounded">Tab</span> to cycle through models mid-conversation
              </p>
            </div>

            {/* Tab Demo Visual */}
            <div className="relative max-w-4xl mx-auto">
              <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-cyan/30 shadow-2xl shadow-neural-cyan/20">
                <div className="bg-[#11111b] px-6 py-4 flex items-center justify-between border-b border-gray-700">
                  <div className="flex items-center gap-3">
                    <div className="flex gap-2">
                      <div className="w-3 h-3 rounded-full bg-[#ff5f56]"></div>
                      <div className="w-3 h-3 rounded-full bg-[#ffbd2e]"></div>
                      <div className="w-3 h-3 rounded-full bg-[#27c93f]"></div>
                    </div>
                    <div className="ml-4 text-sm text-gray-400 font-mono">rycode</div>
                  </div>

                  {/* Animated Model Chips */}
                  <div className="flex gap-2 relative">
                    <div className="absolute -inset-4 bg-neural-cyan/10 rounded-lg animate-pulse pointer-events-none"></div>
                    <div className="px-3 py-1.5 rounded-full bg-[#7aa2f7]/30 border-2 border-[#7aa2f7] text-[#7aa2f7] text-sm font-mono flex items-center gap-2">
                      <span className="w-2 h-2 rounded-full bg-[#7aa2f7] animate-pulse"></span>
                      Claude
                    </div>
                    <div className="px-3 py-1.5 rounded-full bg-gray-700/30 border border-gray-600 text-gray-500 text-sm font-mono">
                      Gemini
                    </div>
                    <div className="px-3 py-1.5 rounded-full bg-gray-700/30 border border-gray-600 text-gray-500 text-sm font-mono hidden md:block">
                      GPT-5
                    </div>
                  </div>
                </div>

                <div className="p-12 font-mono">
                  <div className="p-8 bg-black/30 rounded-lg border border-neural-cyan/20">
                    <div className="flex items-center justify-center gap-12">
                      <div className="text-center">
                        <div className="text-6xl mb-4">⇥</div>
                        <p className="text-neural-cyan font-semibold text-lg">Press Tab</p>
                        <p className="text-sm text-gray-400 mt-2">Cycle models</p>
                      </div>
                      <div className="text-4xl text-gray-600">→</div>
                      <div className="text-center">
                        <div className="text-6xl mb-4">💫</div>
                        <p className="text-neural-magenta font-semibold text-lg">Instant Switch</p>
                        <p className="text-sm text-gray-400 mt-2">No context loss</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* FOLD 3: Live Demo - Simplified */}
        <section id="demo" className="relative py-24 px-6 lg:px-8">
          <div className="mx-auto max-w-5xl w-full">
            <div className="text-center mb-16">
              <h2 className="text-4xl lg:text-5xl font-bold mb-6 text-gradient">
                See It In Action
              </h2>
              <p className="text-xl text-gray-300 max-w-2xl mx-auto">
                Real conversation with multiple AI models
              </p>
            </div>

            <div className="relative max-w-4xl mx-auto">
              <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-magenta/30 shadow-2xl shadow-neural-magenta/20">
                <div className="bg-[#11111b] px-6 py-4 flex items-center justify-between border-b border-gray-700">
                  <div className="flex items-center gap-3">
                    <div className="flex gap-2">
                      <div className="w-3 h-3 rounded-full bg-[#ff5f56]"></div>
                      <div className="w-3 h-3 rounded-full bg-[#ffbd2e]"></div>
                      <div className="w-3 h-3 rounded-full bg-[#27c93f]"></div>
                    </div>
                    <div className="ml-4 text-sm text-gray-400 font-mono">rycode - live session</div>
                  </div>

                  <div className="px-3 py-1.5 rounded-full bg-[#7aa2f7]/20 border border-[#7aa2f7]/40 text-[#7aa2f7] text-sm font-mono flex items-center gap-2">
                    <span className="w-2 h-2 rounded-full bg-[#7aa2f7] animate-pulse"></span>
                    Claude
                  </div>
                </div>

                <div className="p-8 font-mono space-y-6 max-h-[500px] overflow-y-auto">
                  {/* User Prompt */}
                  <div className="flex gap-3">
                    <span className="text-neural-cyan text-lg">❯</span>
                    <span className="text-white">Build a REST API with authentication</span>
                  </div>

                  {/* Claude Response */}
                  <div className="pl-6 space-y-3">
                    <div className="text-gray-300">
                      <span className="text-[#7aa2f7] font-semibold">Claude:</span> I'll create a production-ready REST API...
                    </div>

                    <div className="bg-black/40 rounded-lg p-4 border border-[#7aa2f7]/20">
                      <pre className="text-sm text-matrix-green">
{`const express = require('express');
const jwt = require('jsonwebtoken');

app.post('/register', async (req, res) => {
  const hash = await bcrypt.hash(req.body.password, 10);
  res.status(201).send('User registered');
});`}
                      </pre>
                    </div>
                  </div>

                  <div className="border-t border-gray-700/50"></div>

                  {/* Model Switch */}
                  <div className="flex gap-3">
                    <span className="text-neural-cyan text-lg">❯</span>
                    <span className="text-neural-magenta italic">*Tab → GPT-5*</span>
                  </div>

                  <div className="flex gap-3">
                    <span className="text-neural-cyan text-lg">❯</span>
                    <span className="text-white">Add rate limiting</span>
                  </div>

                  <div className="pl-6">
                    <div className="text-gray-300">
                      <span className="text-[#10a37f] font-semibold">GPT-5:</span> Adding express-rate-limit...
                    </div>
                  </div>
                </div>

                <div className="bg-[#11111b] px-6 py-3 border-t border-gray-700 flex items-center justify-between text-sm">
                  <div className="flex items-center gap-4 text-gray-400">
                    <span className="flex items-center gap-2">
                      <span className="w-2 h-2 rounded-full bg-matrix-green animate-pulse"></span>
                      3 models • Context preserved
                    </span>
                  </div>
                  <div className="text-gray-500 font-mono">
                    45 lines • 0 errors
                  </div>
                </div>
              </div>

              {/* Value Props - Clean Grid */}
              <div className="mt-12 grid md:grid-cols-3 gap-6">
                <div className="text-center p-6 bg-black/40 backdrop-blur-sm rounded-xl border border-[#7aa2f7]/20">
                  <div className="text-4xl mb-3">🎯</div>
                  <h3 className="text-lg font-bold text-[#7aa2f7] mb-2">Context Preserved</h3>
                  <p className="text-sm text-gray-400">
                    Switch models mid-conversation without losing any context
                  </p>
                </div>

                <div className="text-center p-6 bg-black/40 backdrop-blur-sm rounded-xl border border-neural-cyan/20">
                  <div className="text-4xl mb-3">⚡</div>
                  <h3 className="text-lg font-bold text-neural-cyan mb-2">Instant Switching</h3>
                  <p className="text-sm text-gray-400">
                    Get a second opinion in milliseconds
                  </p>
                </div>

                <div className="text-center p-6 bg-black/40 backdrop-blur-sm rounded-xl border border-matrix-green/20">
                  <div className="text-4xl mb-3">🧠</div>
                  <h3 className="text-lg font-bold text-matrix-green mb-2">Best of All Models</h3>
                  <p className="text-sm text-gray-400">
                    Use Claude, Gemini, GPT-5, Grok, and Qwen together
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Built with toolkit-cli - Minimal */}
        <section className="px-6 lg:px-8 py-20 bg-black/20">
          <div className="mx-auto max-w-4xl text-center">
            <h2 className="text-3xl font-bold mb-4 text-gradient">
              Built with toolkit-cli
            </h2>
            <p className="text-lg text-gray-300 mb-8">
              Showcasing{' '}
              <a
                href="https://toolkit-cli.com"
                target="_blank"
                rel="noopener noreferrer"
                className="text-claude-blue hover:text-neural-cyan transition-colors font-semibold"
              >
                toolkit-cli
              </a>
              , commercial software <span className="text-matrix-green font-semibold">patent pending</span>
            </p>

            <div className="grid md:grid-cols-2 gap-6 mt-10">
              <div className="p-6 rounded-xl bg-black/40 border border-claude-blue/30">
                <h3 className="text-xl font-bold mb-3 text-claude-blue">100% AI-Designed</h3>
                <p className="text-gray-300">
                  From concept to completion by Claude AI
                </p>
              </div>
              <div className="p-6 rounded-xl bg-black/40 border border-performance-gold/30">
                <h3 className="text-xl font-bold mb-3 text-performance-gold">Production Ready</h3>
                <p className="text-gray-300">
                  31/31 tests passing, zero known bugs
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* Footer - Clean */}
        <footer className="px-6 lg:px-8 py-12 border-t border-neural-cyan/20">
          <div className="mx-auto max-w-7xl">
            <div className="text-center text-gray-400">
              <p className="font-mono text-sm">
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
