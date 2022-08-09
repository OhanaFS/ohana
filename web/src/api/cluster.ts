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
  startDate: string; // need to format
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
export const getnumOfFiles = () =>
  useQuery(['numOfFiles'], () =>
    APIClient.get<number>(`/api/v1/cluster/stats/num_of_files`)
      .then((res) => res.data)
      .catch(typedError)
  );
//working

// Get the number of historical files count from the cluster
export const getnumOfHistoricalFiles = () =>
  useQuery(['numOfHistoricalFiles'], () =>
    APIClient.get<UsedRes>(`/api/v1/cluster/stats/num_of_files_historical`, {
      headers: {
        rangeType: 1,
        // iso format using toisostring
        startDate: new Date('2022-08-01T00:00:00.000Z').toISOString(),
        endDate: new Date('2022-08-02T00:00:00.000Z').toISOString(),
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the storage used without parity and versioning from the cluster
export const getstorageUsed = () =>
  useQuery(['storageNonReplicaUsed'], () =>
    APIClient.get<number>(`/api/v1/cluster/stats/non_replica_used`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the historical storage used without parity and versioning from the cluster
export const gethistoricalStorageUsed = () =>
  useQuery(['nonReplicaUsedHistorical'], () =>
    APIClient.get<number>(`/api/v1/cluster/stats/non_replica_used_historical`, {
      headers: {
        rangeType: 1,
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the storage used with parity and versioning from the cluster
export const getstorageUsedWithParity = () =>
  useQuery(['replica'], () =>
    APIClient.get<number>(`/api/v1/cluster/stats/replica_used`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the historical storage used with parity and versioning from the cluster
export const gethistoricalStorageUsedWithParity = () =>
  useQuery(['replicaHistorical'], () =>
    APIClient.get<UsedRes>(`/api/v1/cluster/stats/replica_used_historical`, {
      headers: {
        rangeType: 1,
        startDate: '2020-01-01',
        endDate: '2020-01-01',
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

// Get all alerts related to the cluster
export const getAlerts = () =>
  useQuery(['alerts'], () =>
    APIClient.get<Alerts>(`/api/v1/cluster/stats/alerts`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Clear all alerts related to the cluster
export const clearAlerts = () =>
  useMutation(() =>
    APIClient.delete<boolean>(`/api/v1/cluster/stats/alerts`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Get specific alerts of the cluster
export const getAlertsID = (id: number) =>
  useQuery(['alertsID', id], () =>
    APIClient.get<Alerts['log_entries']>(`/api/v1/cluster/stats/alerts/${id}`, {
      headers: {
        id: 1,
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Clear specific alerts of the cluster
export const clearAlertsID = (id: number) =>
  useMutation((id) =>
    APIClient.delete<boolean>(`/api/v1/cluster/stats/alerts/${id}`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Return server logs
export const getserverLogs = () =>
  useQuery(['serverLogs'], () =>
    APIClient.get<Alerts['log_entries']>(`/api/v1/cluster/stats/logs`, {
      headers: {
        startNum: 0,
        startDate: '2020-01-01',
        endDate: '2020-01-01',
        filter: '',
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Clear logs
export const clearserverLogs = () =>
  useQuery(['clearserverLogs'], () =>
    APIClient.delete<boolean>(`/api/v1/cluster/stats/logs`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Return server logs by id
export const getserverLogsID = (id: number) =>
  useQuery(['serverLogsID', id], () =>
    APIClient.get<Alerts['log_entries']>(`/api/v1/cluster/stats/logs/${id}`, {
      headers: {
        startNum: 0,
        startDate: '2020-01-01',
        endDate: '2020-01-01',
        filter: '',
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Return server statuses
export const getserverStatuses = () =>
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
export const getserverStatusesID = (serverName: string) =>
  useQuery(['serverStatusesID', serverName], () =>
    APIClient.get<ServerStatuses>(`/api/v1/cluster/stats/servers/${serverName}`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Delete server specific statuses
export const deleteserverStatusesID = (serverName: string) =>
  useQuery(['deleteserverStatusesID', serverName], () =>
    APIClient.delete<boolean>(`/api/v1/cluster/stats/servers/${serverName}`)
      .then((res) => res.data)
      .catch(typedError)
  );
