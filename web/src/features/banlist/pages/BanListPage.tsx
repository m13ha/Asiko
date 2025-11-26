import { useState } from 'react';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { useAddBan, useBanList, useRemoveBan } from '../hooks';
import { EmptyState, EmptyTitle, EmptyDescription, EmptyAction } from '@/components/EmptyState';
import { Field, FieldLabel } from '@/components/Field';

export function BanListPage() {
  const { data, isLoading, error } = useBanList();
  const add = useAddBan();
  const remove = useRemoveBan();
  const [email, setEmail] = useState('');

  return (
    <div className="min-h-screen bg-gray-50 p-4 sm:p-8">
      <div className="max-w-4xl mx-auto space-y-6">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-800">Ban List</h1>
        
        <Card className="p-4 sm:p-6">
          <CardHeader>
            <CardTitle>Manage Banned Contacts</CardTitle>
          </CardHeader>
          
          <div className="space-y-6">
            <Field>
              <FieldLabel>Email Address</FieldLabel>
              <div className="flex flex-col sm:flex-row gap-3">
                <div className="relative flex-1">
                  <div className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
                    <i className="pi pi-envelope text-sm" aria-hidden="true" />
                  </div>
                  <Input 
                    type="email" 
                    placeholder="email@example.com" 
                    value={email} 
                    onChange={(e) => setEmail(e.target.value)}
                    className="pl-10 w-full"
                  />
                </div>
                <Button 
                  variant="primary" 
                  disabled={!email || add.isPending}
                  onClick={() => add.mutate(email)}
                  className="w-full sm:w-auto"
                >
                  {add.isPending ? 'Adding…' : 'Add'}
                </Button>
              </div>
            </Field>

            {isLoading && (
              <div className="text-center py-4 text-gray-600">Loading...</div>
            )}
            
            {error && (
              <div className="text-red-600 bg-red-50 border border-red-200 rounded-lg p-3 text-sm">
                Failed to load ban list.
              </div>
            )}

            <div className="space-y-3">
              {(data || []).map((entry: any) => (
                <div 
                  key={entry.id} 
                  className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border border-gray-200 rounded-lg p-4 bg-white"
                >
                  <div className="min-w-0 flex-1">
                    <div className="font-semibold text-gray-900 truncate">{entry.bannedEmail}</div>
                    {entry.createdAt && (
                      <div className="text-xs text-gray-500 mt-1">
                        Added {new Date(entry.createdAt).toLocaleString()}
                      </div>
                    )}
                  </div>
                  <Button 
                    variant="ghost" 
                    size="sm"
                    onClick={() => remove.mutate(entry.bannedEmail)}
                    className="w-full sm:w-auto"
                  >
                    Remove
                  </Button>
                </div>
              ))}
              
              {(!data || data.length === 0) && !isLoading && (
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