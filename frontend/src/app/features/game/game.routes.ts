import { Routes } from '@angular/router';
import { GameComponent } from './game.component';

export const GAME_ROUTES: Routes = [
    {
        path: ':id',
        loadComponent: () =>
            import('./game.component').then((m) => m.GameComponent),
    },
];
