import { Component } from '@angular/core';

export interface PlayerCardData {
    id: number;
    name: string;
}

@Component({
    selector: 'app-player-cards',
    imports: [],
    templateUrl: './player-cards.component.html',
    styleUrl: './player-cards.component.scss',
})
export class PlayerCardsComponent {
    players: PlayerCardData[] = [
        {
            id: 0,
            name: 'Алиса',
        },
        {
            id: 1,
            name: 'Борис',
        },
        {
            id: 2,
            name: 'Вера',
        },
        {
            id: 3,
            name: 'Гриша',
        },
    ];
    currentPlayerId = 0;
}
