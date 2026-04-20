import { Component, computed, input } from '@angular/core';

export type CardSize = 'sm' | 'md' | 'lg';
export type CardVariant = 'default' | 'flat' | 'inset';

@Component({
    selector: 'fox-card',
    imports: [],
    templateUrl: './card.component.html',
    styleUrl: './card.component.scss',
})
export class CardComponent {
    readonly size = input<CardSize>('md');
    readonly variant = input<CardVariant>('default');
    readonly classes = computed(() =>
        ['card', `card--${this.variant()}`, `card--${this.size()}`]
            .filter(Boolean)
            .join(' '),
    );
}
