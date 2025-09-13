import React from 'react';
import { TradingSignal } from '../types';
import { api } from '../api';

interface SignalCardProps {
  signal: TradingSignal;
  onStatusUpdate: (id: number, status: string) => void;
}

const SignalCard: React.FC<SignalCardProps> = ({ signal, onStatusUpdate }) => {
  const handleStatusChange = async (status: string) => {
    try {
      await api.updateSignalStatus(signal.id, status);
      onStatusUpdate(signal.id, status);
    } catch (error) {
      console.error('Failed to update signal status:', error);
    }
  };

  const getActionColor = (action: string) => {
    switch (action.toUpperCase()) {
      case 'BUY':
      case 'LONG':
        return '#10b981'; // green
      case 'SELL':
      case 'SHORT':
        return '#ef4444'; // red
      default:
        return '#6b7280'; // gray
    }
  };

  const getConfidenceColor = (confidence: string) => {
    switch (confidence.toUpperCase()) {
      case 'HIGH':
        return '#059669';
      case 'MEDIUM':
        return '#d97706';
      case 'LOW':
        return '#dc2626';
      default:
        return '#6b7280';
    }
  };

  return (
    <div style={{
      border: '1px solid #e5e7eb',
      borderRadius: '8px',
      padding: '16px',
      margin: '8px 0',
      backgroundColor: '#fff',
      boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)'
    }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div style={{ flex: 1 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '8px' }}>
            <span style={{
              backgroundColor: getActionColor(signal.action),
              color: 'white',
              padding: '4px 8px',
              borderRadius: '4px',
              fontSize: '12px',
              fontWeight: 'bold'
            }}>
              {signal.action}
            </span>
            <span style={{ fontSize: '18px', fontWeight: 'bold' }}>{signal.symbol}</span>
            <span style={{ fontSize: '16px', color: '#374151' }}>@ {signal.price}</span>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(120px, 1fr))', gap: '8px', marginBottom: '8px' }}>
            {signal.entry_price && (
              <div>
                <span style={{ fontSize: '12px', color: '#6b7280' }}>Entry: </span>
                <span style={{ fontWeight: 'bold' }}>{signal.entry_price}</span>
              </div>
            )}
            {signal.stop_loss && (
              <div>
                <span style={{ fontSize: '12px', color: '#6b7280' }}>SL: </span>
                <span style={{ fontWeight: 'bold', color: '#ef4444' }}>{signal.stop_loss}</span>
              </div>
            )}
            {signal.take_profit && (
              <div>
                <span style={{ fontSize: '12px', color: '#6b7280' }}>TP: </span>
                <span style={{ fontWeight: 'bold', color: '#10b981' }}>{signal.take_profit}</span>
              </div>
            )}
            <div>
              <span style={{ fontSize: '12px', color: '#6b7280' }}>Confidence: </span>
              <span style={{ 
                fontWeight: 'bold', 
                color: getConfidenceColor(signal.confidence) 
              }}>
                {signal.confidence}
              </span>
            </div>
          </div>

          <div style={{ fontSize: '12px', color: '#6b7280', marginBottom: '8px' }}>
            <div>Source: {signal.source}</div>
            <div>Time: {new Date(signal.timestamp).toLocaleString()}</div>
          </div>

          <details style={{ fontSize: '12px', color: '#6b7280' }}>
            <summary style={{ cursor: 'pointer' }}>Raw Message</summary>
            <div style={{ 
              marginTop: '4px', 
              padding: '8px', 
              backgroundColor: '#f9fafb', 
              borderRadius: '4px',
              fontFamily: 'monospace'
            }}>
              {signal.raw_message}
            </div>
          </details>
        </div>

        <div style={{ marginLeft: '16px' }}>
          <select
            value={signal.status}
            onChange={(e) => handleStatusChange(e.target.value)}
            style={{
              padding: '4px 8px',
              borderRadius: '4px',
              border: '1px solid #d1d5db',
              fontSize: '12px'
            }}
          >
            <option value="ACTIVE">Active</option>
            <option value="CLOSED">Closed</option>
            <option value="EXPIRED">Expired</option>
          </select>
        </div>
      </div>
    </div>
  );
};

export default SignalCard;