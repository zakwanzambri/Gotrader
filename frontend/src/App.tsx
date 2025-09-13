import React, { useState, useEffect, useCallback } from 'react';
import { TradingSignal, Stats } from './types';
import { api } from './api';
import { useWebSocket } from './useWebSocket';
import SignalCard from './components/SignalCard';
import StatsCard from './components/StatsCard';

function App() {
  const [signals, setSignals] = useState<TradingSignal[]>([]);
  const [stats, setStats] = useState<Stats>({ total: 0, active: 0, closed: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState({
    symbol: '',
    status: '',
    page: 1,
    limit: 20
  });

  const handleNewSignal = useCallback((signal: TradingSignal) => {
    setSignals(prev => [signal, ...prev]);
    setStats(prev => ({
      ...prev,
      total: prev.total + 1,
      active: signal.status === 'ACTIVE' ? prev.active + 1 : prev.active
    }));
  }, []);

  const { connected } = useWebSocket(handleNewSignal);

  const loadSignals = useCallback(async () => {
    try {
      setLoading(true);
      const response = await api.getSignals(
        filters.page,
        filters.limit,
        filters.symbol,
        filters.status
      );
      setSignals(response.signals);
      setError(null);
    } catch (err) {
      setError('Failed to load signals');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  const loadStats = useCallback(async () => {
    try {
      const statsData = await api.getStats();
      setStats(statsData);
    } catch (err) {
      console.error('Failed to load stats:', err);
    }
  }, []);

  useEffect(() => {
    loadSignals();
  }, [loadSignals]);

  useEffect(() => {
    loadStats();
  }, [loadStats]);

  const handleStatusUpdate = (id: number, status: string) => {
    setSignals(prev => 
      prev.map(signal => 
        signal.id === id ? { ...signal, status } : signal
      )
    );
    loadStats(); // Refresh stats after status update
  };

  const handleFilterChange = (key: string, value: string) => {
    setFilters(prev => ({
      ...prev,
      [key]: value,
      page: 1 // Reset to first page when filters change
    }));
  };

  return (
    <div style={{
      minHeight: '100vh',
      backgroundColor: '#f8fafc',
      padding: '20px'
    }}>
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto'
      }}>
        {/* Header */}
        <div style={{
          marginBottom: '24px',
          textAlign: 'center'
        }}>
          <h1 style={{
            fontSize: '32px',
            fontWeight: 'bold',
            color: '#1e293b',
            margin: '0 0 8px 0'
          }}>
            GoTrader
          </h1>
          <p style={{
            fontSize: '16px',
            color: '#64748b',
            margin: 0
          }}>
            Trading Signals Dashboard
          </p>
        </div>

        {/* Stats */}
        <StatsCard stats={stats} connected={connected} />

        {/* Filters */}
        <div style={{
          backgroundColor: '#fff',
          border: '1px solid #e2e8f0',
          borderRadius: '8px',
          padding: '16px',
          marginBottom: '24px'
        }}>
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
            gap: '16px'
          }}>
            <div>
              <label style={{
                display: 'block',
                fontSize: '14px',
                fontWeight: '500',
                color: '#374151',
                marginBottom: '4px'
              }}>
                Symbol
              </label>
              <input
                type="text"
                value={filters.symbol}
                onChange={(e) => handleFilterChange('symbol', e.target.value)}
                placeholder="e.g., BTCUSDT"
                style={{
                  width: '100%',
                  padding: '8px 12px',
                  border: '1px solid #d1d5db',
                  borderRadius: '4px',
                  fontSize: '14px'
                }}
              />
            </div>

            <div>
              <label style={{
                display: 'block',
                fontSize: '14px',
                fontWeight: '500',
                color: '#374151',
                marginBottom: '4px'
              }}>
                Status
              </label>
              <select
                value={filters.status}
                onChange={(e) => handleFilterChange('status', e.target.value)}
                style={{
                  width: '100%',
                  padding: '8px 12px',
                  border: '1px solid #d1d5db',
                  borderRadius: '4px',
                  fontSize: '14px'
                }}
              >
                <option value="">All</option>
                <option value="ACTIVE">Active</option>
                <option value="CLOSED">Closed</option>
                <option value="EXPIRED">Expired</option>
              </select>
            </div>

            <div style={{ display: 'flex', alignItems: 'end' }}>
              <button
                onClick={loadSignals}
                style={{
                  padding: '8px 16px',
                  backgroundColor: '#3b82f6',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  fontSize: '14px',
                  cursor: 'pointer'
                }}
              >
                Refresh
              </button>
            </div>
          </div>
        </div>

        {/* Signals */}
        <div style={{
          backgroundColor: '#fff',
          border: '1px solid #e2e8f0',
          borderRadius: '8px',
          padding: '16px'
        }}>
          <h2 style={{
            fontSize: '20px',
            fontWeight: '600',
            color: '#1e293b',
            margin: '0 0 16px 0'
          }}>
            Trading Signals
          </h2>

          {loading && (
            <div style={{
              textAlign: 'center',
              padding: '40px',
              color: '#64748b'
            }}>
              Loading signals...
            </div>
          )}

          {error && (
            <div style={{
              textAlign: 'center',
              padding: '40px',
              color: '#dc2626',
              backgroundColor: '#fef2f2',
              borderRadius: '4px'
            }}>
              {error}
            </div>
          )}

          {!loading && !error && signals.length === 0 && (
            <div style={{
              textAlign: 'center',
              padding: '40px',
              color: '#64748b'
            }}>
              No signals found. Connect your Telegram bot to start receiving signals.
            </div>
          )}

          {!loading && !error && signals.length > 0 && (
            <div>
              {signals.map(signal => (
                <SignalCard
                  key={signal.id}
                  signal={signal}
                  onStatusUpdate={handleStatusUpdate}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;