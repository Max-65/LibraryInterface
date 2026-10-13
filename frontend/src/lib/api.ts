const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

type RequestOptions = RequestInit;

async function request<T>(
    path: string,
    options: RequestOptions = {},
): Promise<T> {
    const response = await fetch(`${API_BASE}${path}`, {
        headers: {
            'Content-Type': 'application/json',
            ...(options.headers ?? {}),
        },
        ...options,
    });

    const text = await response.text();

    const data = text ? JSON.parse(text) : null;

    if (!response.ok) {
        throw new Error(data?.error ?? 'Request failed');
    }

    return data as T;
}

export const api = {
    get<T>(path: string) {
        return request<T>(path);
    },

    post<T>(path: string, body?: unknown) {
        return request<T>(path, {
            method: 'POST',
            body: body ? JSON.stringify(body) : undefined,
        });
    },

    patch<T>(path: string, body?: unknown) {
        return request<T>(path, {
            method: 'PATCH',
            body: body ? JSON.stringify(body) : undefined,
        });
    },
};