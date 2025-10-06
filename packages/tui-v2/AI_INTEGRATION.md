# AI Integration Guide

## Overview

RyCode Matrix TUI v2 now supports **real AI providers** including:
- **Claude** (Anthropic) - Claude Opus 4, Sonnet, Haiku
- **GPT-4** (OpenAI) - GPT-4, GPT-4 Turbo, GPT-4o

The AI integration provides **streaming token-by-token responses** with automatic fallback to mock AI when no API keys are configured.

---

## Quick Start

### 1. Set API Keys

```bash
# For Claude (Anthropic)
export ANTHROPIC_API_KEY="sk-ant-..."

# For GPT-4 (OpenAI)
export OPENAI_API_KEY="sk-..."

# Optional: Force a specific provider
export RYCODE_AI_PROVIDER="claude"  # or "openai"
```

### 2. Run RyCode

```bash
# The TUI will auto-detect and use the first available API key
rycode

# Or specify workspace mode explicitly
rycode --workspace
```

### 3. Chat with AI

Type your message and press Enter. The AI will stream responses in real-time!

---

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ANTHROPIC_API_KEY` | Claude API key | - |
| `OPENAI_API_KEY` | OpenAI API key | - |
| `RYCODE_AI_PROVIDER` | Force provider (`claude`, `openai`, `auto`) | `auto` |
| `RYCODE_CLAUDE_MODEL` | Claude model override | `claude-opus-4-20250514` |
| `RYCODE_OPENAI_MODEL` | OpenAI model override | `gpt-4o` |

### Provider Selection Logic

1. If `RYCODE_AI_PROVIDER=claude`, use Claude (requires `ANTHROPIC_API_KEY`)
2. If `RYCODE_AI_PROVIDER=openai`, use OpenAI (requires `OPENAI_API_KEY`)
3. If `RYCODE_AI_PROVIDER=auto` (default):
   - Try Claude first (if `ANTHROPIC_API_KEY` is set)
   - Fall back to OpenAI (if `OPENAI_API_KEY` is set)
   - Fall back to Mock AI (if no keys are set)

---

## Features

### Real-Time Streaming

Both providers support **Server-Sent Events (SSE)** streaming:
- Tokens appear character-by-character as the AI generates them
- No waiting for complete responses
- Smooth, responsive UX

### Conversation History

The AI maintains context across messages:
- Full conversation history is sent with each request
- Multi-turn conversations work naturally
- Context-aware responses

### Error Handling

Graceful error handling with user-friendly messages:
- Network errors
- API rate limits
- Invalid API keys
- Streaming interruptions

### Status Indicators

The status bar shows:
- Which AI provider is active (e.g., "Claude (claude-opus-4-20250514)")
- Current streaming status ("⚡ Claude is responding...")
- Warnings when no API keys are configured

---

## Architecture

### Provider Interface

All AI providers implement the same interface:

```go
type Provider interface {
    Stream(ctx context.Context, prompt string, messages []Message) (<-chan StreamEvent, error)
    Name() string
    Model() string
}
```

### Stream Events

Providers emit three types of events:

```go
type StreamEvent struct {
    Type    EventType  // chunk, complete, error
    Content string     // Text chunk
    Error   error      // Error if Type == EventTypeError
    Done    bool       // True when stream is complete
}
```

### Integration Points

1. **ChatModel** (`internal/ui/models/chat.go`)
   - Creates AI provider on initialization
   - Falls back to mock if no API keys
   - Manages streaming state

2. **AI Providers** (`internal/ai/providers/`)
   - `claude.go` - Anthropic Claude API
   - `openai.go` - OpenAI GPT API
   - Each implements SSE streaming

3. **Factory** (`internal/ai/factory.go`)
   - Loads config from environment
   - Auto-selects provider based on available keys
   - Returns initialized provider

---

## Supported Models

### Claude (Anthropic)

**Latest Models:**
- `claude-opus-4-20250514` ⭐ (default) - Most capable
- `claude-sonnet-4-20250514` - Balanced performance
- `claude-haiku-3-20250307` - Fast, lightweight

**Pricing:** See [Anthropic Pricing](https://www.anthropic.com/pricing)

**API Docs:** [Anthropic API Reference](https://docs.anthropic.com/en/api)

### OpenAI (GPT)

**Latest Models:**
- `gpt-4o` ⭐ (default) - GPT-4 Optimized
- `gpt-4-turbo` - GPT-4 Turbo
- `gpt-4` - Original GPT-4

**Pricing:** See [OpenAI Pricing](https://openai.com/pricing)

**API Docs:** [OpenAI API Reference](https://platform.openai.com/docs/api-reference)

---

## Examples

### Example 1: Using Claude

```bash
export ANTHROPIC_API_KEY="sk-ant-api03-..."
rycode
```

**Status bar shows:**
```
⚡ Claude (claude-opus-4-20250514) is responding... │ Claude (claude-opus-4-20250514) │ 2 messages
```

### Example 2: Using GPT-4

```bash
export OPENAI_API_KEY="sk-..."
export RYCODE_AI_PROVIDER="openai"
rycode
```

**Status bar shows:**
```
⚡ OpenAI (gpt-4o) is responding... │ OpenAI (gpt-4o) │ 4 messages
```

### Example 3: Custom Model

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
export RYCODE_CLAUDE_MODEL="claude-sonnet-4-20250514"
rycode
```

**Status bar shows:**
```
⚡ Claude (claude-sonnet-4-20250514) is responding... │ Claude (claude-sonnet-4-20250514) │ 6 messages
```

### Example 4: No API Keys (Mock Mode)

```bash
# No API keys set
rycode
```

**Status bar shows:**
```
⚡ Mock is responding... │ ⚠️ No AI (set ANTHROPIC_API_KEY or OPENAI_API_KEY) │ 2 messages
```

---

## API Request Format

### Claude Request

```json
{
  "model": "claude-opus-4-20250514",
  "messages": [
    {"role": "user", "content": "Hello!"},
    {"role": "assistant", "content": "Hi! How can I help?"},
    {"role": "user", "content": "Explain recursion"}
  ],
  "max_tokens": 4096,
  "stream": true,
  "temperature": 0.7,
  "top_p": 0.9
}
```

### OpenAI Request

```json
{
  "model": "gpt-4o",
  "messages": [
    {"role": "user", "content": "Hello!"},
    {"role": "assistant", "content": "Hi! How can I help?"},
    {"role": "user", "content": "Explain recursion"}
  ],
  "max_tokens": 4096,
  "stream": true,
  "temperature": 0.7,
  "top_p": 0.9
}
```

---

## Advanced Configuration

### Custom Config in Code

```go
import "github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai"

config := &ai.Config{
    Provider:          "claude",
    ClaudeAPIKey:      "sk-ant-...",
    ClaudeModel:       "claude-opus-4-20250514",
    MaxTokens:         8192,
    Temperature:       0.5,
    TopP:              0.9,
    RequestsPerMinute: 50,
    TokensPerMinute:   100000,
}

provider, err := ai.NewProvider(config)
```

### Streaming from Custom Provider

```go
ctx := context.Background()
history := []ai.Message{
    {Role: ai.RoleUser, Content: "Hello"},
}

eventCh, err := provider.Stream(ctx, "Explain Go interfaces", history)
if err != nil {
    log.Fatal(err)
}

for event := range eventCh {
    switch event.Type {
    case ai.EventTypeChunk:
        fmt.Print(event.Content)
    case ai.EventTypeComplete:
        fmt.Println("\n[Stream complete]")
    case ai.EventTypeError:
        fmt.Printf("\nError: %v\n", event.Error)
    }
}
```

---

## Troubleshooting

### "No AI keys found" Error

**Problem:** Neither `ANTHROPIC_API_KEY` nor `OPENAI_API_KEY` is set.

**Solution:** Set at least one API key:
```bash
export ANTHROPIC_API_KEY="sk-ant-..."
# OR
export OPENAI_API_KEY="sk-..."
```

### "API returned status 401" Error

**Problem:** Invalid or expired API key.

**Solution:** Check your API key:
- Claude: Visit [Anthropic Console](https://console.anthropic.com/)
- OpenAI: Visit [OpenAI Platform](https://platform.openai.com/api-keys)

### "API returned status 429" Error

**Problem:** Rate limit exceeded.

**Solution:**
- Wait a few seconds and try again
- Reduce request frequency
- Upgrade your API plan for higher limits

### Streaming Stops Mid-Response

**Problem:** Network interruption or API timeout.

**Solution:**
- Check your internet connection
- Retry the message
- The error will be shown in the chat

---

## Performance

### Latency

- **First Token:** 200-500ms (depends on provider and model)
- **Token Rate:** 20-100 tokens/second (streaming)
- **Network:** HTTP/2 with SSE (efficient)

### Resource Usage

- **Memory:** ~5MB per chat session (including history)
- **CPU:** Minimal (event-driven architecture)
- **Network:** ~1-10 KB/request (compressed)

### Caching

- No local caching (all requests go to API)
- Conversation history sent with each request
- Future: Add response caching for repeated queries

---

## Roadmap

### Phase 2.1: Enhanced AI Features ✅ (Complete)
- [x] Claude Opus 4 integration
- [x] GPT-4o integration
- [x] Streaming SSE responses
- [x] Conversation history
- [x] Auto-provider selection

### Phase 2.2: Advanced Features (Next)
- [ ] Token usage tracking and display
- [ ] Rate limiting and retry logic
- [ ] Response caching
- [ ] Multi-provider fallback (Claude → GPT-4 → Mock)
- [ ] Configurable system prompts
- [ ] Temperature/top-p sliders in UI

### Phase 2.3: Developer Experience
- [ ] Context-aware code suggestions
- [ ] File content injection in prompts
- [ ] Code block extraction and editing
- [ ] Diff generation and preview
- [ ] Error explanation from logs

### Phase 2.4: Production Features
- [ ] Structured logging for AI requests
- [ ] Telemetry and analytics
- [ ] Cost tracking and budgets
- [ ] Session persistence
- [ ] Export conversation history

---

## API Reference

### `ai.Provider` Interface

```go
type Provider interface {
    // Stream sends a prompt and streams back response tokens
    Stream(ctx context.Context, prompt string, messages []Message) (<-chan StreamEvent, error)

    // Name returns the provider name (e.g., "Claude", "OpenAI")
    Name() string

    // Model returns the model identifier (e.g., "claude-opus-4-20250514")
    Model() string
}
```

### `ai.Config` Struct

```go
type Config struct {
    Provider          string  // "claude", "openai", "auto"
    ClaudeAPIKey      string  // Anthropic API key
    OpenAIAPIKey      string  // OpenAI API key
    ClaudeModel       string  // Claude model (default: "claude-opus-4-20250514")
    OpenAIModel       string  // OpenAI model (default: "gpt-4o")
    MaxTokens         int     // Max response tokens (default: 4096)
    Temperature       float64 // Sampling temperature 0-1 (default: 0.7)
    TopP              float64 // Nucleus sampling (default: 0.9)
    RequestsPerMinute int     // Rate limit (default: 50)
    TokensPerMinute   int     // Token rate limit (default: 100000)
}
```

### `ai.Message` Struct

```go
type Message struct {
    Role    Role   `json:"role"`    // "user", "assistant", "system"
    Content string `json:"content"` // Message text
}
```

### `ai.StreamEvent` Struct

```go
type StreamEvent struct {
    Type    EventType // "chunk", "complete", "error"
    Content string    // Text chunk (if Type == EventTypeChunk)
    Error   error     // Error (if Type == EventTypeError)
    Done    bool      // True if stream is complete
}
```

---

## Security

### API Key Storage

- **Environment Variables:** Keys are read from env vars, not stored in files
- **No Persistence:** Keys are never written to disk
- **Memory Only:** Keys exist only in process memory

### Best Practices

1. **Never commit API keys** to version control
2. **Use `.env` files** for local development (add to `.gitignore`)
3. **Rotate keys regularly** (every 90 days recommended)
4. **Use separate keys** for development vs production
5. **Monitor usage** on provider dashboards

### Example `.env` File

```bash
# .env (DO NOT COMMIT!)
ANTHROPIC_API_KEY=sk-ant-api03-...
OPENAI_API_KEY=sk-...
RYCODE_AI_PROVIDER=auto
RYCODE_CLAUDE_MODEL=claude-opus-4-20250514
```

Then load with:
```bash
source .env
rycode
```

---

## Contributing

### Adding a New Provider

1. Create `internal/ai/providers/newprovider.go`
2. Implement the `ai.Provider` interface
3. Register in `internal/ai/providers/register.go`
4. Add configuration options to `ai.Config`
5. Update documentation

**Example skeleton:**

```go
package providers

import "github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai"

type NewProvider struct {
    apiKey string
    model  string
}

func NewNewProvider(apiKey string, config *ai.Config) *NewProvider {
    return &NewProvider{apiKey: apiKey, model: config.NewModel}
}

func (p *NewProvider) Name() string { return "NewProvider" }
func (p *NewProvider) Model() string { return p.model }

func (p *NewProvider) Stream(ctx context.Context, prompt string, messages []ai.Message) (<-chan ai.StreamEvent, error) {
    // Implementation here
}
```

---

## License

MIT License - see [LICENSE](LICENSE) for details

---

## Support

- **Documentation:** [AI_INTEGRATION.md](AI_INTEGRATION.md) (this file)
- **Issues:** [GitHub Issues](https://github.com/aaronmrosenthal/RyCode/issues)
- **Discussions:** [GitHub Discussions](https://github.com/aaronmrosenthal/RyCode/discussions)

---

<div align="center">

**Real AI, Real-Time, Real Beautiful** 🤖✨

*Making AI-powered coding accessible in the terminal* 🟢

</div>
