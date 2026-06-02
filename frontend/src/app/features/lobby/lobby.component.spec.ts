import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { LobbyComponent } from './lobby.component';
import { ActiveGameService } from './data/active-game.service';

const mockActiveGameService = {
    startPolling: vi.fn(),
    stopPolling: vi.fn(),
    activeGame: vi.fn(),
};

describe('LobbyComponent', () => {
    beforeEach(async () => {
        vi.clearAllMocks();

        await TestBed.configureTestingModule({
            imports: [LobbyComponent],
            providers: [
                provideRouter([]),
                { provide: ActiveGameService, useValue: mockActiveGameService },
            ],
            schemas: [NO_ERRORS_SCHEMA],
        }).compileComponents();
    });

    it('should create lobby component', () => {
        const fixture = TestBed.createComponent(LobbyComponent);
        expect(fixture.componentInstance).toBeTruthy();
    });

    it('should call startPolling on init', () => {
        const fixture = TestBed.createComponent(LobbyComponent);
        fixture.detectChanges();

        expect(mockActiveGameService.startPolling).toHaveBeenCalled();
    });

    it('should call stopPolling on destroy', () => {
        const fixture = TestBed.createComponent(LobbyComponent);
        fixture.detectChanges();
        fixture.destroy();

        expect(mockActiveGameService.stopPolling).toHaveBeenCalled();
    });

    it('should show return button when active game exist', () => {
        mockActiveGameService.activeGame.mockReturnValue({
            route: '/game/123',
        });

        const fixture = TestBed.createComponent(LobbyComponent);
        fixture.detectChanges();

        const btn = fixture.nativeElement.querySelector('.return-to-game-btn');

        expect(btn).toBeTruthy();
    });

    it('should not show return button when no active game', () => {
        mockActiveGameService.activeGame.mockReturnValue(null);

        const fixture = TestBed.createComponent(LobbyComponent);
        fixture.detectChanges();

        const btn = fixture.nativeElement.querySelector('.return-to-game-btn');

        expect(btn).toBeNull();
    });

    it('should navigate to active game route on returnToGame', () => {
        mockActiveGameService.activeGame.mockReturnValue({
            route: '/game/999',
        });

        const fixture = TestBed.createComponent(LobbyComponent);
        fixture.detectChanges();

        const router = TestBed.inject(Router);
        const navigateSpy = vi.spyOn(router, 'navigateByUrl');

        fixture.componentInstance.returnToGame();

        expect(navigateSpy).toHaveBeenCalledWith('/game/999');
    });

    it('should not navigate if no active game', () => {
        mockActiveGameService.activeGame.mockReturnValue(null);

        const fixture = TestBed.createComponent(LobbyComponent);
        fixture.detectChanges();

        const router = TestBed.inject(Router);
        const navigateSpy = vi.spyOn(router, 'navigateByUrl');

        fixture.componentInstance.returnToGame();

        expect(navigateSpy).not.toHaveBeenCalled();
    });

    it('should render logo with correct src', () => {
        const fixture = TestBed.createComponent(LobbyComponent);
        fixture.detectChanges();

        const img: HTMLImageElement =
            fixture.nativeElement.querySelector('.lobby-logo__img');

        expect(img).toBeTruthy();
        expect(img.getAttribute('src')).toBe('outfoxed_logo.png');
    });
});
