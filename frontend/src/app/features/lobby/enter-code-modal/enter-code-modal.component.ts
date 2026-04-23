import { Component, inject, output, signal } from '@angular/core';
import { InputComponent } from '@fox/ui-kit/input';
import { CardComponent } from '@fox/ui-kit/card';
import { ButtonComponent } from '@fox/ui-kit/button';
import { LobbyStore } from '../data/lobby.store';
import { FormsModule } from '@angular/forms';

@Component({
    selector: 'app-enter-code-modal',
    imports: [InputComponent, CardComponent, ButtonComponent, FormsModule],
    templateUrl: './enter-code-modal.component.html',
    styleUrl: './enter-code-modal.component.scss',
})
export class EnterCodeModalComponent {
    private readonly store = inject(LobbyStore);

    readonly closeModal = output<void>();
    joinCode = '';

    get isLoading(): boolean {
        return this.store.isJoining();
    }

    closeJoinModal(): void {
        this.store.clearError();
        this.closeModal.emit();
    }

    joinByCode(): void {
        const code = this.joinCode.trim();
        if (!code) return;
        this.store.joinByCode(code);
    }
}
