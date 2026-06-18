package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/SupRemaZie/loupGOrou/internal/engine"
)

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9]`)

func sanitize(s string) string {
	return nonAlphaNum.ReplaceAllString(s, "_")
}

// toolID builds a unique, API-safe tool name from a required decision.
func toolID(req engine.RequiredDecision) string {
	return sanitize(string(req.Kind)) + "__" + sanitize(req.ActorID)
}

func ollamaHost() string {
	if h := os.Getenv("OLLAMA_HOST"); h != "" {
		return h
	}
	return "http://localhost:11434"
}

func ollamaModel() string {
	if m := os.Getenv("OLLAMA_MODEL"); m != "" {
		return m
	}
	return "nemotron-3-super:cloud"
}

// --- Types pour l'API native Ollama (/api/chat) ---

type ollamaMessage struct {
	Role      string           `json:"role"`
	Content   string           `json:"content,omitempty"`
	ToolCalls []ollamaToolCall `json:"tool_calls,omitempty"`
}

type ollamaToolCall struct {
	Function ollamaFunctionCall `json:"function"`
}

// Arguments est un objet JSON (pas une string) dans l'API native Ollama.
type ollamaFunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type ollamaTool struct {
	Type     string             `json:"type"`
	Function ollamaToolFunction `json:"function"`
}

type ollamaToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Tools    []ollamaTool    `json:"tools,omitempty"`
	Stream   bool            `json:"stream"`
}

type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
}

// Init creates a GameState with randomly assigned roles.
// nbWolves wolves are assigned; a Seer is added at 5+ players; a Witch at 6+.
func Init(gameID string, playerNames []string, nbWolves int) engine.GameState {
	indices := rand.Perm(len(playerNames))
	roles := make([]string, len(playerNames))

	for i := range nbWolves {
		if i >= len(indices) {
			break
		}
		roles[indices[i]] = "Loup Garou"
	}
	offset := nbWolves
	if len(playerNames) >= 5 && offset < len(indices) {
		roles[indices[offset]] = "Voyante"
		offset++
	}
	if len(playerNames) >= 6 && offset < len(indices) {
		roles[indices[offset]] = "Sorcière"
		offset++
	}

	players := make([]engine.PlayerState, len(playerNames))
	for i, name := range playerNames {
		r := roles[i]
		if r == "" {
			r = "Villageois"
		}
		faction := "Village"
		if r == "Loup Garou" {
			faction = "Loup"
		}
		ps := engine.PlayerState{
			ID:      name,
			Name:    name,
			Role:    r,
			IsAlive: true,
			Faction: faction,
		}
		if r == "Sorcière" {
			ps.HasHeal = true
			ps.HasPoison = true
		}
		players[i] = ps
	}

	return engine.GameState{
		ID:      gameID,
		Round:   1,
		Phase:   engine.PhaseNight,
		Players: players,
	}
}

// Run executes a full game using a local Llama model via Ollama for every decision.
// Set OLLAMA_HOST (default: http://localhost:11434) and OLLAMA_MODEL (default: llama3.1).
func Run(ctx context.Context, state engine.GameState) (engine.GameState, error) {
	var history []engine.Event

	for state.Result == "" {
		all := engine.ComputeRequired(state)

		var decisions []engine.Decision
		if len(all) > 0 {
			fmt.Printf("\n[Tour %d / %s]\n", state.Round, state.Phase)
			d, err := decide(ctx, state, all, history)
			if err != nil {
				return state, fmt.Errorf("tour %d/%s: %w", state.Round, state.Phase, err)
			}
			for _, dec := range d {
				fmt.Printf("  %s (%s) -> %s\n", dec.ActorID, dec.Kind, dec.TargetID)
			}
			decisions = d
		}

		result := engine.Step(state, decisions)
		for _, e := range result.Events {
			if e.Detail != "" {
				fmt.Printf("  [%s] %s (%s)\n", e.Kind, e.PlayerID, e.Detail)
			} else {
				fmt.Printf("  [%s] %s\n", e.Kind, e.PlayerID)
			}
		}
		history = append(history, result.Events...)
		state = result.State
	}

	return state, nil
}

// decide calls the local Llama model via Ollama and returns all decisions for the current phase.
// Mandatory decisions that the model skips fall back to a random candidate.
func decide(
	ctx context.Context,
	state engine.GameState,
	pending []engine.RequiredDecision,
	history []engine.Event,
) ([]engine.Decision, error) {
	req := ollamaRequest{
		Model:  ollamaModel(),
		Stream: false,
		Messages: []ollamaMessage{
			{Role: "system", Content: buildSystem(state)},
			{Role: "user", Content: buildUserMsg(state, pending, history)},
		},
		Tools: buildTools(pending),
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := ollamaHost() + "/api/chat"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var decisions []engine.Decision
	for _, tc := range ollamaResp.Message.ToolCalls {
		if d, ok := parseToolCall(tc, pending); ok {
			decisions = append(decisions, d)
		}
	}

	// Fallback: fill any mandatory decision the model omitted
	for _, req := range pending {
		if req.Optional || len(req.Candidates) == 0 {
			continue
		}
		made := false
		for _, d := range decisions {
			if d.Kind == req.Kind && d.ActorID == req.ActorID {
				made = true
				break
			}
		}
		if !made {
			decisions = append(decisions, engine.Decision{
				Kind:     req.Kind,
				ActorID:  req.ActorID,
				TargetID: req.Candidates[rand.Intn(len(req.Candidates))].ID,
			})
		}
	}

	return decisions, nil
}

func buildTools(pending []engine.RequiredDecision) []ollamaTool {
	tools := make([]ollamaTool, len(pending))
	for i, req := range pending {
		ids := make([]string, len(req.Candidates))
		for j, c := range req.Candidates {
			ids[j] = c.ID
		}
		tools[i] = ollamaTool{
			Type: "function",
			Function: ollamaToolFunction{
				Name:        toolID(req),
				Description: describeDecision(req),
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"target_id": map[string]any{
							"type": "string",
							"enum": ids,
						},
					},
					"required": []string{"target_id"},
				},
			},
		}
	}
	return tools
}

func parseToolCall(tc ollamaToolCall, pending []engine.RequiredDecision) (engine.Decision, bool) {
	for _, req := range pending {
		if toolID(req) != tc.Function.Name {
			continue
		}
		var inp struct {
			TargetID string `json:"target_id"`
		}
		if err := json.Unmarshal(tc.Function.Arguments, &inp); err != nil {
			return engine.Decision{}, false
		}
		return engine.Decision{
			Kind:     req.Kind,
			ActorID:  req.ActorID,
			TargetID: inp.TargetID,
		}, true
	}
	return engine.Decision{}, false
}

func buildSystem(state engine.GameState) string {
	var names []string
	for _, p := range state.Players {
		if p.IsAlive {
			names = append(names, p.Name)
		}
	}
	return fmt.Sprintf(
		"Tu es le maître du jeu Loup-Garou, partie %s. "+
			"Tu joues chaque rôle selon sa connaissance propre : les loups connaissent uniquement leurs alliés, "+
			"le village vote sans savoir qui sont les loups. "+
			"Joueurs en vie: %s. "+
			"Appelle TOUS les outils obligatoires pour avancer la partie.",
		state.ID, strings.Join(names, ", "))
}

func buildUserMsg(state engine.GameState, pending []engine.RequiredDecision, history []engine.Event) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Tour %d — Phase: %s\n", state.Round, state.Phase))

	if state.Victim != "" {
		sb.WriteString(fmt.Sprintf("Victime des loups cette nuit: %s\n", state.Victim))
	}

	if n := len(history); n > 0 {
		start := n - 5
		if start < 0 {
			start = 0
		}
		sb.WriteString("\nEvenements recents:\n")
		for _, e := range history[start:] {
			if e.Detail != "" {
				sb.WriteString(fmt.Sprintf("  [%s] %s (%s)\n", e.Kind, e.PlayerID, e.Detail))
			} else {
				sb.WriteString(fmt.Sprintf("  [%s] %s\n", e.Kind, e.PlayerID))
			}
		}
	}

	sb.WriteString("\nDecisions a prendre:\n")
	for _, req := range pending {
		names := make([]string, len(req.Candidates))
		for i, c := range req.Candidates {
			names[i] = c.Name
		}
		opt := ""
		if req.Optional {
			opt = " [optionnel]"
		}
		sb.WriteString(fmt.Sprintf("  %s par %s%s — cibles: %s\n",
			req.Kind, req.ActorID, opt, strings.Join(names, ", ")))
	}

	sb.WriteString("\nAppelle les outils pour chaque decision obligatoire.")
	return sb.String()
}

func describeDecision(req engine.RequiredDecision) string {
	switch req.Kind {
	case engine.DecisionWerewolfAttack:
		return "Loup-garou " + req.ActorID + " choisit sa victime cette nuit"
	case engine.DecisionSeerInvestigate:
		return "Voyante " + req.ActorID + " enquete sur un joueur"
	case engine.DecisionWitchSave:
		return "Sorciere " + req.ActorID + " peut sauver la victime des loups [optionnel]"
	case engine.DecisionWitchPoison:
		return "Sorciere " + req.ActorID + " peut empoisonner un joueur [optionnel]"
	case engine.DecisionVote:
		return "Joueur " + req.ActorID + " vote pour eliminer un suspect"
	default:
		return string(req.Kind) + " par " + req.ActorID
	}
}
