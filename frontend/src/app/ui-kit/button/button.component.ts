import { Component, input } from '@angular/core';
import { NgClass } from '@angular/common';

@Component({
    selector: 'fox-button',
    imports: [NgClass],
    templateUrl: './button.component.html',
    styleUrl: './button.component.scss',
})
export class ButtonComponent {
    readonly text = input.required<string>();
    readonly isDisabled = input<boolean>(false);
    readonly size = input<string>('md');
    readonly color = input<string>('orange');
    readonly fullWidth = input<boolean>(false);

    protected getButtonClasses(): string {
        return [
            'btn',
            `btn-${this.color()}`,
            `btn-${this.size()}`,
            this.fullWidth() ? 'btn--full' : '',
        ].join(' ');
    }
}
