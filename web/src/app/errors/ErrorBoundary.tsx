import React from 'react';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';

type State = { hasError: boolean; error?: any };

export class ErrorBoundary extends React.Component<React.PropsWithChildren, State> {
  state: State = { hasError: false };

  static getDerivedStateFromError(error: any) {
    return { hasError: true, error } as State;
  }

  componentDidCatch(error: any, info: any) {
    // no-op: could log to monitoring here
  }

  handleReload = () => {
    this.setState({ hasError: false, error: undefined });
    window.location.assign('/');
  };

  render() {
    if (this.state.hasError) {
      return (
        <main className="w-full max-w-5xl mx-auto px-4 py-6" role="main" aria-live="polite">
          <Card>
            <CardHeader>
              <CardTitle>Something went wrong</CardTitle>
            </CardHeader>
            <div style={{ display: 'grid', gap: 12 }}>
              <div>We hit an unexpected error. Please try again.</div>
              <div>
                <Button variant="primary" onClick={this.handleReload}>Go home</Button>
              </div>
            </div>
          </Card>
        </main>
      );
    }
    return this.props.children;
  }
}
