import { test, expect } from '@playwright/test';
import { credsMap, getAuthFilePath } from './auth.config';

const adminUsername = credsMap.admin.username

const aliceCreds = credsMap.alice
const bobCreds = credsMap.bob
const carolCreds = credsMap.carol

const tmpPassword = 'Th1s1s@t3mpP@ss123!'

test.describe(() => {
    test('Add user with immediate password (alice) and do first login', async ({ browser }) => {
        const adminContext = await browser.newContext({ storageState: getAuthFilePath(adminUsername) })
        const adminPage = await adminContext.newPage()

        // Go to general settings
        await adminPage.goto('/admin/general')
        await expect(adminPage.getByRole('heading', { name: 'General Settings' })).toBeVisible()

        // Uncheck require password change
        await adminPage.locator('label').filter({ hasText: 'Require password change on first login' }).locator('input[type=checkbox]').setChecked(false)
        await adminPage.locator('.card').getByRole('button', { name: 'Save' }).click()

        // Go to user page
        await adminPage.goto('/admin/users')
        await expect(adminPage.getByRole('heading', { name: 'User Management' })).toBeVisible()

        const createUserBtn = adminPage.getByRole('button', { name: 'Create User' })

        await expect(createUserBtn).toBeVisible()
        await createUserBtn.click()

        await expect(adminPage.locator('.modal-open').getByRole('heading', { name: 'Create a new user' })).toBeVisible()

        await adminPage.locator('.modal-open').filter({ hasText: 'Username' }).getByPlaceholder('j.smith').fill(aliceCreds.username)
        await adminPage.locator('.modal-open').filter({ hasText: 'Password' }).getByPlaceholder('hunter2').fill(aliceCreds.password)
        await adminPage.locator('.modal-open').filter({ hasText: 'Create a new user' }).getByRole('button', { name: 'Create' }).click()

        await expect(adminPage.getByText('Created new user')).toBeVisible()


        const tmpLoginContext = await browser.newContext({})
        const tmpLoginPage = await tmpLoginContext.newPage()

        await tmpLoginPage.goto('/')

        await expect(tmpLoginPage.getByRole('heading', { name: 'Login to Phatcrack' })).toBeVisible()

        await tmpLoginPage.getByPlaceholder('john.doe').fill(aliceCreds.username)
        await tmpLoginPage.getByPlaceholder('hunter2').fill(aliceCreds.password)
        await tmpLoginPage.getByRole('button').click()

        await expect(tmpLoginPage.getByText('Welcome ' + aliceCreds.username)).toBeVisible()

        return
    })

})

test.describe(() => {
    test('Add user with mandatory password change (bob) and do first login', async ({ browser }) => {
        const adminContext = await browser.newContext({ storageState: getAuthFilePath(adminUsername) })
        const adminPage = await adminContext.newPage()

        // Go to general settings
        await adminPage.goto('/admin/general')
        await expect(adminPage.getByRole('heading', { name: 'General Settings' })).toBeVisible()

        // Uncheck require password change
        await adminPage.locator('label').filter({ hasText: 'Require password change on first login' }).locator('input[type=checkbox]').setChecked(true)
        await adminPage.locator('.card').getByRole('button', { name: 'Save' }).click()

        // Go to user page
        await adminPage.goto('/admin/users')
        await expect(adminPage.getByRole('heading', { name: 'User Management' })).toBeVisible()

        const createUserBtn = adminPage.getByRole('button', { name: 'Create User' })

        await expect(createUserBtn).toBeVisible()
        await createUserBtn.click()

        await expect(adminPage.locator('.modal-open').getByRole('heading', { name: 'Create a new user' })).toBeVisible()

        await adminPage.locator('.modal-open').filter({ hasText: 'Username' }).getByPlaceholder('j.smith').fill(bobCreds.username)
        await adminPage.locator('.modal-open').filter({ hasText: 'Password' }).getByPlaceholder('hunter2').fill(tmpPassword)
        await adminPage.locator('.modal-open').filter({ hasText: 'Create a new user' }).getByRole('button', { name: 'Create' }).click()

        await expect(adminPage.getByText('Created new user')).toBeVisible()

        // Close modal
        await adminPage.locator('.modal-open').getByRole('button', { name: '✕' }).click()
        await expect(adminPage.getByRole('heading', { name: 'Admin' })).toBeVisible()


        const tmpLoginContext = await browser.newContext({})
        const tmpLoginPage = await tmpLoginContext.newPage()

        // Login with temporarily-issued credentials, make sure we're forced to change password
        await tmpLoginPage.goto('/')

        await expect(tmpLoginPage.getByRole('heading', { name: 'Login to Phatcrack' })).toBeVisible()

        await tmpLoginPage.getByPlaceholder('john.doe').fill(bobCreds.username)
        await tmpLoginPage.getByPlaceholder('hunter2').fill(tmpPassword)
        await tmpLoginPage.getByRole('button').click()

        await expect(tmpLoginPage.getByRole('heading', { name: 'Set a new password' })).toBeVisible()
        await tmpLoginPage.locator('div').filter({ hasText: /^New Password$/ }).getByPlaceholder('hunter2').fill(bobCreds.password)
        await tmpLoginPage.getByRole('button', { name: 'Change Password' }).click()
        await expect(tmpLoginPage.getByText('Welcome ' + bobCreds.username)).toBeVisible()


        tmpLoginContext.clearCookies()
        await tmpLoginPage.goto('/')

        await expect(tmpLoginPage.getByRole('heading', { name: 'Login to Phatcrack' })).toBeVisible()

        // Make sure temp password nolonger works
        await tmpLoginPage.getByPlaceholder('john.doe').fill(bobCreds.username)
        await tmpLoginPage.getByPlaceholder('hunter2').fill(tmpPassword)
        await tmpLoginPage.getByRole('button').click()
        await expect(tmpLoginPage.getByText('Invalid credentials')).toBeVisible()

        // Make sure new password works
        await tmpLoginPage.getByPlaceholder('john.doe').fill(bobCreds.username)
        await tmpLoginPage.getByPlaceholder('hunter2').fill(bobCreds.password)
        await tmpLoginPage.getByRole('button').click()
        await expect(tmpLoginPage.getByText('Welcome ' + bobCreds.username)).toBeVisible()

        return
    })

})