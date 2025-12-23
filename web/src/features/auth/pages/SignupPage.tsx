import { useState } from 'react';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useSignup } from '@/features/auth/hooks';
import { useNavigate, Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { HandUnderline } from '@/components/HandUnderline';
import { Spinner } from '@/components/Spinner';

export function SignupPage() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const signup = useSignup();
  const navigate = useNavigate();

  return (
    <div className="max-w-xl mx-auto my-6 text-center">
      <h1 className="flex items-center justify-center gap-2 text-2xl font-semibold text-[var(--text)]">
        Create Account
        <span className="font-['Caveat',cursive] tracking-[0.2px]" style={{ color: 'var(--accent)' }}>
          <HandUnderline />
        </span>
      </h1>
      <Card>
        <CardHeader>
          <CardTitle>Create Account</CardTitle>
        </CardHeader>
        <div className="grid gap-3 text-left">
          <Field>
            <FieldLabel>Name</FieldLabel>
            <FieldRow>
              <div className="relative flex-grow">
                <IconSlot><i className="pi pi-user" aria-hidden="true" /></IconSlot>
                <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Jane Doe" className="pl-9" />
              </div>
            </FieldRow>
          </Field>
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
          <Field>
            <FieldLabel>Confirm password</FieldLabel>
            <FieldRow>
              <div className="relative flex-grow">
                <IconSlot><i className="pi pi-lock" aria-hidden="true" /></IconSlot>
                <Input type={showConfirm ? 'text' : 'password'} value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} placeholder="••••••••" className="pl-9" />
              </div>
              <Button onClick={() => setShowConfirm((s) => !s)} aria-label={showConfirm ? 'Hide password' : 'Show password'} title={showConfirm ? 'Hide password' : 'Show password'} className="flex-shrink-0">
                <i className={`pi ${showConfirm ? 'pi-eye-slash' : 'pi-eye'}`} aria-hidden="true" />
              </Button>
            </FieldRow>
          </Field>
          <Button
            variant="primary"
            onClick={() => {
              if (password !== confirmPassword) {
                toast.error('Passwords do not match');
                return;
              }
              signup.mutate(
                { name, email, password },
                {
                  onSuccess: () => navigate(`/verify?email=${encodeURIComponent(email)}`),
                  onError: (e: any) => {
                    const msg = (e?.message || '').toString().toLowerCase();
                    if (msg.includes('verification') || msg.includes('pending')) {
                      navigate(`/verify?email=${encodeURIComponent(email)}`);
                    }
                  },
                }
              );
            }}
            disabled={signup.isPending}
          >
            {signup.isPending ? (
              <>
                <Spinner /> Creating...
              </>
            ) : (
              <>
                <i className="" aria-hidden="true" /> Sign up
              </>
            )}
          </Button>
          <div className="p-3 font-semibold text-center border border-[var(--border)] rounded-lg bg-[var(--bg-elevated)] text-[var(--text)]">
            <small>Already have an account? <Link to="/login">Login</Link></small>
          </div>
        </div>
      </Card>
      <div className="mt-3 text-xs text-center text-[var(--text-muted)]">
        After signing up, we send a verification code to your email. Enter it on the <Link to={`/verify${email ? `?email=${encodeURIComponent(email)}` : ''}`}>Verify page</Link> to finish creating your account.
      </div>
    </div>
  );
}
