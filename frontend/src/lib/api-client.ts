import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  VaultListResponse,
  Vault,
  EntryListResponse,
  VaultEntry,
  CreateVaultRequest,
  UpdateVaultRequest,
  CreateEntryRequest,
  UpdateEntryRequest,
  User,
} from './api-types'

const API_BASE = '/api/v1'

class ApiClient {
  private accessToken: string | null = null
  private refreshToken: string | null = null

  constructor() {
    if (typeof window !== 'undefined') {
      this.accessToken = localStorage.getItem('access_token')
      this.refreshToken = localStorage.getItem('refresh_token')
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    }

    if (this.accessToken) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${this.accessToken}`
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    })

    if (response.status === 401 && this.refreshToken) {
      const refreshed = await this.refresh()
      if (refreshed) {
        (headers as Record<string, string>)['Authorization'] = `Bearer ${this.accessToken}`
        const retryResponse = await fetch(`${API_BASE}${endpoint}`, {
          ...options,
          headers,
        })
        if (!retryResponse.ok) {
          throw await this.handleError(retryResponse)
        }
        return retryResponse.json()
      }
    }

    if (!response.ok) {
      throw await this.handleError(response)
    }

    return response.json()
  }

  private async handleError(response: Response): Promise<Error> {
    try {
      const data = await response.json()
      return new Error(data.error?.message || 'An error occurred')
    } catch {
      return new Error('An error occurred')
    }
  }

  private async refresh(): Promise<boolean> {
    if (!this.refreshToken) return false

    try {
      const response = await fetch(`${API_BASE}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: this.refreshToken }),
      })

      if (!response.ok) {
        this.logout()
        return false
      }

      const data: AuthResponse = await response.json()
      this.setTokens(data.access_token, data.refresh_token)
      return true
    } catch {
      return false
    }
  }

  setTokens(accessToken: string, refreshToken: string) {
    this.accessToken = accessToken
    this.refreshToken = refreshToken
    if (typeof window !== 'undefined') {
      localStorage.setItem('access_token', accessToken)
      localStorage.setItem('refresh_token', refreshToken)
    }
  }

  clearTokens() {
    this.accessToken = null
    this.refreshToken = null
    if (typeof window !== 'undefined') {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    }
  }

  isAuthenticated(): boolean {
    return !!this.accessToken
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    this.setTokens(response.access_token, response.refresh_token)
    return response
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    this.setTokens(response.access_token, response.refresh_token)
    return response
  }

  async logout(): Promise<void> {
    try {
      await this.request('/auth/logout', { method: 'POST' })
    } finally {
      this.clearTokens()
    }
  }

  async getMe(): Promise<User> {
    return this.request<User>('/auth/me')
  }

  async getVaults(): Promise<VaultListResponse> {
    return this.request<VaultListResponse>('/vaults')
  }

  async createVault(data: CreateVaultRequest): Promise<Vault> {
    return this.request<Vault>('/vaults', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getVault(id: string): Promise<Vault> {
    return this.request<Vault>(`/vaults/${id}`)
  }

  async updateVault(id: string, data: UpdateVaultRequest): Promise<Vault> {
    return this.request<Vault>(`/vaults/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteVault(id: string): Promise<void> {
    await this.request(`/vaults/${id}`, { method: 'DELETE' })
  }

  async getEntries(vaultId: string): Promise<EntryListResponse> {
    return this.request<EntryListResponse>(`/vaults/${vaultId}/entries`)
  }

  async createEntry(vaultId: string, data: CreateEntryRequest): Promise<VaultEntry> {
    return this.request<VaultEntry>(`/vaults/${vaultId}/entries`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getEntry(id: string): Promise<VaultEntry> {
    return this.request<VaultEntry>(`/entries/${id}`)
  }

  async updateEntry(id: string, data: UpdateEntryRequest): Promise<VaultEntry> {
    return this.request<VaultEntry>(`/entries/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteEntry(id: string): Promise<void> {
    await this.request(`/entries/${id}`, { method: 'DELETE' })
  }

  async toggleFavorite(id: string): Promise<VaultEntry> {
    return this.request<VaultEntry>(`/entries/${id}/favorite`, {
      method: 'POST',
    })
  }
}

export const apiClient = new ApiClient()
