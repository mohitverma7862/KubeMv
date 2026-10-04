export type Command = {
  id: string;
  title: string;
  hint: string;
  keywords: string;
};

export function filterCommands(commands: Command[], query: string): Command[] {
  const terms = query.trim().toLowerCase().split(/\s+/).filter(Boolean);
  if (terms.length === 0) {
    return commands;
  }
  return commands.filter((command) => {
    const haystack = `${command.title} ${command.hint} ${command.keywords}`.toLowerCase();
    return terms.every((term) => haystack.includes(term));
  });
}

export const navigationChord: Record<string, string> = {
  o: "/overview",
  c: "/clusters",
  p: "/plugins",
  s: "/settings",
};
