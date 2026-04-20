export type RoomStatus = 'waiting' | 'full' | 'starting';

export type Room = {
    id: string;
    name: string;
    playersNow: number;
    playersMax: number;
    status: RoomStatus;
    isPrivate: boolean;
    hostName: string;
    createdAt: Date;
};
