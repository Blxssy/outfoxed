import { Component } from '@angular/core';
import {
    FormControl,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { ButtonComponent } from '@fox/ui-kit/button';
import { InputComponent } from '@fox/ui-kit/input';
import { RouterLink } from '@angular/router';

@Component({
    selector: 'app-register',
    imports: [ReactiveFormsModule, ButtonComponent, InputComponent, RouterLink],
    templateUrl: './register.component.html',
    styleUrl: './register.component.scss',
})
export class RegisterComponent {
    registerForm = new FormGroup({
        nickName: new FormControl(''),
        email: new FormControl('', [Validators.required, Validators.email]),
        password: new FormControl('', [
            Validators.required,
            Validators.minLength(8),
        ]),
    });

    onSubmit() {
        // мб пригодится для проверки формы
        console.log(this.registerForm.value);
    }
}
