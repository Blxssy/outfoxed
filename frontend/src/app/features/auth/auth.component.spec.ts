import { describe, it, expect, beforeEach } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { AuthComponent } from './auth.component';

describe('AuthComponent', () => {
    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [AuthComponent],
            providers: [provideRouter([])],
        }).compileComponents();
    });

    it('should be created without errors', () => {
        const fixture = TestBed.createComponent(AuthComponent);
        expect(fixture.componentInstance).toBeTruthy();
    });

    it('should show logo with correct src', () => {
        const fixture = TestBed.createComponent(AuthComponent);
        fixture.detectChanges();
        const img: HTMLImageElement =
            fixture.nativeElement.querySelector('.auth-logo__img');
        expect(img).toBeTruthy();
        expect(img.getAttribute('src')).toBe('outfoxed_logo.png');
    });

    it('should show navigation links', () => {
        const fixture = TestBed.createComponent(AuthComponent);
        fixture.detectChanges();
        const links = fixture.nativeElement.querySelectorAll(
            '.auth-switcher__btn',
        );
        expect(links.length).toBe(2);
    });

    it('should first link lead to /auth/login with text "Вход"', () => {
        const fixture = TestBed.createComponent(AuthComponent);
        fixture.detectChanges();
        const links = fixture.nativeElement.querySelectorAll(
            '.auth-switcher__btn',
        );
        expect(links[0].textContent.trim()).toBe('Вход');
        expect(links[0].getAttribute('href')).toBe('/auth/login');
    });

    it('should second link lead to /auth/register with text "Регистрация"', () => {
        const fixture = TestBed.createComponent(AuthComponent);
        fixture.detectChanges();
        const links = fixture.nativeElement.querySelectorAll(
            '.auth-switcher__btn',
        );
        expect(links[1].textContent.trim()).toBe('Регистрация');
        expect(links[1].getAttribute('href')).toBe('/auth/register');
    });

    it('should contain router-outlet', () => {
        const fixture = TestBed.createComponent(AuthComponent);
        fixture.detectChanges();
        const outlet = fixture.nativeElement.querySelector('router-outlet');
        expect(outlet).toBeTruthy();
    });
});
