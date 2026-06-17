// Shown above the RequestEditor when a /parse result is low-confidence
// (needsConfirmation=true). Surfaces source + confidence + notes so the
// user knows whether to trust the parse before sending.

import type { ParseSource } from "../../shared/types";

interface Props {
  source: ParseSource;
  confidence: number;
  notes?: string;
  onConfirm: () => void;
  onReparse: () => void;
}

export function ConfirmBar({ source, confidence, notes, onConfirm, onReparse }: Props) {
  const pct = Math.round(confidence * 100);
  return (
    <div className="confirm-bar">
      <div className="confirm-header">
        <span className={"confirm-pill confirm-" + qualityBucket(confidence)}>
          {pct}% confidence · {source}
        </span>
        <div className="confirm-actions">
          <button type="button" className="add" onClick={onReparse}>
            Re-parse
          </button>
          <button type="button" onClick={onConfirm}>
            Looks right
          </button>
        </div>
      </div>
      {notes && <p className="muted confirm-notes">{notes}</p>}
      <p className="muted confirm-help">
        Review the fields below; nothing is sent until you press <strong>Send</strong>.
      </p>
    </div>
  );
}

function qualityBucket(c: number): "low" | "mid" | "high" {
  if (c >= 0.9) return "high";
  if (c >= 0.75) return "mid";
  return "low";
}
