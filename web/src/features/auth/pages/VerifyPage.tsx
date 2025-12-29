import { useMemo, useState } from 'react';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { useResendVerification, useVerify } from '@/features/auth/hooks';
import { useNavigate, useLocation } from 'react-router-dom';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { HandUnderline } from '@/components/HandUnderline';
import { Spinner } from '@/components/Spinner';
import toast from 'react-hot-toast';

export function VerifyPage() {
  const loc = useLocation();
  const initialEmail = useMemo(() => {
    try {
      const sp = new URLSearchParams(loc.search);
      return sp.get('email') || '';
    } catch {
      return '';
    }
  }, [loc.search]);
  const [email, setEmail] = useState(initialEmail);
  const [code, setCode] = useState('');
  const verify = useVerify();
  const resend = useResendVerification();
  const navigate = useNavigate();

  const onVerify = () =>
    verify.mutate({ email, code }, { onSuccess: () => navigate('/dashboard') });

  const onResend = (e: React.MouseEvent) => {
    e.preventDefault();
    if (!email) {
      toast.error('Enter your email address first');
      return;
    }
    resend.mutate({ email });
  };

  return (
    <div className="flex min-h-[80vh] items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <Card className="w-full max-w-md space-y-8 p-8 shadow-2xl">
        <div className="flex flex-col items-center">
          <h1 className="flex items-center gap-2 text-3xl font-bold text-[var(--text)] text-center">
            Verify your email
            <span className="font-['Caveat',cursive] tracking-[0.2px] text-[var(--accent)]">
              <HandUnderline />
            </span>
          </h1>
          <p className="mt-2 text-center text-sm text-[var(--text-secondary)]">
            We sent a verification code to your email. Enter it below to confirm your account.
          </p>
        </div>

        <div className="space-y-6">
          <div className="space-y-4">
            <Field>
              <FieldLabel>Email Address</FieldLabel>
              <FieldRow>
                <div className="relative w-full">
                  <IconSlot><i className="pi pi-envelope text-[var(--text-secondary)]" aria-hidden="true" /></IconSlot>
                  <Input 
                    type="email" 
                    value={email} 
                    onChange={(e) => setEmail(e.target.value)} 
                    placeholder="you@example.com" 
                    className="pl-9 w-full"
                  />
                </div>
              </FieldRow>
            </Field>

            <Field>
              <FieldLabel>Verification Code</FieldLabel>
              <FieldRow>
                <div className="relative w-full">
                  <IconSlot><i className="pi pi-key text-[var(--text-secondary)]" aria-hidden="true" /></IconSlot>
                  <Input 
                    value={code} 
                    onChange={(e) => setCode(e.target.value)} 
                    placeholder="123456" 
                    className="pl-9 w-full tracking-widest font-mono"
                    autoFocus={!!initialEmail}
                  />
                </div>
              </FieldRow>
            </Field>
          </div>

          <Button 
            variant="primary" 
            onClick={onVerify} 
            disabled={verify.isPending || !code || !email}
            className="w-full py-3 text-lg font-medium shadow-md transition-all hover:shadow-lg focus:ring-2 focus:ring-offset-2 focus:ring-[var(--primary)]"
          >
            {verify.isPending ? (
              <span className="flex items-center justify-center gap-2">
                <Spinner /> Verifying...
              </span>
            ) : (
              <span className="flex items-center justify-center gap-2">
                Verify Email <i className="pi pi-arrow-right" aria-hidden="true" />
              </span>
            )}
          </Button>

          <div className="text-center text-sm">
            <span className="text-[var(--text-secondary)]">Didn't receive the code? </span>
            <button 
              onClick={onResend} 
              disabled={resend.isPending}
              className={`font-medium text-[var(--primary)] hover:text-[var(--primary-hover)] transition-colors focus:outline-none focus:underline ${resend.isPending ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
            >
              {resend.isPending ? 'Sending...' : 'Resend'}
            </button>
          </div>
        </div>
      </Card>
    </div>
  );
}
