import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';
import { TokenService } from '../services/auth/token.service';

@Injectable({ providedIn: 'root' })
export class AuthGuard implements CanActivate {
    constructor(
        private tokenService: TokenService,
        private router: Router,
    ) {}

    canActivate(): boolean {
        if (this.tokenService.isLoggedIn()) {
            return true;
        }

        this.router.navigate(['/auth/login']);
        return false;
    }
}
