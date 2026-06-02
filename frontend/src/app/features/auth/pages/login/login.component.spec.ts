import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { throwError } from 'rxjs';
import { LoginComponent } from './login.component';
import { AuthService } from 'src/app/services/auth/auth.service';
import { TokenService } from 'src/app/services/auth/token.service';

const mockAuthService = { login: vi.fn() };
const mockTokenService = { setTokens: vi.fn() };

describe('LoginComponent', () => {
    beforeEach(async () => {
        vi.clearAllMocks();
        await TestBed.configureTestingModule({
            imports: [LoginComponent],
            providers: [
                provideRouter([]),
                { provide: AuthService, useValue: mockAuthService },
                { provide: TokenService, useValue: mockTokenService },
            ],
        }).compileComponents();
    });

    it('should be created without errors', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        expect(fixture.componentInstance).toBeTruthy();
    });

    it('should form contain email and password fields', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        const { loginForm } = fixture.componentInstance;
        expect(loginForm.contains('email')).toBe(true);
        expect(loginForm.contains('password')).toBe(true);
    });

    it('should form be invalid during initialization', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        expect(fixture.componentInstance.loginForm.invalid).toBe(true);
    });

    it('should errorMessage be empty during initialization', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        expect(fixture.componentInstance.errorMessage).toBe('');
    });

    it('should email field be required', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        const ctrl = fixture.componentInstance.loginForm.get('email')!;
        ctrl.setValue('');
        expect(ctrl.errors?.['required']).toBeTruthy();
    });

    it('should password field be required', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        const ctrl = fixture.componentInstance.loginForm.get('password')!;
        ctrl.setValue('');
        expect(ctrl.errors?.['required']).toBeTruthy();
    });

    it('should form be valid with correct data', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        fixture.componentInstance.loginForm.setValue({
            email: 'user@test.com',
            password: 'pass123',
        });
        expect(fixture.componentInstance.loginForm.valid).toBe(true);
    });

    it('should errorMessage reset when the form is changed', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        fixture.detectChanges();
        const comp = fixture.componentInstance;
        comp.errorMessage = 'Ошибка';
        comp.loginForm.get('email')!.setValue('new@value.com');
        expect(comp.errorMessage).toBe('');
    });

    const errorCases = [
        { status: 400, expected: 'Некорректные данные' },
        { status: 401, expected: 'Неверный e-mail или пароль' },
        { status: 500, expected: 'Ошибка сервера' },
        { status: 0, expected: 'Сервер недоступен' },
        { status: 999, error: 'Custom error', expected: 'Custom error' },
        { status: 999, error: undefined, expected: 'Что-то пошло не так' },
    ];

    errorCases.forEach(({ status, error, expected }) => {
        it(`should set errorMessage "${expected}" during status=${status}`, () => {
            mockAuthService.login.mockReturnValue(
                throwError(() => ({ status, error })),
            );
            const fixture = TestBed.createComponent(LoginComponent);
            fixture.detectChanges();
            const comp = fixture.componentInstance;
            comp.loginForm.setValue({ email: 'u@t.com', password: 'p' });
            comp.onSubmit();
            expect(comp.errorMessage).toBe(expected);
        });
    });

    it('should .error-message block be hidden if errorMessage is empty', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        fixture.detectChanges();
        expect(
            fixture.nativeElement.querySelector('.error-message'),
        ).toBeNull();
    });

    it('should .error-message block displayed with the error text', () => {
        const fixture = TestBed.createComponent(LoginComponent);
        fixture.componentInstance.errorMessage = 'Неверный e-mail или пароль';
        fixture.detectChanges();
        const el = fixture.nativeElement.querySelector('.error-message');
        expect(el).toBeTruthy();
        expect(el.textContent.trim()).toBe('Неверный e-mail или пароль');
    });
});
