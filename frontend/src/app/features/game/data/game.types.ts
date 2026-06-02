export type GameStatus = 'waiting' | 'active' | 'finished';

export type GamePhase =
    | 'choose_goal'
    | 'rolling'
    | 'move_pawn'
    | 'resolve_clue'
    | 'reveal_suspects'
    | 'end_turn';

export type GameResult = 'none' | 'win' | 'lose';

export type GoalType = 'clue' | 'suspect';

export type AvailableAction =
    | 'choose_goal'
    | 'roll_auto'
    | 'reroll_dice'
    | 'finish_roll'
    | 'move_pawn'
    | 'take_clue'
    | 'reveal_suspects'
    | 'end_turn'
    | 'accuse';

export interface PublicPlayer {
    userId: string;
    seat: number;
    name: string;
    pawnCell: number;
    connected: boolean;
}

export interface MePlayer extends PublicPlayer {}

export type CellType = 'start' | 'clue' | 'path' | 'normal';

export interface BoardCell {
    index: number;
    x: number;
    y: number;
    type: CellType;
    hasClue: boolean;
    clueTokenId?: string;
}

export interface Board {
    width: number;
    height: number;
    cells: BoardCell[];
}

export interface MoveState {
    reachableCells: number[];
    stepsRemaining?: number;
}

export interface Fox {
    track: number;
    escapeAt: number;
}

export interface SuspectTraits {
    glasses?: 'yes' | 'no';
    hat?: 'yes' | 'no';
    scarf?: 'yes' | 'no';
    umbrella?: 'yes' | 'no';
    color?: string;
    [key: string]: string | undefined;
}

export interface Suspect {
    id: string;
    revealed: boolean;
    excluded: boolean;
    traits?: SuspectTraits;
}

export interface Clue {
    id: string;
    revealed: boolean;
    trait?: string;
    result?: 'yes' | 'no';
    boardCell?: number;
}

export interface RollState {
    rollsUsed: number;
    maxRolls: number;
    faces: string[];
    kept: boolean[];
    success?: boolean;
}

export interface PublicGameState {
    id: string;
    status: GameStatus;
    phase: GamePhase;
    result: GameResult;
    version: number;
    turn: number;
    activeSeat: number;
    me: MePlayer;
    players: PublicPlayer[];
    board: Board;
    fox: Fox;
    suspects: Suspect[];
    clues: Clue[];
    roll?: RollState;
    move?: MoveState;
    availableActions: AvailableAction[] | null;
    turnDeadlineAt?: string;
    journal: JournalItem[];
}

export interface JournalItem {
    id: string;
    turn: number;
    version: number;
    type: string;
    message: string;
    createdAt: string;
}

export type WsCommandName =
    | 'choose_goal'
    | 'roll_auto'
    | 'reroll_dice'
    | 'finish_roll'
    | 'move_pawn'
    | 'take_clue'
    | 'reveal_suspects'
    | 'end_turn'
    | 'accuse';

export interface ChooseGoalPayload {
    goal: GoalType;
}
export interface RollAutoPayload {}
export interface RerollDicePayload {
    keepIndices: number[];
}
export interface FinishRollPayload {}
export interface MovePawnPayload {
    targetIndex: number;
}
export interface TakeCluePayload {}
export interface RevealSuspectsPayload {
    suspectIds: string[];
}
export interface EndTurnPayload {}
export interface AccusePayload {
    suspectId: string;
}

export interface WsCommandPayloadMap {
    choose_goal: ChooseGoalPayload;
    roll_auto: RollAutoPayload;
    reroll_dice: RerollDicePayload;
    finish_roll: FinishRollPayload;
    move_pawn: MovePawnPayload;
    take_clue: TakeCluePayload;
    reveal_suspects: RevealSuspectsPayload;
    end_turn: EndTurnPayload;
    accuse: AccusePayload;
}

export interface WsClientMessage<T extends WsCommandName = WsCommandName> {
    id: string;
    type: 'command';
    command: T;
    payload: WsCommandPayloadMap[T];
}

export interface WsEvent {
    type: string;
    data?: Record<string, unknown>;
}

export interface WsUpdatePayload {
    state: PublicGameState;
    events: WsEvent[];
}

export interface WsErrorPayload {
    code: string;
    message: string;
}

export interface WsServerUpdate {
    id: string | null;
    type: 'update';
    payload: WsUpdatePayload;
}

export interface WsServerError {
    id: string | null;
    type: 'error';
    payload: WsErrorPayload;
}

export type WsServerMessage = WsServerUpdate | WsServerError;

export type WsConnectionStatus =
    | 'disconnected'
    | 'connecting'
    | 'connected'
    | 'reconnecting'
    | 'error';

export interface PendingRequest {
    requestId: string;
    command: WsCommandName;
    sentAt: number;
}

export type WsErrorCode =
    | 'not_your_turn'
    | 'invalid_phase'
    | 'invalid_move'
    | 'game_not_active'
    | 'game_finished'
    | 'goal_already_set'
    | 'goal_not_set'
    | 'no_pending_action'
    | 'pending_not_clue'
    | 'pending_not_suspect'
    | 'suspect_not_revealed'
    | 'suspect_excluded'
    | 'suspect_not_found'
    | 'invalid_reveal_selection'
    | 'forbidden'
    | 'internal_error';
