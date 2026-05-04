import { Routes } from '@angular/router';
import { AuthGuard } from './guards/auth.guard';

export const routes: Routes = [
    {
        path: 'game',
        loadComponent: () =>
            import('./features/game/game.component').then(
                (m) => m.GameComponent,
            ),
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
