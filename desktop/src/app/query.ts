import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
      staleTime: 5_000,
    },
  },
});

export const sessionKey = ["session"] as const;
export const healthKey = ["health"] as const;
export const clustersKey = ["clusters"] as const;
export const pluginsKey = ["plugins"] as const;
