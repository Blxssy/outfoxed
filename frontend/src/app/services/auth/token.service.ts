import { Injectable } from '@angular/core';

const ACCESS_KEY = 'access_token';
const REFRESH_KEY = 'refresh_token';

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
        return !!this.getAccessToken();
    }
}
