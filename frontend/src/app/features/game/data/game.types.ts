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

export type PublicPlayer = {
    userId: string;
    seat: number;
    name: string;
    pawnCell: number;
    connected: boolean;
};

export interface MePlayer extends PublicPlayer {}

export type CellType = 'start' | 'clue' | 'path' | 'normal';

export type BoardCell = {
    index: number;
    x: number;
    y: number;
    type: CellType;
    hasClue: boolean;
    clueTokenId?: string;
};

export type Board = {
    width: number;
    height: number;
    cells: BoardCell[];
};

export type MoveState = {
    reachableCells: number[];
    stepsRemaining?: number;
};

export type Fox = {
    track: number;
    escapeAt: number;
};

export type SuspectTraits = {
    glasses?: 'yes' | 'no';
    hat?: 'yes' | 'no';
    scarf?: 'yes' | 'no';
    umbrella?: 'yes' | 'no';
    color?: string;
    [key: string]: string | undefined;
};

export type Suspect = {
    id: string;
    revealed: boolean;
    excluded: boolean;
    traits?: SuspectTraits;
};

export type Clue = {
    id: string;
    revealed: boolean;
    trait?: string;
    result?: 'yes' | 'no';
    boardCell?: number;
};

export type RollState = {
    rollsUsed: number;
    maxRolls: number;
    faces: string[];
    kept: boolean[];
    success?: boolean;
};

export type PublicGameState = {
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
};

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

export type ChooseGoalPayload = {
    goal: GoalType;
};
export type RollAutoPayload = {};
export type RerollDicePayload = {
    keepIndices: number[];
};
export type FinishRollPayload = {};
export type MovePawnPayload = {
    targetIndex: number;
};
export type TakeCluePayload = {};
export type RevealSuspectsPayload = {
    suspectIds: string[];
};
export type EndTurnPayload = {};
export type AccusePayload = {
    suspectId: string;
};

export type WsCommandPayloadMap = {
    choose_goal: ChooseGoalPayload;
    roll_auto: RollAutoPayload;
    reroll_dice: RerollDicePayload;
    finish_roll: FinishRollPayload;
    move_pawn: MovePawnPayload;
    take_clue: TakeCluePayload;
    reveal_suspects: RevealSuspectsPayload;
    end_turn: EndTurnPayload;
    accuse: AccusePayload;
};

export interface WsClientMessage<T extends WsCommandName = WsCommandName> {
    id: string;
    type: 'command';
    command: T;
    payload: WsCommandPayloadMap[T];
}

export type WsEvent = {
    type: string;
    data?: Record<string, unknown>;
};

export type WsUpdatePayload = {
    state: PublicGameState;
    events: WsEvent[];
};

export type WsErrorPayload = {
    code: string;
    message: string;
};

export type WsServerUpdate = {
    id: string | null;
    type: 'update';
    payload: WsUpdatePayload;
};

export type WsServerError = {
    id: string | null;
    type: 'error';
    payload: WsErrorPayload;
};

export type WsServerMessage = WsServerUpdate | WsServerError;

export type WsConnectionStatus =
    | 'disconnected'
    | 'connecting'
    | 'connected'
    | 'reconnecting'
    | 'error';

export type PendingRequest = {
    requestId: string;
    command: WsCommandName;
    sentAt: number;
};

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
