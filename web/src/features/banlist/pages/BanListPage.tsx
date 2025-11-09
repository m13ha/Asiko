import { useState } from 'react';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { useAddBan, useBanList, useRemoveBan } from '../hooks';
import { EmptyState, EmptyTitle, EmptyDescription, EmptyAction } from '@/components/EmptyState';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { Mail } from 'lucide-react';

export function BanListPage() {
  const { data, isLoading, error } = useBanList();
  const add = useAddBan();
  const remove = useRemoveBan();
  const [email, setEmail] = useState('');

  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <h1 style={{ margin: 0 }}>Ban List</h1>
      <Card>
        <CardHeader>
          <CardTitle>Manage Banned Contacts</CardTitle>
        </CardHeader>
        <div style={{ display: 'grid', gap: 10 }}>
          <Field>
            <FieldLabel>Email</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><Mail size={16} /></IconSlot>
                <Input type="email" placeholder="email@example.com" value={email} onChange={(e: any) => setEmail(e.target.value)} style={{ paddingLeft: 36 }} />
              </div>
              <Button variant="primary" disabled={!email || add.isPending} onClick={() => add.mutate(email)}>
                {add.isPending ? 'Adding…' : 'Add'}
              </Button>
            </FieldRow>
          </Field>

          {isLoading && <div>Loading...</div>}
          {error && <div style={{ color: 'var(--danger)' }}>Failed to load ban list.</div>}

          <div style={{ display: 'grid', gap: 6 }}>
            {(data || []).map((entry: any) => (
              <div key={entry.id} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', border: '1px solid var(--border)', borderRadius: 8, padding: '8px 12px' }}>
                <div>
                  <div style={{ fontWeight: 600 }}>{entry.bannedEmail}</div>
                  {entry.createdAt && (<div style={{ fontSize: 12, opacity: 0.7 }}>{new Date(entry.createdAt).toLocaleString()}</div>)}
                </div>
                <Button variant="ghost" onClick={() => remove.mutate(entry.bannedEmail)}>Remove</Button>
              </div>
            ))}
            {(!data || data.length === 0) && (
              <EmptyState>
                <EmptyTitle>No banned emails yet</EmptyTitle>
                <EmptyDescription>Add an email to prevent booking attempts.</EmptyDescription>
                <EmptyAction>
                  <Button variant="primary" disabled={!email || add.isPending} onClick={() => add.mutate(email)}>
                    Add current email
                  </Button>
                </EmptyAction>
              </EmptyState>
            )}
          </div>
        </div>
      </Card>
    </div>
  );
}
