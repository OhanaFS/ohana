import { APIClient, typedError } from './api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { WhoamiResponse } from './auth';

export const K_ROOT_FOLDER_ID = '00000000-0000-0000-0000-000000000000';
export const K_HOME_PARENT_FOLDER_ID = '00000000-0000-0000-0000-000000000001';

export type FileUploadRequest = {
  file: File;
  folder_id: string;
  frag_count: number;
  parity_count: number;
};

/**
 * The type of the entry. Folder is `1`, File is `2`.
 */
export enum EntryType {
  Folder = 1,
  File = 2,
}

export type FileMetadata<T = EntryType> = {
  file_id: string;
  file_name: string;
  mime_type: T extends EntryType.File ? string : never;
  entry_type: T;
  parent_folder_id: string;
  version_no: number;
  data_version_no: number;
  size: number;
  actual_size: number;
  created_time: string; // date
  modified_user_user_id: string;
  modified_time: string; // date
  versioning_mode: number;
  checksum: T extends EntryType.File ? string : never;
  frag_count: T extends EntryType.File ? number : never;
  parity_count: T extends EntryType.File ? number : never;
  password_protected: boolean;
  link_file_id: string;
  last_checked: string; // date
  status: number;
};

export const MetadataKeyMap: { [key in keyof FileMetadata]: string } = {
  file_id: 'File ID',
  file_name: 'File Name',
  mime_type: 'MIME Type',
  entry_type: 'Entry Type',
  parent_folder_id: 'Parent Folder ID',
  version_no: 'Version No',
  data_version_no: 'Data Version Number',
  size: 'Size',
  actual_size: 'Actual Size',
  created_time: 'Created Time',
  modified_user_user_id: 'Modified By (UID)',
  modified_time: 'Last Modified',
  versioning_mode: 'Versioning mode',
  checksum: 'Checksum',
  frag_count: 'Frag Count',
  parity_count: 'Parity Count',
  password_protected: 'Password Protected',
  link_file_id: 'Link File ID',
  last_checked: 'Last Checked',
  status: 'Status',
} as const;

/**
 * Uploads a file to a folder.
 */
export const useMutateUploadFile = () => {
  const queryClient = useQueryClient();
  return useMutation(
    async ({ file, ...headers }: FileUploadRequest) => {
      const form = new FormData();
      form.append('file', file);
      return APIClient.post<FileMetadata<EntryType.File>>(
        '/api/v1/file',
        form,
        { headers: { ...headers } }
      )
        .then((res) => res.data)
        .catch(typedError);
    },
    {
      onSuccess: (_, params) => {
        queryClient.invalidateQueries([
          'folder',
          'contents',
          'id',
          params.folder_id,
        ]);
      },
    }
  );
};

export type FileUpdateRequest = {
  file: File;
  file_id: string;
  frag_count: number;
  parity_count: number;
};

/**
 * Upload a newer version of a file.
 */
export const useMutateUpdateFile = () => {
  const queryClient = useQueryClient();
  return useMutation(
    async ({ file, file_id, ...headers }: FileUpdateRequest) => {
      const form = new FormData();
      form.append('file', file);

      return (
        APIClient.post<FileMetadata<EntryType.File>>(
          `/api/v1/file/${file_id}/update`,
          form,
          { headers: { ...headers } }
        )
          .then((res) => res.data)
          .catch(typedError),
        {
          onSuccess: () => {
            queryClient.invalidateQueries(['file']);
            queryClient.invalidateQueries(['folder', 'contents']);
          },
        }
      );
    }
  );
};

/**
 * Get a file's metadata.
 */
export const useQueryFileMetadata = (fileId: string) =>
  useQuery(
    ['file', 'metadata', fileId],
    () =>
      !fileId
        ? Promise.reject()
        : APIClient.get<FileMetadata<EntryType.File>>(
            `/api/v1/file/${fileId}/metadata`
          )
            .then((res) => res.data)
            .catch(typedError),
    { keepPreviousData: true }
  );

export type FileMetadataUpdateRequest = {
  file_id: string;
  file_name?: string;
  mime_type?: string;
  versioning_mode?: number;
  password_modification?: boolean;
  password_protected?: boolean;
  password_hint?: string;
  old_password?: string;
  new_password?: string;
};

/**
 * Update a file's metadata.
 */
export const useMutateUpdateFileMetadata = () => {
  const queryClient = useQueryClient();
  return useMutation(
    ({ file_id, ...body }: FileMetadataUpdateRequest) =>
      APIClient.patch<FileMetadata<EntryType.File>>(
        `/api/v1/file/${file_id}/metadata`,
        body
      )
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['file', 'metadata']);
      },
    }
  );
};

export type MoveFileRequest = {
  /** The ID of the file to move */
  file_id: string;
  /** The ID of the folder to move the file to */
  folder_id: string;
};

/**
 * Move a file to a new folder.
 */
export const useMutateMoveFile = () => {
  const queryClient = useQueryClient();
  return useMutation(
    ({ file_id, folder_id }: MoveFileRequest) =>
      APIClient.post<boolean>(`/api/v1/file/${file_id}/move`, null, {
        headers: { folder_id },
      })
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['folder', 'contents']);
      },
    }
  );
};

/**
 * Copy a file to a new folder.
 */
export const useMutateCopyFile = () => {
  const queryClient = useQueryClient();
  return useMutation(
    ({ file_id, folder_id }: MoveFileRequest) =>
      APIClient.post<boolean>(`/api/v1/file/${file_id}/copy`, null, {
        headers: { folder_id },
      })
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['folder', 'contents']);
      },
    }
  );
};

/**
 * Get the download URL for a file. If a version number is provided, the
 * download URL will point to the specified version.
 */
export const getFileDownloadURL = (
  fileId: string,
  opts?: { versionId?: number; inline?: boolean }
) =>
  [
    window.location.origin,
    `/api/v1/file/${fileId}`,
    opts?.versionId ? `/version/${opts?.versionId}` : '',
    opts?.inline ? '?inline=true' : '',
  ].join('');

/**
 * Delete a file by its ID.
 */
export const useMutateDeleteFile = () => {
  const queryClient = useQueryClient();
  return useMutation(
    (fileId: string) =>
      APIClient.delete<boolean>(`/api/v1/file/${fileId}`)
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['folder', 'contents']);
      },
    }
  );
};

export type Permission = {
  can_read: boolean;
  can_write: boolean;
  can_execute: boolean;
  can_share: boolean;
  can_audit: boolean;
};

export type FilePermission = Permission & {
  file_id: string;
  permission_id: string;
  user_id: string;
  group_id: string;
  version_no: number;
  audit: boolean;
  created_at: string;
  updated_at: string;
  status: number;
  User: WhoamiResponse;
  Group: {
    group_id: string;
    group_name: string;
  };
};

/**
 * Check available permissions on a file.
 */
export const useQueryFilePermissions = (fileId: string) =>
  useQuery(
    ['file', 'permissions', fileId],
    () =>
      !fileId
        ? Promise.reject()
        : APIClient.get<FilePermission[]>(`/api/v1/file/${fileId}/permissions`)
            .then((res) => res.data)
            .catch(typedError),
    { keepPreviousData: true }
  );

export type UpdateFilePermissionsRequest = Permission & {
  file_id: string;
  permission_id: string;
  is_folder?: boolean;
};

/**
 * Modify permissions on a file.
 */
export const useMutateUpdateFilePermissions = () => {
  const queryClient = useQueryClient();
  return useMutation(
    ({
      file_id,
      permission_id,
      is_folder,
      ...perms
    }: UpdateFilePermissionsRequest) =>
      APIClient.patch<boolean>(
        `/api/v1/${
          is_folder ? 'folder' : 'file'
        }/${file_id}/permissions/${permission_id}`,
        null,
        { headers: { ...perms } }
      )
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: (_, params) => {
        queryClient.invalidateQueries(['file', 'permissions', params.file_id]);
        queryClient.invalidateQueries(['file', 'metadata', params.file_id]);
      },
    }
  );
};

export type AddFilePermissionsRequest = Permission & {
  file_id: string;
  users: string[];
  groups: string[];
  is_folder?: boolean;
};

/**
 * Add new permissions to a file with an array of users and groups.
 */
export const useMutateAddFilePermissions = () => {
  const queryClient = useQueryClient();
  return useMutation(
    ({
      file_id,
      users,
      groups,
      is_folder,
      ...perms
    }: AddFilePermissionsRequest) =>
      APIClient.post<boolean>(
        `/api/v1/${is_folder ? 'folder' : 'file'}/${file_id}/permissions`,
        null,
        {
          headers: {
            ...perms,
            users: JSON.stringify(users),
            groups: JSON.stringify(groups),
          },
        }
      )
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: (_, params) => {
        queryClient.invalidateQueries(['file', 'permissions', params.file_id]);
        queryClient.invalidateQueries(['file', 'metadata', params.file_id]);
      },
    }
  );
};

/**
 * Remove a permission from a file
 */
export const useMutateDeleteFilePermissions = () => {
  const queryClient = useQueryClient();
  return useMutation(
    (params: { file_id: string; permission_id: string; is_folder?: boolean }) =>
      APIClient.delete(
        `/api/v1/${params.is_folder ? 'folder' : 'file'}/${
          params.file_id
        }/permissions/${params.permission_id}`
      )
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: (_, params) => {
        queryClient.invalidateQueries(['file', 'permissions', params.file_id]);
        queryClient.invalidateQueries(['file', 'metadata', params.file_id]);
      },
    }
  );
};

/**
 * Get a file version's metadata based on its file ID and version ID.
 */
export const useQueryFileVersionMetadata = (
  fileId: string,
  versionId: number
) =>
  useQuery(['file', 'version', 'metadata', fileId, versionId], () =>
    !fileId || !versionId
      ? Promise.reject()
      : APIClient.get<FileMetadata<EntryType.File>>(
          `/api/v1/file/${fileId}/version/${versionId}/metadata`
        )
          .then((res) => res.data)
          .catch(typedError)
  );

/**
 * Get the version history of a file.
 */
export const useQueryFileVersionHistory = (fileId: string) =>
  useQuery(['file', 'version', 'history', fileId], () =>
    !fileId
      ? Promise.reject()
      : APIClient.get<FileMetadata<EntryType.File>[]>(
          `/api/v1/file/${fileId}/versions`
        )
          .then((res) => res.data)
          .catch(typedError)
  );

export type DeleteFileVersionRequest = {
  file_id: string;
  version_id: number;
};

/**
 * Delete a file version by its ID.
 */
export const useMutateDeleteFileVersion = () => {
  const queryClient = useQueryClient();
  return useMutation(
    ({ file_id, version_id }: DeleteFileVersionRequest) =>
      APIClient.delete<boolean>(
        `/api/v1/file/${file_id}/versions/${version_id}`
      )
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: (_, params) => {
        queryClient.invalidateQueries([
          'file',
          'version',
          'history',
          params.file_id,
        ]);
      },
    }
  );
};

/**
 * Get a file path by its ID.
 */
export const useQueryFilePathById = (fileId: string) =>
  useQuery(['path', 'file', fileId], () =>
    !fileId
      ? Promise.reject()
      : APIClient.get<FileMetadata<EntryType.Folder>[]>(
          `/api/v1/file/${fileId}/path`
        )
          .then((res) => res.data)
          .catch(typedError)
  );

/**
 * Get a folder path by its ID.
 */
export const useQueryFolderPathById = (folderId: string) =>
  useQuery(['path', 'file', folderId], () =>
    !folderId
      ? Promise.reject()
      : APIClient.get<FileMetadata<EntryType.Folder>[]>(
          `/api/v1/folder/${folderId}/path`
        )
          .then((res) => res.data)
          .catch(typedError)
  );

/**
 * Add a file to favorites
 */
export const useMutateFavoriteFile = () => {
  const queryClient = useQueryClient();
  return useMutation(
    (params: { fileId: string; action: 'put' | 'delete' }) =>
      APIClient[params.action]<boolean>(`/api/v1/favorites/${params.fileId}`)
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: (_, params) => {
        queryClient.invalidateQueries(['is-favorite', params.fileId]);
        queryClient.invalidateQueries(['favorites']);
      },
    }
  );
};

/**
 * Check if a file is favorited
 */
export const useQueryIsFavoriteFile = (fileId: string) =>
  useQuery(
    ['is-favorite', fileId],
    () =>
      !fileId
        ? Promise.reject()
        : APIClient.get<FileMetadata>(`/api/v1/favorites/${fileId}`)
            .then((res) => res.data)
            .catch(typedError),
    { retry: false }
  );
