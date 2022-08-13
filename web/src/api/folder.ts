import { APIClient, typedError } from './api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { EntryType, FileMetadata, FilePermission, Permission } from './file';

/**
 * Get contents of a folder by ID.
 */
export const useQueryFolderContents = (folderId: string) =>
  useQuery(['folder', 'contents', 'id', folderId], () =>
    !folderId
      ? Promise.reject()
      : APIClient.get<FileMetadata<EntryType.Folder | EntryType.File>[]>(
          `/api/v1/folder/${folderId}`
        )
          .then((res) => res.data)
          .catch(typedError)
  );

/**
 * Get contents of a folder by its path. Use this if the ID is not known.
 */
export const useQueryFolderContentsByPath = (path: string) =>
  useQuery(['folder', 'contents', 'path', path], () =>
    !path
      ? Promise.reject()
      : APIClient.get<FileMetadata<EntryType.Folder | EntryType.File>[]>(
          `/api/v1/folder`,
          { headers: { path } }
        )
          .then((res) => res.data)
          .catch(typedError)
  );

export type UpdateFolderMetadataRequest = {
  folderId: string;
  newName: string;
  versioningMode: number;
};

/**
 * Update a folder's metadata.
 */
export const useMutateUpdateFolderMetadata = () =>
  useMutation(({ folderId, ...params }: UpdateFolderMetadataRequest) =>
    APIClient.patch<FileMetadata<EntryType.Folder>>(
      `/api/v1/folder/${folderId}`,
      null,
      { headers: { ...params } }
    )
      .then((res) => res.data)
      .catch(typedError)
  );

/**
 * Delete a folder
 */
export const useMutateDeleteFolder = () => {
  const queryClient = useQueryClient();
  return useMutation(
    (folderId: string) =>
      APIClient.delete<boolean>(`/api/v1/folder/${folderId}`)
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: () => {
        queryClient.clear();
      },
    }
  );
};

export type CreateFolderRequest = {
  folder_name: string;
} & ({ parent_folder_id: string } | { parent_folder_path: string });

/**
 * Create a new folder
 */
export const useMutateCreateFolder = () => {
  const queryClient = useQueryClient();
  return useMutation(
    (params: CreateFolderRequest) =>
      APIClient.post<{ file_id: string }>(`/api/v1/folder`, null, {
        headers: { ...params },
      })
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: (_, params) => {
        if ('parent_folder_id' in params)
          queryClient.invalidateQueries([
            'folder',
            'contents',
            'id',
            params.parent_folder_id,
          ]);
      },
    }
  );
};

/**
 * View folder permissions
 */
export const useQueryFolderPermissions = (folderId: string) =>
  useQuery(['folder', 'permissions', 'id', folderId], () =>
    !folderId
      ? Promise.reject()
      : APIClient.get<FilePermission[]>(
          `/api/v1/folder/${folderId}/permissions`
        )
          .then((res) => res.data)
          .catch(typedError)
  );

export type UpdateFolderPermissionsRequest = {
  folder_id: string;
  permission_id: number;
} & Permission;

/**
 * Modify folder permissions
 */
export const useMutateUpdateFolderPermissions = () =>
  useMutation(({ folder_id, ...params }: UpdateFolderPermissionsRequest) =>
    APIClient.patch<boolean>(`/api/v1/folder/${folder_id}/permissions`, null, {
      headers: { ...params },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

export type AddFolderPermissionsRequest = Permission & {
  folder_id: string;
  users: string[];
  groups: string[];
};

/**
 * Add folder permissions
 */
export const useMutateAddFolderPermissions = () =>
  useMutation(
    ({ folder_id, users, groups, ...perms }: AddFolderPermissionsRequest) =>
      APIClient.post<boolean>(`/api/v1/folder/${folder_id}/permissions`, null, {
        headers: {
          ...perms,
          users: JSON.stringify(users),
          groups: JSON.stringify(groups),
        },
      })
        .then((res) => res.data)
        .catch(typedError)
  );

export type MoveFolderRequest = {
  folder_id: string;
  new_folder_id: string;
};

/**
 * Move a folder to a new location
 */
export const useMutateMoveFolder = () =>
  useMutation(({ folder_id, new_folder_id }: MoveFolderRequest) =>
    APIClient.post<boolean>(`/api/v1/folder/${folder_id}/move`, null, {
      headers: { new_folder_id },
    })
      .then((res) => res.data)
      .catch(typedError)
  );

/**
 * View details of a folder
 */
export const useQueryFolderDetails = (folderId: string) =>
  useQuery(['folder', 'details', 'id', folderId], () =>
    !folderId
      ? Promise.reject()
      : APIClient.get<FileMetadata<EntryType.Folder>>(
          `/api/v1/folder/${folderId}/details`
        )
          .then((res) => res.data)
          .catch(typedError)
  );
