import { Injectable } from '@angular/core';
import {
    AuthResponse,
    LoginRequest,
    RefreshResponse,
    RegisterRequest,
    User,
} from './auth.types';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { TokenService } from './token.service';

@Injectable({ providedIn: 'root' })
export class AuthService {
    private api = 'http://localhost:8080/api/v1/auth';

    private userSubject = new BehaviorSubject<User | null>(null);
    user$ = this.userSubject.asObservable();

    constructor(
        private http: HttpClient,
        private tokenService: TokenService,
    ) {}

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
                }),
            );
    }

    logout(): void {
        this.tokenService.clear();
        this.userSubject.next(null);
    }

    private handleAuth(res: AuthResponse): void {
        this.tokenService.setTokens(res.access_token, res.refresh_token);

        this.userSubject.next(res.user);
    }
}
