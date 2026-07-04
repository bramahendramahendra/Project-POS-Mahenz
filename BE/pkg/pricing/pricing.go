// Package pricing menghitung ulang subtotal/total transaksi di server berdasarkan harga
// produk asli di master data, dipakai bersama oleh checkout langsung (transaction service)
// dan payload transaksi via sync offline, supaya aturan rekalkulasi harga tidak terduplikasi.
package pricing

import (
	product_repo "pos_api/domain/product/repo"
	"pos_api/errors"
)

type Item struct {
	ProductID    int
	ProductName  string
	UnitID       *int
	Quantity     float64
	DiscountItem float64
}

type Totals struct {
	ItemPrices    []float64
	ItemSubtotals []float64
	Subtotal      float64
	TotalAmount   float64
}

// Recalculate menghitung ulang harga per item dari master data produk (bukan dari payload
// client), lalu subtotal dan total keseluruhan setelah diskon/pajak.
func Recalculate(productRepo product_repo.ProductRepo, items []Item, discount, tax float64) (*Totals, error) {
	result := &Totals{
		ItemPrices:    make([]float64, len(items)),
		ItemSubtotals: make([]float64, len(items)),
	}

	for i, item := range items {
		unitPrice, err := resolveUnitPrice(productRepo, item.ProductID, item.UnitID, item.Quantity)
		if err != nil {
			return nil, err
		}

		lineGross := unitPrice * item.Quantity
		if item.DiscountItem < 0 || item.DiscountItem > lineGross {
			return nil, &errors.BadRequestError{Message: "Diskon item tidak valid untuk " + item.ProductName}
		}

		lineNet := lineGross - item.DiscountItem
		result.ItemPrices[i] = unitPrice
		result.ItemSubtotals[i] = lineNet
		result.Subtotal += lineNet
	}

	if discount < 0 || discount > result.Subtotal {
		return nil, &errors.BadRequestError{Message: "Diskon transaksi tidak valid"}
	}

	total := result.Subtotal - discount + tax
	if total < 0 {
		total = 0
	}
	result.TotalAmount = total

	return result, nil
}

// resolveUnitPrice mengambil harga jual asli produk dari master data, bukan dari payload client.
func resolveUnitPrice(productRepo product_repo.ProductRepo, productID int, unitID *int, quantity float64) (float64, error) {
	product, err := productRepo.GetByID(productID)
	if err != nil {
		return 0, &errors.InternalServerError{Message: err.Error()}
	}
	if product == nil {
		return 0, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	if unitID != nil && *unitID > 0 {
		packages, err := productRepo.GetPackagesByProduct(productID)
		if err != nil {
			return 0, &errors.InternalServerError{Message: err.Error()}
		}
		for _, pkg := range packages {
			if pkg.ID == *unitID {
				return pkg.SellingPrice, nil
			}
		}
		return 0, &errors.BadRequestError{Message: "Kemasan produk tidak ditemukan untuk " + product.Name}
	}

	prices, err := productRepo.GetPricesByProduct(productID)
	if err != nil {
		return 0, &errors.InternalServerError{Message: err.Error()}
	}

	price := product.SellingPrice
	bestMinQty := -1.0
	for _, tier := range prices {
		if quantity >= tier.MinQty && tier.MinQty > bestMinQty {
			price = tier.Price
			bestMinQty = tier.MinQty
		}
	}
	return price, nil
}
