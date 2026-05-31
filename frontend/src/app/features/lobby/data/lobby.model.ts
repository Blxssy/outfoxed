export type GameStatus = 'waiting' | 'active' | 'finished';
export type GameVisibility = 'public' | 'private';

export interface RoomListItem {
    id: string;
    title: string;
    host_username: string;
    players_count: number;
    max_players: number;
    status: GameStatus;
}

export interface CreateGameRequest {
    title: string;
    visibility: GameVisibility;
}

export interface CreateGameResponse {
    game: { id: string; status: GameStatus };
    player: { user_id: string; seat: number };
    title: string;
    visibility: GameVisibility;
    joinCode?: string;
}

export interface JoinGameResponse {
    game: { id: string; status: GameStatus };
    player: { user_id: string; seat: number };
}

export interface JoinByCodeRequest {
    code: string;
}

export interface LobbyPlayer {
    user_id: string;
    seat: number;
    display_name: string;
    is_me: boolean;
}

export interface LobbySnapshot {
    id: string;
    title: string;
    status: GameStatus;
    visibility: GameVisibility;
    joinCode?: string;
    host_username: string;
    players: LobbyPlayer[];
    can_start: boolean;
    min_players: number;
    max_players: number;
}

export interface LobbySnapshotResponse {
    game: LobbySnapshot;
}

export interface LeaveGameResponse {
    game_deleted: boolean;
    new_host_username?: string;
}

export interface StartGameResponse {
    game: { id: string; status: GameStatus };
    redirect: { route: string };
}

export interface LobbyWsUpdate {
    type: 'update';
    payload: {
        state: LobbySnapshot;
        events: unknown;
    };
}

export interface LobbyWsError {
    type: 'error';
    payload: { code: string; message: string };
}

export type LobbyWsMessage = LobbyWsUpdate | LobbyWsError;

export interface ActiveGameResponse {
    found: boolean;
    game?: { id: string; status: string };
    route?: string;
}
