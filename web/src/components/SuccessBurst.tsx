import styled, { keyframes } from 'styled-components';

const pop = keyframes`
  0% { transform: scale(0.6); opacity: 0; }
  60% { transform: scale(1.06); opacity: 1; }
  100% { transform: scale(1); opacity: 0; }
`;

const Particle = styled.span`
  position: absolute;
  width: 6px;
  height: 6px;
  background: var(--accent);
  border-radius: 50%;
  opacity: 0.9;
`;

const BurstWrap = styled.div`
  position: fixed;
  inset: 0;
  display: grid;
  place-items: center;
  pointer-events: none;
  z-index: 50;
`;

const Tick = styled.div`
  position: relative;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: color-mix(in oklab, var(--primary) 88%, white);
  color: var(--primary-contrast);
  box-shadow: var(--elev-3);
  display: grid;
  place-items: center;
  animation: ${pop} var(--spring-slow) forwards;
  @media (prefers-reduced-motion: reduce) {
    animation: none;
    opacity: 1;
  }
  &:after {
    content: '✓';
    font-size: 28px;
    font-weight: 700;
  }
`;

export function SuccessBurst({ show }: { show: boolean }) {
  if (!show) return null;
  return (
    <BurstWrap aria-hidden>
      <Tick>
        <Particle style={{ top: -8, left: 28 }} />
        <Particle style={{ bottom: -6, left: 10 }} />
        <Particle style={{ right: -6, top: 14 }} />
      </Tick>
    </BurstWrap>
  );
}

