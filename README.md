# Password Manager

A zero-knowledge password manager with Go backend and React frontend.

## Features

- **Secure Authentication**: JWT-based auth with Argon2id password hashing
- **Vault Management**: Organize passwords in multiple vaults
- **Entry Types**: Login, Notes, Cards, Identities
- **Modern UI**: React + Tailwind CSS + Shadcn/ui
- **Dark Mode**: Full dark theme support
- **Encryption**: Client-side + server-side encryption

## Tech Stack

**Backend:**

- Go (Fiber)
- PostgreSQL
- Redis
- Argon2id, AES-256-GCM

**Frontend:**

- React 18
- TypeScript
- Vite
- Tailwind CSS
- Shadcn/ui
- TanStack Query

## Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose

## Getting Started

### 1. Start Infrastructure

```bash
docker-compose up -d
```

This starts PostgreSQL and Redis containers.

### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod tidy

# Create config file
cat > config.yaml << EOF
server:
  address: ":8080"
  allowOrigins: "*"
database:
  host: "localhost"
  port: 5432
  user: "passwordmanager"
  password: "passwordmanager_secret"
  dbName: "passwordmanager"
  sslMode: "disable"
redis:
  host: "localhost"
  port: 6379
jwt:
  secret: "your-super-secret-jwt-key-change-in-production"
  accessExpiry: 15
  refreshExpiry: 7
crypto:
  argon2Memory: 65536
  argon2Iterations: 3
  argon2Parallelism: 4
EOF

# Run the server
go run ./cmd/server
```

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will be available at `http://localhost:3000`.

### 4. Production Build

```bash
cd frontend
npm run build
```

The build output is in the `dist/` directory.

## API Endpoints

### Authentication

| Method | Endpoint | Description |
| -------- | ---------- | ------------- |
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/logout` | Logout |
| POST | `/api/v1/auth/refresh` | Refresh token |
| GET | `/api/v1/auth/me` | Get current user |

### Vaults

| Method | Endpoint | Description |
| -------- | ---------- | ------------- |
| GET | `/api/v1/vaults` | List vaults |
| POST | `/api/v1/vaults` | Create vault |
| GET | `/api/v1/vaults/:id` | Get vault |
| PUT | `/api/v1/vaults/:id` | Update vault |
| DELETE | `/api/v1/vaults/:id` | Delete vault |

### Entries

| Method | Endpoint | Description |
| -------- | ---------- | ------------- |
| GET | `/api/v1/vaults/:id/entries` | List entries |
| POST | `/api/v1/vaults/:id/entries` | Create entry |
| GET | `/api/v1/entries/:id` | Get entry |
| PUT | `/api/v1/entries/:id` | Update entry |
| DELETE | `/api/v1/entries/:id` | Delete entry |
| POST | `/api/v1/entries/:id/favorite` | Toggle favorite |

## Project Structure

```text
PassWordManager/
├── backend/                 # Go backend
│   ├── cmd/server/         # Entry point
│   ├── internal/           # Internal packages
│   │   ├── auth/          # Auth module
│   │   ├── vault/         # Vault module
│   │   ├── crypto/        # Crypto module
│   │   └── common/        # Shared utilities
│   └── docs/              # Module docs
├── frontend/               # React frontend
│   ├── src/                # Source code
│   │   ├── components/    # UI components
│   │   ├── features/      # Feature modules
│   │   └── lib/           # Utilities
│   └── docs/              # Frontend docs
├── docs/                   # Architecture docs
├── docker-compose.yml      # Docker setup
└── README.md
```

## Security

- Master passwords are hashed with Argon2id (not stored in plaintext)
- Vault data is encrypted client-side before transmission
- JWT tokens are short-lived (15 minutes)
- Refresh token rotation
- Rate limiting on login attempts
- Account lockout after failed attempts

## License

MIT
