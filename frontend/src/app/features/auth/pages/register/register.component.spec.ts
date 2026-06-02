import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';

import { RegisterComponent } from './register.component';
import { AuthService } from 'src/app/services/auth/auth.service';
import { provideHttpClient } from '@angular/common/http';

const mockAuthService = { register: vi.fn() };

describe('RegisterComponent', () => {
    beforeEach(async () => {
        vi.clearAllMocks();
        await TestBed.configureTestingModule({
            imports: [RegisterComponent],
            providers: [
                provideRouter([]),
                { provide: AuthService, useValue: mockAuthService },
                provideHttpClient(),
            ],
        }).compileComponents();
    });

    it('should be created without errors', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        expect(fixture.componentInstance).toBeTruthy();
    });

    it('should form contain username, email, password, confirmPassword fields', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        const { registerForm } = fixture.componentInstance;
        expect(registerForm.contains('username')).toBe(true);
        expect(registerForm.contains('email')).toBe(true);
        expect(registerForm.contains('password')).toBe(true);
        expect(registerForm.contains('confirmPassword')).toBe(true);
    });

    it('should form be invalid during initialization', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        expect(fixture.componentInstance.registerForm.invalid).toBe(true);
    });

    it('should errorMessage be empty during initialization', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        expect(fixture.componentInstance.errorMessage).toBe('');
    });

    it('should username field be required', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        const ctrl = fixture.componentInstance.registerForm.get('username')!;
        ctrl.setValue('');
        expect(ctrl.errors?.['required']).toBeTruthy();
    });

    it('should email field be required', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        const ctrl = fixture.componentInstance.registerForm.get('email')!;
        ctrl.setValue('');
        expect(ctrl.errors?.['required']).toBeTruthy();
    });

    it('should email field validate format', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        const ctrl = fixture.componentInstance.registerForm.get('email')!;
        ctrl.setValue('not-an-email');
        expect(ctrl.errors?.['email']).toBeTruthy();
        ctrl.setValue('valid@example.com');
        expect(ctrl.errors).toBeNull();
    });

    it('should password field be required', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        const ctrl = fixture.componentInstance.registerForm.get('password')!;
        ctrl.setValue('');
        expect(ctrl.errors?.['required']).toBeTruthy();
    });

    it('should confirmPassword field be required', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        const ctrl =
            fixture.componentInstance.registerForm.get('confirmPassword')!;
        ctrl.setValue('');
        expect(ctrl.errors?.['required']).toBeTruthy();
    });

    it('should errorMessage be reset when the form is changed', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        fixture.detectChanges();
        const comp = fixture.componentInstance;
        comp.errorMessage = 'Ошибка';
        comp.registerForm.get('username')!.setValue('new');
        expect(comp.errorMessage).toBe('');
    });

    it('should onSubmit not call the service when the form is invalid', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        fixture.detectChanges();
        fixture.componentInstance.onSubmit();
        expect(mockAuthService.register).not.toHaveBeenCalled();
    });

    it('should onSubmit calls markAllAsTouched with invalid form', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        fixture.detectChanges();
        const spy = vi.spyOn(
            fixture.componentInstance.registerForm,
            'markAllAsTouched',
        );
        fixture.componentInstance.onSubmit();
        expect(spy).toHaveBeenCalled();
    });

    it('should .error-message block hiddem if errorMessage is empty', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        fixture.detectChanges();
        expect(
            fixture.nativeElement.querySelector('.error-message'),
        ).toBeNull();
    });

    it('should .error-message block show with error text', () => {
        const fixture = TestBed.createComponent(RegisterComponent);
        fixture.componentInstance.errorMessage = 'E-mail уже используется';
        fixture.detectChanges();
        const el = fixture.nativeElement.querySelector('.error-message');
        expect(el).toBeTruthy();
        expect(el.textContent.trim()).toBe('E-mail уже используется');
    });
});
