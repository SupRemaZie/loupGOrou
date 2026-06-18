import type { PlayerState } from '../types';

const ROLE_EMOJI: Record<string, string> = {
  'Loup Garou': '🐺',
  Voyante: '🔮',
  'Sorcière': '🧙',
  Villageois: '👤',
  '?': '❓',
};

const FACTION_COLOR: Record<string, string> = {
  Loup: 'var(--wolf)',
  Village: 'var(--village)',
  '?': 'var(--text-muted)',
};

interface PlayerCardProps {
  player: PlayerState;
  isMe?: boolean;
  isSelected?: boolean;
  isTarget?: boolean;
  onClick?: () => void;
}

export default function PlayerCard({ player, isMe, isSelected, isTarget, onClick }: PlayerCardProps) {
  const emoji = ROLE_EMOJI[player.role] ?? '❓';
  const factionColor = FACTION_COLOR[player.faction] ?? 'var(--text-muted)';
  const isClickable = !!onClick && player.is_alive;

  return (
    <div
      className={[
        'player-card',
        !player.is_alive && 'dead',
        isMe && 'is-me',
        isSelected && 'selected',
        isTarget && 'is-target',
        isClickable && 'clickable',
      ]
        .filter(Boolean)
        .join(' ')}
      onClick={isClickable ? onClick : undefined}
      style={{ borderColor: isTarget ? 'var(--accent)' : undefined }}
    >
      <div className="player-card-emoji">{player.is_alive ? emoji : '💀'}</div>
      <div className="player-card-name">
        {player.name}
        {isMe && <span className="badge-me">Moi</span>}
      </div>
      {player.role !== '?' && (
        <div className="player-card-role" style={{ color: factionColor }}>
          {player.role}
        </div>
      )}
      {!player.is_alive && <div className="player-card-dead-overlay">Éliminé</div>}
    </div>
  );
}
