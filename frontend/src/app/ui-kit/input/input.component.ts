import {
    ChangeDetectionStrategy,
    Component,
    computed,
    input,
    signal,
} from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

type InputType = 'text' | 'password' | 'email' | 'number';
type InputSize = 'sm' | 'md' | 'lg';

@Component({
    selector: 'fox-input',
    imports: [FormsModule, ReactiveFormsModule],
    templateUrl: './input.component.html',
    styleUrl: './input.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class InputComponent {
    readonly label = input<string>();
    readonly type = input<InputType>();
    readonly size = input<InputSize>('md');
    readonly placeholder = input<string>();
    readonly disabled = input<boolean>(false);
    readonly suffixIcon = input<boolean>(false);

    protected readonly passwordVisible = signal(false);

    protected readonly resolvedType = computed(() =>
        this.type() === 'password'
            ? this.passwordVisible()
                ? 'text'
                : 'password'
            : this.type(),
    );

    protected readonly wrapperClasses = computed(() =>
        ['inp', `inp--${this.size()}`, this.disabled() ? 'inp--disabled' : '']
            .filter(Boolean)
            .join(' '),
    );
    protected readonly inputClasses = computed(() =>
        [
            'inp__native',
            this.suffixIcon() || this.type() === 'password'
                ? 'inp__native--suffix'
                : '',
        ]
            .filter(Boolean)
            .join(' '),
    );

    togglePassword() {
        this.passwordVisible.update((p) => !p);
    }
}
