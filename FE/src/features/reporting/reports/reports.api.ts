import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type {
  CashierPerformance,
  ProfitLossReport,
  SalesReport,
  SalesReportFilter,
  SalesReportSummary,
  StockReport,
  StockReportFilter,
} from './reports.types'

interface SalesReportResponse {
  items: SalesReport[]
  total: number
  summary: SalesReportSummary
}

interface StockReportResponse {
  items: StockReport[]
  total: number
  total_stock_value: number
}

export function useSalesReportQuery(filter?: SalesReportFilter) {
  return useQuery({
    queryKey: queryKeys.reports.sales(filter as Record<string, unknown>),
    queryFn: () => api.get<SalesReportResponse>('/reports/sales', filter),
  })
}

export function useProfitLossReportQuery(filter: { date_from?: string; date_to?: string }) {
  return useQuery({
    queryKey: queryKeys.reports.profitLoss(filter as Record<string, unknown>),
    queryFn: () => api.get<ProfitLossReport>('/reports/profit-loss', filter),
  })
}

export function useStockReportQuery(filter?: StockReportFilter) {
  return useQuery({
    queryKey: queryKeys.reports.stock(filter as Record<string, unknown>),
    queryFn: () => api.get<StockReportResponse>('/reports/stock', filter),
  })
}

export function useCashierPerformanceQuery(filter: { date_from?: string; date_to?: string }) {
  return useQuery({
    queryKey: queryKeys.reports.cashierPerformance(filter as Record<string, unknown>),
    queryFn: () => api.get<CashierPerformance[]>('/reports/cashier-performance', filter),
  })
}
