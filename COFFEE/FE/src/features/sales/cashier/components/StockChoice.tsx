export type StockAction = 'return' | 'waste'

/**
 * Pilihan nasib stok barang jadi saat item dihapus, pesanan dibatalkan, atau transaksi di-void
 * (keputusan 100): kembali ke stok atau dibuang. Tombol segmen besar agar mudah diketuk di tablet.
 */
export function StockChoice({ value, onChange, label }: { value: StockAction; onChange: (v: StockAction) => void; label: string }) {
  return (
    <div className="flex overflow-hidden rounded-md border" role="radiogroup" aria-label={label}>
      {(['return', 'waste'] as const).map((a) => (
        <button
          key={a}
          type="button"
          role="radio"
          aria-checked={value === a}
          onClick={() => onChange(a)}
          className={`px-3 py-1.5 text-xs font-medium ${
            value === a ? (a === 'waste' ? 'bg-red-600 text-white' : 'bg-[#2c3e50] text-white') : 'bg-white text-gray-600 hover:bg-gray-50'
          }`}
        >
          {a === 'return' ? 'Kembali ke stok' : 'Dibuang'}
        </button>
      ))}
    </div>
  )
}
