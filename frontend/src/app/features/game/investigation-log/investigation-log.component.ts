import { DatePipe } from '@angular/common';
import { Component } from '@angular/core';

export interface JournalEntry {
    id: number;
    type: string;
    html: string;
    timestamp: Date;
}

@Component({
    selector: 'app-investigation-log',
    imports: [DatePipe],
    templateUrl: './investigation-log.component.html',
    styleUrl: './investigation-log.component.scss',
})
export class InvestigationLogComponent {
    entries: JournalEntry[] = [];

    get reversed() {
        return [...this.entries].reverse();
    }
}
