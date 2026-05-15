import { Injectable, OnDestroy, signal, computed } from '@angular/core';
import { Subject } from 'rxjs';

export type WsConnectionStatus =
    | 'disconnected'
    | 'connecting'
    | 'connected'
    | 'reconnecting'
    | 'error';

export type WsRawMessage = {
    type: string;
    id?: string | null;
    payload?: unknown;
};

const RECONNECT_DELAYS = [1_000, 2_000, 5_000, 10_000];
const MAX_RECONNECT_ATTEMPTS = 8;

@Injectable({ providedIn: 'root' })
export class WebSocketService implements OnDestroy {
    readonly status = signal<WsConnectionStatus>('disconnected');
    readonly isConnected = computed(() => this.status() === 'connected');
    readonly isReconnecting = computed(() => this.status() === 'reconnecting');

    readonly message$ = new Subject<WsRawMessage>();

    private socket: WebSocket | null = null;
    private currentUrl: string | null = null;
    private reconnectAttempt = 0;
    private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
    private intentionalClose = false;

    connect(url: string): void {
        this.currentUrl = url;
        this.intentionalClose = false;
        this.reconnectAttempt = 0;
        this.openSocket();
    }

    disconnect(): void {
        this.intentionalClose = true;
        this.clearReconnectTimer();
        this.closeSocket();
        this.status.set('disconnected');
    }

    send(message: object): boolean {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
            console.warn(
                '[WS] Попытка отправки при закрытом соединении',
                message,
            );
            return false;
        }
        this.socket.send(JSON.stringify(message));
        return true;
    }

    ngOnDestroy(): void {
        this.disconnect();
        this.message$.complete();
    }

    private openSocket(): void {
        if (
            this.socket &&
            (this.socket.readyState === WebSocket.OPEN ||
                this.socket.readyState === WebSocket.CONNECTING)
        ) {
            console.log('[WS] Сокет уже открыт/подключается');
            return;
        }

        this.closeSocket();
        if (!this.currentUrl) return;

        this.status.set(
            this.reconnectAttempt === 0 ? 'connecting' : 'reconnecting',
        );
        console.log(
            `[WS] Подключение к ${this.currentUrl} (попытка ${this.reconnectAttempt + 1})`,
        );

        const ws = new WebSocket(this.currentUrl);
        this.socket = ws;

        ws.onopen = () => this.handleOpen();
        ws.onmessage = (event) => this.handleMessage(event);
        ws.onclose = (event) => this.handleClose(event);
        ws.onerror = (event) => this.handleError(event);
    }

    private handleOpen(): void {
        console.log('[WS] Соединение установлено');
        this.status.set('connected');
        this.reconnectAttempt = 0;
        this.clearReconnectTimer();
    }

    private handleMessage(event: MessageEvent): void {
        let parsed: WsRawMessage;
        try {
            parsed = JSON.parse(event.data as string) as WsRawMessage;
        } catch {
            console.error('[WS] Не удалось распарсить сообщение', event.data);
            return;
        }
        this.message$.next(parsed);
    }

    private handleClose(event: CloseEvent): void {
        console.warn('[WS] Соединение закрыто', event.code, event.reason);
        this.socket = null;
        if (this.intentionalClose) return;
        this.scheduleReconnect();
    }

    private handleError(event: Event): void {
        console.error('[WS] Ошибка транспорта', event);
    }

    private scheduleReconnect(): void {
        if (this.reconnectAttempt >= MAX_RECONNECT_ATTEMPTS) {
            console.error(
                '[WS] Превышено максимальное число попыток reconnect',
            );
            this.status.set('error');
            return;
        }

        const delay =
            RECONNECT_DELAYS[
                Math.min(this.reconnectAttempt, RECONNECT_DELAYS.length - 1)
            ];

        console.log(`[WS] Переподключение через ${delay}мс...`);
        this.status.set('reconnecting');

        this.reconnectTimer = setTimeout(() => {
            this.reconnectAttempt++;
            this.openSocket();
        }, delay);
    }

    private closeSocket(): void {
        if (this.socket) {
            this.socket.onopen = null;
            this.socket.onmessage = null;
            this.socket.onclose = null;
            this.socket.onerror = null;
            this.socket.close();
            this.socket = null;
        }
    }

    private clearReconnectTimer(): void {
        if (this.reconnectTimer !== null) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }
    }

    static resolveWsBase(): string {
        const proto = location.protocol === 'https:' ? 'wss' : 'ws';
        return `${proto}://${location.host}`;
    }
}
