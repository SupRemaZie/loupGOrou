import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function Home() {
  const navigate = useNavigate();
  const [joinId, setJoinId] = useState('');

  function handleJoin(e: React.FormEvent) {
    e.preventDefault();
    const id = joinId.trim().toLowerCase();
    if (id) navigate(`/game/${id}`);
  }

  return (
    <div className="home-page">
      <div className="home-hero">
        <div className="home-moon">🌕</div>
        <h1 className="home-title">Loup-Garou</h1>
        <p className="home-subtitle">Le jeu de déduction sociale</p>
      </div>

      <div className="home-cards">
        <div className="home-card" onClick={() => navigate('/create')}>
          <div className="home-card-icon">🐺</div>
          <h2>Créer une partie</h2>
          <p>Configurez une partie multijoueur et invitez vos amis</p>
        </div>

        <div className="home-card" onClick={() => navigate('/ai')}>
          <div className="home-card-icon">🤖</div>
          <h2>Regarder l'IA jouer</h2>
          <p>Observez une partie se dérouler automatiquement</p>
        </div>
      </div>

      <div className="home-join">
        <h3>Rejoindre une partie existante</h3>
        <form onSubmit={handleJoin} className="home-join-form">
          <input
            type="text"
            placeholder="Code de la partie (ex: abc123)"
            value={joinId}
            onChange={(e) => setJoinId(e.target.value)}
            className="input"
            maxLength={10}
          />
          <button type="submit" className="btn btn-secondary" disabled={!joinId.trim()}>
            Rejoindre
          </button>
        </form>
      </div>
    </div>
  );
}
