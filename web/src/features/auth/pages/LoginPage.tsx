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
    <div className="max-w-lg mx-auto my-6 text-center">
      <h1 className="flex items-center justify-center gap-2 text-2xl font-semibold text-[var(--text)]">
        Login
        <span className="font-['Caveat',cursive] tracking-[0.2px]" style={{ color: 'var(--accent)' }}>
          <HandUnderline />
        </span>
      </h1>
      {info && (
        <div className="p-3 mb-3 text-center text-[var(--text)] border border-[var(--border)] rounded-lg bg-[var(--bg-elevated)]">
          {info}
        </div>
      )}
      <Card>
        <CardHeader>
          <CardTitle className='w-full text-center'>Welcome back</CardTitle>
        </CardHeader>
        <div className="grid gap-3 text-left">
          <div className="text-xs text-center text-[var(--text-muted)]">
            New here? If you just signed up, check your email for a verification code. Then verify your account on the <Link to={`/verify${email ? `?email=${encodeURIComponent(email)}` : ''}`}>Verify page</Link>.
          </div>
          <Field>
            <FieldLabel>Email</FieldLabel>
            <FieldRow>
              <div className="relative flex-grow">
                <IconSlot><i className="pi pi-envelope" aria-hidden="true" /></IconSlot>
                <Input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@example.com" className="pl-9" />
              </div>
            </FieldRow>
          </Field>
          <Field>
            <FieldLabel>Password</FieldLabel>
            <FieldRow>
              <div className="relative flex-grow">
                <IconSlot><i className="pi pi-lock" aria-hidden="true" /></IconSlot>
                <Input type={showPassword ? 'text' : 'password'} value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" className="pl-9" />
              </div>
              <Button onClick={() => setShowPassword((s) => !s)} aria-label={showPassword ? 'Hide password' : 'Show password'} title={showPassword ? 'Hide password' : 'Show password'} className="flex-shrink-0">
                <i className={`pi ${showPassword ? 'pi-eye-slash' : 'pi-eye'}`} aria-hidden="true" />
              </Button>
            </FieldRow>
          </Field>
      <div className="text-right -mt-2">
        <Link to="/forgot-password" className="text-xs text-[var(--primary)] hover:underline">
          Forgot password?
        </Link>
      </div>
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
                <i className="" aria-hidden="true" /> Login
              </>
            )}
          </Button>
          <div className="p-3 font-semibold text-center border border-[var(--border)] rounded-lg bg-[var(--bg-elevated)] text-[var(--text)]">
            <small>Don’t have an account? <Link to="/signup">Sign up</Link></small>
          </div>
        </div>
      </Card>
    </div>
  );
}
