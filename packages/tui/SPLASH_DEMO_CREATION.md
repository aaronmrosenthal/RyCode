# RyCode Splash Screen Demo Creation Guide

> **Instructions for creating marketing assets (GIF, video, screenshots) for the splash screen**

---

## 🎬 Overview

This guide provides step-by-step instructions for creating high-quality demo assets showcasing the RyCode splash screen's 3D neural cortex animation.

---

## 📦 Prerequisites

### Required Tools

**1. VHS (Charmbracelet's Terminal Recorder)**
```bash
# macOS
brew install vhs

# Linux
go install github.com/charmbracelet/vhs@latest

# Verify installation
vhs --version
```

**2. asciinema (Alternative Terminal Recorder)**
```bash
# macOS
brew install asciinema

# Linux
pip3 install asciinema

# Verify installation
asciinema --version
```

**3. FFmpeg (For video conversion)**
```bash
# macOS
brew install ffmpeg

# Linux
sudo apt-get install ffmpeg

# Verify installation
ffmpeg -version
```

**4. ImageMagick (For GIF optimization)**
```bash
# macOS
brew install imagemagick

# Linux
sudo apt-get install imagemagick

# Verify installation
convert -version
```

---

## 🎥 Method 1: VHS Tape Files (Recommended)

### A. Standard Splash Demo

**File:** `splash_demo.tape`

**Generate GIF:**
```bash
cd /Users/aaron/Code/RyCode/RyCode/packages/tui
vhs splash_demo.tape
```

**Output:** `splash_demo.gif` (~1-3 MB)

**What it shows:**
- Build process
- Launch with `--splash` flag
- Full 3-act animation (Boot → Cortex → Closer)
- Auto-close after 5 seconds

### B. Infinite Donut Mode Demo

**File:** `splash_demo_donut.tape`

**Generate GIF:**
```bash
cd /Users/aaron/Code/RyCode/RyCode/packages/tui
vhs splash_demo_donut.tape
```

**Output:** `splash_demo_donut.gif` (~5-10 MB)

**What it shows:**
- Infinite donut mode (`./rycode donut`)
- Math equations reveal (`?` key)
- Konami code activation (↑↑↓↓←→←→BA)
- Rainbow mode
- Quit command (`q`)

### C. Customize VHS Settings

Edit `.tape` files to adjust:
```tape
Set FontSize 14        # Increase for better readability
Set Width 1200         # Terminal width in pixels
Set Height 800         # Terminal height in pixels
Set Padding 20         # Padding around terminal
Set Theme "Dracula"    # Color theme (Dracula, Nord, Monokai, etc.)
```

**Available themes:**
- Dracula (default, best for cyberpunk aesthetic)
- Nord (cool blues)
- Monokai (warm tones)
- Solarized Dark
- Tomorrow Night

---

## 🎬 Method 2: asciinema (For Web Embedding)

### A. Record Session

```bash
# Build RyCode
cd /Users/aaron/Code/RyCode/RyCode/packages/tui
go build -o rycode ./cmd/rycode

# Record splash screen
asciinema rec splash_demo.cast --overwrite

# Inside recording session:
./rycode --splash
# Wait 6 seconds for splash to complete
# Press Ctrl+D to stop recording
```

### B. Upload to asciinema.org

```bash
# Upload and get shareable URL
asciinema upload splash_demo.cast
```

**Output:** `https://asciinema.org/a/XXXXXX`

### C. Convert to GIF

```bash
# Install agg (asciinema GIF generator)
cargo install --git https://github.com/asciinema/agg

# Convert to GIF
agg splash_demo.cast splash_demo.gif \
  --font-size 14 \
  --theme dracula \
  --speed 1.0
```

### D. Embed in Landing Page

```html
<!-- Using asciinema-player -->
<script src="https://asciinema.org/a.js" id="asciicast-XXXXXX" async></script>

<!-- Or self-hosted with asciinema-player -->
<asciinema-player src="splash_demo.cast" cols="120" rows="30"></asciinema-player>
```

---

## 📸 Method 3: Screenshots (For Documentation)

### A. macOS Built-in Screenshot

```bash
# Build and run
cd /Users/aaron/Code/RyCode/RyCode/packages/tui
go build -o rycode ./cmd/rycode
./rycode --splash

# While splash is running:
# Press Cmd+Shift+4, then Space, then click terminal window
```

**Output:** `Screen Shot YYYY-MM-DD at HH.MM.SS.png` on Desktop

### B. Programmatic Screenshots with `screencapture`

```bash
# Take screenshot after 3 seconds (gives time to focus terminal)
screencapture -w -T 3 splash_screenshot.png
```

### C. Capture Multiple Frames

Create a script to capture animation frames:

```bash
#!/bin/bash
# capture_frames.sh

./rycode --splash &
RYCODE_PID=$!

# Capture at 0.5s intervals
for i in {1..10}; do
  sleep 0.5
  screencapture -w splash_frame_$i.png
done

wait $RYCODE_PID
```

**Combine into GIF:**
```bash
convert -delay 50 -loop 0 splash_frame_*.png splash_animation.gif
```

---

## 🎨 Method 4: High-Quality Video (For Social Media)

### A. Record Terminal with QuickTime (macOS)

1. Open **QuickTime Player**
2. File → New Screen Recording
3. Click red record button
4. Select terminal window area
5. Run: `./rycode --splash`
6. Wait 6 seconds
7. Stop recording (⌘+Control+Esc)
8. File → Export As → 1080p

### B. Convert to Twitter/LinkedIn Format

```bash
# Convert to MP4 with optimal settings
ffmpeg -i splash_demo.mov \
  -vf "scale=1280:720" \
  -c:v libx264 \
  -preset slow \
  -crf 18 \
  -c:a aac \
  -b:a 192k \
  splash_demo_720p.mp4

# Twitter optimized (square format)
ffmpeg -i splash_demo.mov \
  -vf "crop=720:720,scale=720:720" \
  -c:v libx264 \
  -preset slow \
  -crf 18 \
  splash_demo_twitter.mp4
```

### C. Add Subtitles/Captions

Create `subtitles.srt`:
```srt
1
00:00:00,000 --> 00:00:02,000
RyCode - AI-Powered Multi-Agent Terminal

2
00:00:02,000 --> 00:00:05,000
3D Neural Cortex Animation
Real Donut Algorithm Math

3
00:00:05,000 --> 00:00:06,000
30 FPS Smooth Rendering
```

**Burn subtitles into video:**
```bash
ffmpeg -i splash_demo.mp4 \
  -vf "subtitles=subtitles.srt:force_style='FontName=Inter,FontSize=24,PrimaryColour=&H00FFFF'" \
  splash_demo_with_subs.mp4
```

---

## 🚀 Optimized GIF Creation Workflow

### Step 1: Record with VHS
```bash
vhs splash_demo.tape
```

### Step 2: Optimize GIF Size
```bash
# Reduce colors and optimize (from ~3MB to ~1MB)
convert splash_demo.gif \
  -fuzz 10% \
  -layers Optimize \
  -colors 128 \
  splash_demo_optimized.gif

# Further compression with gifsicle
gifsicle -O3 --colors 128 splash_demo_optimized.gif -o splash_demo_final.gif
```

### Step 3: Verify Quality
```bash
# Check file size
ls -lh splash_demo_final.gif

# Preview
open splash_demo_final.gif
```

**Target:** <2 MB for web, <1 MB for GitHub README

---

## 📊 Asset Checklist

Create the following assets for complete coverage:

### GIFs (For Web/GitHub)
- [ ] `splash_demo.gif` - Standard splash (5 seconds) <2MB
- [ ] `splash_demo_donut.gif` - Infinite donut mode (20 seconds) <5MB
- [ ] `splash_konami_code.gif` - Konami code demo (10 seconds) <3MB
- [ ] `splash_math_reveal.gif` - Math equations (5 seconds) <2MB

### Screenshots (For Documentation)
- [ ] `splash_boot.png` - Boot sequence frame
- [ ] `splash_cortex.png` - Rotating cortex frame
- [ ] `splash_closer.png` - Closer screen
- [ ] `splash_rainbow.png` - Rainbow mode
- [ ] `splash_math.png` - Math equations display

### Videos (For Social Media)
- [ ] `splash_demo_1080p.mp4` - Full HD (YouTube)
- [ ] `splash_demo_720p.mp4` - HD (LinkedIn)
- [ ] `splash_demo_twitter.mp4` - Square format (Twitter/Instagram)

### asciinema Casts (For Web Embedding)
- [ ] `splash_demo.cast` - Standard splash
- [ ] `splash_donut.cast` - Infinite donut mode

---

## 🎯 Landing Page Usage

### Hero Fold
**Use:** `splash_demo.gif` or `splash_demo.cast` (embedded player)

**HTML:**
```html
<!-- Option 1: Optimized GIF -->
<img
  src="/assets/splash_demo_optimized.gif"
  alt="RyCode 3D Neural Cortex Splash Screen"
  width="1200"
  height="800"
  class="rounded-lg shadow-2xl"
/>

<!-- Option 2: asciinema player (interactive) -->
<asciinema-player
  src="/assets/splash_demo.cast"
  cols="120"
  rows="30"
  autoplay
  loop
  theme="dracula"
></asciinema-player>
```

### Easter Eggs Section
**Use:** Multiple GIFs showing each easter egg

```html
<div class="grid grid-cols-2 gap-8">
  <div>
    <h3>Infinite Donut Mode</h3>
    <img src="/assets/splash_donut.gif" alt="Infinite Donut Mode" />
  </div>
  <div>
    <h3>Konami Code</h3>
    <img src="/assets/splash_konami.gif" alt="Rainbow Mode" />
  </div>
  <div>
    <h3>Math Reveal</h3>
    <img src="/assets/splash_math.gif" alt="Torus Equations" />
  </div>
  <div>
    <h3>Skip Controls</h3>
    <img src="/assets/splash_skip.gif" alt="Skip and ESC keys" />
  </div>
</div>
```

---

## 🎨 Branding Guidelines

### Colors to Highlight
- **Neural Cyan:** `#00ffff` (cortex gradient start)
- **Neural Magenta:** `#ff00ff` (cortex gradient end)
- **Matrix Green:** `#00ff00` (boot sequence)
- **Claude Blue:** `#7aa2f7` (branding)

### Text Overlays (If Added)
```
"RyCode - Epic 3D Splash Screen"
"Real Donut Algorithm Math"
"30 FPS Smooth Animation"
"5 Hidden Easter Eggs"
"🤖 100% AI-Designed by Claude"
```

### Watermarks (Optional)
Add subtle "RyCode" or "toolkit-cli" watermark:
```bash
ffmpeg -i splash_demo.mp4 \
  -vf "drawtext=text='RyCode.com':x=10:y=10:fontsize=24:fontcolor=white@0.5" \
  splash_demo_watermarked.mp4
```

---

## 🐛 Troubleshooting

### VHS Issues

**Problem:** `vhs: command not found`
**Solution:**
```bash
# Ensure Go bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Reinstall VHS
go install github.com/charmbracelet/vhs@latest
```

**Problem:** GIF is too large (>10MB)
**Solution:**
```bash
# Reduce dimensions or duration in .tape file
Set Width 800    # Reduce from 1200
Set Height 600   # Reduce from 800
Sleep 4s         # Reduce sleep duration
```

### asciinema Issues

**Problem:** Recording is choppy
**Solution:**
```bash
# Record at lower frame rate
asciinema rec --idle-time-limit 1 splash_demo.cast
```

### FFmpeg Issues

**Problem:** `codec not found`
**Solution:**
```bash
# Reinstall FFmpeg with all codecs
brew reinstall ffmpeg
```

---

## 📚 Resources

**Tools:**
- [VHS](https://github.com/charmbracelet/vhs) - Terminal recorder
- [asciinema](https://asciinema.org/) - Terminal session recorder
- [agg](https://github.com/asciinema/agg) - asciinema to GIF converter
- [FFmpeg](https://ffmpeg.org/) - Video conversion
- [ImageMagick](https://imagemagick.org/) - Image manipulation
- [gifsicle](https://www.lcdf.org/gifsicle/) - GIF optimization

**Guides:**
- [VHS Documentation](https://github.com/charmbracelet/vhs#readme)
- [asciinema Documentation](https://docs.asciinema.org/)
- [FFmpeg Guide](https://trac.ffmpeg.org/wiki)

---

## ✅ Quick Start (Recommended)

```bash
# 1. Install VHS
brew install vhs

# 2. Navigate to project
cd /Users/aaron/Code/RyCode/RyCode/packages/tui

# 3. Generate standard splash demo
vhs splash_demo.tape

# 4. Generate donut mode demo
vhs splash_demo_donut.tape

# 5. Optimize GIFs
convert splash_demo.gif -fuzz 10% -layers Optimize -colors 128 splash_demo_optimized.gif
convert splash_demo_donut.gif -fuzz 10% -layers Optimize -colors 128 splash_demo_donut_optimized.gif

# 6. Done! Assets ready for landing page
ls -lh splash_demo*.gif
```

---

**🤖 Created by Claude AI**

*Documentation for creating professional marketing assets*
*Ready for ry-code.com landing page integration*

---

**Status:** Ready for Execution ✅
**Estimated Time:** 30 minutes (includes tool installation)
**Output:** 2-4 high-quality GIF/video assets <2MB each

