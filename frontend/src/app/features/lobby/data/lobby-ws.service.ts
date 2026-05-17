import { Injectable, OnDestroy, inject, signal, computed } from '@angular/core';
import { Subscription } from 'rxjs';
import { Router } from '@angular/router';

import { LobbySnapshot } from './lobby.model';
import {
    WebSocketService,
    WsRawMessage,
} from 'src/app/services/websocket/websocket.service';
import { TokenService } from 'src/app/services/auth/token.service';
const WS_BASE = 'ws://localhost:8080';

@Injectable({ providedIn: 'root' })
export class LobbyWsService implements OnDestroy {
    private readonly ws = inject(WebSocketService);
    private readonly router = inject(Router);
    private readonly tokenService = inject(TokenService);

    private msgSub: Subscription | null = null;
    private currentGameId: string | null = null;

    readonly lobbySnapshot = signal<LobbySnapshot | null>(null);

    readonly wsError = signal<{ code: string; message: string } | null>(null);

    readonly connectionStatus = this.ws.status;
    readonly isConnected = this.ws.isConnected;
    readonly isReconnecting = this.ws.isReconnecting;

    readonly canStart = computed(
        () => this.lobbySnapshot()?.can_start ?? false,
    );

    connect(gameId: string): void {
        const status = this.ws.status();
        if (
            this.currentGameId === gameId &&
            (status === 'connected' || status === 'connecting')
        ) {
            return;
        }

        const token = this.tokenService.getAccessToken() ?? '';
        this.currentGameId = gameId;
        this.reset();
        this.subscribeToMessages(gameId);
        this.ws.connect(`${WS_BASE}/ws/games/${gameId}?token=${token}`);
    }

    disconnect(): void {
        this.currentGameId = null;
        this.ws.disconnect();
        this.msgSub?.unsubscribe();
        this.msgSub = null;
    }

    ngOnDestroy(): void {
        this.disconnect();
    }

    private subscribeToMessages(gameId: string): void {
        this.msgSub?.unsubscribe();

        this.msgSub = this.ws.message$.subscribe((raw: WsRawMessage) => {
            if (raw.type === 'game_update' || raw.type === 'update') {
                const state = (raw.payload as any)?.state;
                if (!state) return;
                this.handleUpdate(state, gameId);
            } else if (raw.type === 'error') {
                const err = raw.payload as { code: string; message: string };
                this.wsError.set(err);
            }
        });
    }

    private handleUpdate(state: any, gameId: string): void {
        const players: {
            userId: string;
            seat: number;
            name: string;
            connected: boolean;
        }[] = state.players ?? [];
        const me: { userId: string; seat: number; name: string } | undefined =
            state.me;

        const host = players.find((p) => p.seat === 0);
        const iAmHost = !!me && me.userId === host?.userId;

        const playerCount = players.length;
        const canStart = iAmHost && playerCount >= 2 && playerCount <= 4;

        const snapshot: LobbySnapshot = {
            id: state.id,
            title: state.title ?? '',
            status: state.status,
            visibility: state.visibility ?? 'public',
            joinCode: state.joinCode,
            host_username: host?.name ?? '',
            players: players.map((p) => ({
                user_id: p.userId,
                seat: p.seat,
                display_name: p.name,
                is_me: p.userId === me?.userId,
            })),
            can_start: canStart,
            min_players: 2,
            max_players: 4,
        };

        this.lobbySnapshot.set(snapshot);
        this.wsError.set(null);

        if (state.status === 'active') {
            this.disconnect();
            this.router.navigate(['/game', gameId]);
        }
    }

    private reset(): void {
        this.lobbySnapshot.set(null);
        this.wsError.set(null);
    }
}
