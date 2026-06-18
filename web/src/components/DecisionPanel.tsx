import { useState } from 'react';
import type { RequiredDecision, Decision, DecisionKind } from '../types';
import PlayerCard from './PlayerCard';

const DECISION_LABEL: Record<DecisionKind, string> = {
  werewolf_attack: '🐺 Choisir une victime',
  seer_investigate: '🔮 Enquêter sur',
  witch_save: '💚 Sauver la victime ?',
  witch_poison: '☠️ Empoisonner',
  vote: '🗳️ Voter pour éliminer',
};

interface DecisionPanelProps {
  pending: RequiredDecision[];
  onDecide: (decisions: Decision[]) => void;
}

export default function DecisionPanel({ pending, onDecide }: DecisionPanelProps) {
  const [optionalChoices, setOptionalChoices] = useState<Record<string, string>>({});

  if (pending.length === 0) return null;

  const mandatory = pending.filter((d) => !d.optional);
  const optional = pending.filter((d) => d.optional);
  const hasOptionalOnly = mandatory.length === 0 && optional.length > 0;

  function handleMandatoryClick(dec: RequiredDecision, targetId: string) {
    onDecide([{ kind: dec.kind, actor_id: dec.actor_id, target_id: targetId }]);
  }

  function handleOptionalToggle(key: string, targetId: string) {
    setOptionalChoices((prev) => ({
      ...prev,
      [key]: prev[key] === targetId ? '' : targetId,
    }));
  }

  function handleConfirmOptional() {
    const decisions: Decision[] = optional
      .map((dec) => {
        const key = dec.kind + dec.actor_id;
        const targetId = optionalChoices[key];
        if (!targetId) return null;
        return { kind: dec.kind, actor_id: dec.actor_id, target_id: targetId };
      })
      .filter(Boolean) as Decision[];
    onDecide(decisions);
    setOptionalChoices({});
  }

  return (
    <div className="decision-panel">
      <h3 className="decision-panel-title">⚡ À toi de jouer</h3>

      {mandatory.map((dec) => (
        <div key={dec.kind + dec.actor_id} className="decision-block mandatory">
          <p className="decision-label">{DECISION_LABEL[dec.kind]}</p>
          <div className="decision-candidates">
            {dec.candidates.map((p) => (
              <PlayerCard
                key={p.id}
                player={p}
                isTarget
                onClick={() => handleMandatoryClick(dec, p.id)}
              />
            ))}
          </div>
        </div>
      ))}

      {optional.map((dec) => {
        const key = dec.kind + dec.actor_id;
        const selected = optionalChoices[key];
        return (
          <div key={key} className="decision-block optional">
            <p className="decision-label">
              {DECISION_LABEL[dec.kind]} <span className="badge-optional">optionnel</span>
            </p>
            <div className="decision-candidates">
              {dec.candidates.map((p) => (
                <PlayerCard
                  key={p.id}
                  player={p}
                  isTarget
                  isSelected={selected === p.id}
                  onClick={() => handleOptionalToggle(key, p.id)}
                />
              ))}
            </div>
          </div>
        );
      })}

      {hasOptionalOnly && (
        <button className="btn btn-primary" onClick={handleConfirmOptional}>
          Terminer mon tour
        </button>
      )}
    </div>
  );
}
