import { Component, input, output } from '@angular/core';
import { GameResult } from '../data/game.types';

@Component({
    selector: 'app-finish-modal',
    imports: [],
    templateUrl: './finish-modal.component.html',
    styleUrl: './finish-modal.component.scss',
})
export class FinishModalComponent {
    result = input.required<GameResult>();
    lobbyClicked = output<void>();
}
