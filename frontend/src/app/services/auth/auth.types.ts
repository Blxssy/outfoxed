export interface User {
    id: string;
    username: string;
    email: string;
    is_guest: boolean;
    role: string;
    created_at: string;
    updated_at: string;
    last_seen_at: string;
}

export interface AuthResponse {
    user: User;
    access_token: string;
    refresh_token: string;
}

export interface RefreshResponse {
    access_token: string;
    refresh_token: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    username: string;
    email: string;
    password: string;
}
