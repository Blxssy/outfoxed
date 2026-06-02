import { Component, input, output, signal, computed } from '@angular/core';

export interface SuspectTraits {
    glasses?: string;
    hat?: string;
    scarf?: string;
    umbrella?: string;
    bag?: string;
    boots?: string;
    gloves?: string;
    badge?: string;
    book?: string;
    camera?: string;
    key?: string;
    watch?: string;
    [key: string]: string | undefined;
}

export interface Suspect {
    id: string;
    name?: string;
    revealed: boolean;
    excluded: boolean;
    traits?: SuspectTraits;
}

export const TRAIT_ICONS: Record<string, string> = {
    glasses: '👓',
    hat: '🎩',
    scarf: '🧣',
    umbrella: '☂️',
    bag: '👜',
    boots: '🥾',
    gloves: '🧤',
    badge: '📛',
    book: '📚',
    camera: '📷',
    key: '🔑',
    watch: '⌚',
};

export const TRAIT_LABELS: Record<string, string> = {
    glasses: 'Очки',
    hat: 'Шляпа',
    scarf: 'Шарф',
    umbrella: 'Зонт',
    bag: 'Сумка',
    boots: 'Ботинки',
};

@Component({
    selector: 'app-suspects',
    imports: [],
    templateUrl: './suspects.component.html',
    styleUrl: './suspects.component.scss',
})
export class SuspectsComponent {
    suspects = input.required<Suspect[]>();
    canReveal = input(false);
    canAccuse = input(false);
    selectedIds = input<string[]>([]);
    disabled = input(false);

    suspectToggled = output<string>();
    revealConfirmed = output<void>();
    accused = output<string>();

    readonly localExcluded = signal<Set<string>>(new Set());

    readonly revealedCount = computed(
        () => this.suspects().filter((s) => s.revealed).length,
    );

    readonly canConfirmReveal = computed(
        () => this.canReveal() && this.selectedIds().length === 2,
    );

    isSelected(id: string): boolean {
        return this.selectedIds().includes(id);
    }

    isExcluded(s: Suspect): boolean {
        return s.excluded || this.localExcluded().has(s.id);
    }

    getActiveTraits(
        traits: SuspectTraits,
    ): { key: string; icon: string; label: string }[] {
        return Object.entries(traits)
            .filter(([, v]) => v === 'yes')
            .map(([key]) => ({
                key,
                icon: TRAIT_ICONS[key] ?? '❓',
                label: TRAIT_LABELS[key] ?? key,
            }));
    }

    getSuspectName(s: Suspect): string {
        return s.name ?? s.id;
    }

    onCardClick(s: Suspect): void {
        if (this.disabled() || this.isExcluded(s)) return;
        if (this.canReveal() && !s.revealed) {
            this.suspectToggled.emit(s.id);
        }
    }

    excludeSuspect(s: Suspect, event: Event): void {
        event.stopPropagation();
        if (this.disabled()) return;
        this.localExcluded.update((set) => new Set([...set, s.id]));
    }

    reinstateSuspect(s: Suspect, event: Event): void {
        event.stopPropagation();
        if (this.disabled()) return;
        this.localExcluded.update((set) => {
            const next = new Set(set);
            next.delete(s.id);
            return next;
        });
    }

    accuseSuspect(s: Suspect, event: Event): void {
        event.stopPropagation();
        if (this.disabled() || !this.canAccuse()) return;
        this.accused.emit(s.id);
    }
}
