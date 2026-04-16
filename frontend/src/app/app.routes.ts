import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'game',
    loadComponent: () =>
      import('./features/game/game.component').then((m) => m.GameComponent),
  },
  {
    path: 'lobby',
    loadComponent: () =>
      import('./features/lobby/lobby.component').then((m) => m.LobbyComponent),
  },
  {
    path: 'auth',
    loadChildren: () =>
      import('./features/auth/auth.routes').then((m) => m.routes),
  },
  {
    path: '',
    redirectTo: '/lobby',
    pathMatch: 'full',
  },
];
