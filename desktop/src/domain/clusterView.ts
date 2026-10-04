import type { Cluster } from "../api/types";

export type ClusterDisplay = {
  id: string;
  name: string;
  provider: string;
  context: string;
  connectionState: string;
  connectionDetail: string;
  createdAt: string;
};

export function displayCluster(record: Cluster): ClusterDisplay {
  return {
    id: record.id,
    name: record.name,
    provider: record.provider,
    context: record.context,
    connectionState: record.connectionState,
    connectionDetail: record.connectionDetail,
    createdAt: record.createdAt,
  };
}
