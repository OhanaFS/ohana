import { APIClient, typedError } from './api';
import { useMutation, useQuery } from '@tanstack/react-query';

export type MaintenanceProgress = {
  id: number;
  start_time: string;
  end_time: string;
  server_name: string;
  in_progress: boolean;
  msg: string;
};

export type Record = {
  id: number;
  start_time: string;
  end_time: string;
  total_time_taken: number;
  missing_shards_check: boolean;
  missing_shards_progress: MaintenanceProgress[];
  orphaned_shards_check: boolean;
  orphaned_shards_progress: MaintenanceProgress[];
  quick_shards_health_check: boolean;
  quick_shards_health_progress: MaintenanceProgress[];
  all_files_shards_health_check: boolean;
  all_files_shards_health_progress: MaintenanceProgress[];
  permission_check: boolean;
  delete_fragments: boolean;
  delete_fragments_results: MaintenanceProgress[];
  orphaned_files_check: boolean;
  orphaned_files_results: MaintenanceProgress[];
  progress: number;
  status_msg: string;
  status: number;
};

// Get all the records
export const useQueryGetMaintenanceRecords = (
  startNum: number,
  startDate: string,
  endDate: string,
  filter: string
) =>
  useQuery(['mainRecords'], () =>
    APIClient.get<Record[]>(`/api/v1/maintenance/all`, {
      headers: {
        start_num: startNum,
        start_date: startDate,
        end_date: endDate,
        filter: filter,
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

// Get the records based on the ID
export const useQueryGetMaintenanceRecordsID = (id: number) => {
  return useQuery(['mainRecordsID', id], () =>
    APIClient.get<Record>(`/api/v1/maintenance/job/${id}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Delete the records based on the ID
const useMutateDeleteMainRecordsID = () => {
  return useMutation((id: number) =>
    APIClient.delete<Record>(`/api/v1/maintenance/job/${id}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Create a job based on the ID
const useMutateCreateMainRecordsID = () => {
  return useMutation((id: number) =>
    APIClient.patch<boolean>(`/api/v1/maintenance/job/${id}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

export type MaintenanceRecordCheck = {
  full_shards_check: boolean;
  quick_Shards_check: boolean;
  missing_shards_check: boolean;
  orphaned_shards_check: boolean;
  orphaned_files_check: boolean;
  permission_check: boolean;
  delete_fragments: boolean;
};

// Start a job based on the ID
export const useMutateStartMainRecordsID = () => {
  return useMutation((params: MaintenanceRecordCheck) =>
    APIClient.post<Record>(`/api/v1/maintenance/start`, null, {
      headers: {
        full_shards_check: String(params.full_shards_check),
        quick_Shards_check: String(params.quick_Shards_check),
        missing_shards_check: String(params.missing_shards_check),
        orphaned_shards_check: String(params.orphaned_shards_check),
        orphaned_files_check: String(params.orphaned_files_check),
        permission_check: String(params.permission_check),
        delete_fragments: String(params.delete_fragments),
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );
};

export type ShardsResults = {
  file_id: string;
  file_name: string;
  data_id: string;
  fragment_id: string;
  server_name: string;
  status: number;
};

// Get full shards job results
const useQueryGetFullShardsResults = (id: number) =>
  useQuery(['fullShardsResults', id], () =>
    APIClient.get<ShardsResults>(`/api/v1/maintenance/jon/${id}/full_shards`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Fix full shards job
const useMutateFixFullShards = (id: number) => {
  return useMutation((id) =>
    APIClient.get<ShardsResults>(`/api/v1/maintenance/jon/${id}/full_shards`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Get quick shards job results
const useQueryGetQuickShardsResults = (id: number) =>
  useQuery(['quickShardsResults', id], () =>
    APIClient.get<ShardsResults>(`/api/v1/maintenance/jon/${id}/quick_shards`)
      .then((res) => res.data)
      .catch(typedError)
  );

export type QuickShardsJob = {
  file_id: string;
  fragment_id: string;
  password: string;
  action: number;
};

// Fix quick shards job
const useMutatefixQuickShards = () => {
  return useMutation((id: number) =>
    APIClient.post<QuickShardsJob>(`/api/v1/maintenance/jon/${id}/quick_shards`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Get missing shards job results
const useQueryGetMissingShardsResults = (id: number) =>
  useQuery(['missingShardsResults', id], () =>
    APIClient.get<ShardsResults>(`/api/v1/maintenance/jon/${id}/missing_shards`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Fix missing shards job
const useMutateFixMissingShards = () => {
  return useMutation((id: number) =>
    APIClient.post<QuickShardsJob>(
      `/api/v1/maintenance/jon/${id}/missing_shards`
    )
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Get orphaned shards job results
const useQueryGetOrphanedShardsResults = (id: number) =>
  useQuery(['orphanedShardsResults', id], () =>
    APIClient.get<ShardsResults>(
      `/api/v1/maintenance/jon/${id}/orphaned_shards`
    )
      .then((res) => res.data)
      .catch(typedError)
  );

// Fix orphaned shards job
const useMutateFixOrphanedShards = () => {
  return useMutation((id: number) =>
    APIClient.post<QuickShardsJob>(
      `/api/v1/maintenance/jon/${id}/orphaned_shards`
    )
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Get orphaned files job results
const useQueryGetOrphanedFilesResults = (id: number) =>
  useQuery(['orphanedFilesResults', id], () =>
    APIClient.get<ShardsResults>(`/api/v1/maintenance/jon/${id}/orphaned_files`)
      .then((res) => res.data)
      .catch(typedError)
  );

// Fix orphaned files job
const useMutateFixOrphanedFiles = () => {
  return useMutation((id: number) =>
    APIClient.post<QuickShardsJob>(
      `/api/v1/maintenance/jon/${id}/orphaned_files`
    )
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Get permission check on jobs
const useQueryGetPermissionCheckResults = (id: number) =>
  useQuery(['permissionCheckResults', id], () =>
    APIClient.get<ShardsResults>(
      `/api/v1/maintenance/jon/${id}/permission_check`
    )
      .then((res) => res.data)
      .catch(typedError)
  );

// Fix permission check job
const useMutateFixPermissionCheck = () => {
  return useMutation((id: number) =>
    APIClient.post<QuickShardsJob>(
      `/api/v1/maintenance/jon/${id}/permission_check`
    )
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Get backup keys
const useQueryGetBackupKeys = () =>
  useQuery(['backupKeys'], () =>
    APIClient.get<boolean>(`/api/v1/maintenance/backup_keys`)
      .then((res) => res.data)
      .catch(typedError)
  );

// needs changes
// Send backup keys
const useMutateSendBackupKeys = () => {
  return useMutation((id: number) =>
    APIClient.put<boolean>(`/api/v1/maintenance/backup_keys`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// needs changes
const useMutatePostFileKey = () => {
  return useMutation((fileID: number) =>
    APIClient.post<boolean>(`/api/v1/maintenance/KEY/${fileID}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Deleting file key
const useMutateDeleteFileKey = () => {
  return useMutation((fileID: number) =>
    APIClient.delete<boolean>(`/api/v1/maintenance/KEY/${fileID}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// needs changes
const useMutatePostFolderKey = () => {
  return useMutation((folderID: number) =>
    APIClient.post<boolean>(`/api/v1/maintenance/key/${folderID}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Deleting folder key
const useMutateDeleteFolderKey = () => {
  return useMutation((folderID: number) =>
    APIClient.delete<any>(`/api/v1/maintenance/key/${folderID}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Deletion master key
const useMutationDeleteMasterKey = () => {
  return useMutation((serverID: number) =>
    APIClient.delete<boolean>(`/api/v1/maintenance/key/${serverID}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Request master key
export const useMutationRequestMasterKey = () => {
  return useMutation((serverID: number) =>
    APIClient.put<boolean>(`/api/v1/maintenance/key/${serverID}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};
