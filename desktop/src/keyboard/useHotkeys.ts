import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { navigationChord } from "./commands";
import { useUiStore } from "../stores/ui";

function isTypingTarget(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) {
    return false;
  }
  const tag = target.tagName;
  return tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT" || target.isContentEditable;
}

export function useHotkeys() {
  const navigate = useNavigate();

  useEffect(() => {
    let chord = false;
    let timer = 0;
    const onKeyDown = (event: KeyboardEvent) => {
      const ui = useUiStore.getState();
      const key = event.key.toLowerCase();
      if ((event.ctrlKey || event.metaKey) && key === "k") {
        event.preventDefault();
        ui.setPaletteOpen(!ui.paletteOpen);
        return;
      }
      if (event.key === "Escape") {
        ui.setPaletteOpen(false);
        ui.setHelpOpen(false);
        return;
      }
      if (ui.paletteOpen || ui.helpOpen || isTypingTarget(event.target)) {
        return;
      }
      if (event.key === "?") {
        ui.setHelpOpen(true);
        return;
      }
      if (key === "g" && !event.ctrlKey && !event.metaKey && !event.altKey) {
        chord = true;
        window.clearTimeout(timer);
        timer = window.setTimeout(() => {
          chord = false;
        }, 800);
        return;
      }
      if (chord) {
        chord = false;
        const path = navigationChord[key];
        if (path) {
          navigate(path);
        }
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => {
      window.clearTimeout(timer);
      window.removeEventListener("keydown", onKeyDown);
    };
  }, [navigate]);
}
