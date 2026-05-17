import { Component, computed, signal, input, output } from '@angular/core';
import { GoalType, RollState } from '../data/game.types';

const FACE_SYMBOLS: Record<string, string> = {
    footprint: '🐾',
    eye: '👁️',
};

@Component({
    selector: 'app-dice-roll',
    standalone: true,
    imports: [],
    templateUrl: './dice-roll.component.html',
    styleUrl: './dice-roll.component.scss',
})
export class DiceRollComponent {
    roll = input<RollState | null>(null);
    goalType = input<GoalType | null>(null);
    isMyTurn = input(false);
    canChooseGoal = input(false);
    canReroll = input(false);
    canFinishRoll = input(false);
    disabled = input(false);

    goalChosen = output<GoalType>();
    rerolled = output<number[]>();
    rollFinished = output<void>();

    readonly keptIndices = signal<number[]>([]);
    readonly isRolling = signal(false);

    readonly faces = computed<string[]>(() => {
        const f = this.roll()?.faces ?? [];
        return [0, 1, 2].map((i) => f[i] ?? 'blank');
    });

    readonly rerollsLeft = computed(() => {
        const r = this.roll();
        return r ? r.maxRolls - r.rollsUsed : 0;
    });

    readonly rollDone = computed(() => !!this.roll());
    readonly goalIsSet = computed(() => !!this.goalType() || this.rollDone());

    readonly diceIndices = [0, 1, 2] as const;

    chooseGoal(goal: GoalType): void {
        if (!this.canChooseGoal() || this.disabled()) return;
        this.keptIndices.set([]);
        this.goalChosen.emit(goal);
    }

    toggleKeep(index: number): void {
        if (!this.canReroll() || this.disabled()) return;
        this.keptIndices.update((prev) =>
            prev.includes(index)
                ? prev.filter((i) => i !== index)
                : [...prev, index],
        );
    }

    reroll(): void {
        if (!this.canReroll() || this.disabled()) return;
        this.isRolling.set(true);
        setTimeout(() => {
            this.isRolling.set(false);
            this.rerolled.emit(this.keptIndices());
            this.keptIndices.set([]);
        }, 350);
    }

    finishRoll(): void {
        if (!this.canFinishRoll() || this.disabled()) return;
        this.rollFinished.emit();
        this.keptIndices.set([]);
    }

    isKept(index: number): boolean {
        return this.keptIndices().includes(index);
    }

    faceSymbol(face: string): string {
        return FACE_SYMBOLS[face] ?? face;
    }

    isFaceMatch(face: string): boolean {
        const goal = this.goalType();
        if (goal === 'clue') return face === 'footprint';
        if (goal === 'suspect') return face === 'eye';
        return false;
    }

    goalLabel(goal: GoalType | null): string {
        if (goal === 'clue') return 'Искать улики';
        if (goal === 'suspect') return 'Опросить свидетелей';
        return '';
    }

    targetFace(): string {
        const goal = this.goalType();
        if (goal === 'clue') return '🐾 след';
        if (goal === 'suspect') return '👁️ глаз';
        return '';
    }
}
