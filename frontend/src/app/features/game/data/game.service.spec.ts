import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { GameService } from './game.service';
import { WebSocketService } from 'src/app/services/websocket/websocket.service';
import { TokenService } from 'src/app/services/auth/token.service';
import { Subject } from 'rxjs';

const message$ = new Subject<any>();

const mockWs = {
    connect: vi.fn(),
    disconnect: vi.fn(),
    send: vi.fn(() => true),
    message$,
    status: {},
    isConnected: vi.fn(() => true),
    isReconnecting: vi.fn(() => false),
    static: {},
};

const mockTokenService = {
    getAccessToken: vi.fn(() => 'token'),
};

describe('GameService', () => {
    let service: GameService;

    beforeEach(() => {
        message$.next(null as any);
        vi.clearAllMocks();

        TestBed.configureTestingModule({
            providers: [
                GameService,
                { provide: WebSocketService, useValue: mockWs },
                { provide: TokenService, useValue: mockTokenService },
            ],
        });

        service = TestBed.inject(GameService);
    });

    it('should start session and connect ws', () => {
        service.startSession('game-1');

        expect(mockWs.connect).toHaveBeenCalled();
    });

    it('should end session and disconnect ws', () => {
        service.startSession('game-1');
        service.endSession();

        expect(mockWs.disconnect).toHaveBeenCalled();
    });

    it('should reset state on startSession', () => {
        service.gameState.set({
            me: {} as any,
        } as any);

        service.startSession('game-1');

        expect(service.gameState()).toBeNull();
    });

    it('should send command and add pending request', () => {
        service.startSession('game-1');

        const result = service.movePawn(5);

        expect(result).toBe(true);
        expect(mockWs.send).toHaveBeenCalled();

        expect(service.pendingRequests().length).toBe(1);
    });

    it('should not add pending if ws send fails', () => {
        mockWs.send.mockReturnValueOnce(false);

        service.startSession('game-1');

        const result = service.movePawn(5);

        expect(result).toBe(false);
        expect(service.pendingRequests().length).toBe(0);
    });

    it('should handle game update message', () => {
        service.startSession('game-1');

        message$.next({
            type: 'game_update',
            payload: {
                state: { activeSeat: 1 } as any,
                events: [{ type: 'test' }],
            },
        });

        expect(service.gameState()?.activeSeat).toBe(1);
        expect(service.eventLog().length).toBe(1);
    });

    it('should handle error message', () => {
        service.startSession('game-1');

        message$.next({
            type: 'error',
            payload: {
                code: 500,
                message: 'fail',
            },
        });

        expect(service.lastError()?.code).toBe(500);
    });

    it('should resolve pending on update with id', () => {
        service.startSession('game-1');

        const reqId = service.movePawn(2);

        message$.next({
            id: (mockWs.send.mock.calls[0][0] as any).id,
            type: 'game_update',
            payload: {
                state: {} as any,
                events: [],
            },
        });

        expect(service.pendingRequests().length).toBe(0);
    });

    it('should compute canDo correctly', () => {
        service.gameState.set({
            activeSeat: 0,
            me: { seat: 0 } as any,
            availableActions: ['move_pawn'] as any,
        } as any);

        expect(service.canDo('move_pawn' as any)).toBe(true);
    });

    it('should detect hasPending', () => {
        service.pendingRequests.set([
            { requestId: '1', command: 'move_pawn', sentAt: 0 } as any,
        ]);

        expect(service.hasPending()).toBe(true);
    });
});
