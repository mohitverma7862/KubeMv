import { describe, expect, it } from "vitest";
import type { Cluster } from "../api/types";
import { displayCluster } from "./clusterView";

describe("displayCluster", () => {
  it("keeps only public registry fields", () => {
    const record = {
      id: "cls_1",
      name: "prod",
      provider: "eks",
      context: "prod",
      connectionState: "not_connected",
      connectionDetail: "metadata only",
      createdAt: "2026-10-04T00:00:00Z",
      kubeconfigRef: "secret://clusters/prod",
    } as Cluster & { kubeconfigRef: string };
    const displayed = displayCluster(record);
    expect(JSON.stringify(displayed)).not.toContain("secret://");
    expect(Object.keys(displayed)).not.toContain("kubeconfigRef");
  });
});
