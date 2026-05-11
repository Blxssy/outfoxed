import { Component } from '@angular/core';
import { GameBoardComponent } from './game-board/game-board.component';
import { FoxTrackComponent } from './fox-track/fox-track.component';

@Component({
    selector: 'app-game',
    imports: [GameBoardComponent, FoxTrackComponent],
    templateUrl: './game.component.html',
    styleUrl: './game.component.scss',
})
export class GameComponent {}
