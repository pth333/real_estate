export interface MenuResponse {
  success: boolean;
  data: Category[];
}

export interface Category {
  ID: number;
  Name: string;
  Slug: string;
  children?: Category[];
}