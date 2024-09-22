import { test as setup, expect } from '@playwright/test'
import { credsMap } from './auth.config'

const defaultCreds = credsMap.default
const adminCreds = credsMap.admin

const badPassword = '1234'

setup.use({
    ignoreHTTPSErrors: true
})

setup('Wipe Database', async ({ request }) => {
    const wipeReq = await request.post('/api/v1/e2e/wipe', {
        headers: {
            'X-E2E-Key': 'authKeyForE2E'
        }
    })
    await expect(wipeReq.status()).toBe(200)
})

setup('First time admin flow', async ({ page }) => {
    await page.goto('/')

    await expect(page.getByRole('heading', { name: 'Login to Phatcrack' })).toBeVisible()

    await page.getByPlaceholder('john.doe').fill(defaultCreds.username)
    await page.getByPlaceholder('hunter2').fill(defaultCreds.password)
    await page.getByRole('button').click()

    // Weak password
    await page.locator('div').filter({ hasText: /^New Password$/ }).getByPlaceholder('hunter2').fill(badPassword)
    await page.getByRole('button', { name: 'Change Password' }).click()
    await expect(page.getByText('Failed to change temporary password')).toBeVisible()

    // Strong password
    await page.locator('div').filter({ hasText: /^New Password$/ }).getByPlaceholder('hunter2').fill(adminCreds.password)
    await page.getByRole('button', { name: 'Change Password' }).click()
    await expect(page.getByText('Success')).toBeVisible()

    await expect(page.getByText('Welcome ' + adminCreds.username)).toBeVisible()

})