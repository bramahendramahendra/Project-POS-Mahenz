import { useState, type ReactNode } from 'react'

import { useBreakpoint } from '@/shared/hooks'

import { Navbar } from './Navbar'
import { Sidebar } from './Sidebar'

// Menu sudah dijamin selesai dimuat oleh ProtectedRoute sebelum AppLayout dirender.
export function AppLayout({ children }: { children: ReactNode }) {
  const isDesktop = useBreakpoint('lg')
  const [sidebarOpen, setSidebarOpen] = useState(false)

  return (
    <div className="min-h-screen bg-[var(--color-bg)]">
      <Navbar onMenuClick={() => setSidebarOpen((v) => !v)} />
      <div style={{ marginTop: 'var(--navbar-height)', display: 'flex' }}>
        <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
        <main
          style={{
            marginLeft: isDesktop ? 'var(--sidebar-width)' : 0,
            flex: 1,
            minWidth: 0,
            padding: '24px',
            minHeight: 'calc(100vh - var(--navbar-height))',
          }}
        >
          {children}
        </main>
      </div>
    </div>
  )
}
