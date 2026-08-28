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
