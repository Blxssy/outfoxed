import { ChangeDetectorRef, Component } from '@angular/core';

export interface GameBoardState {
    cells: Cell[][];
    players: Player[];
    currentPlayerId: number;
    reachableCells: Set<string>;
}

export interface Cell {
    row: number;
    col: number;
}

export interface Player {
    id: number;
    color: string;
    name: string;
    row: number;
    col: number;
}

function getReachable(p: Player): Set<string> {
    const set = new Set<string>();
    const dirs = [
        [-1, 0],
        [1, 0],
        [0, -1],
        [0, 1],
        [-2, 0],
        [2, 0],
        [0, -2],
        [0, 2],
        [-1, -1],
        [-1, 1],
        [1, -1],
        [1, 1],
    ];
    for (const [dr, dc] of dirs) {
        const nr = p.row + dr,
            nc = p.col + dc;
        if (nr >= 0 && nr < 16 && nc >= 0 && nc < 16) set.add(`${nr},${nc}`);
    }
    return set;
}

@Component({
    selector: 'app-game-board',
    imports: [],
    templateUrl: './game-board.component.html',
    styleUrl: './game-board.component.scss',
})
export class GameBoardComponent {
    state: GameBoardState = this.buildDefaultState();

    get currentPlayer(): Player | undefined {
        return this.state.players.find(
            (p) => p.id === this.state.currentPlayerId,
        );
    }

    constructor(private cdr: ChangeDetectorRef) {}

    private buildDefaultState(): GameBoardState {
        const players: Player[] = [
            { id: 0, color: '#e05050', name: 'Алиса', row: 7, col: 7 },
            { id: 1, color: '#4d88e8', name: 'Борис', row: 7, col: 8 },
            { id: 2, color: '#44c060', name: 'Вера', row: 8, col: 7 },
            { id: 3, color: '#e0a030', name: 'Гриша', row: 8, col: 8 },
        ];
        return {
            cells: buildDefaultCells(),
            players,
            currentPlayerId: 0,
            reachableCells: getReachable(players[0]),
        };

        function buildDefaultCells(): Cell[][] {
            const board: Cell[][] = Array.from({ length: 16 }, (_, row) =>
                Array.from({ length: 16 }, (_, col) => ({
                    row,
                    col,
                    clue: null,
                })),
            );

            return board;
        }
    }

    isReachable(row: number, col: number): boolean {
        return this.state.reachableCells.has(`${row},${col}`);
    }

    onCellClick(cell: Cell): void {
        if (!this.isReachable(cell.row, cell.col)) return;

        const current = this.currentPlayer;
        if (!current) return;

        current.row = cell.row;
        current.col = cell.col;

        const players = this.state.players;
        const nextIdx = (players.indexOf(current) + 1) % players.length;
        this.state.currentPlayerId = players[nextIdx].id;
        this.state.reachableCells = getReachable(players[nextIdx]);

        this.cdr.markForCheck();
    }

    getPlayersAt(
        row: number,
        col: number,
    ): { player: Player; total: number }[] {
        const here = this.state.players.filter(
            (p) => p.row === row && p.col === col,
        );
        return here.map((player) => ({ player, total: here.length }));
    }
}
