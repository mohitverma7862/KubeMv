import { useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { logout } from "../api/resources";
import { apiBase } from "../api/client";
import { sessionKey } from "../app/query";
import { useSessionQuery } from "../app/session";
import { useUiStore, type ThemeName } from "../stores/ui";

const shortcuts = [
  ["Ctrl K or ⌘K", "Command palette"],
  ["?", "Shortcut help"],
  ["g then o / c / p / s", "Overview, Clusters, Plugins, Settings"],
  ["Esc", "Close overlays"],
];

export function SettingsPage() {
  const theme = useUiStore((state) => state.theme);
  const setTheme = useUiStore((state) => state.setTheme);
  const session = useSessionQuery();
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const choose = (next: ThemeName) => () => setTheme(next);

  return (
    <section>
      <header className="km-page-head">
        <div>
          <p className="km-kicker">Preferences</p>
          <h1>Settings</h1>
          <p className="km-lead">Theme and the shortcuts that exist in this build. Hotkeys are fixed until a later configuration phase.</p>
        </div>
      </header>
      <div className="km-split">
        <article className="km-panel">
          <header>
            <h2>Appearance</h2>
          </header>
          <div className="km-theme-choice">
            <button type="button" className="km-btn-ghost" aria-pressed={theme === "dark"} data-testid="theme-dark" onClick={choose("dark")}>
              Dark
            </button>
            <button type="button" className="km-btn-ghost" aria-pressed={theme === "light"} data-testid="theme-light" onClick={choose("light")}>
              Light
            </button>
          </div>
          <p className="km-help" style={{ marginTop: 12 }}>
            The choice is stored in this browser as a theme name only.
          </p>
          <header style={{ marginTop: 18 }}>
            <h2>Keyboard</h2>
          </header>
          <dl className="km-shortcut-list">
            {shortcuts.map(([keys, description]) => (
              <span key={keys} style={{ display: "contents" }}>
                <dt>{keys}</dt>
                <dd>{description}</dd>
              </span>
            ))}
          </dl>
        </article>
        <article className="km-panel">
          <header>
            <h2>Session</h2>
          </header>
          <p className="km-muted">API {apiBase()}</p>
          <p>Signed in as {session.data?.principal.username}.</p>
          <button
            type="button"
            className="km-btn-danger"
            data-testid="sign-out"
            onClick={() => {
              const token = session.data?.csrfToken;
              if (!token) {
                return;
              }
              void logout(token).finally(() => {
                queryClient.setQueryData(sessionKey, null);
                navigate("/login", { replace: true });
              });
            }}
          >
            Sign out
          </button>
        </article>
      </div>
    </section>
  );
}
