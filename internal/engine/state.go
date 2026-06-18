package engine

type PlayerState struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	IsAlive   bool   `json:"is_alive"`
	Faction   string `json:"faction"`
	HasHeal   bool   `json:"has_heal"`
	HasPoison bool   `json:"has_poison"`
}

type Phase string

const (
	PhaseNight      Phase = "night"
	PhaseNightWitch Phase = "night_witch"
	PhaseDay        Phase = "day"
)

type GameState struct {
	ID               string        `json:"id"`
	Round            int           `json:"round"`
	Phase            Phase         `json:"phase"`
	Players          []PlayerState `json:"players"`
	Victim           string        `json:"victim,omitempty"`
	Result           string        `json:"result,omitempty"`
	SeerInvestigated []string      `json:"seer_investigated,omitempty"`
}

type EventKind string

const (
	EventKilled      EventKind = "killed"
	EventEliminated  EventKind = "eliminated"
	EventSaved       EventKind = "saved"
	EventRevealed    EventKind = "revealed"
	EventNoConsensus EventKind = "no_consensus"
)

type Event struct {
	Kind     EventKind `json:"kind"`
	PlayerID string    `json:"player_id"`
	Detail   string    `json:"detail"`
}

type DecisionKind string

const (
	DecisionWerewolfAttack  DecisionKind = "werewolf_attack"
	DecisionSeerInvestigate DecisionKind = "seer_investigate"
	DecisionWitchSave       DecisionKind = "witch_save"
	DecisionWitchPoison     DecisionKind = "witch_poison"
	DecisionVote            DecisionKind = "vote"
)

type RequiredDecision struct {
	Kind       DecisionKind  `json:"kind"`
	ActorID    string        `json:"actor_id"`
	Candidates []PlayerState `json:"candidates"`
	Optional   bool          `json:"optional"`
}

type Decision struct {
	Kind     DecisionKind `json:"kind"`
	ActorID  string       `json:"actor_id"`
	TargetID string       `json:"target_id"`
}

type StepResult struct {
	State   GameState          `json:"state"`
	Events  []Event            `json:"events"`
	Pending []RequiredDecision `json:"pending"`
}
