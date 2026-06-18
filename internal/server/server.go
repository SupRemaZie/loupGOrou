package server

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/SupRemaZie/loupGOrou/internal/engine"
	"github.com/SupRemaZie/loupGOrou/internal/orchestrator"
)

// ---- WebSocket message types ----

type ClientMsg struct {
	Type      string            `json:"type"`      // "join", "decide"
	Name      string            `json:"name,omitempty"`
	Decisions []engine.Decision `json:"decisions,omitempty"`
}

type ServerMsg struct {
	Type      string                    `json:"type"`
	State     *engine.GameState         `json:"state,omitempty"`
	Events    []engine.Event            `json:"events,omitempty"`
	Pending   []engine.RequiredDecision `json:"pending,omitempty"`
	Connected []string                  `json:"connected,omitempty"`
	Name      string                    `json:"name,omitempty"`
	Error     string                    `json:"error,omitempty"`
}

// ---- Client ----

type Client struct {
	writeMu sync.Mutex
	conn    *websocket.Conn
	player  string // empty = spectator
}

func (c *Client) write(msg ServerMsg) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	c.conn.WriteJSON(msg) //nolint
}

// ---- Session ----

type Session struct {
	ID      string
	mu      sync.Mutex
	State   engine.GameState
	History []engine.Event
	IsAI    bool
	clients []*Client
	pending []engine.Decision
}

func (s *Session) addClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients = append(s.clients, c)
}

func (s *Session) removeClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, cl := range s.clients {
		if cl == c {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			return
		}
	}
}

// broadcastLocked broadcasts to all clients. Must be called while holding s.mu.
func (s *Session) broadcastLocked(events []engine.Event) {
	if events != nil {
		s.History = append(s.History, events...)
	}
	allPending := engine.ComputeRequired(s.State)
	connected := s.connectedPlayersLocked()

	for _, c := range s.clients {
		msg := ServerMsg{
			Type:      "update",
			State:     s.maskedStateLocked(c.player),
			Events:    events,
			Pending:   s.filterPendingLocked(allPending, c.player),
			Connected: connected,
		}
		go c.write(msg)
	}
}

func (s *Session) maskedStateLocked(playerName string) *engine.GameState {
	cp := s.State
	players := make([]engine.PlayerState, len(cp.Players))
	copy(players, cp.Players)

	// AI spectator games and finished games reveal all roles
	if s.IsAI || cp.Result != "" {
		cp.Players = players
		return &cp
	}

	myFaction := ""
	for _, p := range players {
		if p.Name == playerName {
			myFaction = p.Faction
			break
		}
	}

	for i, p := range players {
		if p.Name == playerName {
			continue
		}
		if myFaction == "Loup" && p.Faction == "Loup" {
			continue // wolves see each other
		}
		players[i].Role = "?"
		players[i].Faction = "?"
	}

	cp.Players = players
	return &cp
}

func (s *Session) filterPendingLocked(pending []engine.RequiredDecision, playerName string) []engine.RequiredDecision {
	if playerName == "" {
		return nil
	}
	var filtered []engine.RequiredDecision
	for _, p := range pending {
		if p.ActorID == playerName {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func (s *Session) connectedPlayersLocked() []string {
	var names []string
	for _, c := range s.clients {
		if c.player != "" {
			names = append(names, c.player)
		}
	}
	return names
}

func (s *Session) handleDecide(c *Client, decisions []engine.Decision) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State.Result != "" {
		return
	}

	// Ensure this player has a role in the current phase
	allRequired := engine.ComputeRequired(s.State)
	hasRole := false
	for _, req := range allRequired {
		if req.ActorID == c.player {
			hasRole = true
			break
		}
	}
	if !hasRole {
		return
	}

	for _, d := range decisions {
		if d.ActorID != c.player {
			continue
		}
		replaced := false
		for i, existing := range s.pending {
			if existing.Kind == d.Kind && existing.ActorID == d.ActorID {
				s.pending[i] = d
				replaced = true
				break
			}
		}
		if !replaced {
			s.pending = append(s.pending, d)
		}
	}

	s.tryAdvanceLocked()
}

func (s *Session) tryAdvanceLocked() {
	required := engine.ComputeRequired(s.State)
	for _, req := range required {
		if req.Optional {
			continue
		}
		found := false
		for _, d := range s.pending {
			if d.Kind == req.Kind && d.ActorID == req.ActorID {
				found = true
				break
			}
		}
		if !found {
			return
		}
	}

	result := engine.Step(s.State, s.pending)
	s.State = result.State
	s.pending = nil
	s.broadcastLocked(result.Events)
}

// ---- Hub ----

type Hub struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

func NewHub() *Hub {
	return &Hub{sessions: make(map[string]*Session)}
}

func (h *Hub) createSession(players []string, wolves int, isAI bool) *Session {
	id := randomID()
	state := orchestrator.Init(id, players, wolves)
	s := &Session{ID: id, State: state, IsAI: isAI}
	h.mu.Lock()
	h.sessions[id] = s
	h.mu.Unlock()
	return s
}

func (h *Hub) get(id string) (*Session, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s, ok := h.sessions[id]
	return s, ok
}

// ---- HTTP handlers ----

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Players []string `json:"players"`
		Wolves  int      `json:"wolves"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Players) < 2 || req.Wolves < 1 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	s := h.createSession(req.Players, req.Wolves, false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": s.ID}) //nolint
}

func (h *Hub) handleCreateAIGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Players []string `json:"players"`
		Wolves  int      `json:"wolves"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Players) < 2 || req.Wolves < 1 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	s := h.createSession(req.Players, req.Wolves, true)
	go h.runAIGame(s)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": s.ID}) //nolint
}

func (h *Hub) handleWS(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/ws/"), "/")
	gameID := parts[0]
	s, ok := h.get(gameID)
	if !ok {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws upgrade error:", err)
		return
	}
	defer conn.Close()

	c := &Client{conn: conn}
	s.addClient(c)
	defer s.removeClient(c)

	// Send current state to the newly connected spectator
	s.mu.Lock()
	initMsg := ServerMsg{
		Type:      "update",
		State:     s.maskedStateLocked(""),
		Events:    s.History,
		Pending:   nil,
		Connected: s.connectedPlayersLocked(),
	}
	s.mu.Unlock()
	c.write(initMsg)

	// Read loop
	for {
		var msg ClientMsg
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		switch msg.Type {
		case "join":
			s.mu.Lock()
			valid := false
			for _, p := range s.State.Players {
				if p.Name == msg.Name {
					valid = true
					break
				}
			}
			taken := false
			if valid {
				for _, cl := range s.clients {
					if cl != c && cl.player == msg.Name {
						taken = true
						break
					}
				}
			}
			if valid && !taken {
				c.player = msg.Name
				allPending := engine.ComputeRequired(s.State)
				joined := ServerMsg{Type: "joined", Name: msg.Name}
				update := ServerMsg{
					Type:      "update",
					State:     s.maskedStateLocked(msg.Name),
					Events:    s.History,
					Pending:   s.filterPendingLocked(allPending, msg.Name),
					Connected: s.connectedPlayersLocked(),
				}
				s.mu.Unlock()
				c.write(joined)
				c.write(update)
			} else {
				s.mu.Unlock()
				c.write(ServerMsg{Type: "error", Error: "nom invalide ou déjà pris"})
			}

		case "decide":
			if !s.IsAI && c.player != "" {
				s.handleDecide(c, msg.Decisions)
			}
		}
	}
}

func (h *Hub) runAIGame(s *Session) {
	time.Sleep(1 * time.Second)

	ctx := context.Background()
	_ = ctx

	for {
		s.mu.Lock()
		if s.State.Result != "" {
			s.broadcastLocked(nil)
			s.mu.Unlock()
			break
		}

		pending := engine.ComputeRequired(s.State)
		if len(pending) == 0 {
			s.mu.Unlock()
			break
		}

		var decisions []engine.Decision
		for _, req := range pending {
			if req.Optional {
				if rand.Float32() < 0.4 && len(req.Candidates) > 0 {
					decisions = append(decisions, engine.Decision{
						Kind:     req.Kind,
						ActorID:  req.ActorID,
						TargetID: req.Candidates[rand.Intn(len(req.Candidates))].ID,
					})
				}
				continue
			}
			if len(req.Candidates) == 0 {
				continue
			}
			decisions = append(decisions, engine.Decision{
				Kind:     req.Kind,
				ActorID:  req.ActorID,
				TargetID: req.Candidates[rand.Intn(len(req.Candidates))].ID,
			})
		}

		result := engine.Step(s.State, decisions)
		s.State = result.State
		s.broadcastLocked(result.Events)
		s.mu.Unlock()

		time.Sleep(2 * time.Second)
	}
}

// ---- Helpers ----

func randomID() string {
	const chars = "abcdefghjkmnpqrstuvwxyz23456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// New creates an HTTP server. staticDir is the path to the built React app.
func New(addr, staticDir string) *http.Server {
	hub := NewHub()
	mux := http.NewServeMux()

	mux.HandleFunc("/api/games", hub.handleCreateGame)
	mux.HandleFunc("/api/ai/games", hub.handleCreateAIGame)
	mux.HandleFunc("/ws/", hub.handleWS)

	if staticDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	}

	return &http.Server{
		Addr:    addr,
		Handler: corsMiddleware(mux),
	}
}
