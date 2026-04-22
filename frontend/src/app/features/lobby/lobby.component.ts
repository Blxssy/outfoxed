import { Component, OnInit, signal } from '@angular/core';
import { ButtonComponent } from '@fox/ui-kit/button';
import { CardComponent } from '@fox/ui-kit/card/card.component';
import { Room } from './lobby.model';
import { CreateGameModalComponent } from './create-game-modal/create-game-modal.component';
import { EnterCodeModalComponent } from './enter-code-modal/enter-code-modal.component';

@Component({
    selector: 'app-lobby',
    imports: [
        ButtonComponent,
        CardComponent,
        CreateGameModalComponent,
        EnterCodeModalComponent,
    ],
    templateUrl: './lobby.component.html',
    styleUrl: './lobby.component.scss',
})
export class LobbyComponent {
    protected readonly rooms = signal<Room[]>([]);
    readonly isLoading = signal(true);
    readonly joinModalOpen = signal(false);
    readonly createModalOpen = signal(false);

    joinCode = '';
    newRoomName = '';

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

    joinRoom(room: Room): void {
        if (room.status !== 'waiting') {
            return;
        }
    }

    joinByCode(): void {}
}
