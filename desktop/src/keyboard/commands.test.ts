import { describe, expect, it } from "vitest";
import { filterCommands, type Command } from "./commands";

const commands: Command[] = [
  { id: "clusters", title: "Go to Clusters", hint: "Cluster registry", keywords: "/clusters" },
  { id: "theme", title: "Toggle theme", hint: "Dark or light", keywords: "appearance" },
];

describe("filterCommands", () => {
  it("returns every command for an empty query", () => {
    expect(filterCommands(commands, "  ")).toHaveLength(2);
  });

  it("matches all terms", () => {
    expect(filterCommands(commands, "cluster registry").map((command) => command.id)).toEqual([
      "clusters",
    ]);
    expect(filterCommands(commands, "theme missing")).toHaveLength(0);
  });
});
