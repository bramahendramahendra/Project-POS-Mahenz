import { useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { PageHeader, ConfirmDialog } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/components/ui/card'
import { useShiftOptionsQuery } from '@/features/operational/shifts'
import { formatRupiah } from '@/shared/utils/currency'
import { CalendarClock, DollarSign, User } from 'lucide-react'

import {
  useBackdateCashDrawerCurrentQuery,
  useCloseBackdateCashDrawerMutation,
  useOpenBackdateCashDrawerMutation,
  useUserOptionsQuery,
} from './backdate-cash.api'
import {
  closeBackdateCashDrawerSchema,
  openBackdateCashDrawerSchema,
  type CloseBackdateCashDrawerFormValues,
  type OpenBackdateCashDrawerFormValues,
} from './backdate-cash.schema'

export function BackdateCashPage() {
  const { data: currentCash, isLoading } = useBackdateCashDrawerCurrentQuery()
  const hasOpenCash = currentCash != null && currentCash.id > 0

  return (
    <div className="space-y-4">
      <PageHeader
        title="Kas Historis"
        breadcrumbs={[{ label: 'Historis' }, { label: 'Kas Historis' }]}
      />

      {isLoading ? (
        <div className="flex items-center justify-center py-12 text-sm text-gray-400">Memuat...</div>
      ) : hasOpenCash ? (
        <ActiveCashView data={currentCash} />
      ) : (
        <OpenCashForm />
      )}
    </div>
  )
}

// ===================== FORM BUKA KAS =====================

function OpenCashForm() {
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<OpenBackdateCashDrawerFormValues | null>(null)

  const { mutate: openCash, isPending } = useOpenBackdateCashDrawerMutation()
  const { data: users } = useUserOptionsQuery()
  const { data: shiftOptions } = useShiftOptionsQuery()

  // Batas maksimal tanggal = kemarin (backdate hanya untuk tanggal lampau).
  // Dihitung sekali lewat lazy initializer useState — tempat yang tepat untuk
  // pembacaan waktu (fungsi impur), sehingga tidak dipanggil setiap render.
  const [maxDate] = useState(() => new Date(Date.now() - 86400000).toISOString().split('T')[0])

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<OpenBackdateCashDrawerFormValues>({
    resolver: zodResolver(openBackdateCashDrawerSchema),
    defaultValues: {
      user_id: 0,
      date: '',
      shift_id: null,
      opening_balance: 0,
      notes: '',
    },
  })

  const onSubmit = (values: OpenBackdateCashDrawerFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirm = () => {
    if (!pendingValues) return
    openCash(
      {
        user_id: pendingValues.user_id,
        date: pendingValues.date,
        shift_id: pendingValues.shift_id || null,
        opening_balance: pendingValues.opening_balance,
        notes: pendingValues.notes || undefined,
      },
      {
        onSuccess: () => {
          setIsConfirming(false)
          setPendingValues(null)
        },
      },
    )
  }

  const selectedUser = users?.find((u) => u.id === pendingValues?.user_id)

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle>Buka Kas Historis</CardTitle>
          <CardDescription>
            Buka kas untuk tanggal lampau. Pilih kasir yang bertugas pada tanggal tersebut.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 max-w-md">
            {/* Kasir */}
            <div className="space-y-1.5">
              <Label>
                Kasir <span className="text-red-500">*</span>
              </Label>
              <Controller
                name="user_id"
                control={control}
                render={({ field }) => (
                  <Select
                    value={field.value ? String(field.value) : ''}
                    onValueChange={(v) => field.onChange(Number(v))}
                  >
                    <SelectTrigger className={errors.user_id ? 'border-red-500' : ''}>
                      <SelectValue placeholder="Pilih kasir..." />
                    </SelectTrigger>
                    <SelectContent>
                      {(users ?? []).map((u) => (
                        <SelectItem key={u.id} value={String(u.id)}>
                          {u.full_name} ({u.username})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.user_id && <p className="text-xs text-red-500">{errors.user_id.message}</p>}
            </div>

            {/* Tanggal */}
            <div className="space-y-1.5">
              <Label>
                Tanggal <span className="text-red-500">*</span>
              </Label>
              <Input
                type="date"
                max={maxDate}
                {...register('date')}
                className={errors.date ? 'border-red-500' : ''}
              />
              {errors.date && <p className="text-xs text-red-500">{errors.date.message}</p>}
            </div>

            {/* Shift */}
            <div className="space-y-1.5">
              <Label>Shift (opsional)</Label>
              <Controller
                name="shift_id"
                control={control}
                render={({ field }) => (
                  <Select
                    value={field.value ? String(field.value) : ''}
                    onValueChange={(v) => field.onChange(v ? Number(v) : null)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Pilih shift..." />
                    </SelectTrigger>
                    <SelectContent>
                      {(shiftOptions ?? []).map((s) => (
                        <SelectItem key={s.id} value={String(s.id)}>
                          {s.name} ({s.start_time} – {s.end_time})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
            </div>

            {/* Saldo Awal */}
            <div className="space-y-1.5">
              <Label>
                Saldo Awal (Rp) <span className="text-red-500">*</span>
              </Label>
              <Controller
                name="opening_balance"
                control={control}
                render={({ field }) => (
                  <RupiahInput
                    placeholder="0"
                    value={field.value}
                    onChange={field.onChange}
                    className={errors.opening_balance ? 'border-red-500' : ''}
                  />
                )}
              />
              {errors.opening_balance && (
                <p className="text-xs text-red-500">{errors.opening_balance.message}</p>
              )}
            </div>

            {/* Catatan */}
            <div className="space-y-1.5">
              <Label>Catatan (opsional)</Label>
              <Input {...register('notes')} placeholder="Catatan pembukaan kas historis..." />
            </div>

            <Button type="submit" className="w-full">
              Buka Kas Historis
            </Button>
          </form>
        </CardContent>
      </Card>

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => {
          if (!val) {
            setIsConfirming(false)
            setPendingValues(null)
          }
        }}
        title="Buka Kas Historis"
        description={`Buka kas historis untuk kasir "${selectedUser?.full_name ?? ''}" pada tanggal ${pendingValues?.date ?? ''}?`}
        confirmLabel="Ya, Buka Kas"
        isLoading={isPending}
        onConfirm={handleConfirm}
      />
    </>
  )
}

// ===================== VIEW KAS AKTIF =====================

function ActiveCashView({ data }: { data: NonNullable<ReturnType<typeof useBackdateCashDrawerCurrentQuery>['data']> }) {
  const [isClosing, setIsClosing] = useState(false)

  return (
    <div className="space-y-4">
      {/* Info Card */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <CalendarClock className="h-5 w-5 text-blue-500" />
            Kas Historis Aktif
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <InfoItem icon={<User className="h-4 w-4" />} label="Kasir" value={data.user_name} />
            <InfoItem
              icon={<CalendarClock className="h-4 w-4" />}
              label="Tanggal"
              value={new Date(data.open_time).toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
            />
            <InfoItem icon={<DollarSign className="h-4 w-4" />} label="Saldo Awal" value={formatRupiah(data.opening_balance)} />
            <InfoItem icon={<DollarSign className="h-4 w-4" />} label="Total Penjualan" value={formatRupiah(data.total_sales)} />
          </div>
          {data.open_notes && (
            <p className="mt-3 text-sm text-gray-500">Catatan: {data.open_notes}</p>
          )}
          <p className="mt-2 text-xs text-gray-400">
            Dibuka oleh: {data.created_by_name}
          </p>
        </CardContent>
      </Card>

      {/* Actions */}
      <div className="flex gap-2">
        <Button variant="destructive" onClick={() => setIsClosing(true)}>
          Tutup Kas Historis
        </Button>
      </div>

      {/* Close Form */}
      {isClosing && <CloseCashForm cashId={data.id} expectedBalance={data.expected_balance} onCancel={() => setIsClosing(false)} />}
    </div>
  )
}

// ===================== FORM TUTUP KAS =====================

function CloseCashForm({ cashId, expectedBalance, onCancel }: { cashId: number; expectedBalance: number; onCancel: () => void }) {
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<CloseBackdateCashDrawerFormValues | null>(null)

  const { mutate: closeCash, isPending } = useCloseBackdateCashDrawerMutation()

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<CloseBackdateCashDrawerFormValues>({
    resolver: zodResolver(closeBackdateCashDrawerSchema),
    defaultValues: { closing_balance: 0, notes: '' },
  })

  const onSubmit = (values: CloseBackdateCashDrawerFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirm = () => {
    if (!pendingValues) return
    closeCash(
      { id: cashId, closing_balance: pendingValues.closing_balance, notes: pendingValues.notes },
      { onSuccess: () => { setIsConfirming(false); onCancel() } },
    )
  }

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle>Tutup Kas Historis</CardTitle>
          <CardDescription>Saldo yang diharapkan: {formatRupiah(expectedBalance)}</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 max-w-md">
            <div className="space-y-1.5">
              <Label>
                Saldo Akhir (Rp) <span className="text-red-500">*</span>
              </Label>
              <Controller
                name="closing_balance"
                control={control}
                render={({ field }) => (
                  <RupiahInput
                    placeholder="0"
                    value={field.value}
                    onChange={field.onChange}
                    className={errors.closing_balance ? 'border-red-500' : ''}
                  />
                )}
              />
              {errors.closing_balance && (
                <p className="text-xs text-red-500">{errors.closing_balance.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <Label>Catatan (opsional)</Label>
              <Input {...register('notes')} placeholder="Catatan penutupan kas..." />
            </div>

            <div className="flex gap-2">
              <Button type="submit" variant="destructive">Tutup Kas</Button>
              <Button type="button" variant="outline" onClick={onCancel}>Batal</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => { if (!val) setIsConfirming(false) }}
        title="Tutup Kas Historis"
        description="Yakin ingin menutup kas historis ini?"
        confirmLabel="Ya, Tutup"
        isLoading={isPending}
        onConfirm={handleConfirm}
      />
    </>
  )
}

// ===================== HELPER COMPONENTS =====================

function InfoItem({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div className="flex items-start gap-2">
      <span className="mt-0.5 text-gray-400">{icon}</span>
      <div>
        <p className="text-xs text-gray-500">{label}</p>
        <p className="text-sm font-medium">{value}</p>
      </div>
    </div>
  )
}
