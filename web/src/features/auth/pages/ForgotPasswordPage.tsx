import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Field, FieldLabel, FieldRow, FieldError, IconSlot } from '@/components/Field';
import { useForgotPassword } from '../hooks';

const schema = z.object({
  email: z.string().email('Please enter a valid email address'),
});

type FormValues = z.infer<typeof schema>;

export function ForgotPasswordPage() {
  const { register, handleSubmit, formState: { errors } } = useForm<FormValues>({
    resolver: zodResolver(schema),
  });
  const forgotPassword = useForgotPassword();
  const [isSent, setIsSent] = useState(false);

  const onSubmit = (data: FormValues) => {
    forgotPassword.mutate(data.email, {
      onSuccess: () => setIsSent(true),
    });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--bg)] p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-center">Reset Password</CardTitle>
        </CardHeader>

        {isSent ? (
          <div className="text-center space-y-4 py-4">
            <div className="w-12 h-12 bg-green-100 text-green-600 rounded-full flex items-center justify-center mx-auto">
              <i className="pi pi-check text-xl" />
            </div>
            <h3 className="text-lg font-semibold text-[var(--text)]">Check your email</h3>
            <p className="text-[var(--text-muted)]">
              We have sent a password reset code to your email address.
            </p>
            <div className="pt-4">
              <Link to="/reset-password">
                <Button className="w-full">Enter Code</Button>
              </Link>
            </div>
            <div className="pt-2">
              <button 
                onClick={() => setIsSent(false)}
                className="text-sm text-[var(--primary)] hover:underline"
              >
                Try a different email
              </button>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <p className="text-sm text-[var(--text-muted)] text-center">
              Enter your email address and we'll send you a code to reset your password.
            </p>

            <Field>
              <FieldLabel>Email Address</FieldLabel>
              <FieldRow>
                <div className="relative flex-1">
                  <IconSlot>
                    <i className="pi pi-envelope" />
                  </IconSlot>
                  <Input
                    {...register('email')}
                    type="email"
                    placeholder="you@example.com"
                    className="pl-9"
                    autoComplete="email"
                  />
                </div>
              </FieldRow>
              {errors.email && <FieldError>{errors.email.message}</FieldError>}
            </Field>

            <Button
              type="submit"
              className="w-full"
              disabled={forgotPassword.isPending}
            >
              {forgotPassword.isPending ? 'Sending...' : 'Send Reset Code'}
            </Button>

            <div className="text-center text-sm">
              <Link to="/login" className="text-[var(--primary)] hover:underline font-medium">
                Back to Login
              </Link>
            </div>
          </form>
        )}
      </Card>
    </div>
  );
}
