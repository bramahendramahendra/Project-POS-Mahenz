package service

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"pos_api/config"
	"pos_api/domain/print/dto"
	"pos_api/errors"
)

type PrintServiceInterface interface {
	PrintReceipt(req *dto.PrintReceiptRequest) error
}

type printService struct{}

func NewPrintService() PrintServiceInterface {
	return &printService{}
}

// Perintah ESC/POS dasar.
var (
	escInit        = []byte{0x1B, 0x40}       // ESC @  — reset
	escAlignLeft   = []byte{0x1B, 0x61, 0x00} // ESC a 0
	escAlignCenter = []byte{0x1B, 0x61, 0x01} // ESC a 1
	escBoldOn      = []byte{0x1B, 0x45, 0x01} // ESC E 1
	escBoldOff     = []byte{0x1B, 0x45, 0x00} // ESC E 0
	escFeed3       = []byte{0x1B, 0x64, 0x03} // ESC d 3 — feed 3 baris
	cutPaper       = []byte{0x1D, 0x56, 0x00} // GS V 0  — potong (diabaikan bila tak ada cutter)
)

// PrintReceipt menyusun struk sebagai ESC/POS lalu menulisnya langsung ke COM
// port printer. Karena memakai font internal printer, hasilnya tajam dan lebar
// otomatis pas (tidak terpotong seperti lewat driver/dialog browser).
func (s *printService) PrintReceipt(req *dto.PrintReceiptRequest) error {
	port := config.Cfg.Printer.Port
	if strings.TrimSpace(port) == "" {
		return &errors.BadRequestError{Message: "Port printer belum dikonfigurasi (Printer.Port kosong)."}
	}
	width := config.Cfg.Printer.CharPerLine
	if width <= 0 {
		width = 32
	}

	payload := buildEscPos(req, width)

	if err := writeToSerial(port, payload); err != nil {
		return &errors.BadRequestError{
			Message: fmt.Sprintf("Gagal mencetak ke printer di %s: %v. Pastikan printer menyala & terhubung.", port, err),
		}
	}
	return nil
}

// buildEscPos menyusun seluruh byte struk.
func buildEscPos(req *dto.PrintReceiptRequest, width int) []byte {
	var b bytes.Buffer
	line := strings.Repeat("-", width) + "\n"

	b.Write(escInit)

	// Header toko — center, nama tebal.
	b.Write(escAlignCenter)
	b.Write(escBoldOn)
	b.WriteString(wrap(req.StoreName, width))
	b.Write(escBoldOff)
	if strings.TrimSpace(req.StoreSub) != "" {
		b.WriteString(wrap(req.StoreSub, width))
	}
	b.Write(escAlignLeft)

	b.WriteString(line)
	b.WriteString(twoCol("No", req.TransactionCode, width))
	b.WriteString(twoCol("Tgl", formatDate(req.TransactionDate), width))
	if strings.TrimSpace(req.CustomerName) != "" {
		b.WriteString(twoCol("Plg", req.CustomerName, width))
	}
	b.WriteString(twoCol("Bayar", paymentLabel(req.PaymentMethod), width))
	b.WriteString(line)

	// Item: nama di baris sendiri, lalu "qty x harga" kiri & subtotal kanan.
	for _, it := range req.Items {
		b.WriteString(wrap(it.Name, width))
		left := fmt.Sprintf("  %s x%s @ %s", it.UnitName, trimNum(it.Qty), rupiah(it.Price))
		b.WriteString(twoCol(left, rupiah(it.Subtotal), width))
	}
	b.WriteString(line)

	b.WriteString(twoCol("Subtotal", rupiah(req.Subtotal), width))
	if req.DiscountAmount > 0 {
		b.WriteString(twoCol("Diskon", "-"+rupiah(req.DiscountAmount), width))
	}
	if req.TaxAmount > 0 {
		b.WriteString(twoCol(fmt.Sprintf("Pajak (%s%%)", trimNum(req.TaxPercent)), "+"+rupiah(req.TaxAmount), width))
	}
	b.Write(escBoldOn)
	b.WriteString(twoCol("TOTAL", rupiah(req.GrandTotal), width))
	b.Write(escBoldOff)
	b.WriteString(twoCol("Dibayar", rupiah(req.AmountPaid), width))
	b.WriteString(twoCol("Kembalian", rupiah(req.Change), width))
	b.WriteString(line)

	if strings.TrimSpace(req.Footer) != "" {
		b.Write(escAlignCenter)
		b.WriteString(wrap(req.Footer, width))
		b.Write(escAlignLeft)
	}

	b.Write(escFeed3)
	b.Write(cutPaper)
	return b.Bytes()
}

// writeToSerial menulis byte mentah ke COM port Windows dengan membukanya sebagai
// device file (\\.\COMx). Untuk Bluetooth SPP, baud rate diabaikan sehingga
// menulis langsung sudah cukup (terbukti dari pengujian manual).
func writeToSerial(port string, data []byte) error {
	devPath := port
	if !strings.HasPrefix(strings.ToUpper(port), `\\.\`) {
		devPath = `\\.\` + port
	}
	f, err := os.OpenFile(devPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}
	// Beri jeda singkat agar buffer BT terkirim sebelum port ditutup.
	time.Sleep(400 * time.Millisecond)
	return nil
}

// ---- helper layout ----

// twoCol menata teks kiri & kanan dalam satu baris selebar width. Bila gabungan
// terlalu panjang, teks kiri dipotong aman agar nilai kanan tetap tampil penuh.
func twoCol(left, right string, width int) string {
	if len(left)+len(right)+1 > width {
		maxLeft := width - len(right) - 1
		if maxLeft < 0 {
			maxLeft = 0
		}
		if len(left) > maxLeft {
			left = left[:maxLeft]
		}
	}
	space := width - len(left) - len(right)
	if space < 1 {
		space = 1
	}
	return left + strings.Repeat(" ", space) + right + "\n"
}

// wrap memecah teks panjang menjadi beberapa baris selebar width.
func wrap(text string, width int) string {
	text = strings.TrimRight(text, "\n")
	if len(text) <= width {
		return text + "\n"
	}
	var out strings.Builder
	for len(text) > width {
		out.WriteString(text[:width])
		out.WriteString("\n")
		text = text[width:]
	}
	if len(text) > 0 {
		out.WriteString(text)
		out.WriteString("\n")
	}
	return out.String()
}

// rupiah memformat angka jadi "Rp 1.234.567".
func rupiah(v float64) string {
	n := int64(v)
	neg := n < 0
	if neg {
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	res := "Rp " + strings.Join(parts, ".")
	if neg {
		res = "-" + res
	}
	return res
}

// trimNum menghilangkan desimal .0 agar qty/persen tampil ringkas (2 bukan 2.0).
func trimNum(v float64) string {
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func paymentLabel(m string) string {
	switch m {
	case "cash":
		return "Tunai"
	case "transfer":
		return "Transfer"
	case "qris":
		return "QRIS"
	case "card":
		return "Kartu"
	case "kredit":
		return "Hutang"
	case "balance":
		return "Saldo"
	default:
		return m
	}
}

// formatDate menampilkan tanggal transaksi secara ringkas. Menerima ISO string;
// bila gagal diparse, tampilkan apa adanya.
func formatDate(s string) string {
	if s == "" {
		return ""
	}
	layouts := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			loc := config.Location
			if loc != nil {
				t = t.In(loc)
			}
			return t.Format("02/01/2006 15:04")
		}
	}
	return s
}
