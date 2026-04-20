import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { Validators } from '@angular/forms';
import { InputComponent } from '@fox/ui-kit/input';
import { ButtonComponent } from '@fox/ui-kit/button';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from 'src/app/services/auth/auth.service';
import { TokenService } from 'src/app/services/auth/token.service';
import { CardComponent } from '@fox/ui-kit/card';

@Component({
    selector: 'app-login',
    imports: [
        ReactiveFormsModule,
        ButtonComponent,
        InputComponent,
        RouterLink,
        CardComponent,
    ],
    templateUrl: './login.component.html',
    styleUrl: './login.component.scss',
})
export class LoginComponent {
    private readonly fb = inject(FormBuilder);
    private readonly authService = inject(AuthService);
    private readonly router = inject(Router);
    private readonly tokenService = inject(TokenService);

    errorMessage = '';

    readonly loginForm = this.fb.group({
        email: ['', [Validators.required]],
        password: ['', [Validators.required]],
    });
    // to do: валидаторы на сильный пароль

    ngOnInit() {
        this.loginForm.valueChanges.subscribe(() => {
            this.errorMessage = '';
        });
    }

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
                    this.errorMessage = this.getLoginErrorMessage(err);
                    console.error(this.errorMessage);
                },
            });
    }

    private getLoginErrorMessage(err: any): string {
        if (err.status === 400) {
            return 'Некорректные данные';
        }

        if (err.status === 401) {
            return 'Неверный e-mail или пароль';
        }

        if (err.status === 500) {
            return 'Ошибка сервера';
        }

        if (err.status === 0) {
            return 'Сервер недоступен';
        }

        return err.error || 'Что-то пошло не так';
    }
}
