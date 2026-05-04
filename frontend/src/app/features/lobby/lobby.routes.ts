import { Routes } from '@angular/router';
import { LobbyComponent } from './lobby.component';

export const LOBBY_ROUTES: Routes = [
    {
        path: '',
        component: LobbyComponent,
        children: [
            {
                path: '',
                loadComponent: () =>
                    import('./rooms-list/rooms-list.component').then(
                        (m) => m.RoomsListComponent,
                    ),
            },
            {
                path: ':id',
                loadComponent: () =>
                    import('./lobby-room/lobby-room.component').then(
                        (m) => m.LobbyRoomComponent,
                    ),
            },
        ],
    },
];
