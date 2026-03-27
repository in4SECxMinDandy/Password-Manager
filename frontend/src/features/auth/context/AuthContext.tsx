import React, { createContext, useContext, useState, useEffect, useCallback } from 'react'
import { apiClient } from '@/lib/api-client'
import type { User } from '@/lib/api-types'

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const checkAuth = useCallback(async () => {
    if (!apiClient.isAuthenticated()) {
      setIsLoading(false)
      return
    }

    try {
      const userData = await apiClient.getMe()
      setUser(userData)
    } catch {
      apiClient.clearTokens()
      setUser(null)
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    checkAuth()
  }, [checkAuth])

  const login = async (email: string, password: string) => {
    const response = await apiClient.login({ email, password })
    setUser(response.user)
  }

  const register = async (email: string, password: string) => {
    const response = await apiClient.register({ email, password })
    setUser(response.user)
  }

  const logout = async () => {
    await apiClient.logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
