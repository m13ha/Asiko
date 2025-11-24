import { useState } from 'react';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useLogin } from '@/features/auth/hooks';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import * as API from '@appointment-master/api-client';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { HandUnderline } from '@/components/HandUnderline';
import { Spinner } from '@/components/Spinner';

export function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [info, setInfo] = useState<string | null>(null);
  const login = useLogin();
  const navigate = useNavigate();
  const loc = useLocation() as any;
  const from = loc.state?.from?.pathname || '/dashboard';

  return (
    <div style={{ maxWidth: 420, margin: '24px auto', textAlign: 'center' }}>
      <h1 style={{ display: 'flex', alignItems: 'center', gap: 8, justifyContent: 'center' }}>
        Login
        <span className="hand" style={{ color: 'var(--accent)' }}>
          <HandUnderline />
        </span>
      </h1>
      {info && (
        <div style={{
          border: '1px solid var(--border)',
          background: 'var(--bg-elevated)',
          borderRadius: 'var(--radius)',
          padding: 12,
          marginBottom: 12,
          color: 'var(--text)',
          textAlign: 'center'
        }}>
          {info}
        </div>
      )}
      <Card>
        <CardHeader>
          <CardTitle>Welcome back</CardTitle>
        </CardHeader>
        <div style={{ display: 'grid', gap: 12, textAlign: 'left' }}>
          <div style={{ fontSize: 12, color: 'var(--text-muted)', textAlign: 'center' }}>
            New here? If you just signed up, check your email for a verification code. Then verify your account on the <Link to={`/verify${email ? `?email=${encodeURIComponent(email)}` : ''}`}>Verify page</Link>.
          </div>
          <Field>
            <FieldLabel>Email</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-envelope" aria-hidden="true" /></IconSlot>
                <Input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@example.com" style={{ paddingLeft: 36 }} />
              </div>
            </FieldRow>
          </Field>
          <Field>
            <FieldLabel>Password</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-lock" aria-hidden="true" /></IconSlot>
                <Input type={showPassword ? 'text' : 'password'} value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" style={{ paddingLeft: 36 }} />
              </div>
              <Button onClick={() => setShowPassword((s) => !s)} aria-label={showPassword ? 'Hide password' : 'Show password'} title={showPassword ? 'Hide password' : 'Show password'}>
                <i className={`pi ${showPassword ? 'pi-eye-slash' : 'pi-eye'}`} aria-hidden="true" />
              </Button>
            </FieldRow>
          </Field>
          <Button
            variant="primary"
            onClick={() =>
              login.mutate(
                { email, password },
                {
                  onSuccess: () => navigate(from, { replace: true }),
                  onError: async (e) => {
                    // If backend signals verification required via message, redirect to verify
                    try {
                      if (e instanceof API.ResponseError) {
                      const body = await e.response.json();
                      const msg = (body?.message || '').toString().toLowerCase();
                      if (msg.includes('verify')) {
                        setInfo('Your account is not verified yet. We\'ve redirected you to the verification page.');
                        navigate(`/verify?email=${encodeURIComponent(email)}`);
                        return;
                      }
                      }
                    } catch {}
                  },
                }
              )
            }
            disabled={login.isPending}
          >
            {login.isPending ? (
              <>
                <Spinner /> Logging in...
              </>
            ) : (
              <>
                <i className="pi pi-sign-in" aria-hidden="true" /> Login
              </>
            )}
          </Button>
          <div style={{
            textAlign: 'center',
            padding: '12px 16px',
            border: '1px solid var(--border)',
            borderRadius: 'var(--radius)',
            background: 'var(--bg-elevated)',
            fontWeight: 600
          }}>
            <small>Don’t have an account? <Link to="/signup">Sign up</Link></small>
          </div>
        </div>
      </Card>
    </div>
  );
}
