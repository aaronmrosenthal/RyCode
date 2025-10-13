const { test, expect } = require('@playwright/test');
const { Terminal } = require('playwright-terminal');

test('Visual inspection of RyCode TUI', async ({ page }) => {
  // Create terminal instance
  const terminal = new Terminal(page, {
    rows: 40,
    cols: 120,
  });

  await page.goto('about:blank');
  await terminal.load();

  // Start the TUI
  await terminal.write('../bin/rycode\r');

  // Wait for TUI to load
  await page.waitForTimeout(2000);

  // Take screenshot of initial state (home view)
  await page.screenshot({
    path: 'screenshots/home-view.png',
    fullPage: true
  });

  // Type some text to test the cursor
  await terminal.write('test input');
  await page.waitForTimeout(500);

  // Take screenshot with text
  await page.screenshot({
    path: 'screenshots/with-input.png',
    fullPage: true
  });

  // Clear and take final screenshot
  await terminal.write('\u0015'); // Ctrl+U to clear
  await page.waitForTimeout(500);

  await page.screenshot({
    path: 'screenshots/cleared.png',
    fullPage: true
  });

  // Keep browser open for manual inspection
  await page.waitForTimeout(60000);
});
