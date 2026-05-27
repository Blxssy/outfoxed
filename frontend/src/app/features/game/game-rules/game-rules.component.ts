import { Component, output } from '@angular/core';

@Component({
    selector: 'app-game-rules',
    imports: [],
    templateUrl: './game-rules.component.html',
    styleUrl: './game-rules.component.scss',
})
export class GameRulesComponent {
    closed = output<void>();

    close(): void {
        this.closed.emit();
    }

    onBackdropClick(event: MouseEvent): void {
        if (
            (event.target as HTMLElement).classList.contains('modal-backdrop')
        ) {
            this.close();
        }
    }
}
