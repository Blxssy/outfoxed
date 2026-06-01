import { Component, inject, signal } from '@angular/core';
import {
    Router,
    RouterLink,
    RouterLinkActive,
    RouterOutlet,
} from '@angular/router';
import { AuthService } from 'src/app/services/auth/auth.service';
import { ButtonComponent } from '@fox/ui-kit/button';

@Component({
    selector: 'app-auth',
    imports: [RouterOutlet, RouterLink, RouterLinkActive, ButtonComponent],
    templateUrl: './auth.component.html',
    styleUrl: './auth.component.scss',
})
export class AuthComponent {
    private readonly authService = inject(AuthService);
    private readonly router = inject(Router);

    readonly isGuestLoading = signal(false);
    readonly guestError = signal('');

    loginAsGuest(): void {
        this.isGuestLoading.set(true);
        this.guestError.set('');

        this.authService.loginAsGuest().subscribe({
            next: () => {
                this.isGuestLoading.set(false);
                this.router.navigate(['/lobby']);
            },
            error: (err) => {
                this.isGuestLoading.set(false);
                this.guestError.set(
                    err.status === 0
                        ? 'Сервер недоступен'
                        : 'Не удалось войти как гость',
                );
            },
        });
    }
}
