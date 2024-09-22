import { defineConfig, devices } from '@playwright/test';

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// require('dotenv').config();

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './tests',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  workers: 1,

  reporter: [
    ['html'],
    [process.env.CI ? 'github' : 'list']
  ],

  use: {
    trace: 'on-first-retry',
    baseURL: process.env.E2E_BASE_URL || 'https://localhost:8443',
    ignoreHTTPSErrors: true
  },

  projects: [
    {
      name: 'initial.setup',
      testMatch: /initial\.setup\.ts/,
    },

    {
      name: 'auth.setup',
      testMatch: /auth\.setup\.ts/,
      dependencies: ['initial.setup']
    },

    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
      // testMatch: /.*\/spec\.ts/,
      dependencies: ['auth.setup']
    },
  ]
});