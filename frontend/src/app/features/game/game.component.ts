import { Component } from '@angular/core';
import { GameBoardComponent } from './game-board/game-board.component';
import { FoxTrackComponent } from './fox-track/fox-track.component';
import { PlayerCardsComponent } from './player-cards/player-cards.component';
import { SuspectsComponent } from './suspects/suspects.component';
import { InvestigationLogComponent } from './investigation-log/investigation-log.component';
import { DiceRollComponent } from './dice-roll/dice-roll.component';

@Component({
    selector: 'app-game',
    imports: [
        GameBoardComponent,
        FoxTrackComponent,
        PlayerCardsComponent,
        SuspectsComponent,
        InvestigationLogComponent,
        DiceRollComponent,
    ],
    templateUrl: './game.component.html',
    styleUrl: './game.component.scss',
})
export class GameComponent {}
