import { Component, computed, effect, inject, signal } from '@angular/core';
import { GameBoardComponent } from './game-board/game-board.component';
import { FoxTrackComponent } from './fox-track/fox-track.component';
import { PlayerCardsComponent } from './player-cards/player-cards.component';
import { SuspectsComponent } from './suspects/suspects.component';
import { InvestigationLogComponent } from './investigation-log/investigation-log.component';
import { DiceRollComponent } from './dice-roll/dice-roll.component';
import { GameService } from './data/game.service';
import { ActivatedRoute, Router } from '@angular/router';
import { GoalType } from './data/game.types';
import { CluesListComponent } from './clues-list/clues-list.component';
import { FinishModalComponent } from './finish-modal/finish-modal.component';

@Component({
    selector: 'app-game',
    imports: [
        GameBoardComponent,
        FoxTrackComponent,
        PlayerCardsComponent,
        SuspectsComponent,
        InvestigationLogComponent,
        DiceRollComponent,
        CluesListComponent,
        FinishModalComponent,
    ],
    templateUrl: './game.component.html',
    styleUrl: './game.component.scss',
})
export class GameComponent {
    protected readonly game = inject(GameService);
    private readonly route = inject(ActivatedRoute);
    private readonly router = inject(Router);

    readonly selectedSuspects = signal<string[]>([]);

    readonly isLeaving = signal(false);

    private errorTimer: ReturnType<typeof setTimeout> | null = null;
    private gameId = '';

    constructor() {
        effect(() => {
            if (this.game.lastError()) {
                if (this.errorTimer) clearTimeout(this.errorTimer);
                this.errorTimer = setTimeout(
                    () => this.game.lastError.set(null),
                    4_000,
                );
            }
        });
        effect(() => {
            if (
                this.game.phase() === 'end_turn' &&
                this.game.canDo('end_turn')
            ) {
                this.game.endTurn();
            }
        });
        effect(() => {
            if (
                this.game.phase() === 'resolve_clue' &&
                this.game.canDo('take_clue') &&
                !this.game.hasPending()
            ) {
                this.game.takeClue();
            }
        });
    }

    ngOnInit(): void {
        this.gameId = this.route.snapshot.paramMap.get('id') ?? '';
        if (this.gameId) this.game.startSession(this.gameId);
    }

    ngOnDestroy(): void {
        this.game.endSession();
        if (this.errorTimer) clearTimeout(this.errorTimer);
    }

    readonly state = this.game.gameState;
    readonly isMyTurn = this.game.isMyTurn;
    readonly activeSeat = this.game.activeSeat;
    readonly myPlayer = this.game.myPlayer;
    readonly phase = this.game.phase;
    readonly players = this.game.players;
    readonly suspects = this.game.suspects;
    readonly clues = this.game.clues;
    readonly fox = this.game.foxTrack;
    readonly eventLog = this.game.eventLog;
    readonly isReconnecting = this.game.isReconnecting;
    readonly errorToast = computed(
        () => this.game.lastError()?.message ?? null,
    );

    readonly activePlayer = computed(() => {
        const s = this.state();
        return s
            ? (s.players.find((p) => p.seat === s.activeSeat) ?? null)
            : null;
    });

    readonly canChooseGoal = computed(() => this.game.canDo('choose_goal'));
    readonly canRoll = computed(() => this.game.canDo('roll_auto'));
    readonly boardCells = this.game.boardCells;
    readonly reachableCells = this.game.reachableCells;
    readonly canMovePawn = computed(() => this.game.canDo('move_pawn'));
    readonly canTakeClue = computed(() => this.game.canDo('take_clue'));
    readonly canRevealSuspects = computed(() =>
        this.game.canDo('reveal_suspects'),
    );
    readonly canEndTurn = computed(() => this.game.canDo('end_turn'));
    readonly canAccuse = computed(() => this.game.canDo('accuse'));
    readonly canReroll = computed(() => this.game.canDo('reroll_dice'));
    readonly canFinishRoll = computed(() => this.game.canDo('finish_roll'));
    readonly hasPending = computed(() => this.game.hasPending());

    readonly roll = this.game.roll;
    readonly faces = computed(() => this.game.roll()?.faces ?? []);
    readonly rerollsLeft = computed(() => {
        const r = this.game.roll();
        return r ? r.maxRolls - r.rollsUsed : 0;
    });

    readonly goalType = computed<GoalType | null>(() => {
        const events = this.game.eventLog();
        for (let i = events.length - 1; i >= 0; i--) {
            const e = events[i];
            if (e.type === 'goal_chosen' && e.data?.['goal']) {
                return e.data['goal'] as GoalType;
            }
        }
        return null;
    });

    readonly isFinished = computed(() => this.state()?.status === 'finished');
    readonly gameResult = computed(() => this.state()?.result ?? 'none');

    readonly canConfirmReveal = computed(
        () => this.canRevealSuspects() && this.selectedSuspects().length === 2,
    );

    readonly revealedClues = computed(() =>
        this.clues().filter((c) => c.revealed),
    );
    readonly totalCluesCount = computed(() => this.clues().length);
    readonly collectedClueIndices = computed(() =>
        this.clues()
            .filter((c) => c.revealed)
            .map((c) => c.boardCell)
            .filter((idx): idx is number => idx !== undefined),
    );

    stepsOptions(): number[] {
        return [1, 2, 3];
    }

    leaveGame(): void {
        this.game.endSession();
        this.router.navigate(['/lobby']);
    }

    chooseGoal(goal: GoalType): void {
        this.game.chooseGoal(goal);
    }

    onCellClicked(targetIndex: number): void {
        this.game.movePawn(targetIndex);
    }

    rerollDice(keepIndices: number[]): void {
        this.game.rerollDice(keepIndices);
    }

    finishRoll(): void {
        this.game.finishRoll();
    }

    movePawn(steps: number): void {
        this.game.movePawn(steps);
    }

    takeClue(): void {
        this.game.takeClue();
    }

    endTurn(): void {
        this.game.endTurn();
    }

    toggleSuspect(id: string): void {
        if (!this.canRevealSuspects()) return;
        this.selectedSuspects.update((sel) => {
            if (sel.includes(id)) return sel.filter((s) => s !== id);
            if (sel.length >= 2) return sel;
            return [...sel, id];
        });
    }

    confirmReveal(): void {
        if (!this.canConfirmReveal()) return;
        this.game.revealSuspects(this.selectedSuspects());
        this.selectedSuspects.set([]);
    }

    accuse(suspectId: string): void {
        this.game.accuse(suspectId);
    }

    isSuspectSelected(id: string): boolean {
        return this.selectedSuspects().includes(id);
    }
}
