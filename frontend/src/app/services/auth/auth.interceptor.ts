import { Injectable } from '@angular/core';
import {
    HttpEvent,
    HttpHandler,
    HttpInterceptor,
    HttpRequest,
    HttpErrorResponse,
} from '@angular/common/http';

import { Observable, catchError, switchMap, throwError } from 'rxjs';

import { TokenService } from './token.service';
import { AuthService } from './auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
    constructor(
        private tokenService: TokenService,
        private authService: AuthService,
    ) {}

    intercept(
        req: HttpRequest<any>,
        next: HttpHandler,
    ): Observable<HttpEvent<any>> {
        const token = this.tokenService.getAccessToken();

        let authReq = req;

        if (token) {
            authReq = req.clone({
                setHeaders: {
                    Authorization: `Bearer ${token}`,
                },
            });
        }

        return next.handle(authReq).pipe(
            catchError((error: HttpErrorResponse) => {
                if (error.status === 401) {
                    return this.authService.refresh().pipe(
                        switchMap((tokens) => {
                            const retryReq = req.clone({
                                setHeaders: {
                                    Authorization: `Bearer ${tokens.access_token}`,
                                },
                            });

                            return next.handle(retryReq);
                        }),
                    );
                }

                return throwError(() => error);
            }),
        );
    }
}
