import { test, expect } from '@playwright/test';

test('projects dashboard renders', async ({ page }) => {
  await page.goto('/#/');
  await expect(page.locator('.app-layout, .login-container').first()).toBeVisible();
});

