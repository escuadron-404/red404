export interface HttpClient {
  get<T>(path: string): Promise<T>;
  post<T>(path: string, body: Record<string, unknown>): Promise<T>;
  put<T>(path: string, body: Record<string, unknown>): Promise<T>;
  delete<T>(path: string): Promise<T>;
}

export interface TokenProvider {
  getToken(): string;
  setToken(token: string): void;
  removeToken(): void;
}

export class AuthClient implements HttpClient {
  private baseUrl: string;
  private tokenProvider?: TokenProvider;

  constructor(baseUrl: string = import.meta.env.VITE_API_URL) {
    this.baseUrl = baseUrl;
  }

  private headers(): HeadersInit {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    const token = this.tokenProvider?.getToken();
    if (token) headers.Authorization = `Bearer ${token}`;
    return headers;
  }

  public async get<T>(path: string): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: "GET",
      headers: this.headers(),
    });
    return res.json();
  }

  public async post<T, B>(path: string, body: B): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: "POST",
      headers: this.headers(),
      body: JSON.stringify(body),
    });
    return res.json();
  }

  public async put<T, B>(path: string, body: B): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: "PUT",
      headers: this.headers(),
      body: JSON.stringify(body),
    });
    return res.json();
  }

  public async patch<T, B>(path: string, body: B): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: "PATCH",
      headers: this.headers(),
      body: JSON.stringify(body),
    });
    return res.json();
  }

  public async delete<T>(path: string): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: "DELETE",
      headers: this.headers(),
    });
    return res.json();
  }
}

export const authClient = new AuthClient(
  import.meta.env.VITE_API_URL || "http://localhost:8080",
);
