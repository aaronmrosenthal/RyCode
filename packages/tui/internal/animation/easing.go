package animation

import (
	"math"
	"time"
)

// EasingFunction defines a function that maps progress [0,1] to eased value [0,1]
type EasingFunction func(float64) float64

// Common easing functions for smooth animations
var (
	// Linear - no easing
	Linear EasingFunction = func(t float64) float64 {
		return t
	}

	// EaseInQuad - accelerating from zero velocity
	EaseInQuad EasingFunction = func(t float64) float64 {
		return t * t
	}

	// EaseOutQuad - decelerating to zero velocity
	EaseOutQuad EasingFunction = func(t float64) float64 {
		return t * (2 - t)
	}

	// EaseInOutQuad - acceleration until halfway, then deceleration
	EaseInOutQuad EasingFunction = func(t float64) float64 {
		if t < 0.5 {
			return 2 * t * t
		}
		return -1 + (4-2*t)*t
	}

	// EaseInCubic - accelerating from zero velocity (cubic)
	EaseInCubic EasingFunction = func(t float64) float64 {
		return t * t * t
	}

	// EaseOutCubic - decelerating to zero velocity (cubic)
	EaseOutCubic EasingFunction = func(t float64) float64 {
		t--
		return t*t*t + 1
	}

	// EaseInOutCubic - acceleration until halfway, then deceleration (cubic)
	EaseInOutCubic EasingFunction = func(t float64) float64 {
		if t < 0.5 {
			return 4 * t * t * t
		}
		t = 2*t - 2
		return (t*t*t + 2) / 2
	}

	// EaseInQuart - accelerating from zero velocity (quartic)
	EaseInQuart EasingFunction = func(t float64) float64 {
		return t * t * t * t
	}

	// EaseOutQuart - decelerating to zero velocity (quartic)
	EaseOutQuart EasingFunction = func(t float64) float64 {
		t--
		return 1 - t*t*t*t
	}

	// EaseInOutQuart - acceleration until halfway, then deceleration (quartic)
	EaseInOutQuart EasingFunction = func(t float64) float64 {
		if t < 0.5 {
			return 8 * t * t * t * t
		}
		t--
		return 1 - 8*t*t*t*t
	}

	// EaseInExpo - exponential accelerating from zero velocity
	EaseInExpo EasingFunction = func(t float64) float64 {
		if t == 0 {
			return 0
		}
		return math.Pow(2, 10*(t-1))
	}

	// EaseOutExpo - exponential decelerating to zero velocity
	EaseOutExpo EasingFunction = func(t float64) float64 {
		if t == 1 {
			return 1
		}
		return 1 - math.Pow(2, -10*t)
	}

	// EaseInOutExpo - exponential acceleration until halfway, then deceleration
	EaseInOutExpo EasingFunction = func(t float64) float64 {
		if t == 0 || t == 1 {
			return t
		}
		if t < 0.5 {
			return math.Pow(2, 20*t-10) / 2
		}
		return (2 - math.Pow(2, -20*t+10)) / 2
	}

	// EaseOutBack - overshoots target and comes back
	EaseOutBack EasingFunction = func(t float64) float64 {
		const c1 = 1.70158
		const c3 = c1 + 1
		return 1 + c3*math.Pow(t-1, 3) + c1*math.Pow(t-1, 2)
	}

	// EaseInOutBack - overshoots on both ends
	EaseInOutBack EasingFunction = func(t float64) float64 {
		const c1 = 1.70158
		const c2 = c1 * 1.525
		if t < 0.5 {
			return (math.Pow(2*t, 2) * ((c2+1)*2*t - c2)) / 2
		}
		return (math.Pow(2*t-2, 2)*((c2+1)*(t*2-2)+c2) + 2) / 2
	}

	// EaseOutElastic - elastic snap back (like a spring)
	EaseOutElastic EasingFunction = func(t float64) float64 {
		const c4 = (2 * math.Pi) / 3
		if t == 0 || t == 1 {
			return t
		}
		return math.Pow(2, -10*t)*math.Sin((t*10-0.75)*c4) + 1
	}

	// EaseInBounce - bounce at start
	EaseInBounce EasingFunction = func(t float64) float64 {
		return 1 - EaseOutBounce(1-t)
	}

	// EaseOutBounce - bounce at end
	EaseOutBounce EasingFunction = func(t float64) float64 {
		const n1 = 7.5625
		const d1 = 2.75
		if t < 1/d1 {
			return n1 * t * t
		} else if t < 2/d1 {
			t -= 1.5 / d1
			return n1*t*t + 0.75
		} else if t < 2.5/d1 {
			t -= 2.25 / d1
			return n1*t*t + 0.9375
		}
		t -= 2.625 / d1
		return n1*t*t + 0.984375
	}
)

// Animation represents an ongoing animation
type Animation struct {
	start      time.Time
	duration   time.Duration
	easing     EasingFunction
	from       float64
	to         float64
	onUpdate   func(float64)
	onComplete func()
	running    bool
}

// NewAnimation creates a new animation
func NewAnimation(duration time.Duration, from, to float64, easing EasingFunction) *Animation {
	return &Animation{
		duration: duration,
		easing:   easing,
		from:     from,
		to:       to,
		running:  false,
	}
}

// Start begins the animation
func (a *Animation) Start() {
	a.start = time.Now()
	a.running = true
}

// Stop stops the animation
func (a *Animation) Stop() {
	a.running = false
}

// IsRunning returns whether the animation is currently running
func (a *Animation) IsRunning() bool {
	return a.running
}

// OnUpdate sets a callback for each animation frame
func (a *Animation) OnUpdate(fn func(float64)) *Animation {
	a.onUpdate = fn
	return a
}

// OnComplete sets a callback for when animation finishes
func (a *Animation) OnComplete(fn func()) *Animation {
	a.onComplete = fn
	return a
}

// Update updates the animation and returns the current value
func (a *Animation) Update() float64 {
	if !a.running {
		return a.to
	}

	elapsed := time.Since(a.start)
	if elapsed >= a.duration {
		a.running = false
		if a.onComplete != nil {
			a.onComplete()
		}
		return a.to
	}

	// Calculate progress [0,1]
	progress := float64(elapsed) / float64(a.duration)

	// Apply easing
	eased := a.easing(progress)

	// Interpolate between from and to
	value := a.from + (a.to-a.from)*eased

	if a.onUpdate != nil {
		a.onUpdate(value)
	}

	return value
}

// Value returns the current animation value
func (a *Animation) Value() float64 {
	return a.Update()
}

// Progress returns the current progress [0,1]
func (a *Animation) Progress() float64 {
	if !a.running {
		return 1.0
	}

	elapsed := time.Since(a.start)
	if elapsed >= a.duration {
		return 1.0
	}

	return float64(elapsed) / float64(a.duration)
}

// Animator manages multiple animations
type Animator struct {
	animations []*Animation
}

// NewAnimator creates a new animator
func NewAnimator() *Animator {
	return &Animator{
		animations: make([]*Animation, 0),
	}
}

// Add adds an animation to the animator
func (a *Animator) Add(anim *Animation) {
	a.animations = append(a.animations, anim)
}

// Update updates all animations
func (a *Animator) Update() {
	// Update all animations
	active := make([]*Animation, 0, len(a.animations))
	for _, anim := range a.animations {
		anim.Update()
		if anim.IsRunning() {
			active = append(active, anim)
		}
	}
	a.animations = active
}

// HasRunning returns true if any animations are running
func (a *Animator) HasRunning() bool {
	for _, anim := range a.animations {
		if anim.IsRunning() {
			return true
		}
	}
	return false
}

// StopAll stops all animations
func (a *Animator) StopAll() {
	for _, anim := range a.animations {
		anim.Stop()
	}
	a.animations = make([]*Animation, 0)
}

// ColorInterpolate interpolates between two colors (hex format)
// Returns a hex color string at the given progress [0,1]
func ColorInterpolate(from, to string, progress float64) string {
	fromR, fromG, fromB := hexToRGB(from)
	toR, toG, toB := hexToRGB(to)

	r := interpolateValue(fromR, toR, progress)
	g := interpolateValue(fromG, toG, progress)
	b := interpolateValue(fromB, toB, progress)

	return rgbToHex(r, g, b)
}

// Helper functions
func hexToRGB(hex string) (r, g, b uint8) {
	// Remove # if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	// Parse hex - simple manual parsing to avoid fmt import issues
	if len(hex) >= 6 {
		r = parseHexByte(hex[0:2])
		g = parseHexByte(hex[2:4])
		b = parseHexByte(hex[4:6])
	}
	return
}

func parseHexByte(s string) uint8 {
	var result uint8
	for _, c := range s {
		result *= 16
		if c >= '0' && c <= '9' {
			result += uint8(c - '0')
		} else if c >= 'A' && c <= 'F' {
			result += uint8(c-'A') + 10
		} else if c >= 'a' && c <= 'f' {
			result += uint8(c-'a') + 10
		}
	}
	return result
}

func rgbToHex(r, g, b uint8) string {
	// Manual hex formatting
	hex := "#"
	hex += byteToHex(r)
	hex += byteToHex(g)
	hex += byteToHex(b)
	return hex
}

func byteToHex(b uint8) string {
	const hexChars = "0123456789ABCDEF"
	return string([]byte{hexChars[b>>4], hexChars[b&0xF]})
}

func interpolateValue(from, to uint8, progress float64) uint8 {
	diff := int(to) - int(from)
	return uint8(int(from) + int(float64(diff)*progress))
}
