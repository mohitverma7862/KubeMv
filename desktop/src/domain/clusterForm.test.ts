import { describe, expect, it } from "vitest";
import { validateClusterDraft } from "./clusterForm";

const valid = {
  name: "prod-eks",
  provider: "eks" as const,
  context: "prod",
  kubeconfigRef: "secret://clusters/prod",
};

describe("validateClusterDraft", () => {
  it("accepts registry metadata", () => {
    expect(validateClusterDraft(valid)).toBeNull();
  });

  it("rejects credential material in the reference", () => {
    expect(
      validateClusterDraft({ ...valid, kubeconfigRef: "token:abc" }),
    ).toMatch(/credential/i);
    expect(validateClusterDraft({ ...valid, name: "bad name" })).toMatch(/Name/);
  });
});
