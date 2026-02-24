import { Component, input } from '@angular/core';
import { NgClass } from "@angular/common";

@Component({
    selector: 'fox-button',
    imports: [NgClass],
    templateUrl: './button.component.html',
    styleUrl: './button.component.scss',
})
export class ButtonComponent {
    readonly text = input.required<string>();
    readonly isDisabled = input<boolean>(false);
    readonly size = input<string>('default');
    readonly color = input<string>('orange');

    protected getButtonClasses(): string {
        return ['btn', `color-${this.color()}`, this.size()].join(' ');
    }
}
