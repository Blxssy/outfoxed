import { Component, input } from '@angular/core';

export type PlayerCardData = {
    userId: string;
    seat: number;
    name: string;
    pawnCell: number;
    connected: boolean;
};

@Component({
    selector: 'app-player-cards',
    imports: [],
    templateUrl: './player-cards.component.html',
    styleUrl: './player-cards.component.scss',
})
export class PlayerCardsComponent {
    players = input.required<PlayerCardData[]>();

    currentPlayerSeat = input.required();
}
