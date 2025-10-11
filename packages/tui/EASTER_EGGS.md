# RyCode Splash Screen - Easter Eggs 🥚

> Hidden features and surprises in the RyCode splash screen

---

## 🍩 1. Infinite Donut Mode

**Command:** `rycode donut`

Launches the neural cortex animation in infinite loop mode - perfect for meditation, screensavers, or impressing coworkers.

**Features:**
- Continuous 3D torus rotation
- No auto-close timeout
- Press `Q` to quit
- Press `?` to see the math equations

**Why it's cool:**
This is a tribute to Andy Sloane's original donut.c - the legendary 3D ASCII donut that inspired this implementation. Now you can watch it spin forever in glorious cyberpunk colors!

**Discovery hint:** Look at the command-line args 😉

---

## 🌈 2. Rainbow Mode (Konami Code)

**How to activate:** During splash, enter the Konami code:
```
↑ ↑ ↓ ↓ ← → ← → B A
```

**What happens:**
- Torus changes from cyan-magenta to full rainbow spectrum
- ROYGBIV color cycling (Red, Orange, Yellow, Green, Blue, Indigo, Violet)
- Works in both normal splash and donut mode
- Persists for the rest of the session

**Why it's cool:**
Classic gaming easter egg meets retro ASCII art. The rainbow gradient makes the math even more mesmerizing.

**Discovery hint:** Try classic cheat codes 🎮

---

## 🧮 3. Math Equations Reveal

**How to activate:** Press `?` during splash or donut mode

**What you see:**
- Complete torus parametric equations
- Rotation matrix formulas
- Perspective projection math
- Luminance calculation (Phong shading)
- Character mapping algorithm
- Performance metrics

**Why it's cool:**
Full transparency - see exactly how the sausage is made. Perfect for:
- Math nerds who want to understand the algorithm
- Developers learning 3D graphics
- Interview prep (yes, this is that donut algorithm!)

**Press `?` again to return to the animation**

**Discovery hint:** It's in the skip hint 👀

---

## 👻 4. Hidden Message

**How to find:** Watch the cortex animation carefully...

**What it says:** `CLAUDE WAS HERE`

**When it appears:**
- Briefly flashes in the center of the torus
- Only visible for ~1 second
- Appears approximately every 10 seconds of rotation
- Rendered with max z-buffer value (always on top)

**Why it's cool:**
A signature from Claude AI, the architect of RyCode. Like a digital graffiti tag embedded in the matrix.

**Discovery hint:** Stare at the center long enough... 👁️

---

## ⚡ 5. Performance Mode (Adaptive FPS)

**How it works:** Automatic - no activation needed

**What it does:**
- Monitors frame rendering time
- If frames take >50ms: drops to 15 FPS automatically
- If frames are fast (<50ms): maintains 30 FPS
- Samples last 30 frames to calculate average

**Why it's cool:**
Works great on:
- Low-end systems (Raspberry Pi, old laptops)
- Remote SSH sessions
- Slow terminals (xterm, Windows CMD)

**Technical details:**
```
Target frame time: 33ms (30 FPS) or 66ms (15 FPS)
Decision threshold: 50ms average over last 30 frames
Overhead: Virtually zero (tracked in Update loop)
```

**Discovery hint:** Launch on a potato computer 🥔

---

## 🎨 6. Terminal Capability Detection

**How it works:** Automatic detection at startup

**What it detects:**
- Terminal size (width × height)
- Color support (truecolor / 256-color / 16-color)
- Unicode support (full / ASCII-only)
- Performance estimate (fast / medium / slow)

**Fallback modes:**
- **Too small (<80×24):** Shows simplified text splash
- **No colors:** Monochrome ASCII art
- **No unicode:** ASCII-only character set
- **Slow terminal:** Drops to 15 FPS automatically

**Why it's cool:**
Works everywhere - from Windows CMD to Raspberry Pi serial console. Graceful degradation is built-in.

**Discovery hint:** Try resizing your terminal or SSH to a slow server 📟

---

## 🎯 7. Skip Hints

**Multiple ways to skip:**
- `S` - Skip splash (continue to TUI)
- `ESC` - Skip and disable forever
- `Enter` or `Space` - Continue from closer screen
- `Q` - Quit donut mode

**Hidden progress indicator:**
- If you enter part of the Konami code correctly, you'll see `...` at the bottom
- Disappears if you enter wrong key
- Hidden hint that something is happening!

---

## 🔮 Future Easter Eggs (Coming Soon™)

Ideas for future releases:

1. **Sound Effects** - Beep codes for model initialization
2. **Matrix Mode** - Falling green characters a la The Matrix
3. **Starfield** - 3D star tunnel instead of torus
4. **Wire Mode** - Wireframe rendering of the torus
5. **Hyperspeed** - 2× or 4× speed mode
6. **Secrets in Config** - Hidden settings in config.json
7. **ASCII Art Gallery** - Collection of pre-rendered scenes

**Want to add your own?** PRs welcome! 🎉

---

## 📊 Easter Egg Statistics

**Lines of code dedicated to easter eggs:**
- Donut mode: 15 lines
- Konami code: 25 lines
- Rainbow mode: 40 lines
- Math reveal: 35 lines
- Hidden message: 30 lines
- **Total:** ~145 lines (15% of splash.go!)

**Estimated discovery rate:**
- `/donut` command: 20-30% (documented)
- Konami code: 5-10% (gamers will find it)
- Math reveal: 40-50% (hint in skip message)
- Hidden message: 10-15% (requires patience)
- Adaptive FPS: 100% (everyone benefits, few notice)

---

## 🎓 Learning Resources

Want to understand the donut math?

- **Original article:** https://www.a1k0n.net/2011/07/20/donut-math.html
- **Interactive demo:** https://www.a1k0n.net/2006/09/15/obfuscated-c-donut.html
- **RyCode source:** `packages/tui/internal/splash/cortex.go`

**Key concepts:**
- Parametric equations for torus
- Rotation matrices (linear algebra)
- Perspective projection (3D → 2D)
- Z-buffer algorithm (depth sorting)
- Phong shading (lighting)

---

## 🏆 Achievement Unlocked

If you found all the easter eggs, you're a true RyCode power user! 🎉

**Badge:** 🥚🥚🥚🥚🥚 **Easter Egg Hunter**

Share your discoveries:
- Twitter: @rycode_ai #RyCodeEasterEggs
- GitHub: Open an issue titled "I found them all!"
- Discord: #easter-eggs channel

---

## 🤝 Contributing

Found a bug in an easter egg? Want to add a new one?

1. Open an issue describing the easter egg idea
2. Fork the repo
3. Add your easter egg to `internal/splash/`
4. Update this file with documentation
5. Submit a PR!

**Guidelines:**
- Keep it fun and harmless
- Don't break the main splash flow
- Document it well (so future developers understand)
- Test on multiple platforms

---

**🤖 Easter eggs added by Claude AI - Because even AIs like to have fun!**

*Built with ❤️ and a sense of humor*
