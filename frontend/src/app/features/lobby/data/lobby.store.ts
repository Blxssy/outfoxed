import { inject, Injectable, signal } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Router } from '@angular/router';
import { LobbyApiService } from './lobby-api.service';
import { LobbyWsService } from './lobby-ws.service';
import { RoomListItem } from './lobby.model';

@Injectable({ providedIn: 'root' })
export class LobbyStore {
    private readonly router = inject(Router);
    private readonly api = inject(LobbyApiService);
    private readonly lobbyWs = inject(LobbyWsService);

    readonly rooms = signal<RoomListItem[]>([]);
    readonly isLoadingList = signal(false);

    readonly currentRoom = this.lobbyWs.lobbySnapshot;
    readonly canStart = this.lobbyWs.canStart;
    readonly wsError = this.lobbyWs.wsError;
    readonly wsStatus = this.lobbyWs.connectionStatus;
    readonly isReconnecting = this.lobbyWs.isReconnecting;

    readonly isCreating = signal(false);
    readonly isJoining = signal(false);
    readonly isLeaving = signal(false);
    readonly isStarting = signal(false);

    readonly privateJoinCode = signal<string | null>(null);

    readonly error = signal<{ kind: string; message: string } | null>(null);

    loadRooms(): void {
        this.isLoadingList.set(true);
        this.clearError();

        this.api.getPublicGames().subscribe({
            next: (res) => {
                this.rooms.set(res.games);
                this.isLoadingList.set(false);
            },
            error: (err: HttpErrorResponse) => {
                this.isLoadingList.set(false);
                this.setError(
                    'load_failed',
                    this.httpMessage(err, 'Не удалось загрузить список комнат'),
                );
            },
        });
    }

    createGame(title: string, isPrivate: boolean): void {
        this.isCreating.set(true);
        this.clearError();

        this.api
            .createGame({
                title: title.trim() || 'Моя игра',
                visibility: isPrivate ? 'private' : 'public',
            })
            .subscribe({
                next: (res) => {
                    this.isCreating.set(false);

                    if (res.joinCode) {
                        this.privateJoinCode.set(res.joinCode);
                    }

                    this.enterRoom(res.game.id);
                },
                error: (err: HttpErrorResponse) => {
                    this.isCreating.set(false);
                    this.setError(
                        'create_failed',
                        this.httpMessage(err, 'Не удалось создать комнату'),
                    );
                },
            });
    }

    joinRoom(id: string): void {
        this.isJoining.set(true);
        this.clearError();

        this.api.joinGame(id).subscribe({
            next: () => {
                this.isJoining.set(false);
                this.enterRoom(id);
            },
            error: (err: HttpErrorResponse) => {
                this.isJoining.set(false);
                this.setError(
                    'join_failed',
                    this.httpMessage(err, 'Не удалось войти в комнату'),
                );
            },
        });
    }

    joinByCode(code: string): void {
        this.isJoining.set(true);
        this.clearError();

        this.api.joinByCode({ code: code.trim() }).subscribe({
            next: (res) => {
                this.isJoining.set(false);
                this.enterRoom(res.game.id);
            },
            error: (err: HttpErrorResponse) => {
                this.isJoining.set(false);
                this.setError('join_code_failed', this.codeErrorMessage(err));
            },
        });
    }

    leaveRoom(): void {
        const id = this.currentRoom()?.id;
        if (!id) return;

        this.isLeaving.set(true);
        this.clearError();

        this.api.leaveGame(id).subscribe({
            next: () => {
                this.isLeaving.set(false);

                this.lobbyWs.disconnect();
                this.privateJoinCode.set(null);

                this.router.navigate(['/lobby']);
            },
            error: (err: HttpErrorResponse) => {
                this.isLeaving.set(false);
                this.setError(
                    'leave_failed',
                    this.httpMessage(err, 'Не удалось выйти из комнаты'),
                );
            },
        });
    }

    startGame(): void {
        const id = this.currentRoom()?.id;
        if (!id) return;

        this.isStarting.set(true);
        this.clearError();

        this.api.startGame(id).subscribe({
            next: (res) => {
                this.isStarting.set(false);
                this.lobbyWs.disconnect();
                this.router.navigateByUrl(res.redirect.route);
            },
            error: (err: HttpErrorResponse) => {
                this.isStarting.set(false);
                this.setError(
                    'start_failed',
                    this.httpMessage(err, 'Не удалось запустить игру'),
                );
            },
        });
    }

    initRoom(gameId: string): void {
        const status = this.lobbyWs.connectionStatus();
        if (status !== 'connected' && status !== 'connecting') {
            this.lobbyWs.connect(gameId);
        }
    }

    clearError(): void {
        this.error.set(null);
    }

    private enterRoom(gameId: string): void {
        this.lobbyWs.connect(gameId);
        this.router.navigate(['/lobby', gameId]);
    }

    private setError(kind: string, message: string): void {
        this.error.set({ kind, message });
    }

    private httpMessage(err: HttpErrorResponse, fallback: string): string {
        if (err.status === 401) return 'Необходима авторизация';
        if (err.status === 403) return 'Доступ запрещён';
        if (err.status === 404) return 'Игра не найдена';
        if (err.status === 500) return 'Ошибка сервера, попробуйте ещё раз';
        return fallback;
    }

    private codeErrorMessage(err: HttpErrorResponse): string {
        if (err.status === 400) return 'Неверный формат кода';
        if (err.status === 404) return 'Комната с таким кодом не найдена';
        if (err.status === 409) return 'Игра уже началась или заполнена';
        return this.httpMessage(err, 'Не удалось войти по коду');
    }
}
