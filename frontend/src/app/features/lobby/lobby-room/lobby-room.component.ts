import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ButtonComponent } from '@fox/ui-kit/button';
import { LobbyStore } from '../data/lobby.store';
import { CardComponent } from '@fox/ui-kit/card';
import { LobbyWsService } from '../data/lobby-ws.service';

@Component({
    selector: 'app-lobby-room',
    imports: [ButtonComponent, CardComponent],
    templateUrl: './lobby-room.component.html',
    styleUrl: './lobby-room.component.scss',
})
export class LobbyRoomComponent implements OnInit, OnDestroy {
    protected readonly store = inject(LobbyStore);
    protected readonly lobbyWs = inject(LobbyWsService);
    private readonly route = inject(ActivatedRoute);

    protected readonly room = this.store.currentRoom;
    protected readonly canStart = this.store.canStart;
    protected readonly isLeaving = this.store.isLeaving;
    protected readonly isStarting = this.store.isStarting;
    protected readonly joinCode = this.store.privateJoinCode;
    protected readonly httpError = this.store.error;
    protected readonly wsError = this.store.wsError;
    protected readonly isReconnecting = this.store.isReconnecting;
    protected readonly wsStatus = this.store.wsStatus;

    ngOnInit(): void {
        const id = this.route.snapshot.paramMap.get('id');
        if (id) {
            this.store.initRoom(id);
        }
    }

    getEmptySlots(filled: number, max: number): number[] {
        const count = Math.max(0, max - filled);
        return Array.from({ length: count }, (_, i) => i + 1);
    }

    ngOnDestroy(): void {
        if (this.lobbyWs.connectionStatus() !== 'disconnected') {
            this.lobbyWs.disconnect();
        }
    }

    leaveRoom(): void {
        this.store.leaveRoom();
    }

    startGame(): void {
        this.store.startGame();
    }
}
