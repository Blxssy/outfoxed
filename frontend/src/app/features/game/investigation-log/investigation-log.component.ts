import { DatePipe } from '@angular/common';
import {
    Component,
    computed,
    effect,
    ElementRef,
    inject,
    viewChild,
} from '@angular/core';

import { GameService } from '../data/game.service';

export interface JournalItem {
    id: string;
    turn: number;
    version: number;
    type: string;
    message: string;
    createdAt: string;
}

@Component({
    selector: 'app-investigation-log',
    imports: [DatePipe],
    templateUrl: './investigation-log.component.html',
    styleUrl: './investigation-log.component.scss',
})
export class InvestigationLogComponent {
    private readonly game = inject(GameService);

    private readonly scrollContainer =
        viewChild<ElementRef<HTMLElement>>('scrollContainer');

    readonly entries = computed<JournalItem[]>(() => {
        const state = this.game.gameState();

        return [...(state?.journal ?? [])].reverse();
    });

    constructor() {
        effect(() => {
            this.entries();

            queueMicrotask(() => {
                const el = this.scrollContainer()?.nativeElement;

                if (el) {
                    el.scrollTop = 0;
                }
            });
        });
    }
}
