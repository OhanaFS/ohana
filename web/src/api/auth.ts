import { useQuery } from '@tanstack/react-query';
import { APIClient, typedError } from './api';

export type WhoamiResponse = {
  user_id: string;
  name: string;
  email: string;
  home_folder_id: string;
};

export const useQueryUser = () =>
  useQuery(
    ['user'],
    () =>
      APIClient.get<WhoamiResponse>('/auth/whoami')
        .then((res) => res.data)
        .catch(typedError),
    {
      retry: 3,
      retryDelay: 100,
    }
  );

export const useQueryUsers = () =>
  useQuery(['users'], () =>
    APIClient.get<WhoamiResponse[]>('/api/v1/users')
      .then((res) => res.data)
      .catch(typedError)
  );
