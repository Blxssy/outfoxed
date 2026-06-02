import { Injectable, inject, signal, OnDestroy } from '@angular/core';
import { LobbyApiService } from './lobby-api.service';

@Injectable({ providedIn: 'root' })
export class ActiveGameService implements OnDestroy {
    private readonly api = inject(LobbyApiService);

    readonly activeGame = signal<{ id: string; route: string } | null>(null);
    readonly isChecking = signal(false);

    private pollTimer: ReturnType<typeof setInterval> | null = null;

    startPolling(intervalMs = 10_000): void {
        this.check();
        this.pollTimer = setInterval(() => this.check(), intervalMs);
    }

    stopPolling(): void {
        if (this.pollTimer) {
            clearInterval(this.pollTimer);
            this.pollTimer = null;
        }
    }

    check(): void {
        this.isChecking.set(true);
        this.api.getActiveGame().subscribe({
            next: (res) => {
                this.isChecking.set(false);
                if (
                    res.found &&
                    res.game &&
                    res.route &&
                    res.game.status === 'active'
                ) {
                    this.activeGame.set({ id: res.game.id, route: res.route });
                } else {
                    this.activeGame.set(null);
                }
            },
            error: () => {
                this.isChecking.set(false);
                this.activeGame.set(null);
            },
        });
    }

    clear(): void {
        this.activeGame.set(null);
    }

    ngOnDestroy(): void {
        this.stopPolling();
    }
}
