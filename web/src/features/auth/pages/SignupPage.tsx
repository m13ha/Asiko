import { useState } from 'react';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useSignup } from '@/features/auth/hooks';
import { useNavigate, Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { User, Mail, Lock, Eye, EyeOff, UserPlus } from 'lucide-react';
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
    <div style={{ maxWidth: 480, margin: '24px auto', textAlign: 'center' }}>
      <h1 style={{ display: 'flex', alignItems: 'center', gap: 8, justifyContent: 'center' }}>
        Create Account
        <span className="hand" style={{ color: 'var(--accent)' }}>
          <HandUnderline />
        </span>
      </h1>
      <Card>
        <CardHeader>
          <CardTitle>Create Account</CardTitle>
        </CardHeader>
        <div style={{ display: 'grid', gap: 12, textAlign: 'left' }}>
          <Field>
            <FieldLabel>Name</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><User size={16} /></IconSlot>
                <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Jane Doe" style={{ paddingLeft: 36 }} />
              </div>
            </FieldRow>
          </Field>
          <Field>
            <FieldLabel>Email</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><Mail size={16} /></IconSlot>
                <Input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@example.com" style={{ paddingLeft: 36 }} />
              </div>
            </FieldRow>
          </Field>
          <Field>
            <FieldLabel>Password</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><Lock size={16} /></IconSlot>
                <Input type={showPassword ? 'text' : 'password'} value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" style={{ paddingLeft: 36 }} />
              </div>
              <Button onClick={() => setShowPassword((s) => !s)} aria-label={showPassword ? 'Hide password' : 'Show password'} title={showPassword ? 'Hide password' : 'Show password'}>
                {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
              </Button>
            </FieldRow>
          </Field>
          <Field>
            <FieldLabel>Confirm password</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><Lock size={16} /></IconSlot>
                <Input type={showConfirm ? 'text' : 'password'} value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} placeholder="••••••••" style={{ paddingLeft: 36 }} />
              </div>
              <Button onClick={() => setShowConfirm((s) => !s)} aria-label={showConfirm ? 'Hide password' : 'Show password'} title={showConfirm ? 'Hide password' : 'Show password'}>
                {showConfirm ? <EyeOff size={16} /> : <Eye size={16} />}
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
                <UserPlus size={16} /> Sign up
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
            <small>Already have an account? <Link to="/login">Login</Link></small>
          </div>
        </div>
      </Card>
      <div style={{ marginTop: 12, fontSize: 12, color: 'var(--text-muted)', textAlign: 'center' }}>
        After signing up, we send a verification code to your email. Enter it on the <Link to={`/verify${email ? `?email=${encodeURIComponent(email)}` : ''}`}>Verify page</Link> to finish creating your account.
      </div>
    </div>
  );
}
