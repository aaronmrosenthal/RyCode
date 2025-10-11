# Provider Authentication System - Architecture

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Interface (TUI)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │ Model Dialog │  │ Status Bar   │  │ Cost Display │         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘         │
└─────────┼──────────────────┼──────────────────┼─────────────────┘
          │                  │                  │
          └──────────────────┼──────────────────┘
                             │
┌─────────────────────────────▼─────────────────────────────────┐
│                      Auth Manager (High-Level API)             │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  authenticate() | getStatus() | getRecommendations()     │ │
│  │  recordUsage() | getCostSummary() | healthCheck()        │ │
│  └──────────────────────────────────────────────────────────┘ │
└──────┬────────┬────────┬────────┬────────┬───────────────────┘
       │        │        │        │        │
       ▼        ▼        ▼        ▼        ▼
┌──────────┐ ┌────┐ ┌────┐ ┌────┐ ┌─────────────────────┐
│ Provider │ │Auto│ │Cost│ │Model│ │    Storage Layer    │
│ Registry │ │Det.│ │Trk.│ │Rec. │ │ ┌─────────────────┐ │
│          │ │    │ │    │ │     │ │ │ Credential      │ │
│ ┌──────┐ │ │    │ │    │ │     │ │ │ Store           │ │
│ │Claude│ │ │    │ │    │ │     │ │ └────────┬────────┘ │
│ │GPT-4 │ │ │    │ │    │ │     │ │ ┌────────▼────────┐ │
│ │Gemini│ │ │    │ │    │ │     │ │ │ Audit Log       │ │
│ │Grok  │ │ │    │ │    │ │     │ │ └─────────────────┘ │
│ │Qwen  │ │ │    │ │    │ │     │ │ ┌─────────────────┐ │
│ └──┬───┘ │ │    │ │    │ │     │ │ │ Existing Auth   │ │
│    │     │ │    │ │    │ │     │ │ │ Namespace       │ │
└────┼─────┘ └────┘ └────┘ └─────┘ │ └─────────────────┘ │
     │                               └─────────────────────┘
     │
     ▼
┌─────────────────────────────────────────────────────────────┐
│                    Security Layer                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ Rate Limiter │  │   Circuit    │  │    Input     │     │
│  │              │  │   Breaker    │  │  Validator   │     │
│  │ 5 auth/min   │  │              │  │              │     │
│  │ 60 req/min   │  │ Auto-recover │  │ Sanitization │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
     │                  │                  │
     ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   Provider APIs                              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │
│  │Anthropic │ │  OpenAI  │ │  Google  │ │Grok/Qwen │      │
│  │   API    │ │   API    │ │   API    │ │   APIs   │      │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 🔄 Authentication Flow

```
User Action
    │
    ├─ Manual Auth ──────────────┐
    │                             │
    ├─ Auto-Detect ──────────────┤
    │                             │
    └─ OAuth Flow ───────────────┤
                                  │
                                  ▼
                         ┌────────────────┐
                         │  Auth Manager  │
                         └────────┬───────┘
                                  │
                    ┌─────────────┼─────────────┐
                    │             │             │
                    ▼             ▼             ▼
            ┌──────────┐  ┌─────────────┐  ┌──────────┐
            │Rate Limit│  │  Validate   │  │ Circuit  │
            │  Check   │  │   Input     │  │ Breaker  │
            └────┬─────┘  └──────┬──────┘  └────┬─────┘
                 │                │               │
                 └────────────────┼───────────────┘
                                  │
                                  ▼
                         ┌────────────────┐
                         │   Provider     │
                         │ Authenticate   │
                         └────────┬───────┘
                                  │
                    ┌─────────────┼─────────────┐
                    │             │             │
                    ▼             ▼             ▼
            ┌──────────┐  ┌─────────────┐  ┌──────────┐
            │  Store   │  │ Audit Log   │  │  Return  │
            │   Creds  │  │   Event     │  │  Status  │
            └──────────┘  └─────────────┘  └──────────┘
```

## 💰 Cost Tracking Flow

```
API Request
    │
    ▼
┌─────────────────┐
│ Provider API    │
│ Returns Usage   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Auth Manager    │
│ recordUsage()   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Cost Tracker    │
│ - Calculate $   │
│ - Store history │
│ - Update totals │
└────────┬────────┘
         │
    ┌────┼────┐
    │    │    │
    ▼    ▼    ▼
┌────┐ ┌──┐ ┌──────────┐
│Save│ │UI│ │Generate  │
│Data│ │  │ │Tips      │
└────┘ └──┘ └──────────┘
         │
         ▼
  Status Bar Update
  "💰 $0.12 today"
```

## 🎯 Model Recommendation Flow

```
User Context
   │
   ├─ Task: code_generation
   ├─ Complexity: medium
   ├─ Speed: balanced
   └─ Cost: cheapest
         │
         ▼
┌─────────────────────┐
│ Model Recommender   │
│ - Score all models  │
│ - Apply filters     │
│ - Rank by fit       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Scoring Algorithm   │
│ + Task match: 30pts │
│ + Quality: 10pts    │
│ + Speed: 8pts       │
│ + Cost: 20pts       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Top 3 Results       │
│ 1. Haiku (92pts)    │
│ 2. Sonnet (88pts)   │
│ 3. GPT-3.5 (75pts)  │
└─────────────────────┘
```

## 🔒 Security Layers

```
Request
  │
  ├──► Rate Limiter ───► Allows?
  │                       │
  │                       ├─ Yes ──┐
  │                       └─ No ───┤
  │                                │
  ├──► Input Validator ─► Valid?  │
  │                       │        │
  │                       ├─ Yes ──┤
  │                       └─ No ───┤
  │                                │
  ├──► Circuit Breaker ─► Open?   │
  │                       │        │
  │                       ├─ No ───┤
  │                       └─ Yes ──┤
  │                                │
  └──► Audit Logger ────► Record  │
                          Event   │
                                  │
                          ┌───────┴────────┐
                          │                │
                          ▼                ▼
                    ┌──────────┐    ┌──────────┐
                    │ Proceed  │    │  Reject  │
                    │  to API  │    │  Return  │
                    └──────────┘    │  Error   │
                                    └──────────┘
```

## 📊 Data Flow

```
┌──────────────────────────────────────────────────────────┐
│                    Memory (Fast)                          │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐        │
│  │Rate Limit  │  │Circuit     │  │Recent      │        │
│  │Counters    │  │States      │  │Audit Events│        │
│  └────────────┘  └────────────┘  └────────────┘        │
└──────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────┐
│                  Disk (Persistent)                        │
│  ┌──────────────────────────────────────────────────┐   │
│  │ ~/.rycode/data/                                   │   │
│  │  ├─ auth.json (encrypted credentials)            │   │
│  │  ├─ auth-audit.log (audit events)                │   │
│  │  └─ costs.json (usage history)                   │   │
│  └──────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────┘
```

## 🔄 Error Handling Hierarchy

```
Error Occurs
    │
    ▼
Is AuthenticationError?
    │
    ├─ Yes ──► Check Type
    │           │
    │           ├─ RateLimitError ──► Wait & Retry
    │           ├─ NetworkError ───► Retry with backoff
    │           ├─ InvalidKeyError ─► Show help URL
    │           ├─ ExpiredError ────► Re-authenticate
    │           └─ ValidationError ─► Show hint
    │
    └─ No ───► Unknown Error
                │
                └─► Log & Show generic message
```

## 🎨 Component Relationships

```
┌─────────────────────────────────────────────────────┐
│                   AuthManager                        │
│  (Orchestrates everything)                           │
└────────┬────────────────────────────────────────────┘
         │
    ┌────┴────┬─────────┬─────────┬──────────┐
    │         │         │         │          │
    ▼         ▼         ▼         ▼          ▼
┌────────┐ ┌────┐ ┌────┐ ┌──────┐ ┌──────────┐
│Provider│ │Auto│ │Cost│ │Model │ │ Storage  │
│Registry│ │Det.│ │Trk.│ │Rec.  │ │          │
└───┬────┘ └────┘ └────┘ └──────┘ └────┬─────┘
    │                                   │
    ├─ anthropicProvider                ├─ credentialStore
    ├─ openaiProvider                   └─ auditLog
    ├─ googleProvider
    ├─ grokProvider
    └─ qwenProvider
         │
         └─ Uses Security Layer
              │
              ├─ rateLimiter
              ├─ circuitBreaker
              └─ inputValidator
```

## 📦 Module Dependencies

```
AuthManager
    │
    ├── depends on ──► ProviderRegistry
    │                   │
    │                   └── depends on ──► Individual Providers
    │                                       │
    │                                       └── depends on ──► Security Layer
    │
    ├── depends on ──► CredentialStore
    │                   │
    │                   └── depends on ──► Existing Auth Namespace
    │
    ├── depends on ──► AuditLog
    │
    ├── depends on ──► CostTracker
    │
    ├── depends on ──► ModelRecommender
    │
    └── depends on ──► SmartSetup (Auto-detect)
```

## 🚀 Deployment Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    RyCode CLI/TUI                        │
│                                                           │
│  ┌────────────────────────────────────────────────┐     │
│  │  Auth Manager (Single Instance)                 │     │
│  │  - Initialized on startup                       │     │
│  │  - Auto-detects credentials                     │     │
│  │  - Loads from disk                              │     │
│  └────────────────────────────────────────────────┘     │
│                        │                                 │
│                        ▼                                 │
│  ┌────────────────────────────────────────────────┐     │
│  │  Provider Instances (Lazy-loaded)               │     │
│  │  - Created on first use                         │     │
│  │  - Cached for performance                       │     │
│  └────────────────────────────────────────────────┘     │
│                        │                                 │
│                        ▼                                 │
│  ┌────────────────────────────────────────────────┐     │
│  │  Security Layer (Always Active)                 │     │
│  │  - Rate limiting counters                       │     │
│  │  - Circuit breaker states                       │     │
│  │  - Validation rules                             │     │
│  └────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
```

## 💡 Key Design Principles

1. **Single Responsibility** - Each module does one thing well
2. **Strategy Pattern** - Providers are interchangeable
3. **Circuit Breaker** - Fail fast, recover automatically
4. **Defense in Depth** - Multiple security layers
5. **User-Centric** - Errors are helpful, not cryptic
6. **Observable** - Everything is logged and trackable
7. **Extensible** - Easy to add new providers
8. **Testable** - Pure functions, clear interfaces

---

**This architecture ensures:**
- ✅ Security at every layer
- ✅ User-friendly error messages
- ✅ Automatic cost tracking
- ✅ Smart model recommendations
- ✅ Resilient to provider outages
- ✅ Easy to extend and maintain
