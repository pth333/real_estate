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
