import axios from 'axios';

export const APIClient = axios.create({
  baseURL: '/',
  withCredentials: true,
});

export type ErrorResponse = {
  code: string;
  message: string;
};

/**
 * typedError re-throws an error as the ErrorResponse type.
 */
export const typedError = (err: any) => {
  throw axios.isAxiosError(err) && err.response
    ? (err.response.data as ErrorResponse)
    : { code: '-1', message: err.toString() };
};

export const qs = (obj: any) =>
  Object.keys(obj)
    .map((key) => [key, obj[key]].map(encodeURIComponent).join('='))
    .join('&');
