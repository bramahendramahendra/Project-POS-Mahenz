export { ProductsPage } from './ProductsPage'
export {
  useProductListQuery,
  useProductSearchQuery,
  useProductDetailQuery,
  useProductPackagesQuery,
  fetchProductDetail,
  fetchProductPackages,
  fetchProductPrices,
  useCreateProductMutation,
  useUpdateProductMutation,
  useDeleteProductMutation,
  useToggleProductStatusMutation,
  useMarkProductReviewedMutation,
} from './products.api'
export type {
  Product,
  ProductSearchOption,
  ProductPackage,
  PriceTier,
  ProductListFilter,
  ProductFilter,
  CreateProductPayload,
  UpdateProductPayload,
} from './products.types'
export { productSchema, grosirSchema } from './products.schema'
export type { ProductFormValues, GrosirFormValues } from './products.schema'
export { formatPackageBreakdown, formatStockNumber, calcMargin } from './products.utils'
