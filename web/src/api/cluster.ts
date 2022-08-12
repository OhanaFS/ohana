import { APIClient, typedError } from './api';
import { useMutation, useQuery } from '@tanstack/react-query';

export type Server = {
  name: string;
  hostname: string;
  port: string;
  status: number;
  free_space: number;
  used_space: number;
};

export type ServerID = {
  server_id: string;
};

export type DateTime = {
  rangeType: number;
  startDate: string; // need to format 2022-08-05T15:57:05.279Z
  endDate: string;
};

export type Logs = {
  startNum: number;
  startDate: string; // need to format
  endDate: string;
  filter: string;
};

export type UsedRes = {
  date: string;
  value: number;
};

// Get the number of files from the cluster
export const useQueryGetnumOfFiles = () =>
  useQuery(['numOfFiles'], () =>
    APIClient.get<number>(`/api/v1/cluster/stats/num_of_files`)
      .then((res) => res.data)
      .catch(typedError)
  );
//working

// Get the number of historical files count from the cluster
export const useQueryGetnumOfHistoricalFiles = (
  range: number,
  startD: string,
  endD: string
) =>
  useQuery(['numOfHistoricalFiles', range, startD, endD], () =>
    APIClient.get<UsedRes>(`/api/v1/cluster/stats/num_of_files_historical`, {
      headers: {
        rangeType: 1,
        startDate: new Date(startD).toISOString(),
        endDate: new Date(endD).toISOString(),
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the storage used without parity and versioning from the cluster
export const useQueryGetstorageUsed = () =>
  useQuery(['storageNonReplicaUsed'], () =>
    APIClient.get<number>(`/api/v1/cluster/stats/non_replica_used`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the historical storage used without parity and versioning from the cluster
export const useQueryGethistoricalStorageUsed = (
  range: number,
  startD: string,
  endD: string
) =>
  useQuery(['nonReplicaUsedHistorical', range, startD, endD], () =>
    APIClient.get<number>(`/api/v1/cluster/stats/non_replica_used_historical`, {
      headers: {
        rangeType: range,
        // iso format using toisostring
        startDate: new Date(startD).toISOString(),
        endDate: new Date(endD).toISOString(),
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the storage used with parity and versioning from the cluster
export const useQueryGetstorageUsedWithParity = () =>
  useQuery(['replica'], () =>
    APIClient.get<number>(`/api/v1/cluster/stats/replica_used`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the historical storage used with parity and versioning from the cluster
export const useQueryGethistoricalStorageUsedWithParity = (
  range: number,
  startD: string,
  endD: string
) =>
  useQuery(['replicaHistorical'], () =>
    APIClient.get<UsedRes>(`/api/v1/cluster/stats/replica_used_historical`, {
      headers: {
        rangeType: range,
        startDate: new Date(startD).toISOString(),
        endDate: new Date(endD).toISOString(),
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

export type Alerts = {
  log_entries: [
    {
      id: number;
      log_type: number;
      server_name: string;
      message: string;
      timestamp: string;
    }
  ];
  fatal_count: number;
  error_count: number;
  warning_count: number;
};

export type LogEntries = {
  id: number;
  log_type: number;
  server_name: string;
  message: string;
  timestamp: string;
};

// Get all alerts related to the cluster
export const useQueryGetAlerts = () =>
  useQuery(['alerts'], () =>
    APIClient.get<Alerts>(`/api/v1/cluster/stats/alerts`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Clear all alerts related to the cluster
export const useQueryClearAlerts = () =>
  useMutation(() =>
    APIClient.delete<boolean>(`/api/v1/cluster/stats/alerts`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Get specific alerts of the cluster
export const useQueryGetAlertsID = (Id: number) =>
  useQuery(['alertsID', Id], () =>
    APIClient.get<LogEntries>(`/api/v1/cluster/stats/alerts/${Id}`, {
      headers: {
        id: Id,
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Clear specific alerts of the cluster
export const useQueryClearAlertsID = (id: number) =>
  useMutation((id) =>
    APIClient.delete<boolean>(`/api/v1/cluster/stats/alerts/${id}`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Return server logs
export const useQueryGetserverLogs = (
  start: number,
  startD: string,
  endD: string,
  fil: string
) =>
  useQuery(['serverLogs'], () =>
    APIClient.get<LogEntries>(`/api/v1/cluster/stats/logs`, {
      headers: {
        startNum: start,
        startDate: new Date(startD).toISOString(),
        endDate: new Date(endD).toISOString(),
        filter: fil,
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Clear logs
export const useQueryClearserverLogs = () =>
  useQuery(['clearserverLogs'], () =>
    APIClient.delete<boolean>(`/api/v1/cluster/stats/logs`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Return server logs by id
export const useQueryGetserverLogsID = (
  id: number,
  start: number,
  startD: string,
  endD: string,
  fil: string
) =>
  useQuery(['serverLogsID', id, start, startD, endD, fil], () =>
    APIClient.get<LogEntries>(`/api/v1/cluster/stats/logs/${id}`, {
      headers: {
        startNum: start,
        startDate: new Date(startD).toISOString(),
        endDate: new Date(endD).toISOString(),
        filter: fil,
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Return server statuses
export const useQueryGetserverStatuses = () =>
  useQuery(['serverStatuses'], () =>
    APIClient.get<Server>(`/api/v1/cluster/stats/servers`)
      .then((res) => res.data)
      .catch(typedError)
  );

export type ServerStatuses = {
  name: string;
  hostname: string;
  port: string;
  status: number;
  free_space: number;
  used_space: number;
  load_avg: string;
  uptime: number;
  cpu: number;
  memory_used: number;
  memory_free: number;
  network_rx_bytes: number;
  network_tx_bytes: number;
  reads: number;
  writes: number;
  warnings: number;
  errors: number;
  smart_good: boolean;
  smart_status: string;
};

// Return server specific statuses
export const useQueryGetserverStatusesID = (serverName: string) =>
  useQuery(['serverStatusesID', serverName], () =>
    APIClient.get<ServerStatuses>(`/api/v1/cluster/stats/servers/${serverName}`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Delete server specific statuses
export const useQueryDeleteserverStatusesID = (serverName: string) =>
  useQuery(['deleteserverStatusesID', serverName], () =>
    APIClient.delete<boolean>(`/api/v1/cluster/stats/servers/${serverName}`)
      .then((res) => res.data)
      .catch(typedError)
  );