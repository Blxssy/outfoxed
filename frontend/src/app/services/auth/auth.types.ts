export type User = {
    id: string;
    username: string;
    email: string;
    is_guest: boolean;
    role: string;
    created_at: string;
    updated_at: string;
    last_seen_at: string;
};

export type AuthResponse = {
    user: User;
    access_token: string;
    refresh_token: string;
};

export type RefreshResponse = {
    access_token: string;
    refresh_token: string;
};

export type LoginRequest = {
    email: string;
    password: string;
};

export type RegisterRequest = {
    username: string;
    email: string;
    password: string;
};
