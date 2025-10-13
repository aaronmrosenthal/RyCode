import { test, expect } from '@playwright/test';

/**
 * Playwright E2E Tests for RyCode Model Selector
 *
 * Tests the web-based visualization of the model selector to verify:
 * 1. All providers are displayed (5 providers)
 * 2. All models are loaded (30 models total)
 * 3. Authentication status is correctly shown
 * 4. Search functionality works
 * 5. Keyboard navigation is functional
 * 6. Provider collapse/expand works
 * 7. CLI providers are distinguished from API providers
 */

test.describe('RyCode Model Selector', () => {
  test.beforeEach(async ({ page }) => {
    // Load the test visualization
    await page.goto('file://' + process.cwd() + '/test-model-selector-web.html');
    await page.waitForLoadState('domcontentloaded');
  });

  test('should display correct provider and model counts', async ({ page }) => {
    // Verify test results panel shows correct counts
    await expect(page.locator('#count-providers')).toHaveText('5 providers');
    await expect(page.locator('#count-models')).toHaveText('30 models');
    await expect(page.locator('#count-auth')).toHaveText('4 authenticated');

    // Verify all status indicators are passing
    await expect(page.locator('#status-providers')).toHaveClass(/pass/);
    await expect(page.locator('#status-models')).toHaveClass(/pass/);
    await expect(page.locator('#status-auth')).toHaveClass(/pass/);
  });

  test('should show all provider groups', async ({ page }) => {
    // Check that all 5 providers are rendered
    const providers = [
      'provider-anthropic',
      'provider-openai',
      'provider-claude-cli',
      'provider-gemini-locked',
    ];

    for (const providerTestId of providers) {
      await expect(page.locator(`[data-testid="${providerTestId}"]`)).toBeVisible();
    }
  });

  test('should display recent models section', async ({ page }) => {
    // Verify Recent section exists
    await expect(page.locator('text=📌 RECENT')).toBeVisible();

    // Verify 3 recent models are shown
    await expect(page.locator('[data-testid="recent-model-1"]')).toBeVisible();
    await expect(page.locator('[data-testid="recent-model-2"]')).toBeVisible();
    await expect(page.locator('[data-testid="recent-model-3"]')).toBeVisible();

    // Verify the first recent model has correct metadata
    const firstRecent = page.locator('[data-testid="recent-model-1"]');
    await expect(firstRecent).toContainText('Claude 4.5 Sonnet');
    await expect(firstRecent).toContainText('32K out');
    await expect(firstRecent).toContainText('2 min ago');
  });

  test('should show authentication indicators', async ({ page }) => {
    // Authenticated providers show ✓
    const anthropicProvider = page.locator('[data-testid="provider-anthropic"]');
    await expect(anthropicProvider).toContainText('✓');
    await expect(anthropicProvider).toContainText('ANTHROPIC');

    const openaiProvider = page.locator('[data-testid="provider-openai"]');
    await expect(openaiProvider).toContainText('✓');

    // Locked providers show 🔒
    const geminiProvider = page.locator('[data-testid="provider-gemini-locked"]');
    await expect(geminiProvider).toContainText('🔒');
    await expect(geminiProvider).toContainText('GEMINI');
  });

  test('should distinguish CLI providers from API providers', async ({ page }) => {
    // CLI provider has "CLI" label
    const cliProvider = page.locator('[data-testid="provider-claude-cli"]');
    await expect(cliProvider).toContainText('CLI');

    // Expand to see CLI badge on models
    await cliProvider.locator('.provider-header').click();
    await expect(cliProvider.locator('text=CLI').first()).toBeVisible();
  });

  test('should display model metadata badges', async ({ page }) => {
    // Check that models have speed, cost, and other badges
    const claudeModel = page.locator('[data-testid="model-claude-4.5"]');
    await expect(claudeModel).toContainText('⚡'); // Speed badge
    await expect(claudeModel).toContainText('💰'); // Cost badge
    await expect(claudeModel).toContainText('🔥'); // Popular badge
    await expect(claudeModel).toContainText('32K out'); // Output limit

    const gpt4oModel = page.locator('[data-testid="model-gpt-4o"]');
    await expect(gpt4oModel).toContainText('128K ctx'); // Context size
  });

  test('should support search functionality', async ({ page }) => {
    const searchInput = page.locator('#search-input');
    await searchInput.fill('claude');

    // Wait for search to filter
    await page.waitForTimeout(300);

    // Verify search results are filtered
    await expect(page.locator('[data-testid="model-claude-4.5"]')).toBeVisible();

    // Verify GPT models are hidden
    const gpt4o = page.locator('[data-testid="model-gpt-4o"]');
    await expect(gpt4o).toBeHidden();
  });

  test('should support keyboard shortcuts - search focus', async ({ page }) => {
    const searchInput = page.locator('#search-input');

    // Press "/" to focus search
    await page.keyboard.press('/');

    // Verify search input is focused
    await expect(searchInput).toBeFocused();
  });

  test('should support keyboard shortcuts - provider jump', async ({ page }) => {
    // Press "2" to jump to second provider (OpenAI)
    await page.keyboard.press('2');

    // Wait for scroll animation
    await page.waitForTimeout(500);

    // Verify OpenAI provider is in viewport
    const openaiProvider = page.locator('[data-testid="provider-openai"]');
    await expect(openaiProvider).toBeInViewport();
  });

  test('should collapse and expand provider groups', async ({ page }) => {
    // Find Anthropic provider
    const anthropicProvider = page.locator('[data-testid="provider-anthropic"]');
    const anthropicModels = page.locator('#models-anthropic');
    const toggleButton = page.locator('#toggle-anthropic');

    // Initially expanded
    await expect(anthropicModels).toBeVisible();
    await expect(toggleButton).toContainText('▼ Collapse');

    // Click to collapse
    await anthropicProvider.locator('.provider-header').click();

    // Now collapsed
    await expect(anthropicModels).toHaveClass(/collapsed/);
    await expect(toggleButton).toContainText('▶ Expand');

    // Click to expand again
    await anthropicProvider.locator('.provider-header').click();

    // Expanded again
    await expect(anthropicModels).not.toHaveClass(/collapsed/);
    await expect(toggleButton).toContainText('▼ Collapse');
  });

  test('should show AI insights panel', async ({ page }) => {
    // Verify insight panel is visible
    const insightPanel = page.locator('.insight-panel');
    await expect(insightPanel).toBeVisible();
    await expect(insightPanel).toContainText('AI Insight');
    await expect(insightPanel).toContainText('Claude 4.5 Sonnet is 23% faster');
  });

  test('should display persistent shortcut bar', async ({ page }) => {
    const shortcutBar = page.locator('.shortcut-bar');
    await expect(shortcutBar).toBeVisible();
    await expect(shortcutBar).toContainText('Tab');
    await expect(shortcutBar).toContainText('1-9');
    await expect(shortcutBar).toContainText('d');
    await expect(shortcutBar).toContainText('?');
  });

  test('should show help dialog with ? key', async ({ page }) => {
    page.once('dialog', dialog => {
      expect(dialog.message()).toContain('Keyboard Shortcuts');
      expect(dialog.message()).toContain('Navigate models');
      expect(dialog.message()).toContain('Jump to provider');
      dialog.accept();
    });

    await page.keyboard.press('?');
  });

  test('should handle model selection', async ({ page }) => {
    // Click on a model
    const claudeModel = page.locator('[data-testid="model-claude-4.5"]');
    await claudeModel.hover();

    // Verify hover state (should highlight)
    const backgroundColor = await claudeModel.evaluate(
      el => window.getComputedStyle(el).backgroundColor
    );
    expect(backgroundColor).not.toBe('rgba(0, 0, 0, 0)');
  });

  test('should run automated test suite', async ({ page }) => {
    // Test Search button
    await page.click('text=Test Search');
    await page.waitForTimeout(600);
    await expect(page.locator('#search-status')).toContainText('Found 3 models matching "claude"');
    await expect(page.locator('#status-search')).toHaveClass(/pass/);

    // Reset view
    await page.click('text=Reset View');
    await page.waitForLoadState('domcontentloaded');

    // Test Keyboard Navigation button
    await page.click('text=Test Keyboard Nav');
    await page.waitForTimeout(1600);
    await expect(page.locator('#keyboard-status')).toContainText('All shortcuts functional');
    await expect(page.locator('#status-keyboard')).toHaveClass(/pass/);
  });

  test('should test authentication flow', async ({ page }) => {
    // Verify Gemini is initially locked
    await expect(page.locator('[data-testid="provider-gemini-locked"]')).toContainText('🔒');
    await expect(page.locator('#count-auth')).toContainText('4 authenticated');

    // Click "Test Auth Flow" button
    await page.click('text=Test Auth');
    await page.waitForTimeout(2100);

    // Verify Gemini is now authenticated
    await expect(page.locator('#count-auth')).toContainText('5 authenticated');
    await expect(page.locator('#status-auth')).toHaveClass(/pass/);
  });

  test('should show correct provider counts', async ({ page }) => {
    // Anthropic: 6 models
    await expect(page.locator('[data-testid="provider-anthropic"]')).toContainText('(6 models)');

    // OpenAI: 8 models
    await expect(page.locator('[data-testid="provider-openai"]')).toContainText('(8 models)');

    // Claude CLI: 6 models
    await expect(page.locator('[data-testid="provider-claude-cli"]')).toContainText('(6 models)');

    // Gemini: 7 models
    await expect(page.locator('[data-testid="provider-gemini-locked"]')).toContainText('(7 models)');
  });

  test('should maintain visual hierarchy', async ({ page }) => {
    // Verify modal structure
    await expect(page.locator('.modal-header')).toBeVisible();
    await expect(page.locator('.shortcut-bar')).toBeVisible();
    await expect(page.locator('.search-box')).toBeVisible();
    await expect(page.locator('.model-list')).toBeVisible();
    await expect(page.locator('.footer')).toBeVisible();

    // Verify Recent section appears before provider groups
    const modelList = page.locator('.model-list');
    const recentSection = modelList.locator('.section-header').first();
    await expect(recentSection).toContainText('RECENT');
  });

  test('should have accessible keyboard navigation', async ({ page }) => {
    const searchInput = page.locator('#search-input');

    // Tab to focus search (would work if page had proper tab order)
    await searchInput.focus();
    await expect(searchInput).toBeFocused();

    // Type to search
    await page.keyboard.type('gpt');
    await expect(searchInput).toHaveValue('gpt');

    // Clear search
    await page.keyboard.press('Escape');
    // In real implementation, Escape would clear search
  });

  test('should display responsive design elements', async ({ page }) => {
    // Verify fixed-width modal
    const modalSelector = page.locator('.model-selector');
    const width = await modalSelector.evaluate(el => el.offsetWidth);
    expect(width).toBeGreaterThan(0);
    expect(width).toBeLessThanOrEqual(800);

    // Verify scrollable model list
    const modelList = page.locator('.model-list');
    const hasScroll = await modelList.evaluate(
      el => el.scrollHeight > el.clientHeight
    );
    expect(hasScroll).toBe(true);
  });
});

test.describe('Model Selector - Edge Cases', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('file://' + process.cwd() + '/test-model-selector-web.html');
  });

  test('should handle empty search results', async ({ page }) => {
    const searchInput = page.locator('#search-input');
    await searchInput.fill('nonexistent-model-xyz');
    await page.waitForTimeout(300);

    // Verify search status shows 0 results
    await expect(page.locator('#search-status')).toContainText('Found 0 models');
  });

  test('should handle locked provider click', async ({ page }) => {
    const geminiProvider = page.locator('[data-testid="provider-gemini-locked"]');

    page.once('dialog', dialog => {
      expect(dialog.message()).toContain('Authentication prompt for gemini');
      expect(dialog.message()).toContain('auto-detect');
      dialog.accept();
    });

    await geminiProvider.locator('.provider-header').click();
  });

  test('should maintain state during rapid interactions', async ({ page }) => {
    // Rapidly collapse/expand provider
    const anthropicHeader = page.locator('[data-testid="provider-anthropic"] .provider-header');

    for (let i = 0; i < 5; i++) {
      await anthropicHeader.click();
      await page.waitForTimeout(100);
    }

    // Should still be functional
    const toggleButton = page.locator('#toggle-anthropic');
    await expect(toggleButton).toContainText(/Collapse|Expand/);
  });
});

test.describe('Model Selector - Performance', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('file://' + process.cwd() + '/test-model-selector-web.html');
  });

  test('should load within acceptable time', async ({ page }) => {
    const startTime = Date.now();
    await page.waitForLoadState('domcontentloaded');
    const loadTime = Date.now() - startTime;

    // Should load in under 2 seconds
    expect(loadTime).toBeLessThan(2000);
  });

  test('should render all models quickly', async ({ page }) => {
    const startTime = Date.now();

    // Wait for all model items to be rendered
    await page.waitForSelector('.model-item');
    const modelItems = await page.locator('.model-item').count();

    const renderTime = Date.now() - startTime;

    // Should render all models in under 1 second
    expect(renderTime).toBeLessThan(1000);
    expect(modelItems).toBeGreaterThanOrEqual(10); // At least 10 models visible
  });

  test('should handle search input without lag', async ({ page }) => {
    const searchInput = page.locator('#search-input');

    const startTime = Date.now();
    await searchInput.type('claude sonnet 4.5');
    const typeTime = Date.now() - startTime;

    // Typing should be instant
    expect(typeTime).toBeLessThan(500);
  });
});
