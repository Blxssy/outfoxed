import { Component } from '@angular/core';

@Component({
    selector: 'app-fox-track',
    imports: [],
    templateUrl: './fox-track.component.html',
    styleUrl: './fox-track.component.scss',
})
export class FoxTrackComponent {
    state = { foxPosition: 16, totalSteps: 16 };

    get steps(): number[] {
        return Array.from(
            { length: this.state.totalSteps - 1 },
            (_, i) => this.state.totalSteps - 1 - i,
        );
    }
}
