# Auth Module (Go)

## 1. Purpose

Handles user authentication and authorization. Belongs to the **Application & Domain layers**.

## 2. Scope & Bounded Context

**Entities owned:**
- `User` - User account with credentials and MFA settings

**Referenced:**
- Redis (token storage, rate limiting)
- Crypto module (password hashing)

## 3. Public API

### HTTP Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login (returns tokens) |
| POST | `/api/v1/auth/logout` | Logout (revoke tokens) |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/api/v1/auth/me` | Get current user |
| POST | `/api/v1/auth/check-password` | Check password strength |

### Services

```go
type AuthService interface {
    Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
    Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
    Logout(ctx context.Context, userID uuid.UUID) error
    RefreshTokens(ctx context.Context, refreshToken string) (*AuthResponse, error)
    GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
}

type TokenService interface {
    GenerateTokens(userID uuid.UUID, email string) (accessToken, refreshToken string, err error)
    ValidateAccessToken(token string) (*Claims, error)
    RefreshTokens(refreshToken string) (newAccess, newRefresh string, err error)
    RevokeToken(tokenID string) error
    RevokeAllUserTokens(userID uuid.UUID) error
}
```

## 4. Business Flows

### Registration Flow

```
1. Validate email format
2. Validate password strength (min 12 chars, upper/lower/digit required)
3. Check if email already exists
4. Generate salt, derive key using Argon2id
5. Store user with password hash
6. Generate JWT tokens
7. Return tokens + user info
```

### Login Flow

```
1. Check rate limiting (max 5 attempts → 15min lock)
2. Find user by email
3. Check if account is locked
4. Derive key from password + stored salt
5. Compare with stored hash
6. On failure: increment failed count, lock if threshold reached
7. On success: reset failed count, generate new tokens
8. Return tokens + user info
```

## 5. Business Rules

| Rule | Description |
|------|-------------|
| Password minimum | 12 characters |
| Password complexity | Uppercase, lowercase, digit required |
| Login attempts | Max 5 → 15 minute lockout |
| Access token expiry | 15 minutes |
| Refresh token expiry | 7 days |

## 6. Security & Privacy

- **Authentication**: JWT with HMAC-SHA256 signing
- **Password storage**: Argon2id with 64MB memory, 3 iterations
- **Token storage**: Refresh tokens stored in Redis
- **Rate limiting**: Redis-based login attempt tracking
- **No PII logging**: Email not logged, secrets masked

## 7. Integration

**Dependencies:**
- `common/cache` - Redis client
- `common/errors` - Error types
- `crypto` - Password hashing

**Used by:**
- Middleware for route protection
- Vault module for authorization checks

## 8. Roadmap

- [ ] MFA (TOTP) support
- [ ] Password reset flow
- [ ] Session management UI
- [ ] OAuth providers (Google, GitHub)
