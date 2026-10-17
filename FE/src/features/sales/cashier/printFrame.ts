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
  const widthCss = paperSize === '58mm' ? '58mm' : '80mm'

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
  @page { size: ${widthCss} auto; margin: 2mm; }
  * { box-sizing: border-box; }
  html, body { margin: 0; padding: 0; background: #fff; color: #000; }
  body {
    width: ${widthCss};
    font-family: 'Courier New', ui-monospace, monospace;
    font-size: 12px;
    line-height: 1.35;
    color: #000;
  }
  .rcpt { width: 100%; padding: 0 2mm; }
  .rcpt * { color: #000 !important; }
  .center { text-align: center; }
  .bold { font-weight: 700; }
  .row { display: flex; justify-content: space-between; gap: 8px; }
  .muted { color: #000; }
  hr { border: none; border-top: 1px dashed #000; margin: 6px 0; }
  .sep { border-top: 1px solid #000; margin: 6px 0; }
  .item-sub { font-size: 11px; padding-left: 2px; }
  .total { font-size: 14px; font-weight: 700; }
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
