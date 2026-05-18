import { Injectable, OnDestroy, signal, computed, inject } from '@angular/core';
import { Subscription } from 'rxjs';

import {
    PublicGameState,
    WsServerUpdate,
    WsServerError,
    WsErrorPayload,
    WsEvent,
    WsClientMessage,
    WsCommandName,
    WsCommandPayloadMap,
    PendingRequest,
    GoalType,
    AvailableAction,
} from './game.types';
import {
    WebSocketService,
    WsRawMessage,
} from 'src/app/services/websocket/websocket.service';
import { TokenService } from 'src/app/services/auth/token.service';

const MAX_EVENT_LOG = 100;
const WS_BASE = 'ws://localhost:8080';

@Injectable({ providedIn: 'root' })
export class GameService implements OnDestroy {
    private readonly ws = inject(WebSocketService);
    private readonly tokenService = inject(TokenService);
    private sub: Subscription | null = null;

    readonly gameState = signal<PublicGameState | null>(null);
    readonly lastError = signal<WsErrorPayload | null>(null);
    readonly eventLog = signal<WsEvent[]>([]);
    readonly pendingRequests = signal<PendingRequest[]>([]);

    readonly connectionStatus = this.ws.status;
    readonly isConnected = this.ws.isConnected;
    readonly isReconnecting = this.ws.isReconnecting;

    readonly myPlayer = computed(() => this.gameState()?.me ?? null);
    readonly activeSeat = computed(() => this.gameState()?.activeSeat ?? -1);
    readonly isMyTurn = computed(
        () => this.gameState()?.me.seat === this.gameState()?.activeSeat,
    );
    readonly availableActions = computed<AvailableAction[]>(
        () => this.gameState()?.availableActions ?? [],
    );
    readonly phase = computed(() => this.gameState()?.phase ?? null);
    readonly foxTrack = computed(() => this.gameState()?.fox ?? null);
    readonly suspects = computed(() => this.gameState()?.suspects ?? []);
    readonly clues = computed(() => this.gameState()?.clues ?? []);
    readonly players = computed(() => this.gameState()?.players ?? []);
    readonly roll = computed(() => this.gameState()?.roll ?? null);
    readonly move = computed(() => this.gameState()?.move ?? null);
    readonly reachableCells = computed(
        () => this.gameState()?.move?.reachableCells ?? [],
    );
    readonly boardCells = computed(() => this.gameState()?.board?.cells ?? []);

    startSession(gameId: string): void {
        const token = this.tokenService.getAccessToken() ?? '';
        this.reset();
        this.subscribeToMessages();

        const wsBase = WebSocketService.resolveWsBase();
        this.ws.connect(
            `${wsBase}/ws/games/${gameId}?token=${encodeURIComponent(token)}`,
        );
    }

    endSession(): void {
        this.ws.disconnect();
        this.sub?.unsubscribe();
        this.sub = null;
    }

    ngOnDestroy(): void {
        this.endSession();
    }

    chooseGoal(goal: GoalType): boolean {
        return this.sendCommand('choose_goal', { goal });
    }
    rollAuto(): boolean {
        return this.sendCommand('roll_auto', {});
    }
    rerollDice(keepIndices: number[]): boolean {
        return this.sendCommand('reroll_dice', { keepIndices });
    }
    finishRoll(): boolean {
        return this.sendCommand('finish_roll', {});
    }
    movePawn(targetIndex: number): boolean {
        return this.sendCommand('move_pawn', { targetIndex });
    }
    takeClue(): boolean {
        return this.sendCommand('take_clue', {});
    }
    revealSuspects(ids: string[]): boolean {
        return this.sendCommand('reveal_suspects', { suspectIds: ids });
    }
    endTurn(): boolean {
        return this.sendCommand('end_turn', {});
    }
    accuse(suspectId: string): boolean {
        return this.sendCommand('accuse', { suspectId });
    }

    canDo(action: AvailableAction): boolean {
        const actions = this.gameState()?.availableActions;
        return (
            this.isMyTurn() &&
            Array.isArray(actions) &&
            actions.includes(action)
        );
    }

    hasPending(): boolean {
        return this.pendingRequests().length > 0;
    }

    private subscribeToMessages(): void {
        this.sub?.unsubscribe();
        this.sub = this.ws.message$.subscribe((raw: WsRawMessage) => {
            if (raw.type === 'game_update' || raw.type === 'update') {
                this.handleUpdate(raw as unknown as WsServerUpdate);
            } else if (raw.type === 'error') {
                this.handleServerError(raw as unknown as WsServerError);
            }
        });
    }

    private handleUpdate(msg: WsServerUpdate): void {
        const { state, events } = msg.payload;

        this.gameState.set(state);
        this.lastError.set(null);

        if (events?.length) {
            this.eventLog.update((log) =>
                [...log, ...events].slice(-MAX_EVENT_LOG),
            );
        }

        if (msg.id) this.resolvePending(msg.id);
    }

    private handleServerError(msg: WsServerError): void {
        console.warn(
            '[Game] Ошибка от сервера:',
            msg.payload.code,
            msg.payload.message,
        );
        this.lastError.set(msg.payload);
        if (msg.id) this.resolvePending(msg.id);
    }

    private sendCommand<T extends WsCommandName>(
        command: T,
        payload: WsCommandPayloadMap[T],
    ): boolean {
        const requestId = this.generateRequestId();
        const message: WsClientMessage<T> = {
            id: requestId,
            type: 'command',
            command,
            payload,
        };

        const sent = this.ws.send(message);
        if (sent) {
            this.pendingRequests.update((prev) => [
                ...prev,
                { requestId, command, sentAt: Date.now() },
            ]);
        }
        return sent;
    }

    private resolvePending(requestId: string): void {
        this.pendingRequests.update((prev) =>
            prev.filter((r) => r.requestId !== requestId),
        );
    }

    private reset(): void {
        this.gameState.set(null);
        this.lastError.set(null);
        this.eventLog.set([]);
        this.pendingRequests.set([]);
    }

    private generateRequestId(): string {
        return `req-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`;
    }
}
