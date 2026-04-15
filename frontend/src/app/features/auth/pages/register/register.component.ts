import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ButtonComponent } from '@fox/ui-kit/button';
import { InputComponent } from '@fox/ui-kit/input';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from 'src/app/services/auth/auth.service';

@Component({
    selector: 'app-register',
    imports: [ReactiveFormsModule, ButtonComponent, InputComponent, RouterLink],
    templateUrl: './register.component.html',
    styleUrl: './register.component.scss',
})
export class RegisterComponent {
    private readonly fb = inject(FormBuilder);
    private readonly authService = inject(AuthService);
    private readonly router = inject(Router);

    readonly registerForm = this.fb.group({
        username: ['', [Validators.required]],
        email: ['', [Validators.required, Validators.email]],
        password: ['', [Validators.required]],
        confirmPassword: ['', [Validators.required]],
    });
    // to do: валидаторы на сильный пароль

    onSubmit() {
        if (this.registerForm.invalid) {
            this.registerForm.markAllAsTouched();
            // to do: настроить ошибки формы
            console.log('submit fired');
            return;
        }

        const { username, email, password, confirmPassword } =
            this.registerForm.getRawValue();

        if (password !== confirmPassword) {
            this.registerForm.get('confirmPassword')?.setErrors({
                mismatch: true,
            });
            return;
        }

        this.authService
            .register({
                username: username!,
                email: email!,
                password: password!,
            })
            .subscribe({
                next: () => {
                    this.router.navigate(['/lobby']);
                },
                error: (err) => {
                    console.error('register error:', err);
                },
            });
    }
}
