import { Injectable } from '@angular/core';
import { LoginRequest, RegisterRequest } from './auth.types';

@Injectable({ providedIn: 'root' })
export class AuthService {
    login(data: LoginRequest) {}

    register(data: RegisterRequest) {}

    logout() {}
}
