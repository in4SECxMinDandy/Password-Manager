# Password Manager - Architecture Overview

## 1. System Objectives

Zero-knowledge password manager with multi-device support, modern UI (web, mobile, browser extension), high-performance backend, and extensible interface.

- **Backend**: Go (Fiber) - REST API, business logic, crypto, sync
- **Frontend Web**: React + TypeScript + Vite, Tailwind CSS + Shadcn/ui for modern, customizable UI
- **Database**: PostgreSQL (primary), Redis (cache/session)

## 2. Stack & Architecture

### Backend Architecture

```
backend/
├── cmd/server/           # Entry point
├── internal/
│   ├── auth/            # Module: Authentication & Authorization
│   ├── vault/          # Module: Vault & Entries
│   ├── crypto/         # Module: Encryption utilities
│   └── common/         # Shared: config, errors, middleware, logger
├── migrations/          # SQL migrations
└── docs/               # Module documentation
```

**Backend Modules:**
- **Auth Module**: User registration, login, JWT tokens, MFA support
- **Vault Module**: Vault & entry CRUD operations, authorization
- **Crypto Module**: Argon2id KDF, AES-256-GCM encryption

### Frontend Architecture

```
frontend/
├── src/
│   ├── app/             # Routing/pages
│   ├── components/     # UI components
│   │   ├── ui/         # Shadcn/ui components
│   │   └── layout/     # Layout components
│   ├── features/       # Feature modules
│   │   ├── auth/       # Auth features
│   │   └── vault/      # Vault features
│   └── lib/            # Utilities, API client
└── docs/               # Frontend documentation
```

## 3. Module List

| Module | Description | Documentation |
|--------|-------------|---------------|
| [Auth](./module-auth.md) | Authentication, authorization, JWT, MFA | [Details](./module-auth.md) |
| [Crypto](./module-crypto.md) | Encryption, key derivation, hashing | [Details](./module-crypto.md) |
| [Vault](./module-vault.md) | Vault & entry management | [Details](./module-vault.md) |
| [Frontend Architecture](./frontend-architecture.md) | Frontend structure & design | [Details](./frontend-architecture.md) |

## 4. End-to-End Flows

### User Registration & Login Flow

```
1. User submits registration form (email, master password)
2. Backend validates password strength (min 12 chars, entropy)
3. Backend derives key using Argon2id, stores hash + salt
4. Backend creates user record, returns JWT tokens
5. Frontend stores tokens, redirects to vault list

Login Flow:
1. User submits login form
2. Backend verifies credentials
3. Backend checks rate limiting (5 attempts → 15min lock)
4. Backend returns new JWT tokens
5. Frontend redirects to vault list
```

### Vault & Entry Management Flow

```
1. User creates vault
2. User adds entry (login, note, card, identity)
3. Entry data encrypted client-side + server-side envelope
4. User can favorite, search, filter entries
5. Only vault owner can access entries (authorization)
```

## 5. Security Principles

### Zero-Knowledge Architecture
- Master password never stored in plaintext
- Password hash (Argon2id) + salt stored
- Vault data encrypted client-side before transmission
- Server-side encryption envelope for additional protection
- Even admin cannot decrypt user data

### Defense in Depth
- JWT access tokens (15min expiry)
- Refresh tokens with rotation
- Rate limiting on login attempts
- Account lockout after failed attempts
- MFA support (TOTP)

### Least Privilege
- Only vault owner can access vault data
- Authorization checks on all vault operations
- Sensitive actions require re-authentication

## 6. Version & Changes

Current Architecture Version: 1.0.0

### Changelog
- Initial MVP implementation
- Core modules: Auth, Vault, Crypto
- Frontend: React + TypeScript + Tailwind + Shadcn/ui
