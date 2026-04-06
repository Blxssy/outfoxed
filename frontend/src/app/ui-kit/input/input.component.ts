import {
    ChangeDetectionStrategy,
    ChangeDetectorRef,
    Component,
    computed,
    forwardRef,
    input,
    output,
    signal,
} from '@angular/core';
import {
    ControlValueAccessor,
    FormsModule,
    NG_VALUE_ACCESSOR,
    ReactiveFormsModule,
} from '@angular/forms';

type InputType = 'text' | 'password' | 'email' | 'number';
type InputSize = 'sm' | 'md' | 'lg';

@Component({
    selector: 'fox-input',
    imports: [FormsModule, ReactiveFormsModule],
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: forwardRef(() => InputComponent),
            multi: true,
        },
    ],
    templateUrl: './input.component.html',
    styleUrl: './input.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class InputComponent implements ControlValueAccessor {
    readonly label = input<string>();
    readonly type = input<InputType>('text');
    readonly size = input<InputSize>('md');
    readonly placeholder = input<string>();
    readonly disabled = input<boolean>(false);
    readonly suffixIcon = input<boolean>(false);
    readonly id = input<string>('');

    readonly valueChange = output<string>();
    readonly focused = output<void>();
    readonly blurred = output<void>();

    value = '';
    inputId = '';
    private static idCounter = 0;

    protected readonly passwordVisible = signal(false);
    protected readonly isFocused = signal(false);

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

    constructor(private cdr: ChangeDetectorRef) {}

    ngOnInit() {
        this.inputId = this.id() || `fox-input-${InputComponent.idCounter++}`;
    }

    togglePassword(): void {
        this.passwordVisible.update((p) => !p);
    }

    onInput(event: Event): void {
        this.value = (event.target as HTMLInputElement).value;
        this.onChange(this.value);
        this.valueChange.emit(this.value);
    }

    onFocus(): void {
        this.isFocused.set(true);
        this.focused.emit();
    }
    onBlur(): void {
        this.isFocused.set(false);
        this.onTouched();
        this.blurred.emit();
    }

    private onChange: (v: string) => void = () => {};
    private onTouched: () => void = () => {};

    writeValue(val: string): void {
        this.value = val ?? '';
        this.cdr.markForCheck();
    }
    registerOnChange(fn: (v: string) => void): void {
        this.onChange = fn;
    }
    registerOnTouched(fn: () => void): void {
        this.onTouched = fn;
    }
}
