import { useQuery } from '@tanstack/react-query';
import { APIClient } from './api';

type WhoamiResponse = {
  user_id: string;
  name: string;
  email: string;
};

export const useQueryUser = () =>
  useQuery(['user'], async () =>
    APIClient.get<WhoamiResponse>('/auth/whoami').then((res) => res.data)
  );
