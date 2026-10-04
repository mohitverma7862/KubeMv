import { expect, test } from "@playwright/test";

async function signIn(page: import("@playwright/test").Page) {
  await page.goto("/login");
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("foundation-test-password");
  await page.getByRole("button", { name: "Sign in" }).click();
  await expect(page.getByRole("heading", { name: "Foundation" })).toBeVisible();
}

test("rejects a bad password without leaving the login screen", async ({ page }) => {
  await page.goto("/login");
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("not-the-password");
  await page.getByRole("button", { name: "Sign in" }).click();
  await expect(page.getByTestId("login-error")).toContainText("Invalid username or password");
  await expect(page.getByTestId("login-error")).not.toContainText("not-the-password");
});

test("operator can register metadata, use the keyboard, and sign out", async ({ page }) => {
  await page.goto("/clusters");
  await expect(page).toHaveURL(/\/login$/);
  await signIn(page);
  await expect(page.getByTestId("current-user")).toHaveText("admin");
  await expect(page.getByTestId("stat-api")).toHaveText("Up");
  await expect(page.getByTestId("stat-plugins")).toHaveText("0");

  await page.keyboard.press("Control+k");
  await expect(page.getByTestId("command-palette")).toBeVisible();
  await page.getByRole("textbox", { name: "Command" }).fill("clusters");
  await page.keyboard.press("Enter");
  await expect(page.getByRole("heading", { name: "Clusters" })).toBeVisible();

  await page.getByLabel("Name").fill("prod-eks");
  await page.getByLabel("Provider").selectOption("eks");
  await page.getByLabel("Context").fill("prod");
  await page.getByLabel("Server-side reference").fill("secret://clusters/prod");
  await page.getByRole("button", { name: "Register cluster" }).click();
  const list = page.getByTestId("cluster-list");
  await expect(list).toContainText("prod-eks");
  await expect(list).toContainText("Not connected");
  await expect(list).not.toContainText("secret://");
  await expect(page.getByLabel("Server-side reference")).toHaveValue("");

  await page.getByRole("button", { name: "Remove" }).click();
  await page.getByTestId("confirm-remove").click();
  await expect(list).toContainText("No clusters are registered");

  await page.getByTestId("nav-plugins").click();
  await expect(page.getByText("No plugins are loaded.")).toBeVisible();

  await page.keyboard.press("?");
  await expect(page.getByTestId("help-dialog")).toBeVisible();
  await page.keyboard.press("Escape");
  await expect(page.getByTestId("help-dialog")).toBeHidden();

  await page.keyboard.press("g");
  await page.keyboard.press("s");
  await expect(page.getByRole("heading", { name: "Settings" })).toBeVisible();
  await page.getByTestId("theme-light").click();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
  await page.reload();
  await expect(page.getByRole("heading", { name: "Settings" })).toBeVisible();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");

  await page.getByTestId("sign-out").click();
  await expect(page.getByTestId("login-form")).toBeVisible();
  await page.goto("/overview");
  await expect(page).toHaveURL(/\/login$/);
});

test("narrow layout keeps navigation available", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await signIn(page);
  await expect(page.getByTestId("nav-clusters")).toBeVisible();
  await page.getByTestId("nav-clusters").click();
  await expect(page.getByRole("heading", { name: "Clusters" })).toBeVisible();
  await expect(page.getByTestId("nav-settings")).toBeVisible();
  const overflow = await page.evaluate(
    () => document.documentElement.scrollWidth > document.documentElement.clientWidth + 8,
  );
  expect(overflow).toBe(false);
});
