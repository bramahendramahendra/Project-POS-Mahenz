import { PageHeader } from '@/shared/components'

import { DirectLoginCard } from './components/DirectLoginCard'
import { UserTable } from './components/UserTable'

// Akun baru dan password diatur di Portal (keputusan 69). User yang pertama kali masuk ke Toko
// muncul di sini dengan role "Menunggu Akses"; admin Toko mengatur role-nya.
export function UserManagementPage() {
  return (
    <div className="space-y-4">
      <PageHeader
        title="Manajemen User"
        description="Akun baru dan password diatur di Portal. User yang pertama kali masuk ke Toko muncul di sini dengan role Menunggu Akses."
        breadcrumbs={[{ label: 'Sistem' }, { label: 'Manajemen User' }]}
      />
      <DirectLoginCard />
      <UserTable />
    </div>
  )
}
