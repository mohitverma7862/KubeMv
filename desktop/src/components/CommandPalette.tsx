import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { logout } from "../api/resources";
import { sessionKey } from "../app/query";
import { useSessionQuery } from "../app/session";
import { filterCommands, type Command } from "../keyboard/commands";
import { useUiStore } from "../stores/ui";

type PaletteCommand = Command & { run: () => void };

export function CommandPalette() {
  const open = useUiStore((state) => state.paletteOpen);
  const setOpen = useUiStore((state) => state.setPaletteOpen);
  const toggleTheme = useUiStore((state) => state.toggleTheme);
  const navigate = useNavigate();
  const session = useSessionQuery();
  const queryClient = useQueryClient();
  const [query, setQuery] = useState("");
  const [active, setActive] = useState(0);

  const commands = useMemo<PaletteCommand[]>(() => {
    const go = (id: string, title: string, path: string, hint: string): PaletteCommand => ({
      id,
      title,
      hint,
      keywords: path,
      run: () => navigate(path),
    });
    return [
      go("overview", "Go to Overview", "/overview", "Foundation status"),
      go("clusters", "Go to Clusters", "/clusters", "Cluster registry"),
      go("plugins", "Go to Plugins", "/plugins", "Plugin manifests"),
      go("settings", "Go to Settings", "/settings", "Theme and shortcuts"),
      {
        id: "theme",
        title: "Toggle theme",
        hint: "Dark or light",
        keywords: "appearance",
        run: toggleTheme,
      },
      {
        id: "signout",
        title: "Sign out",
        hint: "End this session",
        keywords: "logout",
        run: () => {
          const token = session.data?.csrfToken;
          if (!token) {
            return;
          }
          void logout(token).finally(() => {
            queryClient.setQueryData(sessionKey, null);
            navigate("/login", { replace: true });
          });
        },
      },
    ];
  }, [navigate, queryClient, session.data?.csrfToken, toggleTheme]);

  const visible = filterCommands(commands, query);
  const activeCommand = commands.find((command) => command.id === visible[active]?.id);

  useEffect(() => {
    if (!open) {
      setQuery("");
      setActive(0);
    }
  }, [open]);

  useEffect(() => {
    setActive(0);
  }, [query]);

  if (!open) {
    return null;
  }

  return (
    <div className="km-modal-back" onMouseDown={() => setOpen(false)}>
      <div
        className="km-modal"
        role="dialog"
        aria-modal="true"
        aria-label="Command palette"
        data-testid="command-palette"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <input
          className="km-command-input"
          autoFocus
          placeholder="Type a command"
          value={query}
          aria-label="Command"
          onChange={(event) => setQuery(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Escape") {
              setOpen(false);
            } else if (event.key === "ArrowDown") {
              event.preventDefault();
              setActive((index) => Math.min(index + 1, Math.max(visible.length - 1, 0)));
            } else if (event.key === "ArrowUp") {
              event.preventDefault();
              setActive((index) => Math.max(index - 1, 0));
            } else if (event.key === "Enter" && activeCommand) {
              activeCommand.run();
              setOpen(false);
            }
          }}
        />
        <ul className="km-command-list">
          {visible.length === 0 ? <li className="km-empty">No matching commands</li> : null}
          {visible.map((command, index) => (
            <li key={command.id}>
              <button
                type="button"
                data-active={index === active}
                onMouseEnter={() => setActive(index)}
                onClick={() => {
                  commands.find((item) => item.id === command.id)?.run();
                  setOpen(false);
                }}
              >
                <span>{command.title}</span>
                <span className="km-muted">{command.hint}</span>
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
