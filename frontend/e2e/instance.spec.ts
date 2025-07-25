import { expect, test } from '@playwright/test';

test.describe('Instances', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('Application should have an instance', async ({ page }) => {
    await expect(page.getByRole('link', { name: 'instances' })).toContainText('1instances');
    await page.getByRole('link', { name: 'Flatcar Container Linux' }).click();

    await expect(page.locator('#main')).toContainText('1');

    await page.evaluate(() => window.scrollTo(0, 0));
    await expect(page).toHaveScreenshot('in-app-with-a-node-instance.png');

    await page.getByRole('link', { name: 'Alpha (AMD64)' }).click();

    await page.evaluate(() => window.scrollTo(0, 0));
    await expect(page).toHaveScreenshot('in-group-with-a-node-instance.png');

    await expect(page.locator('#main')).toContainText('See all instances');
    await expect(page.getByLabel('Downloaded: 1 instances')).toContainText('Downloaded');
    await expect(page.locator('#main')).toContainText('100.0%');
    await expect(page.locator('#main')).toContainText('Version Breakdown');
    await expect(page.locator('tbody')).toContainText('4081.2.0');

    await page.getByRole('link', { name: 'See all instances' }).click();

    // Fix the width of columns with date values,
    // because date values can change in this test from execution to execution,
    // therefore the width of columns can also stretch more or less based on these values.
    await page.addStyleTag({
      content: `
      table tr td:first-child,
        table tr td:last-child {
        min-width: 300px !important;
      }
    `,
    });

    // maxDiffPixels set due to displaying date times that can change
    await expect(page).toHaveScreenshot('instances-list.png', {
      mask: [page.locator('tbody tr:first-child td:last-child')],
      maxDiffPixels: 200,
    });

    await expect(page.locator('tbody')).toContainText('2c517ad881474ec6b5ab928df2a7b5f4');
    await expect(page.locator('tbody')).toContainText('Updating: downloaded');
    await expect(page.locator('tbody')).toContainText('4081.2.0');

    await page.locator('tbody tr.MuiTableRow-root').getByRole('button').click();

    // mask elements that are: cells where we can find timedate values, and nebraska version at the
    // bottom
    await expect(page).toHaveScreenshot('instance-history.png', {
      mask: [
        page.locator('//*[contains(text(), "/")]'),
        page.locator('td:has(button):last-of-type'),
        page.locator('#main > div:last-child'),
      ],
      maxDiffPixels: 200,
    });

    await expect(page.locator('#main')).toContainText('Downloaded');
    await expect(page.locator('#main')).toContainText('Downloading');
    await expect(page.locator('#main')).toContainText('Granted');

    const searchInput = page.locator('div[aria-label="Search"]').getByRole('textbox');

    await searchInput.click();
    await searchInput.fill('89');
    await searchInput.press('Enter');

    await expect(page.locator('#main')).toContainText('0/1');

    await searchInput.click();
    await searchInput.fill('4081');
    await searchInput.press('Enter');
    await searchInput.fill('517');
    await searchInput.press('Enter');

    await expect(page.locator('tbody')).toContainText('2c517ad881474ec6b5ab928df2a7b5f4');

    await page.getByRole('link', { name: '2c517ad881474ec6b5ab928df2a7b5f4' }).click();

    // mask elements that are: cells where we can find timedate values,
    // and nebraska version at the bottom
    await expect(page).toHaveScreenshot('instance-info.png', {
      mask: [
        page.locator('//*[contains(text(), "/")]'),
        page.locator('#main > div:last-child'),
        page.locator('text=Last Update Check').locator('xpath=following-sibling::div'),
      ],
      maxDiffPixels: 200,
    });

    await expect(page.getByRole('heading')).toContainText('Instance Information');
    await expect(page.locator('#main')).toContainText('2c517ad881474ec6b5ab928df2a7b5f4');
    await expect(page.locator('tbody')).toContainText('5261.0.0');
  });
});
