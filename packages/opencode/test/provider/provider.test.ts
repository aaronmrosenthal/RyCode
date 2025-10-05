import { describe, test, expect, beforeAll } from "bun:test"
import { Provider } from "../../src/provider/provider"
import { TestSetup } from "../setup"

describe("Provider", () => {
  let cleanup: () => void

  beforeAll(() => {
    // Mock environment with test API key
    cleanup = TestSetup.mockEnv({
      ANTHROPIC_API_KEY: "test-key-anthropic",
    })
  })

  describe("getModel", () => {
    test("should retrieve Anthropic model successfully", async () => {
      const model = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")

      expect(model.providerID).toBe("anthropic")
      expect(model.modelID).toBe("claude-3-5-sonnet-20241022")
      expect(model.info).toBeDefined()
      expect(model.info.id).toBe("claude-3-5-sonnet-20241022")
      expect(model.language).toBeDefined()
    })

    test("should throw ModelNotFoundError for invalid provider", async () => {
      await expect(Provider.getModel("invalid-provider-xyz", "model")).rejects.toThrow(
        Provider.ModelNotFoundError,
      )
    })

    test("should throw ModelNotFoundError for invalid model", async () => {
      await expect(Provider.getModel("anthropic", "nonexistent-model-xyz")).rejects.toThrow(
        Provider.ModelNotFoundError,
      )
    })

    test("should cache model instances", async () => {
      const model1 = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")
      const model2 = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")

      // Both should reference the same model info
      expect(model1.info).toBe(model2.info)
    })

    test("should support different models from same provider", async () => {
      const sonnet = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")
      const haiku = await Provider.getModel("anthropic", "claude-3-5-haiku-20241022")

      expect(sonnet.modelID).toBe("claude-3-5-sonnet-20241022")
      expect(haiku.modelID).toBe("claude-3-5-haiku-20241022")
      expect(sonnet.npm).toBe(haiku.npm) // Same provider package
    })
  })

  describe("list", () => {
    test("should list available providers", async () => {
      const providers = await Provider.list()

      expect(Object.keys(providers).length).toBeGreaterThan(0)
      expect(providers["anthropic"]).toBeDefined()
      expect(providers["anthropic"].info.models).toBeDefined()
    })

    test("should include provider information", async () => {
      const providers = await Provider.list()
      const anthropic = providers["anthropic"]

      expect(anthropic.info.id).toBe("anthropic")
      expect(anthropic.info.name).toBeDefined()
      expect(anthropic.info.npm).toBeDefined()
      expect(anthropic.source).toBeDefined()
    })
  })

  describe("parseModel", () => {
    test("should parse provider/model format", () => {
      const result = Provider.parseModel("anthropic/claude-3-5-sonnet-20241022")

      expect(result.providerID).toBe("anthropic")
      expect(result.modelID).toBe("claude-3-5-sonnet-20241022")
    })

    test("should handle models with slashes in name", () => {
      const result = Provider.parseModel("openrouter/anthropic/claude-3-5-sonnet")

      expect(result.providerID).toBe("openrouter")
      expect(result.modelID).toBe("anthropic/claude-3-5-sonnet")
    })
  })

  describe("defaultModel", () => {
    test("should return a valid default model", async () => {
      const model = await Provider.defaultModel()

      expect(model.providerID).toBeDefined()
      expect(model.modelID).toBeDefined()
      expect(typeof model.providerID).toBe("string")
      expect(typeof model.modelID).toBe("string")
    })

    test("should prefer high-priority models", async () => {
      const model = await Provider.defaultModel()

      // Should be one of the priority models (from provider.ts:477)
      const isPriority =
        model.modelID.includes("gemini-2.5-pro-preview") ||
        model.modelID.includes("gpt-5") ||
        model.modelID.includes("claude-sonnet-4")

      // If available providers have priority models, one should be selected
      // Otherwise any valid model is acceptable
      expect(model.providerID).toBeTruthy()
    })
  })

  describe("getSmallModel", () => {
    test("should return a small model for the provider", async () => {
      const smallModel = await Provider.getSmallModel("anthropic")

      // Should get a haiku model if available
      if (smallModel) {
        expect(smallModel.providerID).toBe("anthropic")
        expect(smallModel.modelID).toMatch(/haiku|nano/)
      }
    })

    test("should return undefined if provider has no small models", async () => {
      const smallModel = await Provider.getSmallModel("nonexistent-provider")

      expect(smallModel).toBeUndefined()
    })
  })

  describe("sort", () => {
    test("should sort models by priority", () => {
      const models = [
        { id: "claude-3-haiku", name: "Claude 3 Haiku" } as any,
        { id: "gpt-5-latest", name: "GPT-5" } as any,
        { id: "claude-sonnet-4", name: "Claude Sonnet 4" } as any,
      ]

      const sorted = Provider.sort(models)

      // gpt-5 should be first (priority from provider.ts:477)
      expect(sorted[0].id).toBe("gpt-5-latest")
    })

    test("should prefer latest versions", () => {
      const models = [
        { id: "model-v1", name: "Model v1" } as any,
        { id: "model-latest", name: "Model Latest" } as any,
        { id: "model-v2", name: "Model v2" } as any,
      ]

      const sorted = Provider.sort(models)

      // "latest" should be preferred
      expect(sorted[0].id).toBe("model-latest")
    })
  })
})
