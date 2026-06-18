export interface PlayerState {
  id: string;
  name: string;
  role: string;
  is_alive: boolean;
  faction: string;
  has_heal: boolean;
  has_poison: boolean;
}

export type Phase = 'night' | 'night_witch' | 'day';

export interface GameState {
  id: string;
  round: number;
  phase: Phase;
  players: PlayerState[];
  victim?: string;
  result?: string;
}

export type EventKind = 'killed' | 'eliminated' | 'saved' | 'revealed' | 'no_consensus';

export interface GameEvent {
  kind: EventKind;
  player_id: string;
  detail?: string;
}

export type DecisionKind =
  | 'werewolf_attack'
  | 'seer_investigate'
  | 'witch_save'
  | 'witch_poison'
  | 'vote';

export interface RequiredDecision {
  kind: DecisionKind;
  actor_id: string;
  candidates: PlayerState[];
  optional: boolean;
}

export interface Decision {
  kind: DecisionKind;
  actor_id: string;
  target_id: string;
}

export interface ServerMsg {
  type: 'update' | 'joined' | 'error';
  state?: GameState;
  events?: GameEvent[];
  pending?: RequiredDecision[];
  connected?: string[];
  name?: string;
  error?: string;
}

export interface ClientMsg {
  type: 'join' | 'decide';
  name?: string;
  decisions?: Decision[];
}
