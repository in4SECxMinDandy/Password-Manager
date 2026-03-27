import { Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { useAuth } from '@/features/auth/context/AuthContext'

export function AppLayout() {
  const { user, logout } = useAuth()

  return (
    <div className="flex h-screen bg-background">
      <Sidebar user={user} onLogout={logout} />
      <main className="flex-1 overflow-auto">
        <Outlet />
      </main>
    </div>
  )
}
