import { Component, input } from '@angular/core';

export type Suspect = {
    id: string;
    revealed: boolean;
    excluded: boolean;
};

@Component({
    selector: 'app-suspects',
    imports: [],
    templateUrl: './suspects.component.html',
    styleUrl: './suspects.component.scss',
})
export class SuspectsComponent {
    suspects = input.required<Suspect[]>();

    get revealedCount(): number {
        return this.suspects().filter((s) => s.revealed === true).length;
    }
}
