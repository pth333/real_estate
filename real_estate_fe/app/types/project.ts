import type { ImageResponse } from "~/types/real_estate";

export interface ProjectDetail {
  id: number;
  name: string;
  slug: string;
  status: string;
  full_address: string;
  total_area_ha?: number;
  total_units?: number;
  price_min?: number;
  price_max?: number;
  view_count?: number;
  thumbnail?: string;
  description?: string;
  images?: ImageResponse[]
}
