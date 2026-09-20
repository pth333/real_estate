import type { ImageResponse } from "~/types/real_estate";

/**
 * Dự án rút gọn dùng cho danh sách / khối dự án nổi bật.
 * Trước đây mỗi component tự khai báo 1 interface `Project` cục bộ
 * (ProjectList.vue, FeaturedProjects.vue) — nay gộp về đây để dùng chung.
 */
export interface ProjectSummary {
  id: number;
  name: string;
  slug: string;
  status: string;
  full_address: string;
  description?: string;
  thumbnail?: string;
  total_area_ha?: number;
  total_units?: number;
  price_min?: number;
  price_max?: number;
}

export interface ProjectDetail {
  id: number;
  name: string;
  slug: string;
  status: string;
  full_address: string;
  total_area_ha?: number;
  total_units?: number;
  /** Số căn đã chốt bán — tồn kho còn lại = total_units - sold_units */
  sold_units?: number;
  price_min?: number;
  price_max?: number;
  view_count?: number;
  thumbnail?: string;
  description?: string;
  images?: ImageResponse[]
}
