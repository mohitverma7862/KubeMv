import { defineConfig, devices } from "@playwright/test";

function processEnv(overrides: Record<string, string>): Record<string, string> {
  const env: Record<string, string> = {};
  for (const [key, value] of Object.entries(process.env)) {
    if (value !== undefined) {
      env[key] = value;
    }
  }
  return { ...env, ...overrides };
}

const apiEnv = processEnv({
  KUBEMV_BOOTSTRAP_USERNAME: "admin",
  KUBEMV_BOOTSTRAP_PASSWORD: "foundation-test-password",
  KUBEMV_ALLOWED_ORIGINS: "http://127.0.0.1:1420,http://localhost:1420",
});

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false,
  workers: 1,
  use: {
    ...devices["Desktop Chrome"],
    baseURL: "http://127.0.0.1:1420",
    viewport: { width: 1280, height: 800 },
  },
  webServer: [
    {
      command: "go run ./cmd/kubemv",
      cwd: "../backend",
      url: "http://127.0.0.1:8787/api/v1/health",
      reuseExistingServer: false,
      timeout: 120_000,
      env: apiEnv,
    },
    {
      command: "npm run dev",
      url: "http://127.0.0.1:1420",
      reuseExistingServer: false,
      timeout: 120_000,
      env: processEnv({}),
    },
  ],
});
