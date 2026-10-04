import { useQuery } from "@tanstack/react-query";
import { Navigate, Outlet } from "react-router-dom";
import { fetchSession } from "../api/resources";
import { sessionKey } from "./query";

export function useSessionQuery() {
  return useQuery({
    queryKey: sessionKey,
    queryFn: fetchSession,
  });
}

export function RequireSession() {
  const session = useSessionQuery();
  if (session.isLoading) {
    return (
      <div className="km-screen" role="status">
        Checking session
      </div>
    );
  }
  if (session.isError) {
    return (
      <div className="km-screen" role="alert">
        The API is unreachable. Start the KubeMv backend and reload.
      </div>
    );
  }
  if (!session.data) {
    return <Navigate to="/login" replace />;
  }
  return <Outlet />;
}
