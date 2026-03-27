# Vault Module (Go)

## 1. Purpose

Handles vault and entry management. Belongs to the **Application & Domain layers**.

## 2. Scope & Bounded Context

**Entities owned:**
- `Vault` - Container for entries, owned by user
- `VaultEntry` - Individual password entries, cards, notes, identities

**Referenced:**
- Auth module (authorization)
- Crypto module (data encryption)

## 3. Public API

### HTTP Endpoints

**Vaults:**
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/vaults` | List user's vaults |
| POST | `/api/v1/vaults` | Create vault |
| GET | `/api/v1/vaults/:id` | Get vault details |
| PUT | `/api/v1/vaults/:id` | Update vault |
| DELETE | `/api/v1/vaults/:id` | Delete vault |

**Entries:**
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/vaults/:id/entries` | List entries in vault |
| POST | `/api/v1/vaults/:id/entries` | Create entry |
| GET | `/api/v1/entries/:id` | Get entry details |
| PUT | `/api/v1/entries/:id` | Update entry |
| DELETE | `/api/v1/entries/:id` | Delete entry |
| POST | `/api/v1/entries/:id/favorite` | Toggle favorite |

### Services

```go
type VaultService interface {
    CreateVault(ctx context.Context, userID uuid.UUID, req *CreateVaultRequest) (*Vault, error)
    GetVault(ctx context.Context, userID, vaultID uuid.UUID) (*Vault, error)
    ListVaults(ctx context.Context, userID uuid.UUID) (*VaultListResponse, error)
    UpdateVault(ctx context.Context, userID, vaultID uuid.UUID, req *UpdateVaultRequest) (*Vault, error)
    DeleteVault(ctx context.Context, userID, vaultID uuid.UUID) error
    
    CreateEntry(ctx context.Context, userID, vaultID uuid.UUID, req *CreateEntryRequest) (*VaultEntry, error)
    GetEntry(ctx context.Context, userID, entryID uuid.UUID) (*VaultEntry, error)
    ListEntries(ctx context.Context, userID, vaultID uuid.UUID) (*EntryListResponse, error)
    UpdateEntry(ctx context.Context, userID, entryID uuid.UUID, req *UpdateEntryRequest) (*VaultEntry, error)
    DeleteEntry(ctx context.Context, userID, entryID uuid.UUID) error
    ToggleFavorite(ctx context.Context, userID, entryID uuid.UUID) (*VaultEntry, error)
}
```

## 4. Business Flows

### Create Entry Flow

```
1. Verify user owns the vault
2. Validate entry data (title required)
3. Create entry with current timestamp
4. Store entry in database
5. Return created entry
```

### Delete Entry Flow (Soft Delete)

```
1. Verify user owns the entry's vault
2. Set deleted_at timestamp (soft delete)
3. Entry retained for 30-day recovery
4. Return success
```

## 5. Business Rules

| Rule | Description |
|------|-------------|
| Vault ownership | Only owner can access vault |
| Entry ownership | Inherited from vault |
| Soft delete | 30-day retention |
| Entry types | login, note, card, identity |

## 6. Security & Privacy

- **Authorization**: All operations verify user owns the vault
- **Encryption**: Entry data encrypted client-side
- **Server envelope**: Additional encryption layer
- **No secret logging**: Entry data never logged

## 7. Integration

**Dependencies:**
- `common/database` - PostgreSQL
- `common/errors` - Error types
- `crypto` - Encryption utilities

**Used by:**
- Frontend for vault management
- Future modules (sharing, sync)

## 8. Roadmap

- [ ] Entry versioning
- [ ] Bulk operations
- [ ] Import/Export
- [ ] Entry templates
- [ ] Attachments support
