import { useEffect, useRef } from 'react';
import type { GameEvent } from '../types';

const EVENT_ICON: Record<string, string> = {
  killed: '🔪',
  eliminated: '🗳️',
  saved: '💚',
  revealed: '👁️',
  no_consensus: '⚖️',
};

const EVENT_LABEL: Record<string, (e: GameEvent) => string> = {
  killed: (e) => `${e.player_id} a été tué cette nuit`,
  eliminated: (e) =>
    `${e.player_id} a été éliminé par le village${e.detail ? ` (c'était un ${e.detail})` : ''}`,
  saved: (e) => `${e.player_id} a été sauvé par la sorcière`,
  revealed: (e) => `La voyante a enquêté sur ${e.player_id} : c'est un ${e.detail}`,
  no_consensus: () => 'Pas de consensus — personne n\'est éliminé',
};

interface EventLogProps {
  events: GameEvent[];
}

export default function EventLog({ events }: EventLogProps) {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [events.length]);

  return (
    <div className="event-log">
      <h3 className="event-log-title">📜 Journal</h3>
      <div className="event-log-list">
        {events.length === 0 && (
          <p className="event-log-empty">La partie commence…</p>
        )}
        {events.map((e, i) => (
          <div key={i} className={`event-item event-${e.kind}`}>
            <span className="event-icon">{EVENT_ICON[e.kind] ?? '•'}</span>
            <span>{EVENT_LABEL[e.kind]?.(e) ?? `${e.kind}: ${e.player_id}`}</span>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>
    </div>
  );
}
