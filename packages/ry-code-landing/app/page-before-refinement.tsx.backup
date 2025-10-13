import Image from 'next/image'

export default function Home() {
  return (
    <main className="min-h-screen">
      {/* FOLD 1: Hero - Splash Screen Reveals Model Selector */}
      <section className="relative min-h-screen flex items-center justify-center px-6 py-24">
        <div className="mx-auto max-w-7xl w-full">
          <div className="text-center">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-neural-cyan/10 border border-neural-cyan/30 mb-8">
              <span className="w-2 h-2 rounded-full bg-neural-cyan animate-pulse"></span>
              <span className="text-neural-cyan text-sm font-mono font-semibold">Open Source • Production Ready • 31/31 Tests Passing</span>
            </div>

            <h1 className="text-6xl lg:text-8xl font-bold mb-6 text-gradient glow-cyan">
              RyCode
            </h1>
            <p className="text-2xl lg:text-4xl mb-4 text-neural-cyan font-mono font-bold">
              The World's Most Advanced Open Source Coding Agent
            </p>
            <p className="text-xl lg:text-2xl text-gray-300 mb-6 max-w-3xl mx-auto font-light">
              Switch between <span className="text-neural-magenta font-semibold">multiple state-of-the-art AI models</span> with a single keystroke. <span className="text-neural-cyan font-semibold">Zero context loss.</span> Infinite possibilities.
            </p>

            {/* Key Differentiators */}
            <div className="flex flex-wrap justify-center gap-4 mb-12 max-w-4xl mx-auto">
              <div className="flex items-center gap-2 text-sm text-gray-400">
                <span className="text-matrix-green text-lg">✓</span>
                <span>Production-Grade TUI</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-400">
                <span className="text-matrix-green text-lg">✓</span>
                <span>Instant Model Switching</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-400">
                <span className="text-matrix-green text-lg">✓</span>
                <span>Context Preservation</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-400">
                <span className="text-matrix-green text-lg">✓</span>
                <span>60 FPS Terminal UI</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-400">
                <span className="text-matrix-green text-lg">✓</span>
                <span>19MB Binary</span>
              </div>
            </div>

            {/* Installation */}
            <div className="mb-16 flex flex-col items-center gap-6">
              <div className="bg-black/40 backdrop-blur-sm border border-neural-cyan/30 rounded-lg px-8 py-4 font-mono text-base sm:text-lg">
                <span className="text-gray-400">$</span>{' '}
                <span className="text-neural-cyan">curl -fsSL https://ry-code.com/install | sh</span>
              </div>
              <div className="flex flex-col sm:flex-row items-center gap-4">
                <button className="bg-neural-cyan hover:bg-neural-magenta text-black font-bold py-4 px-10 rounded-lg transition-all duration-200 transform hover:scale-105 shadow-lg shadow-neural-cyan/50">
                  Get Started - It's Free →
                </button>
                <a
                  href="https://github.com/aaronmrosenthal/RyCode"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 text-gray-400 hover:text-neural-cyan transition-colors"
                >
                  <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                  </svg>
                  <span className="font-mono text-sm">Star on GitHub</span>
                </a>
              </div>
            </div>

            {/* SOTA Models Showcase */}
            <div className="mb-12">
              <p className="text-xs sm:text-sm text-gray-500 uppercase tracking-wide mb-4">Latest frontier models (2025)</p>
              <div className="flex flex-wrap justify-center gap-2 sm:gap-3">
                <div className="px-3 sm:px-4 py-1.5 sm:py-2 rounded-full bg-[#7aa2f7]/20 border-2 border-[#7aa2f7]/40 text-[#7aa2f7] font-mono text-xs sm:text-sm font-semibold">
                  Claude Sonnet 4.5
                </div>
                <div className="px-3 sm:px-4 py-1.5 sm:py-2 rounded-full bg-[#ea4aaa]/20 border-2 border-[#ea4aaa]/40 text-[#ea4aaa] font-mono text-xs sm:text-sm font-semibold">
                  Gemini 2.5 Pro
                </div>
                <div className="px-3 sm:px-4 py-1.5 sm:py-2 rounded-full bg-[#ff6b35]/20 border-2 border-[#ff6b35]/40 text-[#ff6b35] font-mono text-xs sm:text-sm font-semibold">
                  GPT-5
                </div>
                <div className="px-3 sm:px-4 py-1.5 sm:py-2 rounded-full bg-neural-cyan/20 border-2 border-neural-cyan/40 text-neural-cyan font-mono text-xs sm:text-sm font-semibold">
                  Grok 4 Fast
                </div>
                <div className="px-3 sm:px-4 py-1.5 sm:py-2 rounded-full bg-neural-magenta/20 border-2 border-neural-magenta/40 text-neural-magenta font-mono text-xs sm:text-sm font-semibold">
                  Qwen3-Coder
                </div>
              </div>
            </div>

            {/* Model Selector Mockup - Shows the actual selector interface */}
            <div className="relative max-w-6xl mx-auto">
              <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-cyan/30 shadow-2xl shadow-neural-cyan/20 p-4 sm:p-6 md:p-8">
                {/* Terminal Header */}
                <div className="bg-[#11111b] -mx-4 sm:-mx-6 md:-mx-8 -mt-4 sm:-mt-6 md:-mt-8 px-3 sm:px-4 py-2 sm:py-3 mb-4 sm:mb-6 flex items-center justify-between border-b border-gray-700">
                  <div className="flex items-center gap-2">
                    <div className="flex gap-1.5 sm:gap-2">
                      <div className="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-[#ff5f56]"></div>
                      <div className="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-[#ffbd2e]"></div>
                      <div className="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-[#27c93f]"></div>
                    </div>
                    <div className="ml-2 sm:ml-4 text-[10px] sm:text-xs text-gray-400 font-mono">RyCode - Model Selector</div>
                  </div>
                </div>

                {/* Prompt */}
                <div className="flex gap-2 sm:gap-3 mb-4 sm:mb-6 font-mono text-xs sm:text-sm">
                  <span className="text-neural-cyan">❯</span>
                  <span className="text-neural-magenta">/model</span>
                </div>

                {/* Model Selector Interface */}
                <div className="space-y-2 sm:space-y-3 font-mono text-xs sm:text-sm">
                  <div className="text-gray-400 text-xs mb-4 uppercase tracking-wide">Select AI Model:</div>

                  {/* Claude Sonnet 4.5 - Active */}
                  <div className="flex items-center gap-2 sm:gap-4 px-2 sm:px-4 py-2 sm:py-3 rounded-lg bg-[#7aa2f7]/10 border-l-2 sm:border-l-4 border-[#7aa2f7]">
                    <span className="text-[#7aa2f7] text-base sm:text-xl">▶</span>
                    <div className="flex-1 min-w-0">
                      <div className="text-[#7aa2f7] font-bold text-sm sm:text-base">Claude Sonnet 4.5</div>
                      <div className="text-[10px] sm:text-xs text-gray-400 mt-0.5 sm:mt-1 hidden sm:block">Best coding model, 77.2% on SWE-bench Verified</div>
                    </div>
                    <div className="px-1.5 sm:px-2 py-0.5 sm:py-1 rounded bg-[#7aa2f7]/20 text-[#7aa2f7] text-[10px] sm:text-xs font-semibold shrink-0">ACTIVE</div>
                  </div>

                  {/* Gemini 2.5 Pro */}
                  <div className="flex items-center gap-2 sm:gap-4 px-2 sm:px-4 py-2 sm:py-3 rounded-lg hover:bg-gray-700/20 transition-colors">
                    <span className="text-gray-600 text-base sm:text-xl">○</span>
                    <div className="flex-1 min-w-0">
                      <div className="text-[#ea4aaa] font-bold text-sm sm:text-base">Gemini 2.5 Pro</div>
                      <div className="text-[10px] sm:text-xs text-gray-500 mt-0.5 sm:mt-1 hidden sm:block">Most intelligent with advanced thinking capabilities</div>
                    </div>
                  </div>

                  {/* GPT-5 */}
                  <div className="flex items-center gap-2 sm:gap-4 px-2 sm:px-4 py-2 sm:py-3 rounded-lg hover:bg-gray-700/20 transition-colors">
                    <span className="text-gray-600 text-base sm:text-xl">○</span>
                    <div className="flex-1 min-w-0">
                      <div className="text-[#ff6b35] font-bold text-sm sm:text-base">GPT-5</div>
                      <div className="text-[10px] sm:text-xs text-gray-500 mt-0.5 sm:mt-1 hidden sm:block">OpenAI's smartest model with built-in reasoning</div>
                    </div>
                  </div>

                  {/* Grok 4 Fast */}
                  <div className="flex items-center gap-2 sm:gap-4 px-2 sm:px-4 py-2 sm:py-3 rounded-lg hover:bg-gray-700/20 transition-colors">
                    <span className="text-gray-600 text-base sm:text-xl">○</span>
                    <div className="flex-1 min-w-0">
                      <div className="text-neural-cyan font-bold text-sm sm:text-base">Grok 4 Fast</div>
                      <div className="text-[10px] sm:text-xs text-gray-500 mt-0.5 sm:mt-1 hidden sm:block">2M context window, 98% cost reduction with reasoning</div>
                    </div>
                  </div>

                  {/* Qwen3-Coder */}
                  <div className="flex items-center gap-2 sm:gap-4 px-2 sm:px-4 py-2 sm:py-3 rounded-lg hover:bg-gray-700/20 transition-colors">
                    <span className="text-gray-600 text-base sm:text-xl">○</span>
                    <div className="flex-1 min-w-0">
                      <div className="text-neural-magenta font-bold text-sm sm:text-base">Qwen3-Coder</div>
                      <div className="text-[10px] sm:text-xs text-gray-500 mt-0.5 sm:mt-1 hidden sm:block">480B parameter agentic coder, 256K-1M context</div>
                    </div>
                  </div>
                </div>

                {/* Help Text */}
                <div className="mt-4 sm:mt-6 pt-3 sm:pt-4 border-t border-gray-700/50 text-[10px] sm:text-xs text-gray-500 font-mono">
                  <div className="flex flex-wrap items-center gap-3 sm:gap-6">
                    <span><span className="text-neural-cyan">↑↓</span> Navigate</span>
                    <span><span className="text-neural-cyan">Tab</span> Quick Switch</span>
                    <span><span className="text-neural-cyan">Enter</span> Select</span>
                    <span className="hidden sm:inline"><span className="text-neural-cyan">Esc</span> Cancel</span>
                  </div>
                </div>
              </div>

              <div className="mt-8 text-center space-y-4">
                <p className="text-neural-cyan text-lg font-semibold">
                  Type <span className="bg-black/40 px-3 py-1 rounded font-mono">/model</span> to switch between the latest frontier models
                </p>
                <div className="flex flex-wrap justify-center gap-6 text-sm text-gray-400">
                  <div className="flex items-center gap-2">
                    <span className="text-2xl">⚡</span>
                    <span>Instant switching with <span className="text-neural-cyan font-mono">Tab</span></span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-2xl">🧠</span>
                    <span>Context preserved across models</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-2xl">🚀</span>
                    <span>Zero configuration required</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* FOLD 2: Tab Switching + Model Chips Update */}
      <section className="relative min-h-screen flex items-center justify-center px-6 py-24 bg-black/20">
        <div className="mx-auto max-w-7xl w-full">
          <div className="text-center mb-16">
            <h2 className="text-5xl font-bold mb-4 text-gradient">
              Switch Models Instantly
            </h2>
            <p className="text-xl text-gray-400 max-w-2xl mx-auto">
              Press <span className="text-neural-cyan font-mono bg-black/40 px-3 py-1 rounded">Tab</span> to cycle through the latest frontier models. Watch the chips update in real-time.
            </p>
          </div>

          {/* Terminal showing Tab switching */}
          <div className="relative max-w-5xl mx-auto">
            <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-cyan/30 shadow-2xl shadow-neural-cyan/20">
              {/* Terminal Header with Model Chips */}
              <div className="bg-[#11111b] px-4 py-3 flex items-center justify-between border-b border-gray-700">
                <div className="flex items-center gap-2">
                  <div className="flex gap-2">
                    <div className="w-3 h-3 rounded-full bg-[#ff5f56]"></div>
                    <div className="w-3 h-3 rounded-full bg-[#ffbd2e]"></div>
                    <div className="w-3 h-3 rounded-full bg-[#27c93f]"></div>
                  </div>
                  <div className="ml-4 text-xs text-gray-400 font-mono">
                    RyCode - Multi-Agent CLI
                  </div>
                </div>

                {/* Model Chips - HIGHLIGHTED */}
                <div className="flex gap-2 relative">
                  <div className="absolute -inset-4 bg-neural-cyan/10 rounded-lg animate-pulse pointer-events-none"></div>
                  <div className="px-3 py-1 rounded-full bg-[#7aa2f7]/20 border border-[#7aa2f7]/40 text-[#7aa2f7] text-xs font-mono flex items-center gap-1">
                    <span className="w-1.5 h-1.5 rounded-full bg-[#7aa2f7] animate-pulse"></span>
                    Sonnet 4.5
                  </div>
                  <div className="px-3 py-1 rounded-full bg-gray-700/30 border border-gray-600 text-gray-500 text-xs font-mono">
                    Gemini 2.5
                  </div>
                  <div className="px-3 py-1 rounded-full bg-gray-700/30 border border-gray-600 text-gray-500 text-xs font-mono">
                    GPT-5
                  </div>
                </div>
              </div>

              {/* Terminal Content - Showing Tab key */}
              <div className="p-8 font-mono text-sm">
                <div className="flex gap-3 mb-4">
                  <span className="text-neural-cyan">❯</span>
                  <div className="flex-1">
                    <span className="text-gray-300">Type your prompt...</span>
                  </div>
                </div>

                <div className="mt-8 p-6 bg-black/30 rounded-lg border border-neural-cyan/20">
                  <div className="flex items-center justify-center gap-8">
                    <div className="text-center">
                      <div className="text-6xl mb-3">⇥</div>
                      <p className="text-neural-cyan font-semibold">Press Tab</p>
                      <p className="text-xs text-gray-400 mt-1">Cycle through models</p>
                    </div>
                    <div className="text-4xl text-gray-600">→</div>
                    <div className="text-center">
                      <div className="text-6xl mb-3">💫</div>
                      <p className="text-neural-magenta font-semibold">Chips Update</p>
                      <p className="text-xs text-gray-400 mt-1">See active model instantly</p>
                    </div>
                  </div>
                </div>

                <div className="mt-6 text-center text-xs text-gray-500">
                  <p>Sonnet 4.5 → Gemini 2.5 Pro → GPT-5 → Grok 4 Fast → Qwen3-Coder</p>
                </div>
              </div>

              {/* Status Bar */}
              <div className="bg-[#11111b] px-6 py-3 border-t border-gray-700 flex items-center justify-between text-xs">
                <div className="flex items-center gap-4 text-gray-400">
                  <span className="flex items-center gap-2">
                    <span className="w-2 h-2 rounded-full bg-matrix-green"></span>
                    Latest frontier models ready
                  </span>
                  <span>Tab to switch</span>
                </div>
                <div className="text-gray-500 font-mono">
                  Session: active
                </div>
              </div>
            </div>

            {/* Annotation */}
            <div className="mt-6 flex justify-center">
              <div className="bg-neural-cyan/10 backdrop-blur-sm border border-neural-cyan/30 rounded-lg px-6 py-3">
                <p className="text-sm text-neural-cyan font-mono">
                  ← Model chips update in real-time as you Tab through options
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* FOLD 3: Console in Action - THE DEAL CLOSER */}
      <section className="relative min-h-screen flex items-center justify-center px-6 py-24">
        <div className="mx-auto max-w-7xl w-full">
          <div className="text-center mb-16">
            <h2 className="text-5xl font-bold mb-4 text-gradient">
              See It In Action
            </h2>
            <p className="text-xl text-gray-400 max-w-2xl mx-auto">
              Real conversation. Real code generation. Real multi-agent switching.
            </p>
          </div>

          {/* Console Demo - Full Conversation Flow */}
          <div className="relative max-w-6xl mx-auto">
            <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-magenta/30 shadow-2xl shadow-neural-magenta/20">
              {/* Terminal Header */}
              <div className="bg-[#11111b] px-4 py-3 flex items-center justify-between border-b border-gray-700">
                <div className="flex items-center gap-2">
                  <div className="flex gap-2">
                    <div className="w-3 h-3 rounded-full bg-[#ff5f56]"></div>
                    <div className="w-3 h-3 rounded-full bg-[#ffbd2e]"></div>
                    <div className="w-3 h-3 rounded-full bg-[#27c93f]"></div>
                  </div>
                  <div className="ml-4 text-xs text-gray-400 font-mono">
                    RyCode - Live Session
                  </div>
                </div>

                <div className="flex gap-2">
                  <div className="px-3 py-1 rounded-full bg-[#7aa2f7]/20 border border-[#7aa2f7]/40 text-[#7aa2f7] text-xs font-mono flex items-center gap-1">
                    <span className="w-1.5 h-1.5 rounded-full bg-[#7aa2f7] animate-pulse"></span>
                    Claude
                  </div>
                </div>
              </div>

              {/* Terminal Content - Real conversation */}
              <div className="p-8 font-mono text-sm space-y-6 max-h-[600px] overflow-y-auto">
                {/* User prompt */}
                <div className="flex gap-3">
                  <span className="text-neural-cyan">❯</span>
                  <div className="flex-1">
                    <span className="text-white">Build me a REST API for a todo app with auth</span>
                  </div>
                </div>

                {/* Claude response */}
                <div className="pl-6 space-y-3">
                  <div className="text-gray-300">
                    <span className="text-[#7aa2f7] font-semibold">Claude:</span> I'll create a production-ready REST API with authentication. Let me set this up:
                  </div>

                  {/* Code output */}
                  <div className="bg-black/40 rounded-lg p-4 border border-[#7aa2f7]/20">
                    <div className="text-xs text-gray-500 mb-2">server.js</div>
                    <pre className="text-xs text-matrix-green overflow-x-auto">
{`const express = require('express');
const jwt = require('jsonwebtoken');
const bcrypt = require('bcrypt');

const app = express();
app.use(express.json());

// Auth middleware
const auth = (req, res, next) => {
  const token = req.header('Authorization');
  if (!token) return res.status(401).send('Access denied');

  try {
    const verified = jwt.verify(token, process.env.JWT_SECRET);
    req.user = verified;
    next();
  } catch (err) {
    res.status(400).send('Invalid token');
  }
};

// Routes
app.post('/register', async (req, res) => {
  const hashedPassword = await bcrypt.hash(req.body.password, 10);
  // Save user to DB...
  res.status(201).send('User registered');
});

app.get('/todos', auth, (req, res) => {
  // Return user's todos...
});`}
                    </pre>
                  </div>

                  <div className="text-gray-400 text-xs">
                    ✓ Express server configured<br/>
                    ✓ JWT authentication implemented<br/>
                    ✓ Password hashing with bcrypt
                  </div>
                </div>

                {/* Divider */}
                <div className="border-t border-gray-700/50 my-6"></div>

                {/* User switches model */}
                <div className="flex gap-3">
                  <span className="text-neural-cyan">❯</span>
                  <div className="flex-1">
                    <span className="text-neural-magenta italic">*presses Tab*</span>
                    <span className="ml-3 text-xs text-gray-500">→ Switched to GPT-4</span>
                  </div>
                </div>

                {/* User asks follow-up */}
                <div className="flex gap-3">
                  <span className="text-neural-cyan">❯</span>
                  <div className="flex-1">
                    <span className="text-white">Add rate limiting and input validation</span>
                  </div>
                </div>

                {/* GPT-4 response */}
                <div className="pl-6 space-y-3">
                  <div className="text-gray-300">
                    <span className="text-[#10a37f] font-semibold">GPT-4:</span> I'll add express-rate-limit and joi validation:
                  </div>

                  <div className="bg-black/40 rounded-lg p-4 border border-[#10a37f]/20">
                    <pre className="text-xs text-matrix-green overflow-x-auto">
{`const rateLimit = require('express-rate-limit');
const Joi = require('joi');

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 100
});

app.use(limiter);

const registerSchema = Joi.object({
  email: Joi.string().email().required(),
  password: Joi.string().min(8).required()
});`}
                    </pre>
                  </div>
                </div>

                {/* User switches again */}
                <div className="border-t border-gray-700/50 my-6"></div>

                <div className="flex gap-3">
                  <span className="text-neural-cyan">❯</span>
                  <div className="flex-1">
                    <span className="text-neural-magenta italic">*presses Tab again*</span>
                    <span className="ml-3 text-xs text-gray-500">→ Switched to Codex</span>
                  </div>
                </div>

                <div className="flex gap-3">
                  <span className="text-neural-cyan">❯</span>
                  <div className="flex-1">
                    <span className="text-white">Write unit tests for the auth middleware</span>
                  </div>
                </div>

                {/* Codex response */}
                <div className="pl-6 space-y-3">
                  <div className="text-gray-300">
                    <span className="text-[#ff6b35] font-semibold">Codex:</span> Testing authentication with Jest:
                  </div>

                  <div className="bg-black/40 rounded-lg p-4 border border-[#ff6b35]/20">
                    <pre className="text-xs text-matrix-green overflow-x-auto">
{`describe('Auth Middleware', () => {
  test('rejects requests without token', async () => {
    const res = await request(app)
      .get('/todos')
      .expect(401);
    expect(res.text).toBe('Access denied');
  });

  test('accepts valid JWT tokens', async () => {
    const token = jwt.sign({ id: 1 }, process.env.JWT_SECRET);
    const res = await request(app)
      .get('/todos')
      .set('Authorization', token)
      .expect(200);
  });
});`}
                    </pre>
                  </div>
                </div>
              </div>

              {/* Status Bar */}
              <div className="bg-[#11111b] px-6 py-3 border-t border-gray-700 flex items-center justify-between text-xs">
                <div className="flex items-center gap-4 text-gray-400">
                  <span className="flex items-center gap-2">
                    <span className="w-2 h-2 rounded-full bg-matrix-green animate-pulse"></span>
                    3 models used in this session
                  </span>
                  <span>Context preserved across switches</span>
                </div>
                <div className="text-gray-500 font-mono">
                  45 lines generated • 0 errors
                </div>
              </div>
            </div>

            {/* Value Proposition */}
            <div className="mt-12 grid md:grid-cols-3 gap-6">
              <div className="text-center p-6 bg-black/40 backdrop-blur-sm rounded-lg border border-[#7aa2f7]/20">
                <div className="text-4xl mb-3">🎯</div>
                <h3 className="text-lg font-bold text-[#7aa2f7] mb-2">Context Preserved</h3>
                <p className="text-sm text-gray-400">
                  Switch models mid-conversation without losing context
                </p>
              </div>

              <div className="text-center p-6 bg-black/40 backdrop-blur-sm rounded-lg border border-neural-cyan/20">
                <div className="text-4xl mb-3">⚡</div>
                <h3 className="text-lg font-bold text-neural-cyan mb-2">Instant Switching</h3>
                <p className="text-sm text-gray-400">
                  Press Tab to get a second opinion in milliseconds
                </p>
              </div>

              <div className="text-center p-6 bg-black/40 backdrop-blur-sm rounded-lg border border-matrix-green/20">
                <div className="text-4xl mb-3">🧠</div>
                <h3 className="text-lg font-bold text-matrix-green mb-2">Best of All Models</h3>
                <p className="text-sm text-gray-400">
                  Claude for reasoning, Gemini for speed, Codex for code
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Built with toolkit-cli Section */}
      <section className="px-6 py-24 lg:px-8 bg-black/20">
        <div className="mx-auto max-w-4xl text-center">
          <h2 className="text-4xl font-bold mb-6 text-gradient">
            Built with toolkit-cli
          </h2>
          <p className="text-xl text-gray-300 mb-8">
            RyCode is a showcase of what's possible with{' '}
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

          <div className="grid md:grid-cols-2 gap-6 mt-12">
            <div className="p-6 rounded-lg bg-black/40 border border-claude-blue/30">
              <h3 className="text-xl font-bold mb-3 text-claude-blue">100% AI-Designed</h3>
              <p className="text-gray-300 text-sm">
                From concept to completion by Claude AI using toolkit-cli's multi-agent architecture.
              </p>
            </div>
            <div className="p-6 rounded-lg bg-black/40 border border-performance-gold/30">
              <h3 className="text-xl font-bold mb-3 text-performance-gold">Production Ready</h3>
              <p className="text-gray-300 text-sm">
                31/31 tests passing, 54.2% coverage, zero known bugs. Built for production.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="px-6 py-12 border-t border-neural-cyan/20">
        <div className="mx-auto max-w-7xl text-center text-gray-400">
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
          <p className="mt-2 text-xs">
            Zero Compromises • Infinite Attention to Detail
          </p>
        </div>
      </footer>
    </main>
  )
}
