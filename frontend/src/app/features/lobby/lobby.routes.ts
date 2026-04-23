import { Routes } from '@angular/router';

export const LOBBY_ROUTES: Routes = [
    {
        path: '',
        loadComponent: () =>
            import('./lobby.component').then((m) => m.LobbyComponent),
    },
    {
        path: ':id',
        loadComponent: () =>
            import('./lobby-room/lobby-room.component').then(
                (m) => m.LobbyRoomComponent,
            ),
    },
];
