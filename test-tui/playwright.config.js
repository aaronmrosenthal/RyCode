const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
  testDir: './tests',
  timeout: 30000,
  use: {
    screenshot: 'on',
    video: 'on',
  },
  reporter: [['html', { outputFolder: 'test-results' }]],
});
