import { inject, Injectable, OnDestroy } from '@angular/core';
import {
    AuthResponse,
    LoginRequest,
    RefreshResponse,
    RegisterRequest,
    User,
} from './auth.types';
import {
    BehaviorSubject,
    Observable,
    Subscription,
    switchMap,
    tap,
    timer,
} from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { TokenService } from './token.service';

const REFRESH_BEFORE_MS = 60_000;

@Injectable({ providedIn: 'root' })
export class AuthService implements OnDestroy {
    private readonly http = inject(HttpClient);
    private readonly tokenService = inject(TokenService);

    private api = '/api/v1/auth';

    private userSubject = new BehaviorSubject<User | null>(null);
    user$ = this.userSubject.asObservable();

    private refreshTimer: Subscription | null = null;

    login(data: LoginRequest): Observable<AuthResponse> {
        return this.http
            .post<AuthResponse>(`${this.api}/login`, data)
            .pipe(tap((res) => this.handleAuth(res)));
    }

    register(data: RegisterRequest): Observable<AuthResponse> {
        return this.http
            .post<AuthResponse>(`${this.api}/register`, data)
            .pipe(tap((res) => this.handleAuth(res)));
    }

    refresh(): Observable<RefreshResponse> {
        const refresh = this.tokenService.getRefreshToken();

        return this.http
            .post<RefreshResponse>(`${this.api}/refresh`, {
                refresh_token: refresh,
            })
            .pipe(
                tap((tokens) => {
                    this.tokenService.setTokens(
                        tokens.access_token,
                        tokens.refresh_token,
                    );

                    this.scheduleTokenRefresh();
                }),
            );
    }

    logout(): void {
        this.cancelTokenRefresh();
        this.tokenService.clear();
        this.userSubject.next(null);
    }

    restoreSession(): void {
        if (this.tokenService.isLoggedIn()) {
            this.scheduleTokenRefresh();
        }
    }

    ngOnDestroy(): void {
        this.cancelTokenRefresh();
    }

    private handleAuth(res: AuthResponse): void {
        this.tokenService.setTokens(res.access_token, res.refresh_token);
        this.userSubject.next(res.user);
        this.scheduleTokenRefresh();
    }

    private scheduleTokenRefresh(): void {
        this.cancelTokenRefresh();

        const msLeft = this.tokenService.msUntilExpiry();
        const delay = Math.max(0, msLeft - REFRESH_BEFORE_MS);

        if (delay === 0 && msLeft === 0) {
            return;
        }

        this.refreshTimer = timer(delay)
            .pipe(switchMap(() => this.refresh()))
            .subscribe({
                error: (err) => {
                    console.warn(err);
                },
            });
    }

    private cancelTokenRefresh(): void {
        this.refreshTimer?.unsubscribe();
        this.refreshTimer = null;
    }
}
