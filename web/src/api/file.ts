import { APIClient, typedError } from './api';
import { useMutation, useQuery } from '@tanstack/react-query';

export type FileUploadRequest = {
  file: File;
  folder_id: string;
  file_name?: string;
  frag_count: number;
  parity_count: number;
};

export type FileMetadata = {
  file_id: string;
  file_name: string;
  mime_type: string;
  entry_type: number;
  parent_folder_id: string;
  version_no: number;
  data_version_no: number;
  size: number;
  actual_size: number;
  created_time: string; // date
  modified_user_user_id: string;
  modified_time: string; // date
  versioning_mode: number;
  checksum: string;
  frag_count: number;
  parity_count: number;
  password_protected: boolean;
  link_file_id: string;
  last_checked: string; // date
  status: number;
};

/**
 * Uploads a file to a folder.
 */
export const useMutateUploadFile = () =>
  useMutation(({ file, ...headers }: FileUploadRequest) => {
    if (!headers.file_name) headers.file_name = file.name;
    return APIClient.post<FileMetadata>('/api/v1/file', file, {
      headers: { ...headers },
    })
      .then((res) => res.data)
      .catch(typedError);
  });

export type FileUpdateRequest = {
  file: File;
  file_id: string;
  frag_count: number;
  parity_count: number;
};

/**
 * Upload a newer version of a file.
 */
export const useMutateUpdateFile = () =>
  useMutation(({ file, file_id, ...headers }: FileUpdateRequest) => {
    return APIClient.post<FileMetadata>(
      `/api/v1/file/${file_id}/update`,
      file,
      { headers: { ...headers } }
    )
      .then((res) => res.data)
      .catch(typedError);
  });

/**
 * Get a file's metadata.
 */
export const useQueryFileMetadata = (fileId: string) =>
  useQuery(['file-metadata', fileId], () =>
    APIClient.get<FileMetadata>(`/api/v1/file/${fileId}/metadata`)
      .then((res) => res.data)
      .catch(typedError)
  );

export type FileMetadataUpdateRequest = {
  file_id: string;
  file_name: string;
  mime_type: string;
  versioning_mode: number;
  password_modification: boolean;
  password_protected: boolean;
  password_hint: string;
  old_password: string;
  new_password: string;
};

/**
 * Update a file's metadata.
 */
export const useMutateUpdateFileMetadata = () =>
  useMutation(({ file_id, ...body }: FileMetadataUpdateRequest) =>
    APIClient.patch<FileMetadata>(`/api/v1/file/${file_id}/metadata`, body)
      .then((res) => res.data)
      .catch(typedError)
  );

export type MoveFileRequest = {
  /** The ID of the file to move */
  file_id: string;
  /** The ID of the folder to move the file to */
  folder_id: string;
};

/**
 * Move a file to a new folder.
 */
export const useMutateMoveFile = () =>
  useMutation(({ file_id, folder_id }: MoveFileRequest) =>
    APIClient.post<boolean>(`/api/v1/file/${file_id}/move`, { folder_id })
      .then((res) => res.data)
      .catch(typedError)
  );

/**
 * Copy a file to a new folder.
 */
export const useMutateCopyFile = () =>
  useMutation(({ file_id, folder_id }: MoveFileRequest) =>
    APIClient.post<boolean>(`/api/v1/file/${file_id}/copy`, { folder_id })
      .then((res) => res.data)
      .catch(typedError)
  );

/**
 * Get the download URL for a file.
 */
export const getFileDownloadURL = (fileId: string) =>
  window.location.origin + `/api/v1/file/${fileId}`;
