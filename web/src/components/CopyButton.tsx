import { useState } from 'react';

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

export function CopyButton({ 
  value, 
  label = 'Copy', 
  copiedLabel = 'Copied', 
  ariaLabel,
  className = ''
}: { 
  value: string; 
  label?: string; 
  copiedLabel?: string; 
  ariaLabel?: string;
  className?: string;
}) {
  const [copied, setCopied] = useState(false);
  
  return (
    <button
      type="button"
      aria-label={ariaLabel || 'Copy to clipboard'}
      className={`inline-flex items-center gap-1.5 px-2 py-1 text-xs border border-gray-300 rounded  text-gray-700 hover:bg-gray-50 hover:border-gray-400 active:scale-95 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-1 transition-all ${className}`}
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
    </button>
  );
}