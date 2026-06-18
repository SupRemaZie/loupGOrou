import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const DEFAULT_NAMES = ['Alice', 'Bob', 'Charlie', 'Diana', 'Eve', 'Frank'];

export default function CreateGame() {
  const navigate = useNavigate();
  const [players, setPlayers] = useState<string[]>(['', '', '', '']);
  const [wolves, setWolves] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  function updatePlayer(i: number, val: string) {
    setPlayers((prev) => prev.map((p, idx) => (idx === i ? val : p)));
  }

  function addPlayer() {
    if (players.length < 10) setPlayers((prev) => [...prev, '']);
  }

  function removePlayer(i: number) {
    if (players.length > 2) setPlayers((prev) => prev.filter((_, idx) => idx !== i));
  }

  function fillDefaults() {
    setPlayers(DEFAULT_NAMES.slice(0, Math.max(players.length, 4)));
  }

  const validPlayers = players.map((p) => p.trim()).filter(Boolean);
  const canCreate = validPlayers.length >= 2 && wolves >= 1 && wolves < validPlayers.length;

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!canCreate) return;
    setLoading(true);
    setError('');
    try {
      const res = await fetch('/api/games', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ players: validPlayers, wolves }),
      });
      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();
      navigate(`/game/${data.id}`);
    } catch (err) {
      setError(String(err));
      setLoading(false);
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <button className="btn btn-ghost" onClick={() => navigate('/')}>← Retour</button>
        <h1>🐺 Nouvelle partie</h1>
      </div>

      <form onSubmit={handleCreate} className="form-card">
        <div className="form-section">
          <div className="form-section-header">
            <h2>Joueurs ({validPlayers.length})</h2>
            <button type="button" className="btn btn-ghost btn-sm" onClick={fillDefaults}>
              Noms par défaut
            </button>
          </div>
          <div className="player-inputs">
            {players.map((name, i) => (
              <div key={i} className="player-input-row">
                <input
                  type="text"
                  placeholder={`Joueur ${i + 1}`}
                  value={name}
                  onChange={(e) => updatePlayer(i, e.target.value)}
                  className="input"
                  maxLength={20}
                />
                {players.length > 2 && (
                  <button
                    type="button"
                    className="btn btn-danger btn-sm"
                    onClick={() => removePlayer(i)}
                  >
                    ✕
                  </button>
                )}
              </div>
            ))}
          </div>
          {players.length < 10 && (
            <button type="button" className="btn btn-ghost btn-sm" onClick={addPlayer}>
              + Ajouter un joueur
            </button>
          )}
        </div>

        <div className="form-section">
          <h2>Loups-garous</h2>
          <div className="wolf-selector">
            {[1, 2, 3].map((n) => (
              <button
                key={n}
                type="button"
                className={`btn ${wolves === n ? 'btn-primary' : 'btn-ghost'}`}
                onClick={() => setWolves(n)}
                disabled={n >= validPlayers.length}
              >
                {'🐺'.repeat(n)}
              </button>
            ))}
          </div>
          <p className="form-hint">
            {validPlayers.length >= 5 && '🔮 Voyante incluse. '}
            {validPlayers.length >= 6 && '🧙 Sorcière incluse.'}
          </p>
        </div>

        {error && <p className="form-error">{error}</p>}

        <button
          type="submit"
          className="btn btn-primary btn-lg"
          disabled={!canCreate || loading}
        >
          {loading ? 'Création…' : '🎮 Créer la partie'}
        </button>
      </form>
    </div>
  );
}
