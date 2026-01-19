import { Component } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import {Validators} from '@angular/forms';

@Component({
    selector: 'app-login',
    imports: [ReactiveFormsModule],
    templateUrl: './login.component.html',
    styleUrl: './login.component.scss',
})
export class LoginComponent {
    loginForm = new FormGroup({
        nickName: new FormControl(''),
        email: new FormControl('', [Validators.required, Validators.email]),
        password: new FormControl('', [Validators.required, Validators.minLength(8)]),
    });

    onSubmit() {
        // мб пригодится для проверки формы
        console.log(this.loginForm.value);
    }
    
}
