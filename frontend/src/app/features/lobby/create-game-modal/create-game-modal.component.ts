import { Component, inject, output, signal } from '@angular/core';
import { InputComponent } from '@fox/ui-kit/input';
import { ButtonComponent } from '@fox/ui-kit/button';
import { CardComponent } from '@fox/ui-kit/card';
import { LobbyStore } from '../data/lobby.store';
import { FormsModule } from '@angular/forms';

@Component({
    selector: 'app-create-game-modal',
    imports: [InputComponent, ButtonComponent, CardComponent, FormsModule],
    templateUrl: './create-game-modal.component.html',
    styleUrl: './create-game-modal.component.scss',
})
export class CreateGameModalComponent {
    private readonly store = inject(LobbyStore);

    readonly closeModal = output<void>();
    readonly newRoomMaxPlayers = signal(4);
    readonly newRoomPrivate = signal(false);

    newRoomName = '';

    get isLoading(): boolean {
        return this.store.isCreating();
    }

    closeCreateModal(): void {
        this.closeModal.emit();
    }

    confirmCreateRoom(): void {
        this.store.createGame(this.newRoomName, this.newRoomPrivate());
        this.closeModal.emit();
    }

    togglePrivate(): void {
        this.newRoomPrivate.update((v) => !v);
    }
}
