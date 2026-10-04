import { create } from "zustand";

export type ThemeName = "dark" | "light";

const storageKey = "kubemv.theme";

function readTheme(): ThemeName {
  try {
    const stored = localStorage.getItem(storageKey);
    if (stored === "light" || stored === "dark") {
      return stored;
    }
  } catch {
    // Private mode can reject storage. The default theme still applies.
  }
  return "dark";
}

function applyTheme(theme: ThemeName) {
  document.documentElement.dataset.theme = theme;
  try {
    localStorage.setItem(storageKey, theme);
  } catch {
    // The in-memory theme still updates for this session.
  }
}

type UiState = {
  theme: ThemeName;
  paletteOpen: boolean;
  helpOpen: boolean;
  setTheme: (theme: ThemeName) => void;
  toggleTheme: () => void;
  setPaletteOpen: (open: boolean) => void;
  setHelpOpen: (open: boolean) => void;
};

export const useUiStore = create<UiState>((set, get) => ({
  theme: readTheme(),
  paletteOpen: false,
  helpOpen: false,
  setTheme: (theme) => {
    applyTheme(theme);
    set({ theme });
  },
  toggleTheme: () => {
    const next = get().theme === "dark" ? "light" : "dark";
    applyTheme(next);
    set({ theme: next });
  },
  setPaletteOpen: (paletteOpen) => set({ paletteOpen, helpOpen: paletteOpen ? false : get().helpOpen }),
  setHelpOpen: (helpOpen) => set({ helpOpen, paletteOpen: helpOpen ? false : get().paletteOpen }),
}));
