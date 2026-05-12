import { Component } from '@angular/core';

export type DiceTarget = 'eyes' | 'paws';

export interface DiceRollResult {
    eyes: number;
    paws: number;
    target: DiceTarget;
}

@Component({
    selector: 'app-dice-roll',
    imports: [],
    templateUrl: './dice-roll.component.html',
    styleUrl: './dice-roll.component.scss',
})
export class DiceRollComponent {
    canRoll = true;
    selectedTarget: DiceTarget = 'eyes';

    selectTarget(t: DiceTarget): void {}
}
