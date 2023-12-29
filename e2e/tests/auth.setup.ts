import { test as setup, expect } from '@playwright/test';
import { credsMap, getAuthFilePath } from './auth.config';

async function doAuth({ page, creds }) {
    await page.goto('https://localhost:8443')

    await page.getByPlaceholder('john.doe').fill(creds.username)
    await page.getByPlaceholder('hunter2').fill(creds.password)
    await page.getByRole('button').click()

    await expect(page.getByText('Welcome ' + creds.username)).toBeVisible()

    await page.context().storageState({ path: getAuthFilePath(creds.username) })
}

setup.use({
    ignoreHTTPSErrors: true
})

setup('authenticate as admin', async ({ page }) => {
    await doAuth({ page, creds: credsMap.admin })
})

