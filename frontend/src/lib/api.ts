export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

export type User = {
  id: string;
  email: string;
  name: string;
  picture?: string;
  provider?: string;
};

async function request(path: string, options: RequestInit = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers ?? {}),
    },
    ...options,
  });
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || 'Request failed');
  }
  if (response.headers.get('content-type')?.includes('application/json')) {
    return response.json();
  }
  return null;
}

export async function fetchCurrentUser(): Promise<User | null> {
  const response = await fetch(`${API_BASE_URL}/me`, {
    credentials: 'include',
  });
  if (!response.ok) {
    return null;
  }
  return response.json();
}

export async function registerWithEmail(payload: { email: string; password: string; name: string }) {
  return request('/auth/email/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function loginWithEmail(payload: { email: string; password: string }) {
  return request('/auth/email/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
