import { Component, inject, OnInit, OnDestroy } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { ActiveGameService } from './data/active-game.service';
import { ButtonComponent } from '@fox/ui-kit/button';

@Component({
    selector: 'app-lobby',
    imports: [RouterOutlet, ButtonComponent],
    templateUrl: './lobby.component.html',
    styleUrl: './lobby.component.scss',
})
export class LobbyComponent implements OnInit, OnDestroy {
    private readonly router = inject(Router);
    readonly activeGameService = inject(ActiveGameService);

    ngOnInit(): void {
        this.activeGameService.startPolling();
    }

    ngOnDestroy(): void {
        this.activeGameService.stopPolling();
    }

    returnToGame(): void {
        const game = this.activeGameService.activeGame();
        if (game) {
            this.router.navigateByUrl(game.route);
        }
    }
}
