import { Component, effect, inject, OnInit, signal } from '@angular/core';
import { ButtonComponent } from '@fox/ui-kit/button';
import { CardComponent } from '@fox/ui-kit/card/card.component';
import { CreateGameModalComponent } from '../create-game-modal/create-game-modal.component';
import { EnterCodeModalComponent } from '../enter-code-modal/enter-code-modal.component';
import { LobbyStore } from '../data/lobby.store';
import { RoomListItem } from '../data/lobby.model';

@Component({
    selector: 'app-rooms-list',
    imports: [
        ButtonComponent,
        CardComponent,
        CreateGameModalComponent,
        EnterCodeModalComponent,
    ],
    templateUrl: './rooms-list.component.html',
    styleUrl: './rooms-list.component.scss',
})
export class RoomsListComponent implements OnInit {
    protected readonly store = inject(LobbyStore);

    readonly joinModalOpen = signal(false);
    readonly createModalOpen = signal(false);

    joinCode = '';
    newRoomName = '';

    protected readonly rooms = this.store.rooms;
    protected readonly isLoading = this.store.isLoadingList;

    constructor() {
        effect(() => {
            console.log('rooms updated:', this.rooms());
        });
    }

    ngOnInit(): void {
        this.store.loadRooms();
        console.log('rooms in lobby', this.rooms());
    }

    openCreateModal(): void {
        this.createModalOpen.set(true);
        console.log('join opened');
    }

    closeCreateModal(): void {
        this.createModalOpen.set(false);
    }

    openJoinModal(): void {
        this.joinModalOpen.set(true);
        console.log('create opened');
    }

    closeJoinModal(): void {
        this.joinModalOpen.set(false);
    }

    createRoom(): void {
        this.openCreateModal();
    }

    joinRoom(room: RoomListItem): void {
        if (room.status !== 'waiting') {
            return;
        }
        this.store.joinRoom(room.id);
    }
}
