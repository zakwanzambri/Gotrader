# GoTrader

A full-stack application that automatically copies trading signals from Telegram groups and displays them in a custom web dashboard with real-time updates.

## Features

- **Telegram Bot Integration**: Automatically monitors Telegram groups for trading signals
- **Real-time Updates**: WebSocket connection for live signal updates
- **Signal Parsing**: Intelligent parsing of trading messages (BUY/SELL, price, stop loss, take profit)
- **Web Dashboard**: Clean, responsive React frontend for viewing and managing signals
- **Signal Management**: Update signal status (Active/Closed/Expired)
- **Filtering & Search**: Filter signals by symbol, status, and other criteria
- **Statistics Dashboard**: Overview of total, active, and closed signals
- **SQLite Database**: Lightweight database for signal storage

## Architecture

- **Backend**: Go with Gin web framework
- **Frontend**: React with TypeScript
- **Database**: SQLite
- **Real-time**: WebSocket for live updates
- **Telegram**: Bot API for monitoring groups

## Quick Start

### Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Telegram Bot Token (from @BotFather)

### 1. Clone and Setup

```bash
git clone https://github.com/zakwanzambri/Gotrader.git
cd Gotrader
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```env
# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_IDS=chat_id_1,chat_id_2,chat_id_3

# Server Configuration
PORT=8080
DB_PATH=./signals.db

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000
```

### 3. Build and Run

```bash
# Build the application
./build.sh

# Run the server
./gotrader
```

### 4. Access the Application

Open your browser and go to `http://localhost:8080`

## Setting up Telegram Bot

### 1. Create a Bot

1. Message @BotFather on Telegram
2. Send `/newbot` and follow the instructions
3. Copy the bot token to your `.env` file

### 2. Get Chat IDs

1. Add your bot to the Telegram groups you want to monitor
2. Send a message in the group
3. Visit `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Find the chat IDs in the response and add them to your `.env` file

### 3. Signal Format

The bot recognizes various trading signal formats:

```
BUY BTCUSDT @ 45000
SL: 44000
TP: 46000

SELL ETHUSDT 3200
STOP LOSS: 3250
TARGET: 3100

LONG BNBUSDT @ 320
Confidence: HIGH
```

## API Endpoints

### Signals
- `GET /api/v1/signals` - Get signals with pagination and filtering
- `PUT /api/v1/signals/:id/status` - Update signal status

### Statistics
- `GET /api/v1/stats` - Get signal statistics

### Health
- `GET /api/v1/health` - Health check

### WebSocket
- `GET /ws` - WebSocket connection for real-time updates

## Development

### Backend Development

```bash
# Install dependencies
go mod tidy

# Run in development mode
go run ./cmd/server/

# Build
go build -o gotrader ./cmd/server/
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm start

# Build for production
npm run build
```

## Docker Deployment

```bash
# Build Docker image
docker build -t gotrader .

# Run container
docker run -p 8080:8080 \
  -e TELEGRAM_BOT_TOKEN=your_token \
  -e TELEGRAM_CHAT_IDS=chat1,chat2 \
  gotrader
```

## Configuration Options

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TELEGRAM_BOT_TOKEN` | Telegram bot token | - |
| `TELEGRAM_CHAT_IDS` | Comma-separated chat IDs | - |
| `PORT` | Server port | 8080 |
| `DB_PATH` | SQLite database path | ./signals.db |
| `ALLOWED_ORIGINS` | CORS allowed origins | http://localhost:3000 |

### Signal Confidence Levels

- **HIGH**: Strong buy/sell signals with high confidence keywords
- **MEDIUM**: Standard signals (default)
- **LOW**: Weak signals with uncertainty keywords

### Signal Status

- **ACTIVE**: Signal is currently active
- **CLOSED**: Signal has been closed/executed
- **EXPIRED**: Signal has expired

## Troubleshooting

### Common Issues

1. **Bot not receiving messages**: Ensure the bot is added to the group and has proper permissions
2. **Database errors**: Check that the application has write permissions to the database path
3. **CORS errors**: Verify `ALLOWED_ORIGINS` includes your frontend URL
4. **WebSocket connection issues**: Ensure the WebSocket URL matches your server configuration

### Logs

The application logs important events including:
- New signals received
- Database operations
- WebSocket connections
- Telegram bot status

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.