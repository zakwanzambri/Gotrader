import { SignalResponse, Stats } from './types';

const API_BASE = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

export const api = {
  async getSignals(page = 1, limit = 20, symbol = '', status = ''): Promise<SignalResponse> {
    const params = new URLSearchParams();
    params.append('page', page.toString());
    params.append('limit', limit.toString());
    if (symbol) params.append('symbol', symbol);
    if (status) params.append('status', status);

    const response = await fetch(`${API_BASE}/signals?${params}`);
    if (!response.ok) {
      throw new Error('Failed to fetch signals');
    }
    return response.json();
  },

  async updateSignalStatus(id: number, status: string): Promise<void> {
    const response = await fetch(`${API_BASE}/signals/${id}/status`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ status }),
    });
    if (!response.ok) {
      throw new Error('Failed to update signal status');
    }
  },

  async getStats(): Promise<Stats> {
    const response = await fetch(`${API_BASE}/stats`);
    if (!response.ok) {
      throw new Error('Failed to fetch stats');
    }
    return response.json();
  },
};