import type { SelectOption } from "naive-ui";
import type { Filter } from "~/types/real_estate";

// Nhóm dấu tiếng Việt → ký tự gốc để build slug không dấu
const VIETNAMESE_MAP: Record<string, string> = {
  á: "a", à: "a", ả: "a", ã: "a", ạ: "a", ă: "a", ắ: "a", ằ: "a", ẳ: "a", ẵ: "a", ặ: "a", â: "a", ấ: "a", ầ: "a", ẩ: "a", ẫ: "a", ậ: "a",
  é: "e", è: "e", ẻ: "e", ẽ: "e", ẹ: "e", ê: "e", ế: "e", ề: "e", ể: "e", ễ: "e", ệ: "e",
  í: "i", ì: "i", ỉ: "i", ĩ: "i", ị: "i",
  ó: "o", ò: "o", ỏ: "o", õ: "o", ọ: "o", ô: "o", ố: "o", ồ: "o", ổ: "o", ỗ: "o", ộ: "o", ơ: "o", ớ: "o", ờ: "o", ở: "o", ỡ: "o", ợ: "o",
  ú: "u", ù: "u", ủ: "u", ũ: "u", ụ: "u", ư: "u", ứ: "u", ừ: "u", ử: "u", ữ: "u", ự: "u",
  ý: "y", ỳ: "y", ỷ: "y", ỹ: "y", ỵ: "y",
  đ: "d", Đ: "d",
};

/**
 * Chuyển chuỗi tiếng Việt có dấu → slug kebab lowercase không dấu.
 * VD: "Hồ Chí Minh" → "ho-chi-minh", "Phường 12" → "phuong-12".
 */
export function toSlugPart(input: string): string {
  const normalized = input
    .trim()
    .toLowerCase()
    .split("")
    .map((ch) => VIETNAMESE_MAP[ch] ?? ch)
    .join("");
  return normalized.replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
}

/** Lấy tên (label) từ list SelectOption theo code (value). */
export function nameFromCode(options: SelectOption[], code?: string): string {
  if (!code) return "";
  const found = options.find((item) => item.value === code);
  return typeof found?.label === "string" ? found.label : "";
}

// Đơn vị quy đổi (khớp seed backend filter_ranges)
const TY = 1_000_000_000;

/**
 * Menu khoảng giá SEO — slug phải khớp CHÍNH XÁC với seed bảng `filter_ranges`
 * bên backend (server-driven: backend lookup slug → min/max). Không tự sinh slug
 * thủ công để tránh lệch dữ liệu.
 */
const PRICE_RANGES: Array<{ slug: string; min?: number; max?: number }> = [
  { slug: "gia-duoi-1-ty", max: 1 * TY },
  { slug: "gia-1-den-3-ty", min: 1 * TY, max: 3 * TY },
  { slug: "gia-3-den-5-ty", min: 3 * TY, max: 5 * TY },
  { slug: "gia-5-den-10-ty", min: 5 * TY, max: 10 * TY },
  { slug: "gia-tren-10-ty", min: 10 * TY },
];

/**
 * Menu khoảng diện tích SEO — cùng nguyên tắc server-driven với PRICE_RANGES.
 * Đơn vị m².
 */
const AREA_RANGES: Array<{ slug: string; min?: number; max?: number }> = [
  { slug: "dien-tich-duoi-30", max: 30 },
  { slug: "dien-tich-30-50", min: 30, max: 50 },
  { slug: "dien-tich-50-100", min: 50, max: 100 },
  { slug: "dien-tich-100-200", min: 100, max: 200 },
  { slug: "dien-tich-tren-200", min: 200 },
];

/**
 * Tìm slug chuẩn cho khoảng đã chọn (server-driven — slug khớp bảng filter_ranges).
 * Khớp chính xác cả 2 đầu (min/max); không khai báo 1 đầu = không giới hạn phía đó.
 * Trả "" khi không có khoảng nào khớp (để FE fallback sang query string).
 */
function findRangeSlug(
  ranges: Array<{ slug: string; min?: number; max?: number }>,
  range: { min?: number; max?: number },
): string {
  const hit = ranges.find((r) => {
    const minOk = range.min !== undefined ? r.min === range.min : r.min === undefined;
    const maxOk = range.max !== undefined ? r.max === range.max : r.max === undefined;
    return minOk && maxOk;
  });
  return hit?.slug ?? "";
}

/** Segment giá cho URL list. Trả "" khi không lọc giá (hoặc khoảng không khớp menu). */
export function buildPriceSegment(min?: number, max?: number): string {
  const hasMin = min != null && min > 0;
  const hasMax = max != null && max > 0;
  if (!hasMin && !hasMax) return "";
  return findRangeSlug(PRICE_RANGES, {
    min: hasMin ? min : undefined,
    max: hasMax ? max : undefined,
  });
}

/** Segment diện tích cho URL list. Trả "" khi không lọc (hoặc không khớp menu). */
export function buildAreaSegment(min?: number, max?: number): string {
  const hasMin = min != null && min > 0;
  const hasMax = max != null && max > 0;
  if (!hasMin && !hasMax) return "";
  return findRangeSlug(AREA_RANGES, {
    min: hasMin ? min : undefined,
    max: hasMax ? max : undefined,
  });
}

/**

 *
 * VD: buildListUrl('nha-dat-ban', { city:'79', min_price:1e9, max_price:3e9, bedrooms:2 }, cityOptions)
 *   → "/nha-dat-ban-ha-noi/gia-1-den-3-ty?bedrooms=2"
 */
export function buildListUrl(
  category: string,
  f: Filter | undefined,
  cityOptions: SelectOption[] = [],
): string {
  if (!f) return `/${category}`;

  // Segment 1: category[-city] (city = tên, slug hóa)
  const cityName = nameFromCode(cityOptions, f.city);
  const seg1 = cityName ? `${category}-${toSlugPart(cityName)}` : category;
  const parts: string[] = [seg1];

  // Segment 2+: giá, diện tích (nếu khớp menu chuẩn → append vào path)
  const priceSeg = buildPriceSegment(f.min_price, f.max_price);
  const hasPriceSeg = priceSeg !== "";
  if (priceSeg) parts.push(priceSeg);

  const areaSeg = buildAreaSegment(f.min_area, f.max_area);
  const hasAreaSeg = areaSeg !== "";
  if (areaSeg) parts.push(areaSeg);

  // Query: advanced + fallback giá/diện tích (khoảng không khớp menu)
  const params: Record<string, string> = {};
  if (!hasPriceSeg) {
    if (f.min_price != null) params.min_price = String(f.min_price);
    if (f.max_price != null) params.max_price = String(f.max_price);
  }
  if (!hasAreaSeg) {
    if (f.min_area != null) params.min_area = String(f.min_area);
    if (f.max_area != null) params.max_area = String(f.max_area);
  }
  if (f.bedrooms != null) params.bedrooms = String(f.bedrooms);
  if (f.bathrooms != null) params.bathrooms = String(f.bathrooms);
  if (f.house_direction) params.house_direction = f.house_direction;
  if (f.balcony_direction) params.balcony_direction = f.balcony_direction;
  if (f.legal_docs) params.legal_docs = f.legal_docs;
  if (f.interior) params.interior = f.interior;

  const qs = objectToQueryString(params);
  return `/${parts.join("/")}${qs}`;
}

/** Chuyển object → query string ("a=b&c=d"), trả "" nếu rỗng. */
export function objectToQueryString(params: Record<string, string>): string {
  const pairs = Object.entries(params)
    .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
    .join("&");
  return pairs ? `?${pairs}` : "";
}

/**
 * Build object params từ toàn bộ Filter phẳng — dùng cho axios (query).
 * (FE giữ nguyên, backend đọc query string riêng; không cần thiết trong API.)
 */
// export function buildListParams(f: Filter | undefined): Record<string, string> {
//   const params: Record<string, string> = {};
//   if (!f) return params;
//   for (const [key, value] of Object.entries(f)) {
//     if (value !== undefined && value !== null && value !== "") {
//       params[key] = String(value);
//     }
//   }
//   return params;
// }

/**
 * Build URL trang chi tiết từ slug listing (đã chứa "-rs{id}").
 * VD "nha-pho-2-tang-cau-giay-rs123" → "/nha-pho-2-tang-cau-giay-rs123".
 */
export function buildDetailUrl(listingSlug: string): string {
  return `/${listingSlug}`;
}