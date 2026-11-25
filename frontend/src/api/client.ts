import type { ApiErrorBody, AuthResponse } from '../types/api';

const BASE_URL = (import.meta.env.VITE_API_URL as string | undefined) || '/api/v1';

export class ApiError extends Error {
  status: number;
  code: string;
  details?: Record<string, string>;

  constructor(status: number, body: ApiErrorBody) {
    super(body.message);
    this.status = status;
    this.code = body.code;
    this.details = body.details;
  }
}

let accessToken: string | null = null;
let refreshPromise: Promise<string | null> | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function getAccessToken(): string | null {
  return accessToken;
}

async function refreshAccessToken(): Promise<string | null> {
  if (refreshPromise) return refreshPromise;

  refreshPromise = (async () => {
    try {
      const res = await fetch(`${BASE_URL}/auth/refresh`, {
        method: 'POST',
        credentials: 'include',
      });
      if (!res.ok) {
        accessToken = null;
        return null;
      }
      const data = (await res.json()) as AuthResponse;
      accessToken = data.access_token;
      return accessToken;
    } catch {
      accessToken = null;
      return null;
    } finally {
      refreshPromise = null;
    }
  })();

  return refreshPromise;
}

interface RequestOptions extends RequestInit {
  skipAuthRefresh?: boolean;
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { skipAuthRefresh, headers, ...rest } = options;

  const makeHeaders = (token: string | null): HeadersInit => ({
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...headers,
  });

  let res = await fetch(`${BASE_URL}${path}`, {
    ...rest,
    credentials: 'include',
    headers: makeHeaders(accessToken),
  });

  if (res.status === 401 && !skipAuthRefresh && !path.startsWith('/auth/')) {
    const newToken = await refreshAccessToken();
    if (newToken) {
      res = await fetch(`${BASE_URL}${path}`, {
        ...rest,
        credentials: 'include',
        headers: makeHeaders(newToken),
      });
    }
  }

  if (!res.ok) {
    let errBody: ApiErrorBody = {
      code: 'HTTP_ERROR',
      message: `Request failed with status ${res.status}`,
    };
    try {
      const parsed = (await res.json()) as { error?: ApiErrorBody };
      if (parsed?.error) {
        errBody = parsed.error;
      }
    } catch {
      // ignore non-JSON error body
    }
    throw new ApiError(res.status, errBody);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return (await res.json()) as T;
}
