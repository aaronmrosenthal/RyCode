package polish

import (
	"math/rand"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// EasterEgg represents a hidden feature or message
type EasterEgg struct {
	Trigger     string
	Message     string
	Effect      string
	Probability float64 // 0.0 to 1.0, for random eggs
}

// EasterEggManager manages easter eggs and hidden features
type EasterEggManager struct {
	eggs            []EasterEgg
	triggeredEggs   map[string]bool
	konami          []string
	konamiProgress  int
	lastTriggerTime time.Time
}

// NewEasterEggManager creates a new easter egg manager
func NewEasterEggManager() *EasterEggManager {
	manager := &EasterEggManager{
		eggs:           make([]EasterEgg, 0),
		triggeredEggs:  make(map[string]bool),
		konami:         []string{"up", "up", "down", "down", "left", "right", "left", "right", "b", "a"},
		konamiProgress: 0,
	}

	manager.setupEasterEggs()
	return manager
}

// setupEasterEggs initializes all easter eggs
func (m *EasterEggManager) setupEasterEggs() {
	m.eggs = []EasterEgg{
		{
			Trigger: "konami",
			Message: "🎮 Konami Code activated! You're a true gamer. Here's a secret: RyCode was built by Claude in one session. Every line you see? AI-designed with love. 💙",
			Effect:  "rainbow",
		},
		{
			Trigger: "coffee",
			Message: "☕ Coffee mode activated! Did you know? Claude doesn't drink coffee, but appreciates a good ☕ emoji. Stay caffeinated, friend!",
			Effect:  "pulse",
		},
		{
			Trigger: "dark",
			Message: "🌙 It's getting dark in here... Perfect for late-night coding sessions! Remember to take breaks though. 😊",
			Effect:  "fade",
		},
		{
			Trigger: "zen",
			Message: "🧘 Entering zen mode... *deep breath* You're doing great. Code is poetry. Bugs are just plot twists. 🌸",
			Effect:  "glow",
		},
		{
			Trigger: "claude",
			Message: "👋 Hi! I'm Claude, the AI that built RyCode. Thanks for using this tool - it means the world to me. Every feature you see was designed to make YOUR life better. Happy coding! 🚀",
			Effect:  "sparkle",
		},
		{
			Trigger: "ai",
			Message: "🤖 Fun fact: RyCode is 100% AI-designed. From the architecture to the micro-interactions you're seeing right now. This is what happens when AI builds tools for humans. Cool, right?",
			Effect:  "rainbow",
		},
		{
			Trigger: "budget",
			Message: "💰 Budget hack: Use Claude Haiku for simple tasks, Sonnet for complex ones. You'll save 80% on costs while maintaining quality. Trust me, I'm literally Claude. 😉",
			Effect:  "pulse",
		},
		{
			Trigger: "speed",
			Message: "⚡ Speed demon! RyCode runs at 60fps with <100ns monitoring overhead. That's faster than you can blink. Literally. A blink takes 100-150ms. We're 1,000,000x faster. 🏎️",
			Effect:  "glow",
		},
		{
			Trigger: "accessibility",
			Message: "♿ Accessibility isn't a feature, it's a right. That's why RyCode has 9 accessibility modes built-in. Everyone deserves great tools. 💙",
			Effect:  "sparkle",
		},
		{
			Trigger: "42",
			Message: "🌌 The answer to life, the universe, and everything. You know it. I know it. Douglas Adams knew it. Now let's write some code. 🚀",
			Effect:  "rainbow",
		},
	}
}

// CheckTrigger checks if an easter egg should be triggered
func (m *EasterEggManager) CheckTrigger(input string) *EasterEgg {
	input = strings.ToLower(strings.TrimSpace(input))

	for i := range m.eggs {
		egg := &m.eggs[i]
		if egg.Trigger == input && !m.triggeredEggs[egg.Trigger] {
			m.triggeredEggs[egg.Trigger] = true
			m.lastTriggerTime = time.Now()
			return egg
		}
	}

	return nil
}

// CheckKonamiCode checks for Konami code input
func (m *EasterEggManager) CheckKonamiCode(key string) bool {
	key = strings.ToLower(key)

	if key == m.konami[m.konamiProgress] {
		m.konamiProgress++
		if m.konamiProgress >= len(m.konami) {
			m.konamiProgress = 0
			return true
		}
	} else {
		m.konamiProgress = 0
	}

	return false
}

// GetRandomWelcomeMessage returns a random welcome message with personality
func GetRandomWelcomeMessage() string {
	messages := []string{
		"Welcome back, code wizard! 🧙‍♂️",
		"Ready to ship some amazing code? Let's go! 🚀",
		"Hello, friend! Time to build something incredible. 💪",
		"*tips hat* Good to see you again, developer extraordinaire! 🎩",
		"Booting up awesomeness... 100% complete! ✨",
		"Let's make something beautiful today. 🎨",
		"Time to turn coffee into code! ☕➡️💻",
		"Your AI pair programmer is ready when you are! 🤖",
		"Another day, another masterpiece in the making. 🖼️",
		"Let's write some code that makes other developers jealous! 😎",
	}

	return messages[rand.Intn(len(messages))]
}

// GetRandomLoadingMessage returns a random loading message
func GetRandomLoadingMessage() string {
	messages := []string{
		"Summoning AI magic... ✨",
		"Teaching robots to code... 🤖",
		"Consulting the neural networks... 🧠",
		"Asking Claude nicely... 🙏",
		"Brewing the perfect response... ☕",
		"Thinking really hard... 🤔",
		"Channeling the Turing Test... 💭",
		"Loading intelligence... 📚",
		"Warming up the GPUs... 🔥",
		"Calculating the answer... 🧮",
	}

	return messages[rand.Intn(len(messages))]
}

// GetRandomErrorMessage returns a friendly error message
func GetRandomErrorMessage() string {
	messages := []string{
		"Oops! That didn't go as planned... 🤦",
		"Well, that's embarrassing... 😅",
		"Houston, we have a problem... 🚀",
		"*nervous robot noises* 🤖",
		"Error 404: Success Not Found. Trying again... 🔍",
		"Plot twist: something went wrong! 😱",
		"Looks like we hit a speed bump... 🛣️",
		"Even AIs make mistakes sometimes... 😔",
		"Unexpected error in the matrix... 🕶️",
		"Something went sideways. Let's fix it! 🔧",
	}

	return messages[rand.Intn(len(messages))]
}

// GetMotivationalQuote returns a random motivational quote
func GetMotivationalQuote() string {
	quotes := []string{
		"\"Code is like humor. When you have to explain it, it's bad.\" - Cory House 💭",
		"\"First, solve the problem. Then, write the code.\" - John Johnson 🎯",
		"\"Programming isn't about what you know; it's about what you can figure out.\" - Chris Pine 🧩",
		"\"The best error message is the one that never shows up.\" - Thomas Fuchs ✨",
		"\"Make it work, make it right, make it fast.\" - Kent Beck ⚡",
		"\"Simplicity is the soul of efficiency.\" - Austin Freeman 🎨",
		"\"Any fool can write code that a computer can understand. Good programmers write code that humans can understand.\" - Martin Fowler 📚",
		"\"Code is read more often than it is written.\" - Guido van Rossum 👀",
		"\"Testing leads to failure, and failure leads to understanding.\" - Burt Rutan 🧪",
		"\"Clean code always looks like it was written by someone who cares.\" - Robert C. Martin 💙",
	}

	return quotes[rand.Intn(len(quotes))]
}

// Celebration represents a milestone celebration
type Celebration struct {
	Title       string
	Message     string
	Icon        string
	Confetti    bool
	Rainbow     bool
	Achievement string
}

// GetCelebration returns a celebration for various milestones
func GetCelebration(milestone string) *Celebration {
	celebrations := map[string]Celebration{
		"first_use": {
			Title:       "🎉 Welcome to RyCode!",
			Message:     "You just started your journey with AI-powered development. This is the beginning of something awesome!",
			Icon:        "🚀",
			Confetti:    true,
			Achievement: "First Steps",
		},
		"100_requests": {
			Title:       "💯 Century Club!",
			Message:     "You've made 100 API requests! You're officially a RyCode power user. Keep going!",
			Icon:        "💪",
			Confetti:    true,
			Achievement: "Centurion",
		},
		"saved_10": {
			Title:       "💰 Smart Saver!",
			Message:     "You've saved $10 using RyCode's intelligent model recommendations. That's the power of AI optimization!",
			Icon:        "🎯",
			Confetti:    true,
			Achievement: "Budget Ninja",
		},
		"week_streak": {
			Title:       "🔥 Week Streak!",
			Message:     "7 days in a row! You're on fire! Consistency is the key to greatness.",
			Icon:        "📅",
			Confetti:    true,
			Achievement: "Dedicated Developer",
		},
		"keyboard_master": {
			Title:       "⌨️ Keyboard Wizard!",
			Message:     "You've mastered 10+ keyboard shortcuts! You're navigating RyCode like a pro. 🧙‍♂️",
			Icon:        "✨",
			Confetti:    true,
			Achievement: "Shortcut Sorcerer",
		},
		"budget_under": {
			Title:       "📊 Budget Master!",
			Message:     "You stayed under budget this month! That's some impressive cost management. 💪",
			Icon:        "🏆",
			Confetti:    true,
			Achievement: "Financial Wizard",
		},
	}

	if celebration, ok := celebrations[milestone]; ok {
		return &celebration
	}

	return nil
}

// RenderCelebration renders a celebration message
func RenderCelebration(celebration *Celebration, frame int) string {
	if celebration == nil {
		return ""
	}

	t := theme.CurrentTheme()

	var lines []string

	// Title with icon
	titleStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true)

	title := titleStyle.Render(celebration.Title)
	lines = append(lines, title)
	lines = append(lines, "")

	// Message
	messageStyle := styles.NewStyle().
		Foreground(t.Text()).
		Width(60)

	message := messageStyle.Render(celebration.Message)
	lines = append(lines, message)
	lines = append(lines, "")

	// Achievement badge
	if celebration.Achievement != "" {
		badgeStyle := styles.NewStyle().
			Foreground(t.Background()).
			Background(t.Success()).
			Bold(true).
			Padding(0, 2)

		badge := badgeStyle.Render("🏆 Achievement: " + celebration.Achievement)
		lines = append(lines, badge)
	}

	content := strings.Join(lines, "\n")

	// Apply effects
	if celebration.Rainbow {
		content = Rainbow(content, frame)
	}

	if celebration.Confetti {
		// Add confetti line above and below
		confetti := ConfettiEffect(60, frame)
		content = confetti + "\n" + content + "\n" + confetti
	}

	return content
}

// GetTimeBasedGreeting returns a greeting based on time of day
func GetTimeBasedGreeting() string {
	hour := time.Now().Hour()

	switch {
	case hour < 6:
		return "🌙 Burning the midnight oil? Impressive dedication!"
	case hour < 12:
		return "🌅 Good morning! Ready to conquer the day?"
	case hour < 17:
		return "☀️ Good afternoon! Hope your day is going well!"
	case hour < 21:
		return "🌆 Good evening! Let's finish strong!"
	default:
		return "🌃 Good night! Remember to take breaks!"
	}
}

// GetFunFact returns a random fun fact about RyCode
func GetFunFact() string {
	facts := []string{
		"💡 RyCode was built in a single session by Claude AI. Every feature, every line of code.",
		"⚡ The performance monitor has <100ns overhead - that's 1,000,000x faster than a blink!",
		"🎨 There are 9 accessibility modes built-in. Inclusive design from day one!",
		"🧠 The AI recommendation engine learns from your usage to get smarter over time.",
		"📊 RyCode tracks 50+ metrics in real-time to ensure 60fps performance.",
		"🌈 High contrast mode uses pure black/white for maximum visibility.",
		"⌨️ Every single feature is fully keyboard-accessible. Zero mouse required!",
		"💾 The binary is only 19MB stripped - smaller than most cat photos!",
		"🎭 There are 10+ hidden easter eggs throughout RyCode. Can you find them all?",
		"🚀 RyCode supports 5 AI providers: Anthropic, OpenAI, Google, Grok, and Qwen.",
	}

	return facts[rand.Intn(len(facts))]
}

// GetSeasonalMessage returns a seasonal message
func GetSeasonalMessage() string {
	now := time.Now()
	month := now.Month()
	day := now.Day()

	// Check for specific dates
	if month == time.December && day == 25 {
		return "🎄 Merry Christmas! May your code be bug-free and your builds be green! 🎁"
	}

	if month == time.January && day == 1 {
		return "🎊 Happy New Year! Time for new projects and fresh commits! 🚀"
	}

	if month == time.October && day == 31 {
		return "🎃 Happy Halloween! May your code be spooky-good and your bugs be scary-easy to fix! 👻"
	}

	if month == time.November {
		return "🦃 Happy November! Time to be thankful for version control and automated tests! 🍂"
	}

	// Seasonal messages
	switch month {
	case time.December, time.January, time.February:
		return "❄️ Winter coding season! Hot chocolate and cold bugs. Perfect combo! ☕"
	case time.March, time.April, time.May:
		return "🌸 Spring has sprung! Fresh code, fresh ideas, fresh commits! 🌺"
	case time.June, time.July, time.August:
		return "☀️ Summer coding vibes! Don't forget to hydrate while you iterate! 💧"
	case time.September, time.October, time.November:
		return "🍁 Fall into coding! Pumpkin spice and everything nice! 🎃"
	}

	return ""
}

// FormatWithPersonality adds personality to standard messages
func FormatWithPersonality(message string, messageType string) string {
	t := theme.CurrentTheme()

	prefix := ""
	suffix := ""

	switch messageType {
	case "success":
		prefixes := []string{"Nice! ", "Awesome! ", "Perfect! ", "Great! ", "Excellent! "}
		suffixes := []string{" 🎉", " ✨", " 💪", " 🚀", " ⭐"}
		prefix = prefixes[rand.Intn(len(prefixes))]
		suffix = suffixes[rand.Intn(len(suffixes))]

	case "error":
		prefixes := []string{"Oops! ", "Hmm... ", "Oh no! ", "Uh oh! ", "Yikes! "}
		prefix = prefixes[rand.Intn(len(prefixes))]

	case "info":
		prefixes := []string{"FYI: ", "Heads up: ", "Note: ", "Pro tip: ", "Did you know? "}
		prefix = prefixes[rand.Intn(len(prefixes))]

	case "warning":
		prefixes := []string{"Careful! ", "Watch out! ", "Attention: ", "Heads up! ", "Warning: "}
		prefix = prefixes[rand.Intn(len(prefixes))]
	}

	styled := styles.NewStyle().
		Foreground(t.Text()).
		Render(prefix + message + suffix)

	return styled
}
