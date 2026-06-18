import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useGameSocket } from '../hooks/useGameSocket';
import PlayerCard from '../components/PlayerCard';
import EventLog from '../components/EventLog';

const DEFAULT_AI_PLAYERS = ['Alice', 'Bob', 'Charlie', 'Diana', 'Eve', 'Frank'];

const PHASE_INFO: Record<string, { label: string; icon: string }> = {
  night: { label: 'Nuit', icon: '🌙' },
  night_witch: { label: 'Nuit — Sorcière', icon: '🧙' },
  day: { label: 'Jour', icon: '☀️' },
};

function WatchGame({ gameId }: { gameId: string }) {
  const { state, events } = useGameSocket(gameId);

  if (!state) {
    return (
      <div className="center">
        <div className="spinner" />
        <p>Chargement…</p>
      </div>
    );
  }

  const phase = PHASE_INFO[state.phase] ?? PHASE_INFO.night;

  return (
    <div className="ai-watch">
      <div className="ai-watch-header">
        <span className="phase-badge">
          {phase.icon} {phase.label} — Tour {state.round}
        </span>
        {state.result && (
          <span className="result-badge">
            {state.result === 'Loups' ? '🐺 Loups gagnent !' : '🏘️ Village gagne !'}
          </span>
        )}
      </div>

      <div className="players-grid">
        {state.players.map((p) => (
          <PlayerCard key={p.id} player={p} />
        ))}
      </div>

      <EventLog events={events} />
    </div>
  );
}

export default function AISpectator() {
  const navigate = useNavigate();
  const [players, setPlayers] = useState(DEFAULT_AI_PLAYERS.join(', '));
  const [wolves, setWolves] = useState(2);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [gameId, setGameId] = useState<string | null>(null);

  async function handleStart(e: React.FormEvent) {
    e.preventDefault();
    const names = players
      .split(/[,\n]+/)
      .map((s) => s.trim())
      .filter(Boolean);
    if (names.length < 2 || wolves < 1 || wolves >= names.length) {
      setError('Minimum 2 joueurs, 1 loup, et moins de loups que de joueurs.');
      return;
    }
    setLoading(true);
    setError('');
    try {
      const res = await fetch('/api/ai/games', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ players: names, wolves }),
      });
      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();
      setGameId(data.id);
    } catch (err) {
      setError(String(err));
      setLoading(false);
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <button className="btn btn-ghost" onClick={() => navigate('/')}>← Retour</button>
        <h1>🤖 Spectateur IA</h1>
      </div>

      {!gameId ? (
        <form onSubmit={handleStart} className="form-card">
          <div className="form-section">
            <h2>Joueurs (séparés par des virgules)</h2>
            <textarea
              className="input textarea"
              rows={3}
              value={players}
              onChange={(e) => setPlayers(e.target.value)}
            />
          </div>

          <div className="form-section">
            <h2>Nombre de loups</h2>
            <div className="wolf-selector">
              {[1, 2, 3].map((n) => (
                <button
                  key={n}
                  type="button"
                  className={`btn ${wolves === n ? 'btn-primary' : 'btn-ghost'}`}
                  onClick={() => setWolves(n)}
                >
                  {'🐺'.repeat(n)}
                </button>
              ))}
            </div>
          </div>

          {error && <p className="form-error">{error}</p>}

          <button type="submit" className="btn btn-primary btn-lg" disabled={loading}>
            {loading ? 'Démarrage…' : '▶ Lancer la partie'}
          </button>
        </form>
      ) : (
        <WatchGame gameId={gameId} />
      )}
    </div>
  );
}
