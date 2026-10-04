import { useQuery } from "@tanstack/react-query";
import { NavLink, Outlet } from "react-router-dom";
import { fetchHealth } from "../api/resources";
import { healthKey } from "../app/query";
import { useSessionQuery } from "../app/session";
import { useHotkeys } from "../keyboard/useHotkeys";
import { useUiStore } from "../stores/ui";
import { CommandPalette } from "./CommandPalette";
import { HelpDialog } from "./HelpDialog";
import { Mark } from "./Mark";

const links = [
  { to: "/overview", label: "Overview", testId: "nav-overview", chord: "g o" },
  { to: "/clusters", label: "Clusters", testId: "nav-clusters", chord: "g c" },
  { to: "/plugins", label: "Plugins", testId: "nav-plugins", chord: "g p" },
  { to: "/settings", label: "Settings", testId: "nav-settings", chord: "g s" },
];

export function AppShell() {
  const session = useSessionQuery();
  const health = useQuery({ queryKey: healthKey, queryFn: fetchHealth });
  const setPaletteOpen = useUiStore((state) => state.setPaletteOpen);
  const setHelpOpen = useUiStore((state) => state.setHelpOpen);
  useHotkeys();
  const principal = session.data?.principal;

  return (
    <div className="km-shell">
      <a className="km-skip" href="#main">
        Skip to content
      </a>
      <header className="km-topbar">
        <NavLink to="/overview" className="km-brand">
          <Mark />
          <span className="km-word">
            Kube<span>Mv</span>
          </span>
        </NavLink>
        <span className="km-phase">Phase 0 · Foundation</span>
        <button
          type="button"
          className="km-palette-trigger"
          data-testid="palette-trigger"
          onClick={() => setPaletteOpen(true)}
        >
          <span>Command</span>
          <kbd className="km-kbd">Ctrl K</kbd>
        </button>
        {principal ? (
          <div className="km-user" data-testid="current-user">
            <span>{principal.username}</span>
          </div>
        ) : null}
      </header>
      <nav className="km-nav" aria-label="Primary">
        {links.map((link) => (
          <NavLink key={link.to} to={link.to} data-testid={link.testId}>
            <span>{link.label}</span>
            <kbd>{link.chord}</kbd>
          </NavLink>
        ))}
      </nav>
      <main className="km-main" id="main">
        <Outlet />
      </main>
      <footer className="km-status">
        <span data-testid="api-status">
          API {health.isSuccess ? "reachable" : health.isError ? "unreachable" : "checking"}
          {health.data ? ` · v${health.data.version}` : ""}
        </span>
        <span>Kubernetes API not connected</span>
        <button type="button" className="km-btn-ghost" onClick={() => setHelpOpen(true)}>
          Shortcuts
        </button>
      </footer>
      <CommandPalette />
      <HelpDialog />
    </div>
  );
}
