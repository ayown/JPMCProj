import { Link, NavLink } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'

function Layout({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth()
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <Link to="/" className="text-xl font-semibold text-primary-700">Fraud Detection</Link>
          <nav className="flex items-center gap-6">
            <NavLink to="/dashboard" className={({isActive})=>`hover:text-primary-700 ${isActive?'text-primary-700':'text-gray-700'}`}>Dashboard</NavLink>
            <NavLink to="/reports" className={({isActive})=>`hover:text-primary-700 ${isActive?'text-primary-700':'text-gray-700'}`}>Reports</NavLink>
            <NavLink to="/settings" className={({isActive})=>`hover:text-primary-700 ${isActive?'text-primary-700':'text-gray-700'}`}>Settings</NavLink>
            <button className="btn btn-secondary" onClick={logout}>Logout</button>
            <span className="text-sm text-gray-600">{user?.email}</span>
          </nav>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-6 py-8">
        {children}
      </main>
    </div>
  )
}

export default Layout
