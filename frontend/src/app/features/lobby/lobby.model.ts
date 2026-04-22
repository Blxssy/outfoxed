export type RoomStatus = 'waiting' | 'active' | 'finished';

export type Room = {
    id: string;
    title: string;
    host_username: string;
    players_count: number;
    max_players: number;
    status: RoomStatus;
    visibility: boolean;
    joinCode?: number;
};
