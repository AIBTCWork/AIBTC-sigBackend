# AI-BTC

AI-BTC is a Go-based blockchain contract service that provides interaction with Ethereum smart contracts and message notifications via Telegram Bot.

## Features

- 🔗 **Smart Contract Interaction** - Support for interacting with Ethereum smart contracts
- ✍️ **EIP-712 Signing** - Support for standard EIP-712 typed data signing
- 🤖 **Telegram Bot Notifications** - Push message notifications via Telegram Bot
- ⏰ **Scheduled Tasks** - Built-in Cron job scheduling
- 🐳 **Docker Deployment** - Support for Docker containerized deployment

## Project Structure

```
AI-BTC
├── conf/                    # Configuration files
│   ├── config.yaml          # Default configuration
│   ├── config_debug.yaml    # Debug environment configuration
│   └── config_prod.yaml     # Production environment configuration
├── config/                  # Configuration structure definitions
├── deployment/              # Deployment related files
│   ├── Dockerfile           # Docker build file
│   ├── docker-compose.yml   # Docker Compose configuration
│   ├── build.sh             # Build script
│   └── deployment.sh        # Deployment script
├── domain/                  # Domain models
├── internal/                # Internal modules
│   ├── bot/                 # Telegram Bot module
│   ├── contract/            # Contract interaction module
│   └── job/                 # Scheduled tasks module
├── ioc/                     # Dependency injection
├── server/                  # Service entry point
├── utils/                   # Utility functions
├── go.mod                   # Go module definition
└── go.sum                   # Go dependency verification
```

## Requirements

- Go 1.25+
- Docker (optional, for containerized deployment)

## Local Development

### 1. Clone the Project

```bash
git clone <repository-url>
cd AI-BTC
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configuration

Copy and modify the configuration file:

```bash
cp conf/config_debug.yaml conf/config.yaml
```

Edit `conf/config.yaml` and configure the following:

| Configuration | Description |
|---------------|-------------|
| `server.port` | Service port |
| `server.mode` | Running mode (local/dev/pre/prod) |
| `bot.api_key` | Telegram Bot API Key |
| `bot.chat_id` | Telegram Chat ID |
| `contract.token_address` | Token contract address |
| `contract.mint_address` | Mint contract address |
| `contract.private_key` | Private key |
| `contract.chain_id` | Chain ID |

### 4. Run the Service

```bash
go run server/main.go
```

## Docker Deployment

### Build Image

Navigate to the `deployment` directory and execute the build script:

```bash
cd deployment
./build.sh --mode=pre --version=latest
```

**Parameter Description:**
- `--mode`: Build mode, supports `prod` or `pre` (default: pre)
- `--version`: Image version tag (default: latest)

### Deploy Service

```bash
./deployment.sh --mode=pre --version=latest
```

### Docker Compose Manual Deployment

```bash
cd deployment
IMAGE_NAME=nebulaicli/ai-btc-contract:latest docker compose up -d
```

## Configuration Reference

### Server Configuration

```yaml
server:
  domain: http://localhost:8080    # Service domain
  port: :8086                       # Service port
  read_timeout: 5s                  # Read timeout
  write_timeout: 10s                # Write timeout
  graceful_shutdown: 30s            # Graceful shutdown duration
  mode: dev                         # Running mode
  whitelist:                        # IP whitelist
    127.0.0.1: true
```

### Bot Configuration

```yaml
bot:
  chat_id: your_chat_id             # Telegram Chat ID
  api_key: your_bot_api_key         # Telegram Bot API Key
  base_url: https://api.telegram.org/bot
```

### Contract Configuration

```yaml
contract:
  token_address: 0x...              # Token contract address
  mint_address: 0x...               # Mint contract address
  private_key: your_private_key     # Private key
  chain_id: 11155111                # Chain ID (Sepolia testnet)
```

## API Endpoints

After the service starts, you can access the following endpoints:

- Health Check: `GET /health`

## Running Modes

| Mode | Description |
|------|-------------|
| `local` | Local development mode, scheduled tasks not running |
| `dev` | Development environment |
| `pre` | Pre-production environment |
| `prod` | Production environment |

## Tech Stack

- [Go](https://golang.org/) - Programming language
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Cron](https://github.com/robfig/cron) - Scheduled tasks
- [Wire](https://github.com/google/wire) - Dependency injection
- [go-ethereum](https://github.com/ethereum/go-ethereum) - Ethereum interaction

## License

MIT License
