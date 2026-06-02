import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { GameBoardComponent } from './game-board.component';

describe('GameBoardComponent', () => {
    beforeEach(async () => {
        vi.clearAllMocks();

        await TestBed.configureTestingModule({
            imports: [GameBoardComponent],
            schemas: [NO_ERRORS_SCHEMA],
        }).compileComponents();
    });

    const create = () => TestBed.createComponent(GameBoardComponent);

    it('should create component', () => {
        const fixture = create();
        expect(fixture.componentInstance).toBeTruthy();
    });

    it('should detect reachable cell when canMove true', () => {
        const fixture = create();

        fixture.componentRef.setInput('canMove', true);
        fixture.componentRef.setInput('reachableCells', [5]);
        fixture.detectChanges();

        expect(fixture.componentInstance.isReachable(5)).toBe(true);
    });

    it('should not detect reachable cell when canMove false', () => {
        const fixture = create();

        fixture.componentRef.setInput('canMove', false);
        fixture.componentRef.setInput('reachableCells', [5]);
        fixture.detectChanges();

        expect(fixture.componentInstance.isReachable(5)).toBe(false);
    });

    it('should mark clue as active', () => {
        const fixture = create();

        fixture.componentRef.setInput('collectedClueIndices', []);
        fixture.detectChanges();

        const cell = { index: 1, hasClue: true } as any;

        expect(fixture.componentInstance.hasActiveClue(cell)).toBe(true);
    });

    it('should mark clue as collected', () => {
        const fixture = create();

        fixture.componentRef.setInput('collectedClueIndices', [1]);
        fixture.detectChanges();

        const cell = { index: 1, hasClue: true } as any;

        expect(fixture.componentInstance.isCollected(cell)).toBe(true);
    });

    it('should emit click on valid cell', () => {
        const fixture = create();

        fixture.componentRef.setInput('canMove', true);
        fixture.componentRef.setInput('reachableCells', [2]);
        fixture.detectChanges();

        const spy = vi.spyOn(fixture.componentInstance.cellClicked, 'emit');

        fixture.componentInstance.onCellClick({
            index: 2,
        } as any);

        expect(spy).toHaveBeenCalledWith(2);
    });

    it('should not emit click on invalid cell', () => {
        const fixture = create();

        fixture.componentRef.setInput('canMove', false);
        fixture.detectChanges();

        const spy = vi.spyOn(fixture.componentInstance.cellClicked, 'emit');

        fixture.componentInstance.onCellClick({
            index: 2,
        } as any);

        expect(spy).not.toHaveBeenCalled();
    });

    it('should detect active player', () => {
        const fixture = create();

        fixture.componentRef.setInput('activeSeat', 1);
        fixture.detectChanges();

        expect(
            fixture.componentInstance.isActivePlayer({
                seat: 1,
            } as any),
        ).toBe(true);
    });

    it('should detect me player', () => {
        const fixture = create();

        fixture.componentRef.setInput('mySeat', 3);
        fixture.detectChanges();

        expect(
            fixture.componentInstance.isMe({
                seat: 3,
            } as any),
        ).toBe(true);
    });

    it('should return correct seat class', () => {
        const fixture = create();

        expect(fixture.componentInstance.seatClass(1)).toBe('player-1');
        expect(fixture.componentInstance.seatClass(99)).toBe('player-0');
    });

    it('should group players at cell', () => {
        const fixture = create();

        fixture.componentRef.setInput('players', [
            { pawnCell: 10, userId: '1' } as any,
            { pawnCell: 10, userId: '2' } as any,
        ]);

        fixture.detectChanges();

        const res = fixture.componentInstance.playersAt(10);

        expect(res.length).toBe(2);
        expect(res[0].total).toBe(2);
    });
});
