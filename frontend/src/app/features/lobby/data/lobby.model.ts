export type GameStatus = 'waiting' | 'active' | 'finished';
export type GameVisibility = 'public' | 'private';

export type RoomListItem = {
    id: string;
    title: string;
    host_username: string;
    players_count: number;
    max_players: number;
    status: GameStatus;
};

export type CreateGameRequest = {
    title: string;
    visibility: GameVisibility;
};

export type CreateGameResponse = {
    game: { id: string; status: GameStatus };
    player: { user_id: string; seat: number };
    title: string;
    visibility: GameVisibility;
    joinCode?: string;
};

export type JoinGameResponse = {
    game: { id: string; status: GameStatus };
    player: { user_id: string; seat: number };
};

export type JoinByCodeRequest = {
    code: string;
};

export type LobbyPlayer = {
    user_id: string;
    seat: number;
    display_name: string;
    is_me: boolean;
};

export type LobbySnapshot = {
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
};

export type LobbySnapshotResponse = {
    game: LobbySnapshot;
};

export type LeaveGameResponse = {
    game_deleted: boolean;
    new_host_username?: string;
};

export type StartGameResponse = {
    game: { id: string; status: GameStatus };
    redirect: { route: string };
};
