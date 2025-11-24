export interface User {
  id: number;
  username: string;
  email: string;
  avatarUrl?: string;
  aboutText?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Forum {
  id: number;
  name: string;
  description: string;
  slug: string;
  threadCount: number;
  postCount: number;
}

export interface Thread {
  id: number;
  forumId: number;
  userId: number;
  title: string;
  slug: string;
  isLocked: boolean;
  isPinned: boolean;
  viewCount: number;
  replyCount: number;
  lastPostAt?: string;
  createdAt: string;
  author?: User;
  forum?: Forum;
}

export interface Post {
  id: number;
  threadId: number;
  userId: number;
  postNumber: number;
  content: string;
  createdAt: string;
  updatedAt?: string;
  author?: User;
}

export interface Pagination {
  page: number;
  limit: number;
  totalPages: number;
  totalItems: number;
}

export interface ThreadListResponse {
  threads: Thread[];
  pagination: Pagination;
}

export interface PostListResponse {
  posts: Post[];
  pagination: Pagination;
}

export interface AuthResponse {
  accessToken: string;
  user: User;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface ThreadCreateRequest {
  forumId: number;
  title: string;
  content: string;
}

export interface PostCreateRequest {
  content: string;
}

