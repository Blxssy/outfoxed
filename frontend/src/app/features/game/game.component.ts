import { Component } from '@angular/core';
import { GameBoardComponent } from './game-board/game-board.component';
import { FoxTrackComponent } from './fox-track/fox-track.component';
import { PlayerCardsComponent } from './player-cards/player-cards.component';

@Component({
    selector: 'app-game',
    imports: [GameBoardComponent, FoxTrackComponent, PlayerCardsComponent],
    templateUrl: './game.component.html',
    styleUrl: './game.component.scss',
})
export class GameComponent {}
