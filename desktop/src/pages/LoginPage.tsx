import { useState, type FormEvent } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Navigate } from "react-router-dom";
import { ApiError } from "../api/client";
import { login } from "../api/resources";
import { sessionKey } from "../app/query";
import { useSessionQuery } from "../app/session";
import { Mark } from "../components/Mark";

export function LoginPage() {
  const session = useSessionQuery();
  const queryClient = useQueryClient();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const mutation = useMutation({
    mutationFn: () => login(username.trim(), password),
    onSuccess: (next) => {
      queryClient.setQueryData(sessionKey, next);
    },
  });

  if (session.isLoading) {
    return (
      <div className="km-screen" role="status">
        Checking session
      </div>
    );
  }
  if (session.data) {
    return <Navigate to="/overview" replace />;
  }

  const onSubmit = (event: FormEvent) => {
    event.preventDefault();
    setFormError(null);
    if (username.trim() === "" || password === "") {
      setFormError("Enter a username and password.");
      return;
    }
    mutation.mutate(undefined, {
      onError: (error) => {
        setFormError(error instanceof ApiError ? error.message : "Could not reach the API.");
      },
    });
  };

  return (
    <div className="km-login">
      <section className="km-login-panel">
        <div className="km-brand" style={{ marginBottom: 28 }}>
          <Mark />
          <span className="km-word">
            Kube<span>Mv</span>
          </span>
        </div>
        <p className="km-kicker">Operator sign in</p>
        <h1>The intelligent command center for Kubernetes</h1>
        <p className="km-lead">
          This build is the architecture foundation. Signing in starts a server session. Kubernetes credentials are not accepted here.
        </p>
        <form data-testid="login-form" onSubmit={onSubmit} style={{ marginTop: 22 }}>
          {session.isError ? (
            <div className="km-alert" role="alert">
              The API is unreachable.
            </div>
          ) : null}
          {formError ? (
            <div className="km-alert" role="alert" data-testid="login-error">
              {formError}
            </div>
          ) : null}
          <label className="km-field">
            <span>Username</span>
            <input
              name="username"
              autoComplete="username"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
            />
          </label>
          <label className="km-field">
            <span>Password</span>
            <input
              name="password"
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
            />
          </label>
          <button className="km-btn" type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "Signing in" : "Sign in"}
          </button>
        </form>
      </section>
      <aside className="km-login-aside">
        <div>
          <blockquote>See. Understand. Troubleshoot. Fix. Validate. Audit.</blockquote>
          <p>Phase 0 establishes the shell, session, cluster registry, and plugin boundary.</p>
        </div>
      </aside>
    </div>
  );
}
