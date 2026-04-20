import { Component, signal } from '@angular/core';
import { ButtonComponent } from '@fox/ui-kit/button';
import { CardComponent } from '@fox/ui-kit/card/card.component';
import { Room } from './lobby.model';

@Component({
    selector: 'app-lobby',
    imports: [ButtonComponent, CardComponent],
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

    openJoinModal(): void {}
    closeJoinModal(): void {}
    openCreateModal(): void {}
    closeCreateModal(): void {}

    createRoom(): void {
        this.openCreateModal();
    }

    joinRoom(room: Room): void {}

    joinByCode(): void {}
}
