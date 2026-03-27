export interface User {
  id: string
  email: string
  mfa_enabled: boolean
  created_at: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: User
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
}

export interface Vault {
  id: string
  user_id: string
  name: string
  created_at: string
  updated_at: string
}

export interface VaultEntry {
  id: string
  vault_id: string
  type: 'login' | 'note' | 'card' | 'identity'
  title: string
  data: string
  favorite: boolean
  deleted_at: string | null
  created_at: string
  updated_at: string
}

export interface VaultListResponse {
  vaults: Vault[]
  total: number
}

export interface EntryListResponse {
  entries: VaultEntry[]
  total: number
}

export interface CreateVaultRequest {
  name: string
}

export interface UpdateVaultRequest {
  name: string
}

export interface CreateEntryRequest {
  type: VaultEntry['type']
  title: string
  data: string
  favorite?: boolean
}

export interface UpdateEntryRequest {
  title?: string
  data?: string
  favorite?: boolean
}

export interface ApiError {
  error: {
    code: string
    message: string
  }
}
