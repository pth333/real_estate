import type { SelectOption } from "naive-ui";
import type { FilterLocation, FilterPriceRange } from "~/types/real_estate";

// Nhóm dấu tiếng Việt → ký tự gốc để build slug không dấu
const VIETNAMESE_MAP: Record<string, string> = {
  á: "a", à: "a", ả: "a", ã: "a", ạ: "a", ă: "a", ắ: "a", ằ: "a", ẳ: "a", ẵ: "a", ặ: "a",
  â: "a", ấ: "a", ầ: "a", ẩ: "a", ẫ: "a", ậ: "a",
  é: "e", è: "e", ẻ: "e", ẽ: "e", ẹ: "e", ê: "e", ế: "e", ề: "e", ể: "e", ễ: "e", ệ: "e",
  í: "i", ì: "i", ỉ: "i", ĩ: "i", ị: "i",
  ó: "o", ò: "o", ỏ: "o", õ: "o", ọ: "o", ô: "o", ố: "o", ồ: "o", ổ: "o", ỗ: "o", ộ: "o",
  ơ: "o", ớ: "o", ờ: "o", ở: "o", ỡ: "o", ợ: "o",
  ú: "u", ù: "u", ủ: "u", ũ: "u", ụ: "u", ư: "u", ứ: "u", ừ: "u", ử: "u", ữ: "u", ự: "u",
  ý: "y", ỳ: "y", ỷ: "y", ỹ: "y", ỵ: "y",
  đ: "d", Đ: "d",
};

/**
 * Chuyển chuỗi tiếng Việt có dấu → slug kebab lowercase không dấu.
 * VD: "Hồ Chí Minh" → "ho-chi-minh", "Phường 12" → "phuong-12".
 */
export function toSlugPart(input: string): string {
  const normalized = input.trim().toLowerCase().split("").map((ch) => VIETNAMESE_MAP[ch] ?? ch).join("");
  return normalized
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

/** Lấy tên (label) từ list SelectOption theo code (value). */
function nameFromCode(options: SelectOption[], code?: string): string {
  if (!code) return "";
  const found = options.find((item) => item.value === code);
  return typeof found?.label === "string" ? found.label : "";
}

/** Lấy code (value) từ list SelectOption theo tên (label). */
function codeFromName(options: SelectOption[], name: string): string {
  const found = options.find((item) => typeof item.label === "string" && item.label === name);
  return typeof found?.value === "string" ? found.value : found?.value?.toString() ?? "";
}

/**
 * Build segment đầu của URL SEO: `{category-slug}-{city-name}-{ward-name}`.
 * location chứa code → đổi sang name qua options. Các đoạn trống sẽ bỏ qua.
 */
export function buildLocationSegment(
  categorySlug: string,
  location: FilterLocation,
  cityOptions: SelectOption[],
  wardOptions: SelectOption[],
): string {
  const parts = [categorySlug];
  const cityName = nameFromCode(cityOptions, location.city);
  const wardName = nameFromCode(wardOptions, location.ward);
  if (cityName) parts.push(toSlugPart(cityName));
  if (wardName) parts.push(toSlugPart(wardName));
  return parts.join("-");
}

/**
 * Segment giá: không lọc → "" (rỗng, bỏ hẳn khỏi URL); có min/max →
 * "gia-tu-{min}-den-{max}". Giá truyền vào là VND, format gọn.
 */
export function buildPriceSegment(range: FilterPriceRange | undefined): string {
  if (!range || (range.min_price === undefined && range.max_price === undefined)) {
    return "";
  }
  const min = range.min_price ?? 0;
  const max = range.max_price;
  const parts = ["gia-tu", String(min)];
  if (max !== undefined && max > min) {
    parts.push("den", String(max));
  }
  return parts.join("-");
}

/**
 * Ghép URL SEO: `/{location-segment}/{price-segment}/{page}`.
 * Segment giá rỗng → bỏ; page <= 1 → bỏ. Giữ phần sau `/` gọn nhất có thể.
 */
export function buildSeoUrl(
  locationSegment: string,
  priceSegment: string,
  page: number,
): string {
  const parts = [locationSegment];
  if (priceSegment) parts.push(priceSegment);
  if (page > 1) parts.push(String(page));
  return `/${parts.join("/")}`;
}

/**
 * Giải mã location segment ngược → trả về code. Duyệt đệ quy các name đã biết
 * (city rồi ward) trong list options, lấy phần đầu khớp token.
 */
export function parseLocationFromSegment(
  segment: string,
  cityOptions: SelectOption[],
  wardOptions: SelectOption[],
): FilterLocation {
  const loc: FilterLocation = {};
  let rest = segment;

  // Sinh "slug-reverse" cho từng option name để so khớp tiền tố
  const tryMatch = (options: SelectOption[]) => {
    for (const opt of options) {
      const name = typeof opt.label === "string" ? opt.label : "";
      if (!name) continue;
      const token = toSlugPart(name);
      if (!token) continue;
      if (rest === token || rest.startsWith(`${token}-`)) {
        const code = typeof opt.value === "string" ? opt.value : opt.value?.toString() ?? "";
        rest = rest.slice(token.length + 1);
        return code;
      }
    }
    return "";
  };

  const cityCode = tryMatch(cityOptions);
  if (cityCode) {
    loc.city = cityCode;
    const wardCode = tryMatch(wardOptions);
    if (wardCode) loc.ward = wardCode;
  }
  return loc;
}

/**
 * Giải price segment "gia-tu-{min}-den-{max}" / "tat-ca" về object giá.
 * Trả về rỗng khi không phải dạng hợp lệ.
 */
export function parsePriceFromSegment(segment?: string): FilterPriceRange {
  if (!segment || segment === "") return {};
  const m = segment.match(/^gia-tu-(\d+)(?:-den-(\d+))?$/);
  if (!m) return {};
  const range: FilterPriceRange = {};
  if (m[1]) range.min_price = Number(m[1]);
  if (m[2]) range.max_price = Number(m[2]);
  return range;
}

export { nameFromCode, codeFromName };