// Mirror từ backend DTO response — field names khớp JSON snake_case
export interface RealEstateResponse {
  id: number;
  title: string;
  slug?: string;
  price_vnd: number;
  address: string;
  district: string;
  city: string;
  acreage: number;
  price_per_m2: number;
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

  // Thông tin bổ sung (mục Đặc điểm BĐS)
  house_direction?: string;
  balcony_direction?: string;
  floors?: number;
  legal_docs?: string;
  interior?: string;
  price_electricity?: number;
  price_water?: number;
  price_internet?: number;
  amenities?: string[];
}

export interface DashboardSummary {
  total_posts: number;
  avg_price_m2: number;
  max_price_m2: number;
  min_price_m2: number;
}

// Filter phẳng — khớp 1-1 với dto.Filter backend. Mọi trường optional,
// giá trị rỗng/undefined bị bỏ qua khi build URL (không ghi lên query).
export interface Filter {
  // Vị trí (lưu tên/code trực tiếp, không lồng location)
  city?: string;
  district?: string;
  ward?: string;
  street?: string;
  // Giá & diện tích (VNĐ / m²)
  min_price?: number;
  max_price?: number;
  min_acreage?: number;
  max_acreage?: number;
  // Bộ lọc nâng cao (query string)
  bedrooms?: number;
  bathrooms?: number;
  house_direction?: string;
  balcony_direction?: string;
  legal_docs?: string;
  interior?: string;
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

  // Ảnh/video đã upload (id + url + thumbnail) — dùng để khôi phục preview
  images: UploadedMediaItem[] = [];

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

// Ảnh/video đã upload thành công — id + url S3 (public_url/thumbnail)
// Dùng để khôi phục preview khi quay lại bản nháp.
// name + type (MIME thật) lưu kèm để hiển thị alt text và giữ đúng định dạng
// file khi khôi phục (không hardcode image/jpeg hay video/mp4).
export interface UploadedMediaItem {
  imageId: number;
  publicUrl: string;
  thumbnailUrl?: string;
  fileType: "image" | "video";
  name: string;
  type: string;
}
