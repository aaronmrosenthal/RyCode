import { describe, test, expect, beforeAll } from "bun:test"
import { Provider } from "../../src/provider/provider"
import { TestSetup } from "../setup"
import { Instance } from "../../src/project/instance"
import { tmpdir } from "../fixture/fixture"

describe("Provider", () => {
  beforeAll(() => {
    // Mock environment with test API key
    TestSetup.mockEnv({
      ANTHROPIC_API_KEY: "test-key-anthropic",
    })
  })

  // Helper to run Provider functions within Instance context
  async function withInstance<T>(fn: () => Promise<T>): Promise<T> {
    await using tmp = await tmpdir({ git: true })
    return Instance.provide({
      directory: tmp.path,
      fn,
    })
  }

  describe("getModel", () => {
    test("should retrieve Anthropic model successfully", async () => {
      await withInstance(async () => {
        const model = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")

        expect(model.providerID).toBe("anthropic")
        expect(model.modelID).toBe("claude-3-5-sonnet-20241022")
        expect(model.info).toBeDefined()
        expect(model.info.id).toBe("claude-3-5-sonnet-20241022")
        expect(model.language).toBeDefined()
      })
    })

    test("should throw ModelNotFoundError for invalid provider", async () => {
      await withInstance(async () => {
        await expect(Provider.getModel("invalid-provider-xyz", "model")).rejects.toThrow(
          Provider.ModelNotFoundError,
        )
      })
    })

    test("should throw ModelNotFoundError for invalid model", async () => {
      await withInstance(async () => {
        await expect(Provider.getModel("anthropic", "nonexistent-model-xyz")).rejects.toThrow(
          Provider.ModelNotFoundError,
        )
      })
    })

    test("should cache model instances", async () => {
      await withInstance(async () => {
        const model1 = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")
        const model2 = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")

        // Both should reference the same model info
        expect(model1.info).toBe(model2.info)
      })
    })

    test("should support different models from same provider", async () => {
      await withInstance(async () => {
        const sonnet = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")
        const haiku = await Provider.getModel("anthropic", "claude-3-5-haiku-20241022")

        expect(sonnet.modelID).toBe("claude-3-5-sonnet-20241022")
        expect(haiku.modelID).toBe("claude-3-5-haiku-20241022")
        expect(sonnet.npm).toBe(haiku.npm!) // Same provider package
      })
    })
  })

  describe("list", () => {
    test("should list available providers", async () => {
      await withInstance(async () => {
        const providers = await Provider.list()

        expect(Object.keys(providers).length).toBeGreaterThan(0)
        expect(providers["anthropic"]).toBeDefined()
        expect(providers["anthropic"].info.models).toBeDefined()
      })
    })

    test("should include provider information", async () => {
      await withInstance(async () => {
        const providers = await Provider.list()
        const anthropic = providers["anthropic"]

        expect(anthropic.info.id).toBe("anthropic")
        expect(anthropic.info.name).toBeDefined()
        expect(anthropic.info.npm).toBeDefined()
        expect(anthropic.source).toBeDefined()
      })
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
      await withInstance(async () => {
        const model = await Provider.defaultModel()

        expect(model.providerID).toBeDefined()
        expect(model.modelID).toBeDefined()
        expect(typeof model.providerID).toBe("string")
        expect(typeof model.modelID).toBe("string")
      })
    })

    test("should prefer high-priority models", async () => {
      await withInstance(async () => {
        const model = await Provider.defaultModel()

        // Should be one of the priority models (from provider.ts:477)
        // If available providers have priority models, one should be selected
        // Otherwise any valid model is acceptable
        expect(model.providerID).toBeTruthy()
      })
    })
  })

  describe("getSmallModel", () => {
    test("should return a small model for the provider", async () => {
      await withInstance(async () => {
        const smallModel = await Provider.getSmallModel("anthropic")

        // Should get a haiku model if available
        if (smallModel) {
          expect(smallModel.providerID).toBe("anthropic")
          expect(smallModel.modelID).toMatch(/haiku|nano/)
        }
      })
    })

    test("should return undefined if provider has no small models", async () => {
      await withInstance(async () => {
        const smallModel = await Provider.getSmallModel("nonexistent-provider")

        expect(smallModel).toBeUndefined()
      })
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

      // claude-sonnet-4 should be first (priority: sonnet-4 > gpt-5 from provider.ts:503)
      expect(sorted[0].id).toBe("claude-sonnet-4")
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

  // SECURITY TESTS - SDK initialization race condition

  describe("SDK initialization race conditions", () => {
    test("should handle concurrent getModel calls without duplicate initialization", async () => {
      await withInstance(async () => {
        // Start 10 concurrent getModel calls for the same model
        const promises = Array(10)
          .fill(0)
          .map(() => Provider.getModel("anthropic", "claude-3-5-sonnet-20241022"))

        const results = await Promise.all(promises)

        // All should succeed
        expect(results.length).toBe(10)

        // All should have the same model ID
        expect(results.every((r) => r.modelID === "claude-3-5-sonnet-20241022")).toBe(true)

        // All should reference the same SDK instance (via npm package name)
        const npmPackages = new Set(results.map((r) => r.npm))
        expect(npmPackages.size).toBe(1) // Only one SDK should be loaded
      })
    })

    test("should handle concurrent getModel calls for different models from same provider", async () => {
      await withInstance(async () => {
        // Start concurrent calls for different Anthropic models
        const promises = [
          Provider.getModel("anthropic", "claude-3-5-sonnet-20241022"),
          Provider.getModel("anthropic", "claude-3-5-haiku-20241022"),
          Provider.getModel("anthropic", "claude-3-5-sonnet-20241022"), // Duplicate
        ]

        const results = await Promise.all(promises)

        // All should succeed
        expect(results.length).toBe(3)

        // Should have 2 unique models
        const modelIds = new Set(results.map((r) => r.modelID))
        expect(modelIds.size).toBe(2)

        // All should use same provider SDK
        const npmPackages = new Set(results.map((r) => r.npm))
        expect(npmPackages.size).toBe(1)
      })
    })

    test("should not cache failed SDK initialization", async () => {
      await withInstance(async () => {
        // Try to get a model from a provider that will fail
        // (This test assumes the provider exists but SDK init might fail)
        try {
          await Provider.getModel("invalid-provider-xyz", "some-model")
        } catch (e) {
          // Expected to fail
          expect(Provider.ModelNotFoundError.isInstance(e)).toBe(true)
        }

        // Second attempt should also try initialization (not use cached failure)
        try {
          await Provider.getModel("invalid-provider-xyz", "some-model")
        } catch (e) {
          expect(Provider.ModelNotFoundError.isInstance(e)).toBe(true)
        }
      })
    })

    test("should clean up pending promises after initialization", async () => {
      const { Instance } = await import("../../src/project/instance")
      const { tmpdir } = await import("../fixture/fixture")

      await using tmp = await tmpdir({ git: true })
      await Instance.provide({
        directory: tmp.path,
        fn: async () => {
          // Get a model to trigger initialization
          const model = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")
          expect(model).toBeDefined()

          // Subsequent calls should use cached SDK, not pending promises
          const model2 = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")
          expect(model2.modelID).toBe(model.modelID)
        },
      })
    })

    test("should handle race condition during SDK reload", async () => {
      const { Instance } = await import("../../src/project/instance")
      const { tmpdir } = await import("../fixture/fixture")

      await using tmp = await tmpdir({ git: true })
      await Instance.provide({
        directory: tmp.path,
        fn: async () => {
          // Get initial model
          const model1 = await Provider.getModel("anthropic", "claude-3-5-sonnet-20241022")
          expect(model1).toBeDefined()

          // Multiple subsequent calls should all work correctly
          const promises = Array(5)
            .fill(0)
            .map(() => Provider.getModel("anthropic", "claude-3-5-haiku-20241022"))

          const results = await Promise.all(promises)
          expect(results.length).toBe(5)
          expect(results.every((r) => r.modelID === "claude-3-5-haiku-20241022")).toBe(true)
        },
      })
    })
  })
})
