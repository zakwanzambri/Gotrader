export interface TradingSignal {
  id: number;
  symbol: string;
  action: string;
  price: number;
  entry_price?: number;
  stop_loss?: number;
  take_profit?: number;
  confidence: string;
  source: string;
  raw_message: string;
  timestamp: string;
  status: string;
}

export interface SignalResponse {
  signals: TradingSignal[];
  total: number;
  page: number;
  limit: number;
}

export interface Stats {
  total: number;
  active: number;
  closed: number;
}

export interface WSMessage {
  type: string;
  payload: any;
}