import { QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { AppShell } from "../components/Shell";
import { LoginPage } from "../pages/LoginPage";
import { ClustersPage } from "../pages/ClustersPage";
import { OverviewPage } from "../pages/OverviewPage";
import { PluginsPage } from "../pages/PluginsPage";
import { SettingsPage } from "../pages/SettingsPage";
import { queryClient } from "./query";
import { RequireSession } from "./session";

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route element={<RequireSession />}>
            <Route element={<AppShell />}>
              <Route path="/overview" element={<OverviewPage />} />
              <Route path="/clusters" element={<ClustersPage />} />
              <Route path="/plugins" element={<PluginsPage />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Route>
          </Route>
          <Route path="/" element={<Navigate to="/overview" replace />} />
          <Route path="*" element={<Navigate to="/overview" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
