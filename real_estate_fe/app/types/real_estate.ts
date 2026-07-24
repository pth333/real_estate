// Mirror từ backend DTO response — field names khớp JSON snake_case
export interface RealEstateResponse {
  id: number
  title: string
  price_vnd: number
  address: string
  district: string
  city: string
  acreage: number
  price_per_m2: number
  type_of_real_estate: string
  created_at: string

  // Optional fields
  images?: string[]
  bedrooms?: number
  bathrooms?: number
  description?: string
  agent_name?: string
  agent_phone?: string
  badge?: string
  is_favorite?: boolean
  source?: string
  source_url?: string
}

export interface DashboardSummary {
  total_posts: number
  avg_price_m2: number
  max_price_m2: number
  min_price_m2: number
}

export interface Filter {
  acreage?: number
  type_of_real_estate?: string
  location?: FilterLocation
  price_range?: FilterPriceRange
}

export interface FilterPriceRange {
  min_price?: number
  max_price?: number
}

export interface FilterLocation {
  district?: string
  city?: string
  ward?: string
  street?: string
}

export interface RealEstateSearchRequest {
  page: number
  cursor_id?: number // keyset cursor: id của bản ghi cuối trang trước
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

export interface CityOption {
  name: string,
  code: string
}

export interface DistrictOption {
  name: string,
  code: string
}

export interface WardOption {
  name: string,
  code: string
}