import { useEffect } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Printer } from 'lucide-react'

import { RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Switch } from '@/shared/components/ui/switch'
import { Textarea } from '@/shared/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import { usePrinterSettingsQuery, useUpdatePrinterSettingsMutation } from '../printer.api'
import type { PrinterSettings } from '../printer.types'
import { printerSettingsSchema, type PrinterSettingsFormValues } from '../printer.schema'

const DEFAULT_SETTINGS: PrinterSettingsFormValues = {
  paper_size: '80mm',
  receipt_header: '',
  receipt_footer: 'Terima kasih telah berbelanja!',
  show_logo: false,
  auto_print: false,
}

function ReceiptPreview({ settings }: { settings: PrinterSettings }) {
  const previewWidth = settings.paper_size === '58mm' ? '210px' : '290px'
  const now = new Date().toLocaleString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })

  return (
    <div
      className="print-root"
      style={{
        width: previewWidth,
        margin: '0 auto',
        fontFamily: "'Courier New', monospace",
        fontSize: 12,
        color: '#111',
        padding: '16px 12px',
      }}
    >
      <div style={{ textAlign: 'center' }}>
        {settings.show_logo && <div style={{ color: '#777', fontSize: 11 }}>[LOGO]</div>}
        <div style={{ fontSize: 15, fontWeight: 700, letterSpacing: '0.5px' }}>
          {settings.receipt_header || 'Nama Toko'}
        </div>
      </div>

      <hr style={{ border: 'none', borderTop: '1px dashed #aaa', margin: '10px 0' }} />

      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 3 }}>
        <span style={{ color: '#777' }}>No. Transaksi</span>
        <span style={{ fontWeight: 700 }}>TRX-TEST-001</span>
      </div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 3 }}>
        <span style={{ color: '#777' }}>Tanggal</span>
        <span>{now}</span>
      </div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 3 }}>
        <span style={{ color: '#777' }}>Kasir</span>
        <span>Test Kasir</span>
      </div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 3 }}>
        <span style={{ color: '#777' }}>Pembayaran</span>
        <span>Tunai</span>
      </div>

      <hr style={{ border: 'none', borderTop: '1px dashed #aaa', margin: '10px 0' }} />

      <div style={{ marginBottom: 8 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 3 }}>
          <span style={{ fontWeight: 600 }}>Produk A</span>
          <span style={{ fontWeight: 700 }}>Rp 20.000</span>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', color: '#777', fontSize: 11 }}>
          <span>pcs &times; 2 @ Rp 10.000</span>
        </div>
      </div>
      <div style={{ marginBottom: 8 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 3 }}>
          <span style={{ fontWeight: 600 }}>Produk B</span>
          <span style={{ fontWeight: 700 }}>Rp 15.000</span>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', color: '#777', fontSize: 11 }}>
          <span>pcs &times; 1 @ Rp 15.000</span>
        </div>
      </div>

      <hr style={{ border: 'none', borderTop: '1px dashed #aaa', margin: '10px 0' }} />

      <div style={{ display: 'flex', justifyContent: 'space-between', color: '#777', fontSize: 11 }}>
        <span>Subtotal</span>
        <span>Rp 35.000</span>
      </div>

      <hr style={{ border: 'none', borderTop: '1px solid #ccc', margin: '8px 0' }} />

      <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 14, fontWeight: 700, margin: '4px 0' }}>
        <span>TOTAL</span>
        <span>Rp 35.000</span>
      </div>
      <div style={{ display: 'flex', justifyContent: 'space-between', color: '#777', fontSize: 11 }}>
        <span>Dibayar (Tunai)</span>
        <span>Rp 50.000</span>
      </div>
      <div style={{ display: 'flex', justifyContent: 'space-between', fontWeight: 600, color: '#16a34a', marginTop: 2 }}>
        <span>Kembalian</span>
        <span>Rp 15.000</span>
      </div>

      <hr style={{ border: 'none', borderTop: '1px dashed #aaa', margin: '10px 0' }} />

      <div style={{ textAlign: 'center', color: '#888', fontSize: 11, marginTop: 8 }}>
        {settings.receipt_footer || ''}
      </div>
    </div>
  )
}

export function PrinterSettingsTab() {
  const { data, isLoading } = usePrinterSettingsQuery()
  const { mutate: save, isPending } = useUpdatePrinterSettingsMutation()

  const {
    register,
    handleSubmit,
    reset,
    control,
    setValue,
  } = useForm<PrinterSettingsFormValues>({
    resolver: zodResolver(printerSettingsSchema),
    defaultValues: DEFAULT_SETTINGS,
  })

  useEffect(() => {
    if (data) reset(data)
  }, [data, reset])

  const liveSettings = useWatch({ control }) as PrinterSettings

  const onSubmit = (values: PrinterSettingsFormValues) => {
    save(values)
  }

  if (isLoading) {
    return (
      <div className="space-y-3 max-w-lg">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="h-10 animate-pulse rounded bg-gray-100" />
        ))}
      </div>
    )
  }

  return (
    <div className="flex flex-col lg:flex-row gap-8 lg:items-start">
      <form onSubmit={handleSubmit(onSubmit)} className="w-full max-w-sm space-y-6 shrink-0">
        <div className="rounded-lg border bg-white p-5 space-y-4">
          <h3 className="text-sm font-semibold text-gray-700">Konfigurasi Struk</h3>

          <div className="space-y-1.5">
            <Label>Ukuran Kertas</Label>
            <Select
              value={liveSettings.paper_size}
              onValueChange={(v) => {
                if (v) setValue('paper_size', v as '58mm' | '80mm')
              }}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="58mm">58mm</SelectItem>
                <SelectItem value="80mm">80mm</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="pr-header">Header Struk</Label>
            <Input
              id="pr-header"
              {...register('receipt_header')}
              placeholder="Nama toko atau teks header"
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="pr-footer">Footer Struk</Label>
            <Textarea
              id="pr-footer"
              {...register('receipt_footer')}
              placeholder="Teks bawah struk..."
              className="resize-none"
              rows={2}
            />
          </div>
        </div>

        <div className="rounded-lg border bg-white px-5 py-2">
          <h3 className="text-sm font-semibold text-gray-700 py-3 border-b border-gray-100">
            Preferensi
          </h3>

          <div className="flex items-center justify-between py-3 border-b border-gray-100">
            <div>
              <p className="text-sm font-medium text-gray-800">Tampilkan Logo</p>
              <p className="text-xs text-gray-500 mt-0.5">Tampilkan logo toko di bagian atas struk</p>
            </div>
            <Switch
              checked={liveSettings.show_logo}
              onCheckedChange={(v) => setValue('show_logo', v)}
            />
          </div>

          <div className="flex items-center justify-between py-3">
            <div>
              <p className="text-sm font-medium text-gray-800">Auto Print</p>
              <p className="text-xs text-gray-500 mt-0.5">Langsung cetak struk setelah transaksi selesai</p>
            </div>
            <Switch
              checked={liveSettings.auto_print}
              onCheckedChange={(v) => setValue('auto_print', v)}
            />
          </div>
        </div>

        <div className="flex gap-3">
          <Button
            type="button"
            variant="outline"
            className="gap-1.5"
            onClick={() => window.print()}
          >
            <Printer size={15} />
            Test Print
          </Button>
          <RoleGuard menuKey="sistem.printer" action="can_edit">
            <Button type="submit" disabled={isPending}>
              {isPending ? 'Menyimpan...' : 'Simpan'}
            </Button>
          </RoleGuard>
        </div>
      </form>

      <div className="shrink-0">
        <p className="text-xs font-medium text-gray-500 mb-2 no-print">Preview Struk</p>
        <div className="rounded-lg border bg-gray-50 p-4 shadow-sm">
          <ReceiptPreview settings={liveSettings} />
        </div>
      </div>
    </div>
  )
}
