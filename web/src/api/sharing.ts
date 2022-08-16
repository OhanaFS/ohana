import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { APIClient, typedError } from './api';
import { FileMetadata } from './file';

export type SharedLink = {
  shortened_link: string;
  file_id: string;
  created_time: string;
};

/**
 * Gets a list of sharing links for a file
 */
export const useQueryFileSharingLinks = (fileId: string) =>
  useQuery(['sharing-links', fileId], () =>
    APIClient.get<SharedLink[]>(`/api/v1/file/${fileId}/share`)
      .then((res) => res.data)
      .catch(typedError)
  );

/**
 * Gets metadata for a sharing link
 */
export const useQuerySharingLinkMetadata = (linkId: string) =>
  useQuery(
    ['sharing-link', 'meta', linkId],
    () =>
      !linkId
        ? Promise.reject()
        : APIClient.get<FileMetadata>(`/api/v1/shared/${linkId}/metadata`)
            .then((res) => res.data)
            .catch(typedError),
    {
      retry: 3,
      retryDelay: 100,
      keepPreviousData: true,
    }
  );

export type CreateSharingLinkParams = {
  fileId: string;
  link?: string;
};

/**
 * Creates a new sharing link
 */
export const useMutateCreateSharingLink = () => {
  const queryClient = useQueryClient();
  return useMutation(
    (params: CreateSharingLinkParams) =>
      APIClient.post<SharedLink>(
        `/api/v1/file/${params.fileId}/share` +
          (params.link ? '/' + encodeURIComponent(params.link) : '')
      )
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: (_, params) => {
        queryClient.invalidateQueries(['sharing-links', params.fileId]);
      },
    }
  );
};

export type DeleteSharingLinkParams = {
  fileId: string;
  link: string;
};

/**
 * Removes an existing sharing link
 */
export const useMutateDeleteSharingLink = () => {
  const queryClient = useQueryClient();
  return useMutation(
    (params: DeleteSharingLinkParams) =>
      APIClient.delete(`/api/v1/file/${params.fileId}/share/${params.link}`)
        .then((res) => res.data)
        .catch(typedError),
    {
      onSuccess: (_, params) => {
        queryClient.invalidateQueries(['sharing-links', params.fileId]);
      },
    }
  );
};

/**
 * Get the publicly-accessible link from a shortened_link
 */
export const getSharingLinkURL = (
  shortenedLink: string,
  linkType: 'preview' | 'inline' | 'download'
) =>
  [
    window.location.origin,
    linkType === 'preview' ? '/share/' : '/api/v1/shared/',
    encodeURIComponent(shortenedLink),
    linkType === 'inline' ? '?inline=true' : '',
  ].join('');
