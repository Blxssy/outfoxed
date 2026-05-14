import { Routes } from '@angular/router';
import { AuthGuard } from './guards/auth.guard';

export const routes: Routes = [
    {
        path: 'game',
        canActivate: [AuthGuard],
        loadChildren: () =>
            import('./features/game/game.routes').then((m) => m.GAME_ROUTES),
    },
    {
        path: 'lobby',
        canActivate: [AuthGuard],
        loadChildren: () =>
            import('./features/lobby/lobby.routes').then((m) => m.LOBBY_ROUTES),
    },
    {
        path: 'auth',
        loadChildren: () =>
            import('./features/auth/auth.routes').then((m) => m.routes),
    },
    {
        path: '**',
        redirectTo: '/auth',
        pathMatch: 'full',
    },
];
