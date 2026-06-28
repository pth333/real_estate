// Mirror từ backend model
export interface RealEstateModel {
  ID: number
  Title: string
  PriceVND: number
  Address: string
  District: string
  City: string
  Acreage: number
  PricePerM2: number
  TypeOfRealEstate: string
  Source: string
  SourceURL: string
  CrawledAt: string
  CreatedAt: string
  UpdatedAt: string
}

export interface DashboardSummary {
  total_posts: number
  avg_price_m2: number
  max_price_m2: number
  min_price_m2: number
}

export interface Filter {
  min_price?: number
  max_price?: number
  district?: string
}

export interface RealEstateSearchRequest {
  page: number
  size: number
  filter?: Filter
}

export interface PaginatedResponse<T> {
  total: number
  data: T[]
}

export interface NotificationItem {
  id: number
  user_id: number
  title: string
  message: string
  is_read: boolean
  created_at: string
}

export interface NotificationSSEPayload {
  type: 'new_listing'
  title: string
  message: string
  price_vnd: number
  acreage: number
  address: string
  source_url: string
}
