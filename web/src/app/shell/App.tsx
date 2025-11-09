import { Outlet, Link, NavLink, useLocation } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { Button } from '@/components/Button';
import { useTheme } from '@/app/providers/ThemeProvider';
import styled from 'styled-components';
import { Toaster } from 'react-hot-toast';
import { useAuth } from '@/features/auth/AuthProvider';
import { useLogout } from '@/features/auth/hooks';
import { useNavigate } from 'react-router-dom';
import { Sun, Moon } from 'lucide-react';

const Header = styled.header`
  position: sticky;
  top: 0;
  z-index: 10;
  background: var(--bg-elevated);
  border-bottom: 1px solid var(--border);
  box-shadow: var(--elev-1);
`;

const HeaderInner = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 12px 20px;
`;

const Nav = styled.nav`
  display: flex;
  gap: 12px;
  a {
    position: relative;
    padding: 8px 12px;
    border-radius: var(--radius);
    color: var(--text);
    text-decoration: none;
    transition: background var(--spring-fast), box-shadow var(--spring-fast), color var(--spring-fast), transform var(--spring-fast);
  }
  a:hover { background: color-mix(in oklab, var(--primary) 6%, transparent); box-shadow: var(--elev-1); }
  a:focus-visible { outline: none; box-shadow: 0 0 0 2px color-mix(in oklab, var(--primary) 25%, transparent), var(--elev-1); }
  a.active { background: var(--bg-elevated); border: 1px solid var(--border); box-shadow: var(--elev-1); font-weight: 600; }

  /* Mobile */
  @media (max-width: 768px) {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: var(--bg-elevated);
    border-bottom: 1px solid var(--border);
    box-shadow: var(--elev-2);
    padding: 8px 12px;
    display: none;
    flex-direction: column;
    gap: 6px;
    &[data-open="true"] { display: flex; }
  }
`;

const MenuToggle = styled.button`
  display: none;
  @media (max-width: 768px) {
    display: inline-flex;
  }
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: var(--bg-elevated);
  color: var(--text);
  cursor: pointer;
  transition: background var(--spring-fast), box-shadow var(--spring-fast), transform var(--spring-fast);
  &:hover { box-shadow: var(--elev-1); }
  &:active { transform: scale(0.98); }
  &:focus-visible { outline: none; box-shadow: 0 0 0 2px color-mix(in oklab, var(--primary) 25%, transparent), var(--elev-1); }
`;

export function App() {
  const [navOpen, setNavOpen] = useState(false);
  const { toggle, mode } = useTheme();
  const { isAuthed } = useAuth();
  const doLogout = useLogout();
  const navigate = useNavigate();
  const location = useLocation();
  // Auto-close mobile nav on route change
  useEffect(() => { setNavOpen(false); }, [location.pathname]);
  return (
    <div>
      <a href="#main" className="skip-link">Skip to content</a>
      <Header>
        <HeaderInner>
          <Link to="/" style={{ fontWeight: 700 }}>Asiko</Link>
          <Nav aria-label="Primary Navigation" id="primary-nav" data-open={navOpen}>
            <NavLink to="/">Book</NavLink>
            {isAuthed ? (
              <>
                <NavLink to="/dashboard">Dashboard</NavLink>
                <NavLink to="/appointments">Appointments</NavLink>
                <NavLink to="/bookings">My Bookings</NavLink>
                <NavLink to="/analytics">Analytics</NavLink>
                <NavLink to="/ban-list">Ban List</NavLink>
                <NavLink to="/notifications">Notifications</NavLink>
                <a href="#" onClick={(e) => { e.preventDefault(); doLogout.mutate(undefined, { onSettled: () => navigate('/') }); }}>Logout</a>
              </>
            ) : (
              <NavLink to="/login">Login</NavLink>
            )}
          </Nav>
          <div style={{ display: 'flex', gap: 8 }}>
            <Button aria-label="Toggle theme" onClick={toggle} title="Toggle theme">
              {mode === 'light' ? <Sun size={16} /> : <Moon size={16} />}
              <span className="hide-sm">Theme</span>
            </Button>
            <MenuToggle aria-label="Toggle navigation menu" aria-controls="primary-nav" aria-expanded={navOpen} onClick={() => setNavOpen(v => !v)}>
              {/* Simple hamburger glyph */}
              <span aria-hidden="true" style={{ display: 'inline-block', width: 18, height: 2, background: 'currentColor', boxShadow: '0 6px currentColor, 0 -6px currentColor' }} />
            </MenuToggle>
          </div>
        </HeaderInner>
      </Header>
      <main id="main" className="container" role="main">
        <Outlet />
      </main>
      <Toaster position="top-center" />
    </div>
  );
}
