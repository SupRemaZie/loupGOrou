import { useEffect, useRef, useState, useCallback } from 'react';
import type { GameState, GameEvent, RequiredDecision, Decision, ServerMsg } from '../types';

export function useGameSocket(gameId: string | undefined) {
  const [state, setState] = useState<GameState | null>(null);
  const [events, setEvents] = useState<GameEvent[]>([]);
  const [pending, setPending] = useState<RequiredDecision[]>([]);
  const [connected, setConnected] = useState<string[]>([]);
  const [playerName, setPlayerName] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [wsReady, setWsReady] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!gameId) return;

    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;
    const ws = new WebSocket(`${protocol}://${host}/ws/${gameId}`);
    wsRef.current = ws;

    ws.onopen = () => setWsReady(true);

    ws.onmessage = (event) => {
      const msg: ServerMsg = JSON.parse(event.data);
      if (msg.type === 'update') {
        if (msg.state) setState(msg.state);
        if (msg.events && msg.events.length > 0) {
          setEvents((prev) => [...prev, ...msg.events!]);
        }
        setPending(msg.pending ?? []);
        setConnected(msg.connected ?? []);
      } else if (msg.type === 'joined') {
        setPlayerName(msg.name ?? null);
      } else if (msg.type === 'error') {
        setError(msg.error ?? 'Erreur inconnue');
      }
    };

    ws.onerror = () => setError('Connexion WebSocket échouée');
    ws.onclose = () => setWsReady(false);

    return () => ws.close();
  }, [gameId]);

  const join = useCallback((name: string) => {
    wsRef.current?.send(JSON.stringify({ type: 'join', name }));
  }, []);

  const decide = useCallback((decisions: Decision[]) => {
    wsRef.current?.send(JSON.stringify({ type: 'decide', decisions }));
  }, []);

  return { state, events, pending, connected, playerName, error, wsReady, join, decide };
}
