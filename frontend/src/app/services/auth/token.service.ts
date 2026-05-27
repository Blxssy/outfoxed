import { Injectable } from '@angular/core';

const ACCESS_KEY = 'access_token';
const REFRESH_KEY = 'refresh_token';

const EXPIRY_BUFFER_SEC = 30;

@Injectable({ providedIn: 'root' })
export class TokenService {
    getAccessToken(): string | null {
        return localStorage.getItem(ACCESS_KEY);
    }

    getRefreshToken(): string | null {
        return localStorage.getItem(REFRESH_KEY);
    }

    setTokens(access: string, refresh: string): void {
        localStorage.setItem(ACCESS_KEY, access);
        localStorage.setItem(REFRESH_KEY, refresh);
    }

    clear(): void {
        localStorage.removeItem(ACCESS_KEY);
        localStorage.removeItem(REFRESH_KEY);
    }

    isLoggedIn(): boolean {
        return !!this.getAccessToken() && !this.isAccessTokenExpired();
    }

    isAccessTokenExpired(): boolean {
        const token = this.getAccessToken();

        if (!token) {
            return true;
        }

        try {
            const payload = JSON.parse(atob(token.split('.')[1]));

            return Date.now() >= payload.exp * 1000 - EXPIRY_BUFFER_SEC * 1000;
        } catch {
            return true;
        }
    }

    msUntilExpiry(): number {
        const token = this.getAccessToken();
        if (!token) return 0;
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return Math.max(0, payload.exp * 1000 - Date.now());
        } catch {
            return 0;
        }
    }
}
