import type { IInformationRealestate, ImageResponse } from "~/types/real_estate";

export interface IManagerPostItem {
  id: number;
  title: string;
  slug: string;
  thumbnail: string;
  type: string;
  price: number; // Đã đổi sang kiểu số thực trả về từ api
  area: number; // Đã đổi sang kiểu số thực
  // status: "approved" | "pending" | "rejected" | "hidden";
  created_at: string;
}

export interface IManagerPostListResponse {
  total: number;
  posts: IManagerPostItem[];
}

export interface ManagerProject {
  id: number;
  name: string;
  slug: string;
  alternative_name: string;
  status: string;
  full_address: string;
  thumbnail: string;
  total_area_ha: number | null;
  total_units: number | null;
  /** Số căn đã chốt bán — chỉ tăng khi admin duyệt tài liệu mua nhà */
  sold_units: number;
  price_min: number | null;
  price_max: number | null;
  created_at: string;
}

// Dữ liệu form tạo/chỉnh sửa dự án (dùng chung giữa trang và component form)
export interface ProjectFormData {
  name: string;
  alternative_name: string;
  province: string | null;
  ward: string | null;
  full_address: string;
  status: string | null;
  category_id: number | null;
  total_area_ha: number | null;
  total_units: number | null;
  price_min: number | null;
  price_max: number | null;
  construction_start_date: string | null;
  handover_date: string | null;
  image_ids: number[];
}

// ── Payload gửi lên API manager ──────────────────────────

/**
 * Payload tạo/sửa tin đăng — lấy đúng những gì trang đăng tin gửi lên.
 * Giống IInformationRealestate, chỉ khác `real_estate_type` là ID loại BĐS
 * (form giữ chuỗi "id-name" để hiển thị select, store đã tách lấy id).
 */
export type CreatePostPayload = Omit<IInformationRealestate, "real_estate_type"> & {
  real_estate_type: string | null;
};

/** Payload cập nhật tin đăng — cùng shape với payload tạo */
export type UpdatePostPayload = CreatePostPayload;

/** Payload tạo dự án — chính là dữ liệu form dự án */
export type CreateProjectPayload = ProjectFormData;

/** Payload cập nhật dự án — cùng shape với payload tạo */
export type UpdateProjectPayload = ProjectFormData;

// ── Response của API manager ─────────────────────────────

/**
 * Chi tiết 1 dự án (GET /manager/projects/:id) để điền form tạo/sửa.
 * province/ward là MÃ tỉnh/phường để select khớp.
 */
export interface ManagerProjectDetailResponse {
  id: number;
  name: string;
  slug: string;
  alternative_name: string;
  status: string;
  full_address: string;
  province: string;
  ward: string;
  total_area_ha: number | null;
  total_units: number | null;
  sold_units: number;
  price_min: number | null;
  price_max: number | null;
  construction_start_date: string;
  handover_date: string;
  category_id: number | null;
  images: ImageResponse[];
}
