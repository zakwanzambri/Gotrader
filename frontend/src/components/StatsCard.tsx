import React from 'react';
import { Stats } from '../types';

interface StatsCardProps {
  stats: Stats;
  connected: boolean;
}

const StatsCard: React.FC<StatsCardProps> = ({ stats, connected }) => {
  return (
    <div style={{
      display: 'grid',
      gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
      gap: '16px',
      marginBottom: '24px'
    }}>
      <div style={{
        backgroundColor: '#f8fafc',
        border: '1px solid #e2e8f0',
        borderRadius: '8px',
        padding: '16px',
        textAlign: 'center'
      }}>
        <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#1e293b' }}>
          {stats.total}
        </div>
        <div style={{ fontSize: '14px', color: '#64748b' }}>Total Signals</div>
      </div>

      <div style={{
        backgroundColor: '#f0fdf4',
        border: '1px solid #bbf7d0',
        borderRadius: '8px',
        padding: '16px',
        textAlign: 'center'
      }}>
        <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#166534' }}>
          {stats.active}
        </div>
        <div style={{ fontSize: '14px', color: '#16a34a' }}>Active Signals</div>
      </div>

      <div style={{
        backgroundColor: '#fafafa',
        border: '1px solid #e5e5e5',
        borderRadius: '8px',
        padding: '16px',
        textAlign: 'center'
      }}>
        <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#525252' }}>
          {stats.closed}
        </div>
        <div style={{ fontSize: '14px', color: '#737373' }}>Closed Signals</div>
      </div>

      <div style={{
        backgroundColor: connected ? '#f0fdf4' : '#fef2f2',
        border: `1px solid ${connected ? '#bbf7d0' : '#fecaca'}`,
        borderRadius: '8px',
        padding: '16px',
        textAlign: 'center'
      }}>
        <div style={{
          width: '12px',
          height: '12px',
          borderRadius: '50%',
          backgroundColor: connected ? '#16a34a' : '#dc2626',
          margin: '0 auto 8px'
        }}></div>
        <div style={{ fontSize: '14px', color: connected ? '#16a34a' : '#dc2626' }}>
          {connected ? 'Connected' : 'Disconnected'}
        </div>
      </div>
    </div>
  );
};

export default StatsCard;