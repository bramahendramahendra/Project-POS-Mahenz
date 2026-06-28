export { DashboardPage } from './DashboardPage'

export {
  useDashboardStatsQuery,
  useSalesTrendQuery,
  useTopProductsQuery,
} from './dashboard.api'

export type {
  DashboardPeriod,
  DashboardStats,
  TodayStats,
  MonthStats,
  SalesTrendItem,
  TopProductItem,
  SummaryExtraResponse,
} from './dashboard.types'
