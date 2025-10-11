'use client'

export function TUIDemo() {
  return (
    <div className="relative max-w-5xl mx-auto">
      {/* Terminal Window */}
      <div className="bg-[#1e1e2e] rounded-xl overflow-hidden border-2 border-neural-cyan/30 shadow-2xl shadow-neural-cyan/20">
        {/* Terminal Header */}
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
          
          {/* Model Chips (Right Corner) */}
          <div className="flex gap-2">
            <div className="px-3 py-1 rounded-full bg-[#7aa2f7]/20 border border-[#7aa2f7]/40 text-[#7aa2f7] text-xs font-mono flex items-center gap-1">
              <span className="w-1.5 h-1.5 rounded-full bg-[#7aa2f7] animate-pulse"></span>
              Claude
            </div>
            <div className="px-3 py-1 rounded-full bg-gray-700/30 border border-gray-600 text-gray-500 text-xs font-mono">
              GPT-4
            </div>
            <div className="px-3 py-1 rounded-full bg-gray-700/30 border border-gray-600 text-gray-500 text-xs font-mono">
              Gemini
            </div>
          </div>
        </div>

        {/* Terminal Content */}
        <div className="p-6 font-mono text-sm space-y-4">
          {/* Previous conversation */}
          <div className="flex gap-3 text-gray-500">
            <span className="text-neural-cyan">❯</span>
            <div>help</div>
          </div>

          <div className="pl-6 text-gray-400 space-y-1 text-xs">
            <div>Available commands:</div>
            <div className="pl-4 space-y-1">
              <div>/model - Switch AI model</div>
              <div>/clear - Clear screen</div>
              <div>/help - Show this help</div>
            </div>
          </div>

          {/* Divider */}
          <div className="border-t border-gray-700/50"></div>

          {/* User typing /model command */}
          <div className="flex gap-3">
            <span className="text-neural-cyan">❯</span>
            <div className="flex-1">
              <span className="text-neural-magenta">/model</span>
              <span className="ml-2 w-2 h-4 bg-neural-cyan inline-block animate-pulse"></span>
            </div>
          </div>

          {/* Model Selection UI */}
          <div className="pl-6 space-y-3 bg-[#181825]/50 rounded-lg p-4 border border-neural-cyan/20">
            <div className="text-gray-400 text-xs mb-3">
              Select model (use Tab or arrow keys):
            </div>
            
            {/* Model Options */}
            <div className="space-y-2">
              <div className="flex items-center gap-3 px-3 py-2 rounded bg-neural-cyan/10 border-l-2 border-neural-cyan">
                <span className="text-neural-cyan">▶</span>
                <div className="flex-1">
                  <div className="text-[#7aa2f7] font-semibold">Claude (Sonnet 4.5)</div>
                  <div className="text-xs text-gray-400">Best for reasoning & code</div>
                </div>
                <div className="text-xs text-neural-cyan">Active</div>
              </div>

              <div className="flex items-center gap-3 px-3 py-2 rounded hover:bg-gray-700/30 text-gray-400">
                <span className="text-gray-600">○</span>
                <div className="flex-1">
                  <div className="text-gray-300">GPT-4 Turbo</div>
                  <div className="text-xs text-gray-500">Fast & capable</div>
                </div>
              </div>

              <div className="flex items-center gap-3 px-3 py-2 rounded hover:bg-gray-700/30 text-gray-400">
                <span className="text-gray-600">○</span>
                <div className="flex-1">
                  <div className="text-gray-300">Gemini 2.0 Flash</div>
                  <div className="text-xs text-gray-500">Ultra fast responses</div>
                </div>
              </div>

              <div className="flex items-center gap-3 px-3 py-2 rounded hover:bg-gray-700/30 text-gray-400">
                <span className="text-gray-600">○</span>
                <div className="flex-1">
                  <div className="text-gray-300">Codex</div>
                  <div className="text-xs text-gray-500">Specialized for code</div>
                </div>
              </div>

              <div className="flex items-center gap-3 px-3 py-2 rounded hover:bg-gray-700/30 text-gray-400">
                <span className="text-gray-600">○</span>
                <div className="flex-1">
                  <div className="text-gray-300">Mixtral 8x22B</div>
                  <div className="text-xs text-gray-500">Open source powerhouse</div>
                </div>
              </div>

              <div className="flex items-center gap-3 px-3 py-2 rounded hover:bg-gray-700/30 text-gray-400">
                <span className="text-gray-600">○</span>
                <div className="flex-1">
                  <div className="text-gray-300">DeepSeek V3</div>
                  <div className="text-xs text-gray-500">Deep reasoning</div>
                </div>
              </div>
            </div>

            <div className="text-xs text-gray-500 pt-2 border-t border-gray-700/50">
              Press Tab to cycle • Enter to select • Esc to cancel
            </div>
          </div>

          {/* Bottom prompt indicator */}
          <div className="flex gap-3 pt-2">
            <span className="text-neural-cyan">❯</span>
            <div className="flex-1 text-gray-600">
              Waiting for selection...
            </div>
          </div>
        </div>

        {/* Status Bar */}
        <div className="bg-[#11111b] px-6 py-3 border-t border-gray-700 flex items-center justify-between text-xs">
          <div className="flex items-center gap-4 text-gray-400">
            <span className="flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-matrix-green"></span>
              6 models available
            </span>
            <span>Tab to switch</span>
            <span>/help for commands</span>
          </div>
          <div className="text-gray-500 font-mono">
            Session: active
          </div>
        </div>
      </div>

      {/* Floating Label - Right */}
      <div className="absolute -right-8 top-12 bg-neural-cyan/10 backdrop-blur-sm border border-neural-cyan/30 rounded-lg px-4 py-2 animate-pulse">
        <p className="text-xs text-neural-cyan font-mono whitespace-nowrap">
          ← Model chips update in real-time
        </p>
      </div>

      {/* Floating Label - Left */}
      <div className="absolute -left-8 top-1/2 bg-neural-magenta/10 backdrop-blur-sm border border-neural-magenta/30 rounded-lg px-4 py-2">
        <p className="text-xs text-neural-magenta font-mono whitespace-nowrap">
          Tab to switch models →
        </p>
      </div>
    </div>
  )
}
