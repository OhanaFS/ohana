import { APIClient, typedError } from './api';
import { useMutation, useQuery } from '@tanstack/react-query';

export type Records = [
  {
    Id: 0;
    date_time_started: string;
    date_time_ended: string;
    total_time_taken: string;
    total_shards_scanned: number;
    total_files_scanned: number;
    tasks: [
      {
        job_type: string;
        id: number;
        status: number;
      }
    ];
    progress: number;
    status_msg: string;
    status: number;
  }
];

/*
  [
    {
      "file_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "file_name": "string",
      "data_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "fragment_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "server_name": "string",
      "status": 0
    }
  ]
*/

// Get all the records
export const getMainRecords = () => {
  return useQuery(['mainRecords'], () =>
    APIClient.get<Records>(`/api/v1/maintenance/all`, {
      headers: {
        startNum: 0,
      },
    })
      .then((res) => res.data)
      .catch(typedError)
  );
};

// Get the records based on the ID
export const getMainRecordsID = (id: number) => {
  return useQuery(['mainRecordsID', id], () =>
    APIClient.get<any>(`/api/v1/maintenance/job/${id}`)
      .then((res) => res.data)
      .catch(typedError)
  );
};

/*
// Delete the records based on the ID
const deleteMainRecordsID = () => {
    return useMutation((id: number) =>
        APIClient.delete<any>(`/api/v1/maintenance/job/${id}`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Create a job based on the ID
const createMainRecordsID = () => {
    return useMutation((id: number) =>
        APIClient.patch<any>(`/api/v1/maintenance/job/${id}`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Start a job based on the ID
const startMainRecordsID = () => {
    return useMutation((id: number) =>
        APIClient.post<any>(`/api/v1/maintenance/start`, {
            headers: {
                full_shards_check: true,
                quick_Shards_check: true,
                missing_shards_check: true,
                orphaned_shards_check: true,
                orphaned_files_check: true,
                permission_check: true,
                delete_fragments: true,
        },
        })
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Get full shards job results
const getFullShardsResults = (id: number) => {
    return useQuery(['fullShardsResults', id], () =>
        APIClient.get<any>(`/api/v1/maintenance/jon/${id}/full_shards`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// FIx full shards job 
const fixFullShards = (id: number) => {
    return useMutation((id: number) =>
        APIClient.post<any>(`/api/v1/maintenance/jon/${id}/full_shards`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Get quick shards job results
const getQuickShardsResults = (id: number) => {
    return useQuery(['quickShardsResults', id], () =>
        APIClient.get<any>(`/api/v1/maintenance/jon/${id}/quick_shards`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Fix quick shards job
const fixQuickShards = () => {
    return useMutation((id: number) =>
        APIClient.post<any>(`/api/v1/maintenance/jon/${id}/quick_shards`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Get missing shards job results
const getMissingShardsResults = (id: number) => {
    return useQuery(['missingShardsResults', id], () =>
        APIClient.get<any>(`/api/v1/maintenance/jon/${id}/missing_shards`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Fix missing shards job
const fixMissingShards = () => {
    return useMutation((id: number) =>
        APIClient.post<any>(`/api/v1/maintenance/jon/${id}/missing_shards`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Get orphaned shards job results
const getOrphanedShardsResults = (id: number) => {
    return useQuery(['orphanedShardsResults', id], () =>
        APIClient.get<any>(`/api/v1/maintenance/jon/${id}/orphaned_shards`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Fix orphaned shards job
const fixOrphanedShards = () => {
    return useMutation((id: number) =>
        APIClient.post<any>(`/api/v1/maintenance/jon/${id}/orphaned_shards`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Get orphaned files job results
const getOrphanedFilesResults = (id: number) => {
    return useQuery(['orphanedFilesResults', id], () =>
        APIClient.get<any>(`/api/v1/maintenance/jon/${id}/orphaned_files`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Fix orphaned files job
const fixOrphanedFiles = () => {
    return useMutation((id: number) =>
        APIClient.post<any>(`/api/v1/maintenance/jon/${id}/orphaned_files`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Get permission check job results
const getPermissionCheckResults = (id: number) => {
    return useQuery(['permissionCheckResults', id], () =>
        APIClient.get<any>(`/api/v1/maintenance/jon/${id}/permission_check`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Fix permission check job
const fixPermissionCheck = () => {
    return useMutation((id: number) =>
        APIClient.post<any>(`/api/v1/maintenance/jon/${id}/permission_check`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Get backup keys
const getBackupKeys = () => {
    return useQuery(['backupKeys'], () =>
        APIClient.get<any>(`/api/v1/maintenance/backup_keys`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Send backup keys
const sendBackupKeys = () => {
    return useMutation((id: number) =>
        APIClient.put<any>(`/api/v1/maintenance/backup_keys`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Deleting file key
const deleteFileKey = () => {
    return useMutation((id: number) =>
        APIClient.delete<any>(`/api/v1/maintenance/backup_keys`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

//


// Deleting folder key
const deleteFolderKey = () => {
    return useMutation((id: number) =>
        APIClient.delete<any>(`/api/v1/maintenance/backup_keys`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

//

// Deletion master key
const deleteMasterKey = () => {
    return useMutation((id: number) =>
        APIClient.delete<any>(`/api/v1/maintenance/backup_keys`)
        .then((res) => res.data)
        .catch(typedError)
    );
}

// Request master key
const requestMasterKey = () => {
    return useMutation((id: number) =>
        APIClient.put<any>(`/api/v1/maintenance/backup_keys`)
        .then((res) => res.data)
        .catch(typedError)
    );
}*/
