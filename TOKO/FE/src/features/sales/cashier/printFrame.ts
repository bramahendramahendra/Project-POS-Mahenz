/**
 * Cetak HTML struk lewat IFRAME tersembunyi — deterministik dan tidak terpengaruh
 * oleh layout aplikasi (#root, modal Radix, @media print global). Hanya isi iframe
 * yang dikirim ke printer, jadi tidak ada halaman kosong tambahan maupun struk yang
 * hilang. Dipakai sebagai jalur "Browser / Dialog Print".
 *
 * @param innerHtml  HTML isi struk (tanpa <html>/<body>).
 * @param paperSize  '58mm' | '80mm' — menentukan lebar & @page.
 */
export function printReceiptHtml(innerHtml: string, paperSize: '58mm' | '80mm'): void {
  // Lebar AREA CETAK efektif, bukan lebar kertas fisik. Printer 58mm (spt EPPOS
  // EP58M) hanya bisa mencetak ~48mm; sisanya margin fisik. Memakai 58mm penuh
  // membuat teks melewati batas kanan dan terpotong. 80mm efektif ~72mm.
  const printWidth = paperSize === '58mm' ? '48mm' : '72mm'
  const pageWidth = paperSize === '58mm' ? '58mm' : '80mm'

  const iframe = document.createElement('iframe')
  iframe.setAttribute('aria-hidden', 'true')
  iframe.style.position = 'fixed'
  iframe.style.right = '0'
  iframe.style.bottom = '0'
  iframe.style.width = '0'
  iframe.style.height = '0'
  iframe.style.border = '0'
  iframe.style.visibility = 'hidden'
  document.body.appendChild(iframe)

  const doc = iframe.contentWindow?.document
  if (!doc) {
    iframe.remove()
    return
  }

  doc.open()
  doc.write(`<!doctype html>
<html>
<head>
<meta charset="utf-8" />
<style>
  @page { size: ${pageWidth} auto; margin: 0; }
  * { box-sizing: border-box; }
  html, body { margin: 0; padding: 0; background: #fff; }
  /* Thermal 203dpi: teks HITAM PEKAT & tebal agar tajam. Lebar dibatasi ke area
     cetak efektif (printWidth) supaya tidak ada tulisan yang terpotong di kanan.
     Semua elemen boleh membungkus — tidak ada yang dipaksa satu baris. */
  body {
    width: ${printWidth};
    max-width: ${printWidth};
    margin: 0;                     /* rata kiri — mulai dari tepi area cetak */
    padding: 0 1mm;
    font-family: 'Courier New', 'Consolas', ui-monospace, monospace;
    font-size: 12px;               /* pas untuk 48mm — lebih besar bikin terpotong */
    font-weight: 700;              /* semua teks tebal — paling terbaca di thermal */
    line-height: 1.35;
    color: #000;
    -webkit-font-smoothing: none;  /* matikan anti-alias halus yang bikin tepi abu */
    text-rendering: geometricPrecision;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
    word-break: break-word;
    overflow-wrap: anywhere;
  }
  .rcpt { width: 100%; }
  .rcpt * {
    color: #000 !important;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
  .center { text-align: center; }
  .bold { font-weight: 900; }
  .store-name { font-size: 14px; font-weight: 900; }
  /* Baris nilai: label kiri, nominal kanan. Boleh membungkus bila sempit. */
  .row { display: flex; justify-content: space-between; gap: 4px; align-items: baseline; }
  .row > span:first-child { flex: 1 1 auto; overflow-wrap: anywhere; }
  .row > span:last-child { flex: 0 0 auto; text-align: right; white-space: nowrap; }
  /* Baris info transaksi: label lalu nilai. Nilai panjang (kode trx/tanggal)
     membungkus rapi tanpa memotong karakter. */
  .info { margin: 1px 0; }
  .info .lbl { font-weight: 700; }
  .info .val { font-weight: 900; overflow-wrap: anywhere; }
  hr { border: none; border-top: 1px solid #000; margin: 4px 0; }
  .sep { border-top: 1px solid #000; margin: 4px 0; }
  .item-name { font-size: 12px; font-weight: 900; overflow-wrap: anywhere; }
  .item-sub { font-size: 11px; font-weight: 700; }
  .total { font-size: 15px; font-weight: 900; }
  img { max-width: 100%; }
</style>
</head>
<body>
  <div class="rcpt">${innerHtml}</div>
</body>
</html>`)
  doc.close()

  const triggerPrint = () => {
    const win = iframe.contentWindow
    if (!win) {
      iframe.remove()
      return
    }
    win.focus()
    win.print()
    // Beri jeda agar dialog print sempat membaca dokumen sebelum iframe dibuang.
    setTimeout(() => iframe.remove(), 1000)
  }

  // Tunggu konten (termasuk gambar logo) selesai dimuat sebelum mencetak.
  if (iframe.contentWindow?.document.readyState === 'complete') {
    setTimeout(triggerPrint, 50)
  } else {
    iframe.onload = () => setTimeout(triggerPrint, 50)
  }
}
