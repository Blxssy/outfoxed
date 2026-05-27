import { Component, input } from '@angular/core';
import { Clue } from '../data/game.types';

const TRAIT_LABELS: Record<string, string> = {
    glasses: 'Очки',
    hat: 'Шляпа',
    scarf: 'Шарф',
    umbrella: 'Зонт',
    bag: 'Сумка',
    boots: 'Ботинки',
};

@Component({
    selector: 'app-clues-list',
    standalone: true,
    imports: [],
    templateUrl: './clues-list.component.html',
    styleUrl: './clues-list.component.scss',
})
export class CluesListComponent {
    clues = input<Clue[]>([]);
    totalCount = input<number>(0);

    traitLabel(trait: string | undefined): string {
        if (!trait) return '—';
        return TRAIT_LABELS[trait] ?? trait;
    }

    resultIcon(result: 'yes' | 'no' | undefined): string {
        return result === 'yes' ? '✅' : result === 'no' ? '❌' : '?';
    }

    resultLabel(result: 'yes' | 'no' | undefined): string {
        return result === 'yes' ? 'Да' : result === 'no' ? 'Нет' : '';
    }
}
