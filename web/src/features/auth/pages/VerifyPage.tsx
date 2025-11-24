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

  const onResend = () => {
    if (!email) {
      toast.error('Enter your email address first');
      return;
    }
    resend.mutate({ email });
  };

  return (
    <div style={{ maxWidth: 480, margin: '24px auto' }}>
      <h1 style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        Verify your email
        <span className="hand" style={{ color: 'var(--accent)' }}>
          <HandUnderline />
        </span>
      </h1>
      <Card>
        <CardHeader>
          <CardTitle>Verification</CardTitle>
        </CardHeader>
        <div style={{ display: 'grid', gap: 12 }}>
          <Field>
            <FieldLabel>Email</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-envelope" aria-hidden="true" /></IconSlot>
                <Input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@example.com" style={{ paddingLeft: 36 }} />
              </div>
            </FieldRow>
          </Field>
          <Field>
            <FieldLabel>Verification code</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-key" aria-hidden="true" /></IconSlot>
                <Input value={code} onChange={(e) => setCode(e.target.value)} placeholder="123456" style={{ paddingLeft: 36 }} />
              </div>
            </FieldRow>
          </Field>
          <Button variant="primary" onClick={onVerify} disabled={verify.isPending}>
            {verify.isPending ? (
              <>
                <Spinner /> Verifying...
              </>
            ) : (
              <>
                <i className="pi pi-check-circle" aria-hidden="true" /> Verify
              </>
            )}
          </Button>
          <Button variant="ghost" onClick={onResend} disabled={resend.isPending}>
            {resend.isPending ? (
              <>
                <Spinner /> Sending...
              </>
            ) : (
              <>
                <i className="pi pi-envelope" aria-hidden="true" /> Resend code
              </>
            )}
          </Button>
        </div>
      </Card>
    </div>
  );
}
