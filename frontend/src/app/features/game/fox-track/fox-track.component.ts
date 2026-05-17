import { Component, computed, input } from '@angular/core';

@Component({
    selector: 'app-fox-track',
    imports: [],
    templateUrl: './fox-track.component.html',
    styleUrl: './fox-track.component.scss',
})
export class FoxTrackComponent {
    foxTrack = input<number>(0);
    escapeAt = input<number>(15);

    startStep = 1;
    finishStep = computed(() => this.escapeAt() + 1);
    foxPosition = computed(() => this.foxTrack() + 1);

    get steps(): number[] {
        return Array.from({ length: this.finishStep() - 1 }, (_, i) => i + 2);
    }
}
