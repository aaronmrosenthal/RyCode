package responsive

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
)

// Platform provides platform-specific implementations
type Platform interface {
	Name() string
	Haptic() HapticProvider
	Voice() VoiceProvider
	Storage() StorageProvider
	Accessibility() AccessibilityProvider
}

// HapticProvider provides haptic feedback
type HapticProvider interface {
	Trigger(event HapticEvent) error
	IsAvailable() bool
}

// VoiceProvider provides voice input
type VoiceProvider interface {
	StartRecording() tea.Cmd
	StopRecording() (string, float64, error) // Returns: text, confidence, error
	IsAvailable() bool
}

// StorageProvider provides persistent storage
type StorageProvider interface {
	Save(key string, value interface{}) error
	Load(key string, target interface{}) error
	Delete(key string) error
	Clear() error
}

// AccessibilityProvider provides platform accessibility features
type AccessibilityProvider interface {
	AnnounceForScreenReader(message string)
	IsScreenReaderActive() bool
	IsHighContrastEnabled() bool
	IsReducedMotionEnabled() bool
}

// DetectPlatform detects and returns appropriate platform implementation
func DetectPlatform(caps *TerminalCapabilities) Platform {
	switch caps.Platform {
	case "ios":
		return NewIOSPlatform()
	case "android":
		return NewAndroidPlatform()
	case "macos":
		return NewMacOSPlatform()
	case "linux":
		return NewLinuxPlatform()
	case "windows":
		return NewWindowsPlatform()
	default:
		return NewGenericPlatform()
	}
}

// ====================
// iOS Platform
// ====================

type IOSPlatform struct {
	haptic        HapticProvider
	voice         VoiceProvider
	storage       StorageProvider
	accessibility AccessibilityProvider
}

func NewIOSPlatform() *IOSPlatform {
	return &IOSPlatform{
		haptic:        &IOSHaptic{},
		voice:         &IOSVoice{},
		storage:       NewFileStorage("ios"),
		accessibility: &IOSAccessibility{},
	}
}

func (p *IOSPlatform) Name() string                             { return "iOS" }
func (p *IOSPlatform) Haptic() HapticProvider                   { return p.haptic }
func (p *IOSPlatform) Voice() VoiceProvider                     { return p.voice }
func (p *IOSPlatform) Storage() StorageProvider                 { return p.storage }
func (p *IOSPlatform) Accessibility() AccessibilityProvider     { return p.accessibility }

type IOSHaptic struct{}

func (h *IOSHaptic) Trigger(event HapticEvent) error {
	// In a real iOS terminal app, this would call native APIs
	// via JavaScript bridge or URL scheme
	// For now, we can try to send system beep
	fmt.Print("\a") // BEL character
	return nil
}

func (h *IOSHaptic) IsAvailable() bool { return true }

type IOSVoice struct{}

func (h *IOSVoice) StartRecording() tea.Cmd {
	// Would integrate with iOS speech recognition
	// For now, simulate
	return func() tea.Msg {
		return VoiceStartMsg{}
	}
}

func (h *IOSVoice) StopRecording() (string, float64, error) {
	// Real implementation would use iOS Speech framework
	return "", 0, fmt.Errorf("voice not implemented for iOS")
}

func (h *IOSVoice) IsAvailable() bool { return false }

type IOSAccessibility struct{}

func (a *IOSAccessibility) AnnounceForScreenReader(message string) {
	// Would use iOS accessibility APIs
	// Could write to a special log file that Blink/Termius monitors
}

func (a *IOSAccessibility) IsScreenReaderActive() bool {
	// Check if VoiceOver is running
	return os.Getenv("VOICEOVER_RUNNING") == "1"
}

func (a *IOSAccessibility) IsHighContrastEnabled() bool {
	return os.Getenv("IOS_HIGH_CONTRAST") == "1"
}

func (a *IOSAccessibility) IsReducedMotionEnabled() bool {
	return os.Getenv("IOS_REDUCED_MOTION") == "1"
}

// ====================
// Android Platform
// ====================

type AndroidPlatform struct {
	haptic        HapticProvider
	voice         VoiceProvider
	storage       StorageProvider
	accessibility AccessibilityProvider
}

func NewAndroidPlatform() *AndroidPlatform {
	return &AndroidPlatform{
		haptic:        &AndroidHaptic{},
		voice:         &AndroidVoice{},
		storage:       NewFileStorage("android"),
		accessibility: &AndroidAccessibility{},
	}
}

func (p *AndroidPlatform) Name() string                         { return "Android" }
func (p *AndroidPlatform) Haptic() HapticProvider               { return p.haptic }
func (p *AndroidPlatform) Voice() VoiceProvider                 { return p.voice }
func (p *AndroidPlatform) Storage() StorageProvider             { return p.storage }
func (p *AndroidPlatform) Accessibility() AccessibilityProvider { return p.accessibility }

type AndroidHaptic struct{}

func (h *AndroidHaptic) Trigger(event HapticEvent) error {
	fmt.Print("\a")
	return nil
}

func (h *AndroidHaptic) IsAvailable() bool { return true }

type AndroidVoice struct{}

func (h *AndroidVoice) StartRecording() tea.Cmd {
	return func() tea.Msg {
		return VoiceStartMsg{}
	}
}

func (h *AndroidVoice) StopRecording() (string, float64, error) {
	return "", 0, fmt.Errorf("voice not implemented for Android")
}

func (h *AndroidVoice) IsAvailable() bool { return false }

type AndroidAccessibility struct{}

func (a *AndroidAccessibility) AnnounceForScreenReader(message string) {}
func (a *AndroidAccessibility) IsScreenReaderActive() bool              { return false }
func (a *AndroidAccessibility) IsHighContrastEnabled() bool             { return false }
func (a *AndroidAccessibility) IsReducedMotionEnabled() bool            { return false }

// ====================
// macOS Platform
// ====================

type MacOSPlatform struct {
	haptic        HapticProvider
	voice         VoiceProvider
	storage       StorageProvider
	accessibility AccessibilityProvider
}

func NewMacOSPlatform() *MacOSPlatform {
	return &MacOSPlatform{
		haptic:        &AudioHaptic{},
		voice:         &MacOSVoice{},
		storage:       NewFileStorage("macos"),
		accessibility: &MacOSAccessibility{},
	}
}

func (p *MacOSPlatform) Name() string                         { return "macOS" }
func (p *MacOSPlatform) Haptic() HapticProvider               { return p.haptic }
func (p *MacOSPlatform) Voice() VoiceProvider                 { return p.voice }
func (p *MacOSPlatform) Storage() StorageProvider             { return p.storage }
func (p *MacOSPlatform) Accessibility() AccessibilityProvider { return p.accessibility }

type AudioHaptic struct{}

func (h *AudioHaptic) Trigger(event HapticEvent) error {
	// Use different beep frequencies for different haptic types
	// macOS: afplay, Linux: beep, Windows: rundll32
	fmt.Print("\a")
	return nil
}

func (h *AudioHaptic) IsAvailable() bool { return true }

type MacOSVoice struct{}

func (h *MacOSVoice) StartRecording() tea.Cmd {
	// Could use macOS `say` command or NSSpeechRecognizer
	return func() tea.Msg {
		return VoiceStartMsg{}
	}
}

func (h *MacOSVoice) StopRecording() (string, float64, error) {
	return "", 0, fmt.Errorf("voice not implemented for macOS")
}

func (h *MacOSVoice) IsAvailable() bool { return false }

type MacOSAccessibility struct{}

func (a *MacOSAccessibility) AnnounceForScreenReader(message string) {
	// Use macOS `say` command
	// exec.Command("say", message).Run()
}

func (a *MacOSAccessibility) IsScreenReaderActive() bool {
	// Check if VoiceOver is running
	// Could check: defaults read com.apple.universalaccess voiceOverOnOffKey
	return false
}

func (a *MacOSAccessibility) IsHighContrastEnabled() bool {
	// Check system preferences
	return false
}

func (a *MacOSAccessibility) IsReducedMotionEnabled() bool {
	// Check system preferences
	return false
}

// ====================
// Linux Platform
// ====================

type LinuxPlatform struct {
	haptic        HapticProvider
	voice         VoiceProvider
	storage       StorageProvider
	accessibility AccessibilityProvider
}

func NewLinuxPlatform() *LinuxPlatform {
	return &LinuxPlatform{
		haptic:        &AudioHaptic{},
		voice:         &GenericVoice{},
		storage:       NewFileStorage("linux"),
		accessibility: &LinuxAccessibility{},
	}
}

func (p *LinuxPlatform) Name() string                         { return "Linux" }
func (p *LinuxPlatform) Haptic() HapticProvider               { return p.haptic }
func (p *LinuxPlatform) Voice() VoiceProvider                 { return p.voice }
func (p *LinuxPlatform) Storage() StorageProvider             { return p.storage }
func (p *LinuxPlatform) Accessibility() AccessibilityProvider { return p.accessibility }

type LinuxAccessibility struct{}

func (a *LinuxAccessibility) AnnounceForScreenReader(message string) {
	// Use espeak or speech-dispatcher
}

func (a *LinuxAccessibility) IsScreenReaderActive() bool {
	// Check if Orca is running
	return false
}

func (a *LinuxAccessibility) IsHighContrastEnabled() bool { return false }
func (a *LinuxAccessibility) IsReducedMotionEnabled() bool { return false }

// ====================
// Windows Platform
// ====================

type WindowsPlatform struct {
	haptic        HapticProvider
	voice         VoiceProvider
	storage       StorageProvider
	accessibility AccessibilityProvider
}

func NewWindowsPlatform() *WindowsPlatform {
	return &WindowsPlatform{
		haptic:        &AudioHaptic{},
		voice:         &GenericVoice{},
		storage:       NewFileStorage("windows"),
		accessibility: &WindowsAccessibility{},
	}
}

func (p *WindowsPlatform) Name() string                         { return "Windows" }
func (p *WindowsPlatform) Haptic() HapticProvider               { return p.haptic }
func (p *WindowsPlatform) Voice() VoiceProvider                 { return p.voice }
func (p *WindowsPlatform) Storage() StorageProvider             { return p.storage }
func (p *WindowsPlatform) Accessibility() AccessibilityProvider { return p.accessibility }

type WindowsAccessibility struct{}

func (a *WindowsAccessibility) AnnounceForScreenReader(message string) {
	// Use Windows SAPI
}

func (a *WindowsAccessibility) IsScreenReaderActive() bool {
	// Check if Narrator or JAWS is running
	return false
}

func (a *WindowsAccessibility) IsHighContrastEnabled() bool { return false }
func (a *WindowsAccessibility) IsReducedMotionEnabled() bool { return false }

// ====================
// Generic Platform (Fallback)
// ====================

type GenericPlatform struct {
	haptic        HapticProvider
	voice         VoiceProvider
	storage       StorageProvider
	accessibility AccessibilityProvider
}

func NewGenericPlatform() *GenericPlatform {
	return &GenericPlatform{
		haptic:        &VisualHapticProvider{},
		voice:         &GenericVoice{},
		storage:       NewFileStorage("generic"),
		accessibility: &GenericAccessibility{},
	}
}

func (p *GenericPlatform) Name() string                         { return "Generic" }
func (p *GenericPlatform) Haptic() HapticProvider               { return p.haptic }
func (p *GenericPlatform) Voice() VoiceProvider                 { return p.voice }
func (p *GenericPlatform) Storage() StorageProvider             { return p.storage }
func (p *GenericPlatform) Accessibility() AccessibilityProvider { return p.accessibility }

type VisualHapticProvider struct{}

func (h *VisualHapticProvider) Trigger(event HapticEvent) error {
	// Visual only (emojis from our original implementation)
	return nil
}

func (h *VisualHapticProvider) IsAvailable() bool { return true }

type GenericVoice struct{}

func (h *GenericVoice) StartRecording() tea.Cmd {
	return func() tea.Msg {
		return VoiceStartMsg{}
	}
}

func (h *GenericVoice) StopRecording() (string, float64, error) {
	return "", 0, fmt.Errorf("voice not available on this platform")
}

func (h *GenericVoice) IsAvailable() bool { return false }

type GenericAccessibility struct{}

func (a *GenericAccessibility) AnnounceForScreenReader(message string) {}
func (a *GenericAccessibility) IsScreenReaderActive() bool              { return false }
func (a *GenericAccessibility) IsHighContrastEnabled() bool             { return false }
func (a *GenericAccessibility) IsReducedMotionEnabled() bool            { return false }

// ====================
// File Storage Implementation
// ====================

type FileStorage struct {
	basePath string
}

func NewFileStorage(platform string) *FileStorage {
	// Determine config directory based on platform
	var basePath string

	switch platform {
	case "ios":
		basePath = filepath.Join(os.Getenv("HOME"), "Documents", ".opencode")
	case "android":
		basePath = filepath.Join(os.Getenv("HOME"), ".opencode")
	case "macos":
		basePath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "opencode")
	case "linux":
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			basePath = filepath.Join(xdgConfig, "opencode")
		} else {
			basePath = filepath.Join(os.Getenv("HOME"), ".config", "opencode")
		}
	case "windows":
		basePath = filepath.Join(os.Getenv("APPDATA"), "opencode")
	default:
		basePath = filepath.Join(os.Getenv("HOME"), ".opencode")
	}

	// Ensure directory exists
	os.MkdirAll(basePath, 0755)

	return &FileStorage{basePath: basePath}
}

func (fs *FileStorage) Save(key string, value interface{}) error {
	filePath := filepath.Join(fs.basePath, key+".json")

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (fs *FileStorage) Load(key string, target interface{}) error {
	filePath := filepath.Join(fs.basePath, key+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}

func (fs *FileStorage) Delete(key string) error {
	filePath := filepath.Join(fs.basePath, key+".json")
	return os.Remove(filePath)
}

func (fs *FileStorage) Clear() error {
	return os.RemoveAll(fs.basePath)
}

// ====================
// Preferences Management
// ====================

type Preferences struct {
	Accessibility AccessibilityConfig
	LastAI        string
	VoiceHistory  []string
	Theme         string
	UpdatedAt     time.Time
}

func LoadPreferences(storage StorageProvider) (*Preferences, error) {
	prefs := &Preferences{}
	err := storage.Load("preferences", prefs)
	if err != nil {
		// Return defaults if not found
		return &Preferences{
			Accessibility: *NewAccessibilityConfig(),
			LastAI:        "claude",
			Theme:         "default",
			UpdatedAt:     time.Now(),
		}, nil
	}
	return prefs, nil
}

func SavePreferences(storage StorageProvider, prefs *Preferences) error {
	prefs.UpdatedAt = time.Now()
	return storage.Save("preferences", prefs)
}
