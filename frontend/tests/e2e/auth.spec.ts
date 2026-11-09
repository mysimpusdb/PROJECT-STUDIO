import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/');
  // Basic check to see if the page loads and has a title
  await expect(page).toHaveTitle(/PROJECT-STUDIO|Vite App/i);
});

test('login page renders correctly', async ({ page }) => {
  await page.goto('/login');
  await expect(page.locator('text=Login to PROJECT-STUDIO')).toBeVisible();
});
