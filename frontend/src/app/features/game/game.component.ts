import { Component } from '@angular/core';
import { GameBoardComponent } from './game-board/game-board.component';

@Component({
    selector: 'app-game',
    imports: [GameBoardComponent],
    templateUrl: './game.component.html',
    styleUrl: './game.component.scss',
})
export class GameComponent {}
