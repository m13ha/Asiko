import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Badge } from '@/components/Badge';
import { Link } from 'react-router-dom';
import { CopyButton } from '@/components/CopyButton';

export function AppointmentCard({ item }: { item: any }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{item.title}</CardTitle>
        {item.type && <Badge>{String(item.type)}</Badge>}
      </CardHeader>
      <div style={{ display: 'grid', gap: 6 }}>
        {item.appCode && (
          <div>
            <small>Code:</small> <strong>{item.appCode}</strong> <CopyButton value={item.appCode} ariaLabel="Copy appointment code" />
          </div>
        )}
        <div>
          <small>
            {item.startDate} {item.startTime} → {item.endDate} {item.endTime}
          </small>
        </div>
        <div>
          <Link to={`/appointments/${item.id}`} state={{ appointment: item }}>View details</Link>
        </div>
      </div>
    </Card>
  );
}
