import axios from 'axios';

export const APIClient = axios.create({
  baseURL: '/',
  withCredentials: true,
});
