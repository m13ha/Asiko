import { useNavigate, useSearchParams, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Field, FieldLabel, FieldRow, FieldError, IconSlot } from '@/components/Field';
import { useResetPassword } from '../hooks';

const schema = z.object({
  code: z.string().min(1, 'Code is required'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
});

type FormValues = z.infer<typeof schema>;

export function ResetPasswordPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const defaultCode = searchParams.get('code') || '';
  
  const { register, handleSubmit, formState: { errors } } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      code: defaultCode,
    },
  });
  
  const resetPassword = useResetPassword();

  const onSubmit = (data: FormValues) => {
    resetPassword.mutate(
      { token: data.code, newPassword: data.password },
      {
        onSuccess: () => {
          navigate('/login');
        },
      }
    );
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--bg)] p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-center">Set New Password</CardTitle>
        </CardHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Field>
            <FieldLabel>Reset Code</FieldLabel>
            <FieldRow>
              <div className="relative flex-1">
                <IconSlot>
                  <i className="pi pi-key" />
                </IconSlot>
                <Input
                  {...register('code')}
                  placeholder="Enter the code from your email"
                  className="pl-9"
                />
              </div>
            </FieldRow>
            {errors.code && <FieldError>{errors.code.message}</FieldError>}
          </Field>

          <Field>
            <FieldLabel>New Password</FieldLabel>
            <FieldRow>
              <div className="relative flex-1">
                <IconSlot>
                  <i className="pi pi-lock" />
                </IconSlot>
                <Input
                  {...register('password')}
                  type="password"
                  placeholder="Min. 8 characters"
                  className="pl-9"
                />
              </div>
            </FieldRow>
            {errors.password && <FieldError>{errors.password.message}</FieldError>}
          </Field>

          <Field>
            <FieldLabel>Confirm Password</FieldLabel>
            <FieldRow>
              <div className="relative flex-1">
                <IconSlot>
                  <i className="pi pi-lock" />
                </IconSlot>
                <Input
                  {...register('confirmPassword')}
                  type="password"
                  placeholder="Re-enter password"
                  className="pl-9"
                />
              </div>
            </FieldRow>
            {errors.confirmPassword && <FieldError>{errors.confirmPassword.message}</FieldError>}
          </Field>

          <Button
            type="submit"
            className="w-full"
            disabled={resetPassword.isPending}
          >
            {resetPassword.isPending ? 'Resetting...' : 'Reset Password'}
          </Button>

          <div className="text-center text-sm">
            <Link to="/login" className="text-[var(--primary)] hover:underline font-medium">
              Back to Login
            </Link>
          </div>
        </form>
      </Card>
    </div>
  );
}
