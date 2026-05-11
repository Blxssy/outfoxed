import { Component, signal } from '@angular/core';

export interface Suspect {
    id: number;
    name: string;
    state: SuspectState;
}

export type SuspectState = 'unknown' | 'revealed' | 'innocent';

@Component({
    selector: 'app-suspects',
    imports: [],
    templateUrl: './suspects.component.html',
    styleUrl: './suspects.component.scss',
})
export class SuspectsComponent {
    suspects: Suspect[] = [
        {
            id: 0,
            name: 'mr Fox',
            state: 'revealed',
        },
        {
            id: 1,
            name: 'ms Foxy',
            state: 'unknown',
        },
        {
            id: 2,
            name: 'mini Fox',
            state: 'innocent',
        },
    ];
    // suspectClicked = output();

    get revealedCount(): number {
        return this.suspects.filter((s) => s.state === 'revealed').length;
    }
}
