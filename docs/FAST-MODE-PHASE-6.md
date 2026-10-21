# Fast Mode Phase 6

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


The selected-resource inspector and the relationship tree read Kubernetes ownership metadata.

## Added

- The inspector shows status, namespace, node, IP, readiness, restarts, age, QoS, CPU, memory, limits, ports, mounts, containers, labels, and the immediate owner when those fields exist.
- Pod rows from the live client carry container images, resource requests, limits, ports, mount names, QoS, and controller references. Mounts name a Secret or ConfigMap. They do not include secret data or environment values.
- The relationship tree walks `Deployment → ReplicaSet → Pod → Container`. Tab focuses the tree, `j`/`k` move, and Enter opens a listable owner. A missing owner is loaded from the API.
- Secret rows keep type and key count. Values, tokens, and private keys are omitted from the inspector and the tree.
- Contextual actions in the inspector come from the shortcut registry.

## Not done

- Logs, exec, describe, YAML editing, and rollout actions remain Phase 7. The action labels are visible. They do not perform those operations yet.
