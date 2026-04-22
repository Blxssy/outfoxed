import { Component, output, signal } from '@angular/core';
import { InputComponent } from '@fox/ui-kit/input';
import { ButtonComponent } from '@fox/ui-kit/button';
import { CardComponent } from '@fox/ui-kit/card';

@Component({
    selector: 'app-create-game-modal',
    imports: [InputComponent, ButtonComponent, CardComponent],
    templateUrl: './create-game-modal.component.html',
    styleUrl: './create-game-modal.component.scss',
})
export class CreateGameModalComponent {
    readonly closeModal = output<void>();
    readonly newRoomMaxPlayers = signal(4);
    readonly newRoomPrivate = signal(false);

    newRoomName = '';

    closeCreateModal(): void {
        this.closeModal.emit();
    }

    confirmCreateRoom(): void {
        const name = this.newRoomName.trim() || 'Моя игра';
        console.log('создаём комнату');
        this.closeCreateModal();
    }

    togglePrivate(): void {
        this.newRoomPrivate.update((v) => !v);
    }
}
