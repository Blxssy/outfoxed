import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {
    CreateGameRequest,
    CreateGameResponse,
    JoinByCodeRequest,
    JoinGameResponse,
    LeaveGameResponse,
    LobbySnapshotResponse,
    RoomListItem,
    StartGameResponse,
} from './lobby.model';

interface GamesListResponse {
    games: RoomListItem[];
}

@Injectable({ providedIn: 'root' })
export class LobbyApiService {
    private readonly http = inject(HttpClient);
    private readonly api = 'http://localhost:8080/api/v1/games';

    getPublicGames(): Observable<GamesListResponse> {
        return this.http.get<GamesListResponse>(this.api);
    }

    createGame(body: CreateGameRequest): Observable<CreateGameResponse> {
        console.log('request on create game');
        return this.http.post<CreateGameResponse>(this.api, body);
    }

    joinGame(id: string): Observable<JoinGameResponse> {
        return this.http.post<JoinGameResponse>(`${this.api}/${id}/join`, {});
    }

    joinByCode(body: JoinByCodeRequest): Observable<JoinGameResponse> {
        return this.http.post<JoinGameResponse>(
            `${this.api}/join-by-code`,
            body,
        );
    }

    getLobbySnapshot(id: string): Observable<LobbySnapshotResponse> {
        return this.http.get<LobbySnapshotResponse>(`${this.api}/${id}`);
    }

    leaveGame(id: string): Observable<LeaveGameResponse> {
        console.log('leave game resp');
        return this.http.post<LeaveGameResponse>(`${this.api}/${id}/leave`, {});
    }

    startGame(id: string): Observable<StartGameResponse> {
        return this.http.post<StartGameResponse>(`${this.api}/${id}/start`, {});
    }
}
