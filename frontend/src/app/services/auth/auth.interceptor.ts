import { inject, Injectable } from '@angular/core';
import {
    HttpEvent,
    HttpHandler,
    HttpInterceptor,
    HttpRequest,
    HttpErrorResponse,
} from '@angular/common/http';

import {
    BehaviorSubject,
    Observable,
    catchError,
    filter,
    switchMap,
    take,
    throwError,
} from 'rxjs';

import { TokenService } from './token.service';
import { AuthService } from './auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
    private isRefreshing = false;
    private refreshToken$ = new BehaviorSubject<string | null>(null);

    tokenService = inject(TokenService);
    authService = inject(AuthService);

    intercept(
        req: HttpRequest<unknown>,
        next: HttpHandler,
    ): Observable<HttpEvent<unknown>> {
        if (req.url.includes('/auth/refresh')) {
            return next.handle(req);
        }

        const token = this.tokenService.getAccessToken();
        const authReq = token ? this.addToken(req, token) : req;

        return next.handle(authReq).pipe(
            catchError((error: HttpErrorResponse) => {
                if (error.status === 401) {
                    return this.handle401(req, next);
                }
                return throwError(() => error);
            }),
        );
    }

    private handle401(
        req: HttpRequest<unknown>,
        next: HttpHandler,
    ): Observable<HttpEvent<unknown>> {
        if (this.isRefreshing) {
            return this.refreshToken$.pipe(
                filter((token) => token !== null),
                take(1),
                switchMap((token) => next.handle(this.addToken(req, token!))),
            );
        }

        this.isRefreshing = true;
        this.refreshToken$.next(null);

        return this.authService.refresh().pipe(
            switchMap((res) => {
                this.isRefreshing = false;
                this.refreshToken$.next(res.access_token);

                return next.handle(this.addToken(req, res.access_token));
            }),
            catchError((err) => {
                this.isRefreshing = false;
                this.authService.logout();
                return throwError(() => err);
            }),
        );
    }

    private addToken(
        req: HttpRequest<unknown>,
        token: string,
    ): HttpRequest<unknown> {
        return req.clone({
            setHeaders: { Authorization: `Bearer ${token}` },
        });
    }
}
