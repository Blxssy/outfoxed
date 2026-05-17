import { Component, computed, input, output } from '@angular/core';
import { BoardCell, PublicPlayer } from '../data/game.types';

export const PLAYER_SEAT_CLASSES = [
    'player-0',
    'player-1',
    'player-2',
    'player-3',
] as const;

@Component({
    selector: 'app-game-board',
    standalone: true,
    imports: [],
    templateUrl: './game-board.component.html',
    styleUrl: './game-board.component.scss',
})
export class GameBoardComponent {
    cells = input<BoardCell[]>([]);
    players = input<PublicPlayer[]>([]);
    activeSeat = input<number>(-1);
    mySeat = input<number | undefined>(-1);
    reachableCells = input<number[]>([]);
    canMove = input(false);

    cellClicked = output<number>();

    readonly rows = computed<BoardCell[][]>(() => {
        const flat = this.cells();
        const grid: BoardCell[][] = Array.from({ length: 16 }, () => []);
        flat.forEach((cell) => {
            const row = Math.floor(cell.index / 16);
            grid[row]?.push(cell);
        });
        return grid;
    });

    private readonly reachableSet = computed(
        () => new Set(this.reachableCells()),
    );

    isReachable(cellIndex: number): boolean {
        return this.canMove() && this.reachableSet().has(cellIndex);
    }

    onCellClick(cell: BoardCell): void {
        if (!this.isReachable(cell.index)) return;
        this.cellClicked.emit(cell.index);
    }

    playersAt(cellIndex: number): { player: PublicPlayer; total: number }[] {
        const here = this.players().filter((p) => p.pawnCell === cellIndex);
        return here.map((player) => ({ player, total: here.length }));
    }

    isActivePlayer(player: PublicPlayer): boolean {
        return player.seat === this.activeSeat();
    }

    isMe(player: PublicPlayer): boolean {
        return player.seat === this.mySeat();
    }

    seatClass(seat: number): string {
        return PLAYER_SEAT_CLASSES[seat] ?? 'player-0';
    }
}
