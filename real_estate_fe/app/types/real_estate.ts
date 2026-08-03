// Mirror từ backend DTO response — field names khớp JSON snake_case
export interface RealEstateResponse {
  id: number;
  title: string;
  price_vnd: number;
  address: string;
  district: string;
  city: string;
  acreage: number;
  price_per_m2: number;
  type_of_real_estate: string;
  created_at: string;

  // Optional fields
  images?: string[];
  bedrooms?: number;
  bathrooms?: number;
  description?: string;
  agent_name?: string;
  agent_phone?: string;
  badge?: string;
  is_favorite?: boolean;
  source?: string;
  source_url?: string;
}

export interface DashboardSummary {
  total_posts: number;
  avg_price_m2: number;
  max_price_m2: number;
  min_price_m2: number;
}

export interface Filter {
  acreage?: number;
  type_of_real_estate?: string;
  location?: FilterLocation;
  price_range?: FilterPriceRange;
}

export interface FilterPriceRange {
  min_price?: number;
  max_price?: number;
}

export interface FilterLocation {
  district?: string;
  city?: string;
  ward?: string;
  street?: string;
}

export interface RealEstateSearchRequest {
  page: number;
  cursor_id?: number; // keyset cursor: id của bản ghi cuối trang trước
  size: number;
  filter?: Filter;
}

export interface PaginatedResponse<T> {
  total: number;
  data: T[];
}

export interface NotificationItem {
  id: number;
  user_id: number;
  title: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

export interface NotificationSSEPayload {
  type: "new_listing";
  title: string;
  message: string;
  price_vnd: number;
  acreage: number;
  address: string;
  source_url: string;
}

export class InformationRealestate {
  // Nhu cầu
  listingType: "sell" | "rent" = "sell";

  // Địa chỉ
  province: string | null = null;
  ward: string | null = null;
  detail_address: string | null = null;

  // Thông tin chính
  real_estate_type: string | null = null;
  area: number | null = null;
  price: number | null = null;
  unit: string | null = "vnd";

  // Thông tin khác
  legal_docs: string | null = null;
  interior: string | null = null;
  bathroom_count: number | null = null;
  bedroom_count: number | null = null;
  house_direction: string | null = null;
  balcony_direction: string | null = null;
  move_in_time: string | null = null;
  price_electricity: number | null = null;
  price_water: number | null = null;
  price_internet: number | null = null;
  amenities: string[] = [];

  // Ảnh / video đã upload — lưu id trả về từ /upload/confirm
  // Dùng để liên kết ảnh với tin đăng khi tạo real_estate
  image_ids: number[] = [];

  // Thông tin liên hệ
  contact_name: string = "";
  contact_email: string = "";
  contact_phone: string = "";

  // Tiêu đề & mô tả
  title: string = "";
  description: string = "";

  tab: "information" | "upload" | "review" = "information";

  isTabInformation() {
    return this.tab === "information";
  }
  isTabUpload() {
    return this.tab === "upload";
  }
  isTabReview() {
    return this.tab === "review";
  }
}

export interface CityOption {
  name: string;
  code: string;
}

export interface DistrictOption {
  name: string;
  code: string;
}

export interface WardOption {
  name: string;
  code: string;
}

export interface OptionTypeRealestate {
  id: number;
  name: string;
}
