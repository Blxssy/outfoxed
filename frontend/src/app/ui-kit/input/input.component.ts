import { Component, computed, input, signal } from '@angular/core';
import { NgClass } from '@angular/common';

@Component({
    selector: 'fox-input',
    imports: [NgClass],
    templateUrl: './input.component.html',
    styleUrl: './input.component.scss',
})
export class InputComponent {
    readonly placeholder = input<string>();
    readonly type = input<string>('text');
    readonly isDisabled = input<boolean>(false);

    readonly isInvalid = input<boolean>(false);
    readonly errorText = input<string>();

    protected isPasswordVisible = signal<boolean>(false);

    protected inputType = computed(() => {
        if (this.type() !== 'password') {
            return this.type();
        }

        return this.isPasswordVisible() ? 'text' : 'password';
    });

    protected getInputClasses(): string {
        return [
            'inp',
            this.isDisabled() ? 'inp--disabled' : '',
            this.isInvalid() ? 'inp--error' : '',
        ].join(' ');
    }

    protected togglePasswordVisibility(): void {
        this.isPasswordVisible.set(!this.isPasswordVisible());
    }
}
