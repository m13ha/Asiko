import { useMemo, useState } from 'react';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { useAddBan, useBanList, useRemoveBan } from '../hooks';
import { EmptyState, EmptyTitle, EmptyDescription, EmptyAction } from '@/components/EmptyState';
import { Field, FieldLabel } from '@/components/Field';
import type { EntitiesBanListEntry } from '@appointment-master/api-client';

export function BanListPage() {
  const { data, isLoading, error } = useBanList();
  const add = useAddBan();
  const remove = useRemoveBan();
  const [email, setEmail] = useState('');
  const [removingEmail, setRemovingEmail] = useState<string>('');

  const entries = useMemo(() => (data ?? []) as EntitiesBanListEntry[], [data]);
  const canSubmit = !!email.trim() && !add.isPending;

  const onSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (!canSubmit) return;
    const nextEmail = email.trim();
    add.mutate(nextEmail, { onSuccess: () => setEmail('') });
  };

  return (
    <div className="grid gap-6">
      <div className="max-w-4xl mx-auto w-full space-y-6">
        <header className="flex items-start justify-between gap-3 flex-wrap">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-[var(--text)] m-0">Ban List</h1>
            <p className="mt-1 text-sm text-[var(--text-muted)] m-0">
              Block abusive or unwanted booking attempts by email.
            </p>
          </div>
        </header>
        
        <Card className="p-4 sm:p-6">
          <CardHeader>
            <CardTitle>Manage Banned Contacts</CardTitle>
          </CardHeader>
          
          <div className="space-y-6">
            <form onSubmit={onSubmit} className="space-y-2">
              <Field>
                <FieldLabel>Email Address</FieldLabel>
                <div className="grid grid-cols-1 sm:grid-cols-[minmax(0,1fr)_auto] gap-3 items-start">
                  <div className="relative min-w-0">
                    <div className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)]">
                      <i className="pi pi-envelope text-sm" aria-hidden="true" />
                    </div>
                    <Input
                      type="email"
                      placeholder="email@example.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="pl-10 w-full"
                      inputMode="email"
                      autoComplete="email"
                    />
                  </div>
                  <Button
                    type="submit"
                    variant="primary"
                    disabled={!canSubmit}
                    className="w-full sm:w-auto"
                  >
                    {add.isPending ? 'Adding…' : 'Add'}
                  </Button>
                </div>
              </Field>
            </form>

            {isLoading && (
              <div className="text-center py-4 text-[var(--text-muted)]">Loading...</div>
            )}
            
            {error && (
              <div className="text-red-600 bg-red-50 border border-red-200 rounded-lg p-3 text-sm">
                Failed to load ban list.
              </div>
            )}

            <div className="space-y-3">
              {entries.map((entry) => (
                <div 
                  key={entry.id ?? entry.bannedEmail} 
                  className="grid grid-cols-1 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center gap-3 border border-[var(--border)] rounded-lg p-4 bg-[var(--bg-elevated)]"
                >
                  <div className="min-w-0 flex-1">
                    <div className="font-semibold text-[var(--text)] break-all">{entry.bannedEmail}</div>
                    {entry.createdAt && (
                      <div className="text-xs text-[var(--text-muted)] mt-1">
                        Added {new Date(entry.createdAt).toLocaleString()}
                      </div>
                    )}
                  </div>
                  <Button 
                    variant="ghost" 
                    size="sm"
                    disabled={remove.isPending && removingEmail === entry.bannedEmail}
                    onClick={() => {
                      const target = entry.bannedEmail ?? '';
                      if (!target) return;
                      setRemovingEmail(target);
                      remove.mutate(target, { onSettled: () => setRemovingEmail('') });
                    }}
                    className="w-full sm:w-auto"
                  >
                    {remove.isPending && removingEmail === entry.bannedEmail ? 'Removing…' : 'Remove'}
                  </Button>
                </div>
              ))}
              
              {(!entries || entries.length === 0) && !isLoading && (
                <EmptyState>
                  <EmptyTitle>No banned emails yet</EmptyTitle>
                  <EmptyDescription>Add an email to prevent booking attempts.</EmptyDescription>
                  {email && (
                    <EmptyAction>
                      <Button 
                        variant="primary" 
                        disabled={add.isPending} 
                        onClick={() => add.mutate(email)}
                      >
                        Add current email
                      </Button>
                    </EmptyAction>
                  )}
                </EmptyState>
              )}
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}
