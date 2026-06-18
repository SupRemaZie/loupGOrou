import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useGameSocket } from '../hooks/useGameSocket';
import PlayerCard from '../components/PlayerCard';
import EventLog from '../components/EventLog';
import DecisionPanel from '../components/DecisionPanel';
import type { Decision } from '../types';

const PHASE_INFO: Record<string, { label: string; icon: string; className: string }> = {
  night: { label: 'Nuit', icon: '🌙', className: 'phase-night' },
  night_witch: { label: 'Nuit — Sorcière', icon: '🧙', className: 'phase-night' },
  day: { label: 'Jour', icon: '☀️', className: 'phase-day' },
};

export default function GamePage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { state, events, pending, connected, playerName, error, wsReady, join, decide } =
    useGameSocket(id);

  const [nameInput, setNameInput] = useState('');
  const [joinError, setJoinError] = useState('');

  function handleJoin(e: React.FormEvent) {
    e.preventDefault();
    const name = nameInput.trim();
    if (!name) return;
    join(name);
    setJoinError('');
  }

  function handleDecide(decisions: Decision[]) {
    decide(decisions);
  }

  if (!wsReady && !state) {
    return (
      <div className="page center">
        <div className="spinner" />
        <p>Connexion en cours…</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="page center">
        <p className="text-danger">⚠️ {error}</p>
        <button className="btn btn-ghost" onClick={() => navigate('/')}>Retour</button>
      </div>
    );
  }

  // Game over screen
  if (state?.result) {
    const isWolf = state.result === 'Loups';
    return (
      <div className="page">
        <div className={`game-over ${isWolf ? 'wolf-wins' : 'village-wins'}`}>
          <div className="game-over-icon">{isWolf ? '🐺' : '🏘️'}</div>
          <h1>{isWolf ? 'Les Loups ont gagné !' : 'Le Village a gagné !'}</h1>
          <div className="players-grid">
            {state.players.map((p) => (
              <PlayerCard key={p.id} player={p} isMe={p.name === playerName} />
            ))}
          </div>
          <EventLog events={events} />
          <button className="btn btn-primary" onClick={() => navigate('/')}>
            Nouvelle partie
          </button>
        </div>
      </div>
    );
  }

  // Join screen (no player name yet)
  if (!playerName) {
    const allPlayers = state?.players ?? [];
    const available = allPlayers
      .filter((p) => p.is_alive && !connected.includes(p.name))
      .map((p) => p.name);

    return (
      <div className="page center">
        <div className="join-card">
          <h1>🐺 Partie {id}</h1>
          {state && (
            <p className="text-muted">
              {connected.length} / {state.players.length} joueurs connectés
            </p>
          )}
          <form onSubmit={handleJoin} className="join-form">
            {available.length > 0 ? (
              <>
                <p>Choisissez votre personnage :</p>
                <div className="name-buttons">
                  {available.map((name) => (
                    <button
                      key={name}
                      type="button"
                      className={`btn ${nameInput === name ? 'btn-primary' : 'btn-ghost'}`}
                      onClick={() => setNameInput(name)}
                    >
                      {name}
                    </button>
                  ))}
                </div>
                <p className="text-muted">ou saisissez un nom :</p>
              </>
            ) : (
              <p>Saisissez votre nom :</p>
            )}
            <div className="join-input-row">
              <input
                type="text"
                placeholder="Votre nom"
                value={nameInput}
                onChange={(e) => setNameInput(e.target.value)}
                className="input"
                maxLength={20}
              />
              <button type="submit" className="btn btn-primary" disabled={!nameInput.trim()}>
                Rejoindre
              </button>
            </div>
            {joinError && <p className="form-error">{joinError}</p>}
          </form>
          <button className="btn btn-ghost btn-sm" onClick={() => navigate('/')}>
            ← Retour
          </button>
        </div>
      </div>
    );
  }

  if (!state) return null;

  const phase = PHASE_INFO[state.phase] ?? PHASE_INFO.night;
  const myPlayer = state.players.find((p) => p.name === playerName);
  const waitingFor = state.players
    .filter((p) => p.is_alive && !connected.includes(p.name))
    .map((p) => p.name);

  return (
    <div className={`game-page ${phase.className}`}>
      <header className="game-header">
        <div className="game-header-phase">
          {phase.icon} {phase.label} — Tour {state.round}
        </div>
        <div className="game-header-id">Partie : {id}</div>
        <div className="game-header-role">
          {myPlayer && (
            <>
              Vous êtes : <strong>{myPlayer.role}</strong>
            </>
          )}
        </div>
      </header>

      <div className="game-layout">
        <main className="game-main">
          <div className="players-grid">
            {state.players.map((p) => (
              <PlayerCard key={p.id} player={p} isMe={p.name === playerName} />
            ))}
          </div>

          {pending.length > 0 && (
            <DecisionPanel pending={pending} onDecide={handleDecide} />
          )}

          {pending.length === 0 && state.result === '' && (
            <div className="waiting-message">
              <div className="spinner" />
              <p>
                {waitingFor.length > 0
                  ? `En attente de : ${waitingFor.join(', ')}…`
                  : 'En attente des autres joueurs…'}
              </p>
            </div>
          )}
        </main>

        <aside className="game-sidebar">
          <div className="connected-list">
            <h4>Connectés ({connected.length})</h4>
            {connected.map((name) => (
              <span key={name} className={`connected-pill ${name === playerName ? 'me' : ''}`}>
                {name}
              </span>
            ))}
          </div>
          <EventLog events={events} />
        </aside>
      </div>
    </div>
  );
}
