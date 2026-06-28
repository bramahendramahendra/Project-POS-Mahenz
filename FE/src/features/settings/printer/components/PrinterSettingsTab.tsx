import { useEffect, useState } from 'react'
import { Printer } from 'lucide-react'

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

import { usePrinterSettingsQuery, useUpdatePrinterSettingsMutation } from '../../settings.api'
import type { PrinterSettings } from '../../settings.types'

const DEFAULT_SETTINGS: PrinterSettings = {
  paper_size: '80mm',
  receipt_header: '',
  receipt_footer: 'Terima kasih telah berbelanja!',
  show_logo: false,
  auto_print: false,
}

function buildReceiptHtml(settings: PrinterSettings, forPrint = false): string {
  const bodyWidth = settings.paper_size === '58mm' ? '210px' : '290px'
  const now = new Date().toLocaleString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })
  const printScript = forPrint
    ? `<script>window.onload=function(){window.print();window.onafterprint=function(){window.close()}}<\/script>`
    : ''
  return `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8" />
  <title>Test Struk</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: 'Courier New', monospace;
      font-size: 12px;
      color: #111;
      width: ${bodyWidth};
      margin: 0 auto;
      padding: 16px 12px;
    }
    .center { text-align: center; }
    .store-name { font-size: 15px; font-weight: 700; letter-spacing: 0.5px; }
    .divider { border: none; border-top: 1px dashed #aaa; margin: 10px 0; }
    .divider-solid { border: none; border-top: 1px solid #ccc; margin: 8px 0; }
    .row { display: flex; justify-content: space-between; margin-bottom: 3px; }
    .label { color: #777; }
    .muted { color: #777; font-size: 11px; }
    .bold { font-weight: 700; }
    .item { margin-bottom: 8px; }
    .item-name { font-weight: 600; }
    .total-row { display: flex; justify-content: space-between; font-size: 14px; font-weight: 700; margin: 4px 0; }
    .kembalian { display: flex; justify-content: space-between; font-weight: 600; color: #16a34a; margin-top: 2px; }
    .footer { text-align: center; color: #888; font-size: 11px; margin-top: 8px; }
    @media print { body { width: 100%; } }
  </style>
</head>
<body>
  <div class="center">
    ${settings.show_logo ? '<div class="muted">[LOGO]</div>' : ''}
    <div class="store-name">${settings.receipt_header || 'Nama Toko'}</div>
  </div>

  <hr class="divider" />

  <div class="row"><span class="label">No. Transaksi</span><span class="bold">TRX-TEST-001</span></div>
  <div class="row"><span class="label">Tanggal</span><span>${now}</span></div>
  <div class="row"><span class="label">Kasir</span><span>Test Kasir</span></div>
  <div class="row"><span class="label">Pembayaran</span><span>Tunai</span></div>

  <hr class="divider" />

  <div class="item">
    <div class="row"><span class="item-name">Produk A</span><span class="bold">Rp 20.000</span></div>
    <div class="row muted"><span>pcs &times; 2 @ Rp 10.000</span></div>
  </div>
  <div class="item">
    <div class="row"><span class="item-name">Produk B</span><span class="bold">Rp 15.000</span></div>
    <div class="row muted"><span>pcs &times; 1 @ Rp 15.000</span></div>
  </div>

  <hr class="divider" />

  <div class="row muted"><span>Subtotal</span><span>Rp 35.000</span></div>

  <hr class="divider-solid" />

  <div class="total-row"><span>TOTAL</span><span>Rp 35.000</span></div>
  <div class="row muted"><span>Dibayar (Tunai)</span><span>Rp 50.000</span></div>
  <div class="kembalian"><span>Kembalian</span><span>Rp 15.000</span></div>

  <hr class="divider" />

  <div class="footer">${settings.receipt_footer || ''}</div>

  ${printScript}
</body>
</html>`
}

function openTestPrint(settings: PrinterSettings) {
  const win = window.open('', '_blank', 'width=400,height=600')
  if (!win) return
  win.document.write(buildReceiptHtml(settings, true))
  win.document.close()
}

function ReceiptPreview({ settings }: { settings: PrinterSettings }) {
  const previewWidth = settings.paper_size === '58mm' ? '230px' : '310px'
  return (
    <iframe
      srcDoc={buildReceiptHtml(settings, false)}
      title="Preview Struk"
      style={{ width: previewWidth, height: '420px', border: 'none' }}
      scrolling="no"
    />
  )
}

export function PrinterSettingsTab() {
  const { data, isLoading } = usePrinterSettingsQuery()
  const { mutate: save, isPending } = useUpdatePrinterSettingsMutation()
  const [form, setForm] = useState<PrinterSettings>(DEFAULT_SETTINGS)

  useEffect(() => {
    if (data) setForm(data)
  }, [data])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    save(form)
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
    <div className="flex gap-8 items-start">
      <form onSubmit={handleSubmit} className="w-full max-w-sm space-y-6 shrink-0">
        <div className="rounded-lg border bg-white p-5 space-y-4">
          <h3 className="text-sm font-semibold text-gray-700">Konfigurasi Struk</h3>

          <div className="space-y-1.5">
            <Label>Ukuran Kertas</Label>
            <Select
              value={form.paper_size}
              onValueChange={(v) => setForm((f) => ({ ...f, paper_size: v as '58mm' | '80mm' }))}
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
              value={form.receipt_header}
              onChange={(e) => setForm((f) => ({ ...f, receipt_header: e.target.value }))}
              placeholder="Nama toko atau teks header"
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="pr-footer">Footer Struk</Label>
            <Textarea
              id="pr-footer"
              value={form.receipt_footer}
              onChange={(e) => setForm((f) => ({ ...f, receipt_footer: e.target.value }))}
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
              checked={form.show_logo}
              onCheckedChange={(v) => setForm((f) => ({ ...f, show_logo: v }))}
            />
          </div>

          <div className="flex items-center justify-between py-3">
            <div>
              <p className="text-sm font-medium text-gray-800">Auto Print</p>
              <p className="text-xs text-gray-500 mt-0.5">Langsung cetak struk setelah transaksi selesai</p>
            </div>
            <Switch
              checked={form.auto_print}
              onCheckedChange={(v) => setForm((f) => ({ ...f, auto_print: v }))}
            />
          </div>
        </div>

        <div className="flex gap-3">
          <Button
            type="button"
            variant="outline"
            className="gap-1.5"
            onClick={() => openTestPrint(form)}
          >
            <Printer size={15} />
            Test Print
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? 'Menyimpan...' : 'Simpan'}
          </Button>
        </div>
      </form>

      <div className="hidden lg:block shrink-0">
        <p className="text-xs font-medium text-gray-500 mb-2">Preview Struk</p>
        <div className="rounded-lg border bg-gray-50 p-4 shadow-sm">
          <ReceiptPreview settings={form} />
        </div>
      </div>
    </div>
  )
}
