# Security Policy

## Supported versions

Security fixes target the default branch and the latest released version.

## Reporting a vulnerability

Please do not open a public GitHub issue for a suspected vulnerability.

Use GitHub's private vulnerability reporting/security advisory mechanism for this repository when available. Include:

- affected version/commit
- affected component or file
- reproduction steps or proof of concept
- security impact
- suggested mitigation, if known

Do not include credentials, tokens, private keys, kubeconfigs, customer data, or other secrets in the report.

## Secret handling

KubeMv must never commit Kubernetes credentials, kubeconfigs, API tokens, cloud credentials, private keys, `.env` files containing secrets, or production configuration containing credentials.

If a secret is accidentally committed:

1. Revoke/rotate it immediately.
2. Remove it from the working tree.
3. Remove it from Git history when appropriate.
4. Review GitHub security alerts and audit logs.

Removing a secret from Git history does not replace credential rotation.
