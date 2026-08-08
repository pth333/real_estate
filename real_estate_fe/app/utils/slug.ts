import type { SelectOption } from "naive-ui";
import type {
  FilterAreaRange,
  FilterLocation,
  FilterPriceRange,
} from "~/types/real_estate";

// Nhóm dấu tiếng Việt → ký tự gốc để build slug không dấu
const VIETNAMESE_MAP: Record<string, string> = {
  á: "a",
  à: "a",
  ả: "a",
  ã: "a",
  ạ: "a",
  ă: "a",
  ắ: "a",
  ằ: "a",
  ẳ: "a",
  ẵ: "a",
  ặ: "a",
  â: "a",
  ấ: "a",
  ầ: "a",
  ẩ: "a",
  ẫ: "a",
  ậ: "a",
  é: "e",
  è: "e",
  ẻ: "e",
  ẽ: "e",
  ẹ: "e",
  ê: "e",
  ế: "e",
  ề: "e",
  ể: "e",
  ễ: "e",
  ệ: "e",
  í: "i",
  ì: "i",
  ỉ: "i",
  ĩ: "i",
  ị: "i",
  ó: "o",
  ò: "o",
  ỏ: "o",
  õ: "o",
  ọ: "o",
  ô: "o",
  ố: "o",
  ồ: "o",
  ổ: "o",
  ỗ: "o",
  ộ: "o",
  ơ: "o",
  ớ: "o",
  ờ: "o",
  ở: "o",
  ỡ: "o",
  ợ: "o",
  ú: "u",
  ù: "u",
  ủ: "u",
  ũ: "u",
  ụ: "u",
  ư: "u",
  ứ: "u",
  ừ: "u",
  ử: "u",
  ữ: "u",
  ự: "u",
  ý: "y",
  ỳ: "y",
  ỷ: "y",
  ỹ: "y",
  ỵ: "y",
  đ: "d",
  Đ: "d",
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
const TRIEU = 1_000_000;

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
 * Tìm slug chuẩn cho 1 khoảng (min/max). Trả "" nếu không khớp menu nào.
 * Khớp chính xác cả 2 đầu (min/max), min/max không khai báo = không giới hạn.
 */
function matchRangeSlug(
  ranges: Array<{ slug: string; min?: number; max?: number }>,
  hasMin: boolean,
  hasMax: boolean,
  min?: number,
  max?: number,
): string {
  const hit = ranges.find((r) => {
    const sameMin = hasMin ? r.min === min : r.min === undefined;
    const sameMax = hasMax ? r.max === max : r.max === undefined;
    return sameMin && sameMax;
  });
  return hit?.slug ?? "";
}

/**
 * Build segment giá SEO từ khoảng đã chọn (server-driven — slug khớp bảng
 * filter_ranges backend, backend quy đổi ngược về VND). Trả về "" khi không lọc giá.
 *   min=1e9 max=3e9   → "gia-1-den-3-ty"
 *   chỉ max=1e9        → "gia-duoi-1-ty"
 *   chỉ min=10e9       → "gia-tren-10-ty"
 */
export function buildPriceSegment(range: FilterPriceRange | undefined): string {
  if (!range) return "";
  const hasMin = range.min_price != null && range.min_price > 0;
  const hasMax = range.max_price != null && range.max_price > 0;
  if (!hasMin && !hasMax) return "";

  return matchRangeSlug(
    PRICE_RANGES,
    hasMin,
    hasMax,
    range.min_price,
    range.max_price,
  );
}

/**
 * Build segment diện tích SEO tương tự buildPriceSegment. Trả "" khi không lọc.
 *   min=50 max=100 → "dien-tich-50-100"
 *   chỉ max=30      → "dien-tich-duoi-30"
 */
export function buildAreaSegment(range: FilterAreaRange | undefined): string {
  if (!range) return "";
  const hasMin = range.min_acreage != null && range.min_acreage > 0;
  const hasMax = range.max_acreage != null && range.max_acreage > 0;
  if (!hasMin && !hasMax) return "";

  return matchRangeSlug(
    AREA_RANGES,
    hasMin,
    hasMax,
    range.min_acreage,
    range.max_acreage,
  );
}

/**
 * Build đường dẫn SEO server-driven: `/{category}/{location?}/{price?}/{area?}`.
 * location chỉ có 1 segment (TÊN thành phố đã slug hóa). page > 1 → thêm query `?page=N`.
 */
export function buildCategoryPath(
  category: string,
  location: FilterLocation,
  cityOptions: SelectOption[],
  price: FilterPriceRange | undefined,
  area?: FilterAreaRange | undefined,
  page = 1,
): string {
  // Segment 1: category[-city]
  const cityName = nameFromCode(cityOptions, location.city);
  const seg1 = cityName
    ? `${category}-${toSlugPart(cityName)}` // 'nha-dat-ban-ha-noi'
    : category;

  const parts: string[] = [seg1];

  // Segment 2+: filters
  const priceSeg = buildPriceSegment(price);
  if (priceSeg) parts.push(priceSeg);

  const areaSeg = buildAreaSegment(area);
  if (areaSeg) parts.push(areaSeg);

  let url = `/${parts.join("/")}`;
  if (page > 1) url += `?page=${page}`;
  return url;
}

/**
 * Build URL trang chi tiết từ slug listing (đã chứa "-rs{id}").
 * VD "nha-pho-2-tang-cau-giay-rs123" → "/nha-pho-2-tang-cau-giay-rs123".
 */
export function buildDetailUrl(listingSlug: string): string {
  return `/${listingSlug}`;
}
