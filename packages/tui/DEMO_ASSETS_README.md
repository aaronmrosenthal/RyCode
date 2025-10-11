# RyCode Splash Screen Demo Assets

> **Quick reference for creating and using splash screen marketing materials**

---

## 📁 Files Created

### VHS Tape Files (Automated Recording)
```
splash_demo.tape              - Standard splash animation (5s)
splash_demo_donut.tape        - Infinite donut mode with easter eggs (20s)
```

### Documentation
```
SPLASH_DEMO_CREATION.md       - Complete guide for creating GIFs, videos, screenshots
```

### Scripts
```
scripts/record_splash_simple.sh  - Simple recording helper (no external tools)
```

---

## 🚀 Quick Start (3 Options)

### Option 1: VHS (Best Quality, Automated)

**Install VHS:**
```bash
brew install vhs
```

**Generate Demos:**
```bash
cd /Users/aaron/Code/RyCode/RyCode/packages/tui

# Standard splash
vhs splash_demo.tape
# Output: splash_demo.gif

# Donut mode with easter eggs
vhs splash_demo_donut.tape
# Output: splash_demo_donut.gif
```

**Result:** Production-ready GIF files optimized for web

---

### Option 2: Manual Recording (macOS Built-in)

**Run helper script:**
```bash
./scripts/record_splash_simple.sh
```

**Then use:**
- **Cmd+Shift+5** - macOS screenshot tool (screen recording)
- **QuickTime Player** - File → New Screen Recording

**Result:** .mov video file (convert to GIF with FFmpeg or online tool)

---

### Option 3: asciinema (Web Embeddable)

**Install:**
```bash
brew install asciinema
```

**Record:**
```bash
asciinema rec splash_demo.cast --overwrite
./rycode --splash
# Wait 6 seconds
# Press Ctrl+D
```

**Result:** .cast file that can be embedded on landing page

---

## 📊 Recommended Assets for Landing Page

### Hero Fold
- **Primary:** `splash_demo.gif` (optimized, <2MB)
- **Alternative:** `splash_demo.cast` (interactive asciinema player)

### Easter Eggs Section
- `splash_demo_donut.gif` - Shows infinite mode + all easter eggs
- Individual screenshots for each egg (capture manually)

### Social Media
- Twitter/LinkedIn: 720p MP4 video (create from GIF using FFmpeg)
- Instagram: Square format MP4 (1:1 aspect ratio)

---

## 🎯 Asset Specifications

### GIF Requirements
- **Dimensions:** 1200×800 (can scale down)
- **File size:** <2MB (optimized)
- **Frame rate:** 30 FPS (matches splash)
- **Colors:** 128-256 (good quality, reasonable size)
- **Loop:** Yes (continuous playback)

### Video Requirements
- **Format:** MP4 (H.264)
- **Resolution:** 1080p or 720p
- **Bitrate:** 2-4 Mbps
- **Duration:** 5-20 seconds
- **Audio:** Optional (can add music/narration)

### Screenshot Requirements
- **Format:** PNG (lossless)
- **Resolution:** Native terminal size
- **Purpose:** Documentation, blog posts

---

## 📚 Full Documentation

See **SPLASH_DEMO_CREATION.md** for:
- Complete tool installation guides
- VHS tape file customization
- FFmpeg video conversion recipes
- GIF optimization techniques
- asciinema web embedding
- Troubleshooting guide
- Advanced workflows

---

## ✅ Checklist for Landing Page

**Assets Needed:**
- [ ] Install VHS (`brew install vhs`)
- [ ] Generate `splash_demo.gif` with VHS
- [ ] Generate `splash_demo_donut.gif` with VHS
- [ ] Optimize GIFs (<2MB each)
- [ ] (Optional) Create .cast file for interactive player
- [ ] (Optional) Convert GIF to MP4 for social media

**Integration:**
- [ ] Upload GIFs to `/public/assets/` in Next.js project
- [ ] Add to Hero fold component
- [ ] Add to Easter Eggs section
- [ ] Add alt text and captions
- [ ] Test loading performance
- [ ] Verify loop behavior

---

## 🎨 Usage Examples

### Next.js Landing Page

```tsx
// Hero Fold Component
import Image from 'next/image';

export function HeroFold() {
  return (
    <div className="relative">
      <Image
        src="/assets/splash_demo_optimized.gif"
        alt="RyCode 3D Neural Cortex Splash Screen - Real donut algorithm math rendering at 30 FPS"
        width={1200}
        height={800}
        priority
        className="rounded-lg shadow-2xl"
      />
    </div>
  );
}
```

### Interactive Player (asciinema)

```tsx
// Install: npm install asciinema-player
import 'asciinema-player/dist/bundle/asciinema-player.css';
import AsciinemaPlayer from 'asciinema-player';

export function InteractiveSplashDemo() {
  useEffect(() => {
    AsciinemaPlayer.create('/assets/splash_demo.cast', document.getElementById('demo'), {
      cols: 120,
      rows: 30,
      autoPlay: true,
      loop: true,
      theme: 'dracula',
    });
  }, []);

  return <div id="demo" className="rounded-lg shadow-2xl" />;
}
```

---

## 🔥 Marketing Copy for Assets

### Social Media Posts

**Twitter/X:**
```
🌀 Just shipped: Epic 3D ASCII splash screen for RyCode!

✨ Real donut algorithm math
⚡ 30 FPS smooth animation
🎮 5 hidden easter eggs
🤖 100% AI-designed by Claude

Built with toolkit-cli → Try it: ry-code.com

[GIF: splash_demo_optimized.gif]
```

**LinkedIn:**
```
Excited to showcase RyCode's new 3D terminal splash screen! 🚀

This isn't just eye candy—it's a technical demonstration of what's possible with modern terminal graphics:

✅ Real torus parametric equations (not fake ASCII art)
✅ Z-buffer depth sorting for proper occlusion
✅ 30 FPS rendering (0.318ms per frame—85× faster than needed!)
✅ Adaptive accessibility (respects PREFERS_REDUCED_MOTION)

RyCode is built entirely with toolkit-cli, Anthropic's official AI toolkit for creating multi-agent CLI tools.

See it in action: ry-code.com

#AI #CLI #Terminal #Developer Tools

[Video: splash_demo_720p.mp4]
```

### Blog Post Hero Image
- Use: `splash_demo.gif` or high-res screenshot
- Caption: "RyCode's 3D Neural Cortex splash screen—real math, real performance"

---

## 📐 Technical Specs (For Reference)

**Current Implementation:**
- Rendering engine: Go + Bubble Tea
- Frame time: 0.318ms (M1 Max)
- Animation: 3-act sequence (Boot → Cortex → Closer)
- Duration: ~5 seconds
- Colors: Cyan-magenta gradient
- Math: Torus parametric equations with rotation matrices

**Demo Targets:**
- Capture all 3 acts
- Show smooth 30 FPS animation
- Highlight cyberpunk color palette
- Demonstrate auto-close behavior

---

## 🎬 Next Steps

1. **Install VHS** (5 min):
   ```bash
   brew install vhs
   ```

2. **Generate GIFs** (5 min):
   ```bash
   vhs splash_demo.tape
   vhs splash_demo_donut.tape
   ```

3. **Optimize** (2 min):
   ```bash
   # Install ImageMagick if needed
   brew install imagemagick

   # Optimize
   convert splash_demo.gif -fuzz 10% -layers Optimize -colors 128 splash_demo_optimized.gif
   ```

4. **Verify** (1 min):
   ```bash
   ls -lh splash_demo*.gif
   open splash_demo_optimized.gif
   ```

5. **Ready for landing page!** ✅

---

**🤖 Asset Creation Guide by Claude AI**

*Ready for ry-code.com landing page integration*

---

**Total Time to Generate Assets:** ~15 minutes
**Output:** 2-4 production-ready GIF/video files
**Status:** Ready for Execution ✅

