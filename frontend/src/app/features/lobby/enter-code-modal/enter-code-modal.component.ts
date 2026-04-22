import { Component, output, signal } from '@angular/core';
import { InputComponent } from '@fox/ui-kit/input';
import { CardComponent } from '@fox/ui-kit/card';
import { ButtonComponent } from '@fox/ui-kit/button';

@Component({
    selector: 'app-enter-code-modal',
    imports: [InputComponent, CardComponent, ButtonComponent],
    templateUrl: './enter-code-modal.component.html',
    styleUrl: './enter-code-modal.component.scss',
})
export class EnterCodeModalComponent {
    readonly closeModal = output<void>();
    readonly codeError = signal('');
    joinCode = '';

    closeJoinModal(): void {
        this.closeModal.emit();
    }

    joinByCode(): void {
        const code = this.joinCode.trim().toUpperCase();

        console.log('вход');
        this.closeJoinModal();
    }
}
