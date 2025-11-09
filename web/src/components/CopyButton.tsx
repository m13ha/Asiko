import styled from 'styled-components';
import { useState } from 'react';

const Small = styled.button`
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  font-size: 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--bg-elevated);
  color: var(--text);
  cursor: pointer;
  transition: background var(--spring-fast), box-shadow var(--spring-fast), transform var(--spring-fast), border-color var(--spring-fast);
  &:hover { box-shadow: var(--elev-1); border-color: color-mix(in oklab, var(--primary) 20%, var(--border)); }
  &:active { transform: scale(0.98); }
  &:focus-visible { outline: none; box-shadow: 0 0 0 2px color-mix(in oklab, var(--primary) 25%, transparent), var(--elev-1); }
`;

async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch {
    try {
      const ta = document.createElement('textarea');
      ta.value = text;
      document.body.appendChild(ta);
      ta.select();
      document.execCommand('copy');
      document.body.removeChild(ta);
      return true;
    } catch {
      return false;
    }
  }
}

export function CopyButton({ value, label = 'Copy', copiedLabel = 'Copied', ariaLabel }: { value: string; label?: string; copiedLabel?: string; ariaLabel?: string; }) {
  const [copied, setCopied] = useState(false);
  return (
    <Small
      type="button"
      aria-label={ariaLabel || 'Copy to clipboard'}
      onClick={async () => {
        if (!value) return;
        const ok = await copyToClipboard(value);
        if (ok) {
          setCopied(true);
          window.setTimeout(() => setCopied(false), 900);
        }
      }}
    >
      {copied ? copiedLabel : label}
    </Small>
  );
}

