import axios from './axios';
import {
  User,
  Forum,
  Thread,
  Post,
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  ThreadCreateRequest,
  PostCreateRequest,
  ThreadListResponse,
  PostListResponse,
} from '../types';

// Auth API
export const authAPI = {
  login: (data: LoginRequest) =>
    axios.post<AuthResponse>('/auth/login', data),
  
  register: (data: RegisterRequest) =>
    axios.post<AuthResponse>('/users/register', data),
  
  getCurrentUser: () =>
    axios.get<User>('/users/me'),
};

// Forum API
export const forumAPI = {
  getAll: () =>
    axios.get<Forum[]>('/forums'),
  
  getThreads: (slug: string, page: number = 1, limit: number = 25) =>
    axios.get<ThreadListResponse>(`/forums/${slug}/threads`, {
      params: { page, limit },
    }),
};

// Thread API
export const threadAPI = {
  getById: (id: number) =>
    axios.get<Thread>(`/threads/${id}`),
  
  create: (data: ThreadCreateRequest) =>
    axios.post<Thread>('/threads', data),
  
  getPosts: (id: number, page: number = 1, limit: number = 25) =>
    axios.get<PostListResponse>(`/threads/${id}/posts`, {
      params: { page, limit },
    }),
};

// Post API
export const postAPI = {
  create: (threadId: number, data: PostCreateRequest) =>
    axios.post<Post>(`/threads/${threadId}/posts`, data),
  
  update: (id: number, data: PostCreateRequest) =>
    axios.patch<Post>(`/posts/${id}`, data),
  
  delete: (id: number) =>
    axios.delete(`/posts/${id}`),
};

// User API
export const userAPI = {
  updateProfile: (data: { avatarUrl?: string; aboutText?: string }) =>
    axios.patch<User>('/users/me', data),
};

