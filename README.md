# IoT Platform Go

IoT Device Management Platform built with Go

A scalable IoT platform for managing and monitoring smart devices in real-time.
Features MQTT integration, time-series data processing, and WebSocket-based dashboards.

## Tech Stack

- **Backend**: Go, Gin (Web Framework)
- **Database**: PostgreSQL (Metadata), InfluxDB (Time-series data)
- **Message Broker**: MQTT (Mosquitto)
- **Real-time Communication**: WebSocket
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose
- **Monitoring**: Grafana (Optional)

## Features

### Phase 1: Core Infrastructure ✅
- [x] RESTful API for device management
- [x] PostgreSQL database integration
- [x] Clean architecture with repository pattern
- [x] Configuration management
- [x] CORS middleware
- [x] Health check endpoint

### Phase 2: MQTT Integration ✅
- [x] MQTT client implementation
- [x] Device data collection
- [x] Real-time message processing
- [x] Device status monitoring

### Phase 3: Time-series Data (Planned)
- [ ] InfluxDB integration
- [ ] Data aggregation and compression
- [ ] Historical data analysis
- [ ] Performance optimization

### Phase 4: Real-time Dashboard (Planned)
- [ ] WebSocket implementation
- [ ] Frontend dashboard
- [ ] Real-time data visualization
- [ ] Device control interface

## Project Structure

```
iot-platform-go/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── api/            # HTTP handlers
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and setup
│   ├── device/         # Device business logic
│   ├── mqtt/           # MQTT client (planned)
│   └── websocket/      # WebSocket hub (planned)
├── pkg/
│   ├── models/         # Data models
│   └── utils/          # Utility functions
├── configs/            # Configuration files
├── docs/               # Documentation
├── scripts/            # Build and deployment scripts
├── docker-compose.yml  # Docker services
├── Makefile           # Development commands
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL (via Docker)

### 1. Clone the repository

```bash
git clone <repository-url>
cd iot-platform-go
```

### 2. Start the infrastructure

```bash
make docker-up
```

This will start:
- PostgreSQL database on port 5432
- MQTT broker (Mosquitto) on port 1883
- Grafana on port 3000 (optional)

### 3. Set up environment variables

```bash
cp configs/env.example .env
# Edit .env with your configuration
```

### 4. Install dependencies and run

```bash
make deps
make run
```

The server will start on `http://localhost:8080`

### 5. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Create a device
curl -X POST http://localhost:8080/api/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Temperature Sensor 1",
    "type": "temperature",
    "location": "Living Room"
  }'

# Get all devices
curl http://localhost:8080/api/devices
```

## API Endpoints

### Devices

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/devices` | Get all devices |
| POST | `/api/devices` | Create a new device |
| GET | `/api/devices/:id` | Get device by ID |
| PUT | `/api/devices/:id` | Update device |
| DELETE | `/api/devices/:id` | Delete device |
| GET | `/api/devices/:id/status` | Get device status |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |

## Development

### Available Commands

```bash
make build      # Build the application
make run        # Run the application
make test       # Run tests
make clean      # Clean build artifacts
make docker-up  # Start Docker services
make docker-down# Stop Docker services
make logs       # Show Docker logs
make fmt        # Format code
make lint       # Lint code
make deps       # Install dependencies
make help       # Show all commands
```

### Adding New Features

1. **Models**: Add new data models in `pkg/models/`
2. **Repository**: Implement database operations in `internal/device/`
3. **Handler**: Add HTTP handlers in `internal/api/`
4. **Routes**: Register new routes in `cmd/server/main.go`

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/device -v
```

## Deployment

### Docker

```bash
# Build Docker image
docker build -t iot-platform-go .

# Run with Docker Compose
docker-compose up -d
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | 8080 |
| `SERVER_HOST` | Server host | localhost |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_NAME` | Database name | iot_platform |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | password |
| `MQTT_BROKER` | MQTT broker URL | tcp://localhost:1883 |
| `JWT_SECRET` | JWT secret key | your-secret-key-here |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Security

### GitHub Actions Security

このプロジェクトでは、セキュリティスキャンと依存関係チェックを自動化しています。

#### 必要な設定

GitHub Actionsでセキュリティチェックを実行するには、以下の設定が必要です：

1. **Dependency Graph の有効化**
   - リポジトリの `Settings` → `Security & analysis` で有効化
   - プライベートリポジトリの場合は GitHub Advanced Security も必要

2. **セキュリティワークフロー**
   - `security.yml`: 包括的なセキュリティスキャン
   - `dependency-check.yml`: 依存関係の詳細チェック

詳細な設定手順は [docs/SECURITY_SETUP.md](docs/SECURITY_SETUP.md) を参照してください。

#### 実行されるチェック

- **Dependency Review**: 依存関係の変更を自動チェック
- **Trivy Vulnerability Scanner**: 既知の脆弱性をスキャン
- **GoSec**: Goコードのセキュリティ問題を検出
- **govulncheck**: Goの脆弱性データベースをチェック

## License

This project is licensed under the MIT License.

## Roadmap

- [ ] MQTT integration for real-time device communication
- [ ] WebSocket support for live dashboard updates
- [ ] InfluxDB integration for time-series data
- [ ] Authentication and authorization
- [ ] Device firmware management
- [ ] Alert system
- [ ] Mobile app support
- [ ] Machine learning for anomaly detection
