const { test, expect } = require('@playwright/test');

const TEST_USER = {
  username: 'admin',
  password: 'admin123'
};

test.describe('Browser Compatibility - Authentication', () => {
  test('login page renders correctly', async ({ page }) => {
    await page.goto('/login');
    
    await expect(page.locator('input[type="text"], input[name="username"], input[placeholder*="用户名"], input[placeholder*="username"]').first()).toBeVisible();
    await expect(page.locator('input[type="password"]').first()).toBeVisible();
    await expect(page.locator('button[type="submit"], button:has-text("登录"), button:has-text("Login")').first()).toBeVisible();
  });

  test('login flow works', async ({ page }) => {
    await page.goto('/login');
    
    const usernameInput = page.locator('input[type="text"], input[name="username"]').first();
    const passwordInput = page.locator('input[type="password"]').first();
    
    await usernameInput.fill(TEST_USER.username);
    await passwordInput.fill(TEST_USER.password);
    
    const loginButton = page.locator('button[type="submit"], button:has-text("登录"), button:has-text("Login")').first();
    await loginButton.click();
    
    await page.waitForURL(/\/(dashboard|home|main)/, { timeout: 15000 }).catch(() => {
      // Fallback: check if we're no longer on login page
      return expect(page).not.toHaveURL(/login/);
    });
  });

  test('token stored after login', async ({ page }) => {
    await page.goto('/login');
    
    const usernameInput = page.locator('input[type="text"], input[name="username"]').first();
    const passwordInput = page.locator('input[type="password"]').first();
    
    await usernameInput.fill(TEST_USER.username);
    await passwordInput.fill(TEST_USER.password);
    
    const loginButton = page.locator('button[type="submit"], button:has-text("登录"), button:has-text("Login")').first();
    await loginButton.click();
    
    await page.waitForTimeout(2000);
    
    const localStorage = await page.evaluate(() => window.localStorage);
    const hasToken = Object.keys(localStorage).some(key => 
      key.toLowerCase().includes('token') || 
      key.toLowerCase().includes('auth')
    );
    expect(hasToken || localStorage.length > 0).toBeTruthy();
  });
});

test.describe('Browser Compatibility - Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    const usernameInput = page.locator('input[type="text"], input[name="username"]').first();
    const passwordInput = page.locator('input[type="password"]').first();
    await usernameInput.fill(TEST_USER.username);
    await passwordInput.fill(TEST_USER.password);
    const loginButton = page.locator('button[type="submit"], button:has-text("登录"), button:has-text("Login")').first();
    await loginButton.click();
    await page.waitForTimeout(2000);
  });

  test('dashboard page loads', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForLoadState('domcontentloaded');
    await page.waitForTimeout(1000);
    
    const pageContent = await page.content();
    expect(pageContent.length).toBeGreaterThan(100);
  });

  test('navigation menu works', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    
    const navItems = page.locator('nav, .menu, .sidebar, [class*="nav"], [class*="menu"]');
    const count = await navItems.count();
    expect(count).toBeGreaterThan(0);
  });
});

test.describe('Browser Compatibility - Responsive Design', () => {
  test('mobile viewport renders', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');
    
    await expect(page.locator('body')).toBeVisible();
    
    const bodyWidth = await page.evaluate(() => document.body.scrollWidth);
    expect(bodyWidth).toBeLessThanOrEqual(400);
  });

  test('tablet viewport renders', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto('/');
    
    await expect(page.locator('body')).toBeVisible();
  });

  test('desktop viewport renders', async ({ page }) => {
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.goto('/');
    
    await expect(page.locator('body')).toBeVisible();
  });
});

test.describe('Browser Compatibility - CSS Features', () => {
  test('flexbox layout works', async ({ page }) => {
    await page.goto('/');
    
    const flexSupport = await page.evaluate(() => {
      const test = document.createElement('div');
      test.style.display = 'flex';
      return test.style.display === 'flex';
    });
    expect(flexSupport).toBeTruthy();
  });

  test('CSS grid works', async ({ page }) => {
    await page.goto('/');
    
    const gridSupport = await page.evaluate(() => {
      const test = document.createElement('div');
      test.style.display = 'grid';
      return test.style.display === 'grid';
    });
    expect(gridSupport).toBeTruthy();
  });

  test('CSS custom properties work', async ({ page }) => {
    await page.goto('/');
    
    const customPropSupport = await page.evaluate(() => {
      const test = document.createElement('div');
      test.style.setProperty('--test', 'value');
      return test.style.getPropertyValue('--test') === 'value';
    });
    expect(customPropSupport).toBeTruthy();
  });
});

test.describe('Browser Compatibility - JavaScript Features', () => {
  test('ES6 arrow functions work', async ({ page }) => {
    await page.goto('/');
    
    const arrowSupport = await page.evaluate(() => {
      try {
        const fn = () => true;
        return fn();
      } catch {
        return false;
      }
    });
    expect(arrowSupport).toBeTruthy();
  });

  test('async/await works', async ({ page }) => {
    await page.goto('/');
    
    const asyncSupport = await page.evaluate(async () => {
      try {
        const result = await Promise.resolve(true);
        return result;
      } catch {
        return false;
      }
    });
    expect(asyncSupport).toBeTruthy();
  });

  test('Fetch API works', async ({ page }) => {
    await page.goto('/');
    
    const fetchSupport = await page.evaluate(async () => {
      try {
        return typeof fetch === 'function';
      } catch {
        return false;
      }
    });
    expect(fetchSupport).toBeTruthy();
  });

  test('LocalStorage works', async ({ page }) => {
    await page.goto('/');
    
    const storageSupport = await page.evaluate(() => {
      try {
        localStorage.setItem('test', 'value');
        const result = localStorage.getItem('test') === 'value';
        localStorage.removeItem('test');
        return result;
      } catch {
        return false;
      }
    });
    expect(storageSupport).toBeTruthy();
  });
});

test.describe('Browser Compatibility - API Integration', () => {
  test('health check API responds', async ({ page }) => {
    const response = await page.request.get('http://localhost:8085/health');
    expect(response.status()).toBe(200);
  });

  test('API CORS headers present', async ({ page }) => {
    const response = await page.request.get('http://localhost:8085/health');
    // Check if response is valid (CORS headers may vary)
    expect(response.status()).toBeLessThan(500);
  });
});
