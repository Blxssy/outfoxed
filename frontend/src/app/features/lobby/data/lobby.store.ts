import { computed, inject, Injectable, signal } from '@angular/core';
import { LobbyApiService } from './lobby-api.service';
import { LobbySnapshot, RoomListItem } from './lobby.model';
import { HttpErrorResponse } from '@angular/common/http';
import { catchError, EMPTY, interval, Subscription, switchMap } from 'rxjs';
import { Router } from '@angular/router';

const POLL_INTERVAL_MS = 3000;

@Injectable({ providedIn: 'root' })
export class LobbyStore {
    private readonly router = inject(Router);
    private readonly api = inject(LobbyApiService);

    readonly rooms = signal<RoomListItem[]>([]);
    readonly currentRoom = signal<LobbySnapshot | null>(null);

    readonly isLoadingList = signal(false);
    readonly isCreating = signal(false);
    readonly isJoining = signal(false);
    readonly isLeaving = signal(false);
    readonly isStarting = signal(false);

    readonly privateJoinCode = signal<string | null>(null);

    readonly error = signal<{} | null>(null);

    readonly canStart = computed(() => this.currentRoom()?.can_start ?? false);

    private pollSub: Subscription | null = null;

    loadRooms(): void {
        this.isLoadingList.set(true);
        this.clearError();

        this.api.getPublicGames().subscribe({
            next: (res) => {
                console.log(res.games);
                this.rooms.set(res.games);
                console.log('rooms', this.rooms());
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
                    console.log('game created');

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
            next: (res) => {
                this.isLeaving.set(false);
                this.stopRoomPolling();
                this.currentRoom.set(null);
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
                this.stopRoomPolling();
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

    startRoomPolling(gameId: string): void {
        this.stopRoomPolling();
        this.pollSub = interval(POLL_INTERVAL_MS)
            .pipe(
                switchMap(() =>
                    this.api
                        .getLobbySnapshot(gameId)
                        .pipe(catchError(() => EMPTY)),
                ),
            )
            .subscribe((res) => {
                this.currentRoom.set(res.game);

                if (res.game.status === 'active') {
                    this.stopRoomPolling();
                    this.router.navigate(['/game', gameId]);
                }
            });
    }

    stopRoomPolling(): void {
        this.pollSub?.unsubscribe();
        this.pollSub = null;
    }

    refreshCurrentRoom(): void {
        const id = this.currentRoom()?.id;
        if (!id) return;
        this.fetchLobbySnapshot(id);
    }

    clearError(): void {
        this.error.set(null);
    }

    private setError(kind: any, message: string): void {
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

    private enterRoom(id: string): void {
        this.fetchLobbySnapshot(id, () => {
            this.router.navigate(['/lobby', id]);
        });
    }

    private fetchLobbySnapshot(id: string, onSuccess?: () => void): void {
        this.api.getLobbySnapshot(id).subscribe({
            next: (res) => {
                this.currentRoom.set(res.game);
                onSuccess?.();
            },
            error: (err: HttpErrorResponse) => {
                this.setError(
                    'load_failed',
                    this.httpMessage(err, 'Не удалось загрузить комнату'),
                );
            },
        });
    }

    ngOnDestroy(): void {
        this.stopRoomPolling();
    }
}
