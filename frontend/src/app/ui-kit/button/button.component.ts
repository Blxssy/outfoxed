import { Component, input } from '@angular/core';
import { NgClass } from '@angular/common';

type ButtonType = 'button' | 'submit' | 'reset';
type ButtonColor = 'orange' | 'green';
type ButtonSize = 'sm' | 'md' | 'lg';

@Component({
    selector: 'fox-button',
    imports: [NgClass],
    templateUrl: './button.component.html',
    styleUrl: './button.component.scss',
})
export class ButtonComponent {
    readonly text = input.required<string>();
    readonly isDisabled = input<boolean>(false);
    readonly size = input<ButtonSize>('md');
    readonly color = input<ButtonColor>('orange');
    readonly fullWidth = input<boolean>(false);
    readonly type = input<ButtonType>('button');

    protected getButtonClasses(): string {
        return [
            'btn',
            `btn-${this.color()}`,
            `btn-${this.size()}`,
            this.fullWidth() ? 'btn--full' : '',
        ].join(' ');
    }
}
