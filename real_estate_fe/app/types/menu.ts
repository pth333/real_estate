export interface MenuResponse {
  success: boolean;
  data: Category[];
  user_id?: number;
}

export interface Category {
  ID: number;
  Name: string;
  Slug: string;
  children?: Category[];
}