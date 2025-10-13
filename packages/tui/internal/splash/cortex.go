package splash

import (
	"math"
	"strings"
)

// CortexRenderer renders a 3D rotating torus (neural cortex) in ASCII
type CortexRenderer struct {
	width       int       // Screen width
	height      int       // Screen height
	A           float64   // Rotation angle around X-axis
	B           float64   // Rotation angle around Z-axis
	screen      []rune    // Character buffer
	zbuffer     []float64 // Depth buffer for z-buffering
	chars       []rune    // Character set for luminance mapping
	rainbowMode bool      // Easter egg: rainbow colors
	brandColor  *RGB      // Optional brand color for provider-themed rendering
}

// NewCortexRenderer creates a new cortex renderer
func NewCortexRenderer(width, height int) *CortexRenderer {
	size := width * height
	return &CortexRenderer{
		width:   width,
		height:  height,
		A:       0.0,
		B:       0.0,
		screen:  make([]rune, size),
		zbuffer: make([]float64, size),
		chars:   []rune{' ', '.', '·', ':', '*', '◉', '◎', '⚡'},
	}
}

// RenderFrame renders a single frame of the torus animation
func (r *CortexRenderer) RenderFrame() {
	// Clear buffers
	for i := range r.screen {
		r.screen[i] = ' '
		r.zbuffer[i] = 0
	}

	// Precompute rotation matrix elements
	sinA, cosA := math.Sin(r.A), math.Cos(r.A)
	sinB, cosB := math.Sin(r.B), math.Cos(r.B)

	// Render torus surface
	const thetaStep = 0.07 // ~90 steps around torus
	const phiStep = 0.02   // ~314 steps around tube

	for theta := 0.0; theta < 2*math.Pi; theta += thetaStep {
		sinTheta, cosTheta := math.Sin(theta), math.Cos(theta)

		for phi := 0.0; phi < 2*math.Pi; phi += phiStep {
			sinPhi, cosPhi := math.Sin(phi), math.Cos(phi)

			// Torus geometry (R=2, r=1)
			const majorRadius = 2.0 // Major radius
			const minorRadius = 1.0 // Minor radius

			circleX := majorRadius + minorRadius*cosPhi
			circleY := minorRadius * sinPhi

			// Apply rotations (Rx then Rz)
			x := circleX*(cosB*cosTheta+sinA*sinB*sinTheta) - circleY*cosA*sinB
			y := circleX*(sinB*cosTheta-sinA*cosB*sinTheta) + circleY*cosA*cosB
			z := 5.0 + cosA*circleX*sinTheta + circleY*sinA // z=5 pushes away from camera

			// Perspective projection
			ooz := 1.0 / z // "one over z"
			xp := int(float64(r.width)*0.5 + 30.0*ooz*x)
			yp := int(float64(r.height)*0.5 - 15.0*ooz*y)

			// Bounds check
			if xp < 0 || xp >= r.width || yp < 0 || yp >= r.height {
				continue
			}

			// Calculate luminance (Phong-style shading)
			L := cosPhi*cosTheta*sinB - cosA*cosTheta*sinPhi - sinA*sinTheta +
				cosB*(cosA*sinPhi-cosTheta*sinA*sinTheta)

			// Z-buffer test
			idx := yp*r.width + xp
			if ooz > r.zbuffer[idx] {
				r.zbuffer[idx] = ooz

				// Map luminance to character (8 levels)
				luminanceIdx := int((L + 1.0) * 3.5) // Map [-1,1] to [0,7]
				if luminanceIdx < 0 {
					luminanceIdx = 0
				}
				if luminanceIdx > 7 {
					luminanceIdx = 7
				}

				r.screen[idx] = r.chars[luminanceIdx]
			}
		}
	}

	// Update rotation angles
	r.A += 0.04 // Rotate around X-axis
	r.B += 0.02 // Rotate around Z-axis

	// Easter egg: hide secret message occasionally
	r.renderSecretMessage()
}

// renderSecretMessage occasionally reveals a hidden message
func (r *CortexRenderer) renderSecretMessage() {
	// Only show every 300 frames (~10 seconds at 30 FPS)
	// and only for 30 frames (1 second)
	frameNumber := int(r.A*100) % 600 // Pseudo-frame based on rotation
	if frameNumber < 300 || frameNumber >= 330 {
		return
	}

	message := "CLAUDE WAS HERE"
	startX := r.width/2 - len(message)/2
	centerY := r.height / 2

	if centerY < 0 || centerY >= r.height {
		return
	}

	for i, char := range message {
		x := startX + i
		if x >= 0 && x < r.width {
			idx := centerY*r.width + x
			if idx >= 0 && idx < len(r.screen) {
				r.screen[idx] = char
				r.zbuffer[idx] = 999.0 // Always on top
			}
		}
	}
}

// SetRainbowMode enables or disables rainbow mode
func (r *CortexRenderer) SetRainbowMode(enabled bool) {
	r.rainbowMode = enabled
}

// SetBrandColor sets a custom brand color for provider-themed rendering
func (r *CortexRenderer) SetBrandColor(color RGB) {
	r.brandColor = &color
}

// ClearBrandColor removes the custom brand color and returns to default gradient
func (r *CortexRenderer) ClearBrandColor() {
	r.brandColor = nil
}

// Render renders the torus with colors and returns the string
func (r *CortexRenderer) Render() string {
	r.RenderFrame()

	var buf strings.Builder
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			idx := y*r.width + x
			char := r.screen[idx]

			if char != ' ' {
				var color RGB
				if r.rainbowMode {
					// Rainbow mode: cycle through all colors
					angle := math.Atan2(float64(y-r.height/2), float64(x-r.width/2))
					color = RainbowColor(angle + r.B + r.A)
				} else if r.brandColor != nil {
					// Brand mode: gradient from brand color to brighter version
					angle := math.Atan2(float64(y-r.height/2), float64(x-r.width/2))
					// Normalize angle to [0, 1]
					t := math.Mod(angle+r.B, 2*math.Pi) / (2 * math.Pi)
					if t < 0 {
						t += 1.0
					}
					// Create gradient from brand color to brighter version
					bright := RGB{
						R: uint8(math.Min(255, float64(r.brandColor.R)*1.5)),
						G: uint8(math.Min(255, float64(r.brandColor.G)*1.5)),
						B: uint8(math.Min(255, float64(r.brandColor.B)*1.5)),
					}
					color = lerpRGB(*r.brandColor, bright, t)
				} else {
					// Normal mode: cyan-magenta gradient
					angle := math.Atan2(float64(y-r.height/2), float64(x-r.width/2))
					color = GradientColor(angle + r.B) // Rotate gradient with torus
				}
				buf.WriteString(Colorize(string(char), color))
			} else {
				buf.WriteRune(' ')
			}
		}
		if y < r.height-1 {
			buf.WriteRune('\n')
		}
	}
	return buf.String()
}

// String returns the rendered frame as a string (without colors)
func (r *CortexRenderer) String() string {
	var buf strings.Builder
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			buf.WriteRune(r.screen[y*r.width+x])
		}
		if y < r.height-1 {
			buf.WriteRune('\n')
		}
	}
	return buf.String()
}
