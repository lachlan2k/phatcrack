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
  /* Run tests in files in parallel */
  fullyParallel: true,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */

  reporter: [
    'html',
    [process.env.CI ? 'github' : 'list']
  ],

  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    // baseURL: 'http://127.0.0.1:3000',

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
    baseURL: process.env.E2E_BASE_URL || 'https://localhost:8443',
    ignoreHTTPSErrors: true
  },

  /* Configure projects for major browsers */
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
  ],

  /* Run your local dev server before starting the tests */
  // webServer: {
  //   command: 'npm run start',
  //   url: 'http://127.0.0.1:3000',
  //   reuseExistingServer: !process.env.CI,
  // },
});