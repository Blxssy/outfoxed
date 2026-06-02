import { inject, Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';
import { TokenService } from '../services/auth/token.service';
import { AuthService } from '../services/auth/auth.service';
import { catchError, map, Observable, of } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AuthGuard implements CanActivate {
    private readonly tokenService = inject(TokenService);
    private readonly authService = inject(AuthService);
    private readonly router = inject(Router);

    canActivate(): Observable<boolean> | boolean {
        if (this.tokenService.isLoggedIn()) {
            return true;
        }

        const refresh = this.tokenService.getRefreshToken();
        if (refresh) {
            return this.authService.refresh().pipe(
                map(() => true),
                catchError(() => {
                    this.tokenService.clear();
                    this.router.navigate(['/auth/login']);
                    return of(false);
                }),
            );
        }

        this.router.navigate(['/auth/login']);
        return false;
    }
}
