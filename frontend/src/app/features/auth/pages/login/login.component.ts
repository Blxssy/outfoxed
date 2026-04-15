import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { Validators } from '@angular/forms';
import { InputComponent } from '@fox/ui-kit/input';
import { ButtonComponent } from '@fox/ui-kit/button';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from 'src/app/services/auth/auth.service';
import { TokenService } from 'src/app/services/auth/token.service';

@Component({
    selector: 'app-login',
    imports: [ReactiveFormsModule, ButtonComponent, InputComponent, RouterLink],
    templateUrl: './login.component.html',
    styleUrl: './login.component.scss',
})
export class LoginComponent {
    private readonly fb = inject(FormBuilder);
    private readonly authService = inject(AuthService);
    private readonly router = inject(Router);
    private readonly tokenService = inject(TokenService);

    readonly loginForm = this.fb.group({
        email: ['', [Validators.required]],
        password: ['', [Validators.required]],
    });
    // to do: валидаторы на сильный пароль

    onSubmit() {
        if (this.loginForm.invalid) {
            this.loginForm.markAllAsTouched();
            // to do: настроить ошибки формы
            return;
        }

        const { email, password } = this.loginForm.getRawValue();

        this.authService
            .login({ email: email!, password: password! })
            .subscribe({
                next: (res: any) => {
                    this.tokenService.setTokens(
                        res.accessToken,
                        res.refreshToken,
                    );

                    this.router.navigate(['/lobby']);
                },
                error: (err) => {
                    console.log('login error:', err);
                },
            });
    }
}
