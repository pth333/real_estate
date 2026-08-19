export interface ImageResponse {
  id: number;
  url: string;
  thumbnail_url: string;
  file_name: string;
  file_type: string;
  file_size: number;
  status?: string
}

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

  real_estate_type: string;

  // Optional fields
  images?: ImageResponse[];
  image_urls?: string[];
  bedrooms?: number;
  bathrooms?: number;
  description?: string;
  agent_name?: string;
  agent_email?: string;
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
  min_area?: number;
  max_area?: number;
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
  listing_id?: number;
  type: string;
  payload: unknown;
  created_at: string;
}

export interface NotificationSSEPayload {
  title: string;
  address: string;
  price: number;
  acreage: number;
  url: string;
}

export interface OptionTypeRealestate {
  id: number;
  name: string;
}

export interface IInformationRealestate {
  // Nhu cầu
  listingType: "sell" | "rent";

  // Địa chỉ
  province: string | null;
  ward: string | null;
  detail_address: string | null;
  project_id: number | null;

  // Thông tin chính
  real_estate_type: string | null;
  area: number | null;
  price_per_m2: number | null;
  unit: string | null;

  // Thông tin khác
  legal_docs: string | null;
  interior: string | null;
  bathroom_count: number | null;
  bedroom_count: number | null;
  floor_count: number | null;
  house_direction: string | null;
  balcony_direction: string | null;
  move_in_time: string | null;
  price_electricity: number | null;
  price_water: number | null;
  price_internet: number | null;
  amenities: string[];

  // Ảnh/video đã upload (id + url + thumbnail) — dùng để khôi phục preview
  images: UploadedMediaItem[];

  // Thông tin liên hệ
  contact_name: string;
  contact_email: string;
  contact_phone: string;

  // Tiêu đề & mô tả
  title: string;
  description: string;

  tab: "information" | "upload" | "review";
}

export class InformationRealestate {
  // Hàm tạo object mặc định
  static createEmpty(): IInformationRealestate {
    return {
      listingType: "sell",
      province: null,
      ward: null,
      detail_address: null,
      project_id: null,
      real_estate_type: null,
      area: null,
      price_per_m2: null,
      unit: "vnd",
      legal_docs: null,
      interior: null,
      bathroom_count: null,
      bedroom_count: null,
      floor_count: null,
      house_direction: null,
      balcony_direction: null,
      move_in_time: null,
      price_electricity: null,
      price_water: null,
      price_internet: null,
      amenities: [],
      images: [],
      contact_name: "",
      contact_email: "",
      contact_phone: "",
      title: "",
      description: "",
      tab: "information",
    };
  }

  static isTabInformation(form: IInformationRealestate) {
    return form.tab === "information";
  }

  static isTabUpload(form: IInformationRealestate) {
    return form.tab === "upload";
  }

  static isTabReview(form: IInformationRealestate) {
    return form.tab === "review";
  }

  static fromResponse(res: RealEstateResponse): IInformationRealestate {
    const form = this.createEmpty();
    if (!res) return form;

    form.title = res.title || "";
    form.description = res.description || "";
    form.area = res.acreage || null;
    form.price_per_m2 = res.price_per_m2 || null;
    form.unit = "vnd";

    form.real_estate_type = res.real_estate_type || null;

    form.legal_docs = res.legal_docs || null;
    form.interior = res.interior || null;
    form.bedroom_count = res.bedrooms || null;
    form.bathroom_count = res.bathrooms || null;
    form.floor_count = res.floors || null;
    form.house_direction = res.house_direction || null;
    form.balcony_direction = res.balcony_direction || null;

    form.price_electricity = res.price_electricity || null;
    form.price_water = res.price_water || null;
    form.price_internet = res.price_internet || null;

    if (res.amenities) {
      form.amenities = res.amenities;
    }

    form.contact_name = res.agent_name || "";
    form.contact_email = res.agent_email || "";
    form.contact_phone = res.agent_phone || "";

    form.province = res.city;

    form.ward = res.district || null;
    form.detail_address = res.address || null;

    if (res.images && res.images.length > 0) {
      form.images = res.images.map((img) => ({
        id: img.id,
        url: img.url,
        thumnail_url: img.thumbnail_url,
        file_type: img.file_type,
        file_name: img.file_name,
        file_size: img.file_size
      }));
    }

    return form;
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

// Ảnh/video đã upload thành công — id + url S3 (public_url/thumbnail)
// Dùng để khôi phục preview khi quay lại bản nháp.
// name + type (MIME thật) lưu kèm để hiển thị alt text và giữ đúng định dạng
// file khi khôi phục (không hardcode image/jpeg hay video/mp4).
export interface UploadedMediaItem {
  id: number;
  url: string;
  thumnail_url?: string;
  file_type: string;
  file_name: string;
  file_size: number;
  status?: string
}

export interface ProjectOption {
  id: number;
  name: string;
}

export interface CreatePostResponse {
  success: boolean;
  message?: string;
  data?: {
    id: number;
  };
  error?: string;
}

export interface UpdatePostResponse {
  success: boolean;
  message?: string;
  error?: string;
}
