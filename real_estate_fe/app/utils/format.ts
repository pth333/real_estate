/**
 * Format giá BĐS VND: tỷ → triệu → đ
 * VD: 2.500.000.000 → "2.5 tỷ", 800.000 → "800.000 đ"
 */
export function formatPrice(priceVND: number | undefined): string {
  if (!priceVND) return "Thỏa thuận";

  if (priceVND >= 1_000_000_000) {
    return `${(priceVND / 1_000_000_000).toLocaleString("vi-VN", { maximumFractionDigits: 2 })} tỷ`;
  }
  if (priceVND >= 1_000_000) {
    return `${(priceVND / 1_000_000).toLocaleString("vi-VN", { maximumFractionDigits: 2 })} triệu`;
  }
  return `${priceVND.toLocaleString("vi-VN")} đ`;
}

/**
 * Format giá/m²
 * VD: 25.000.000 → "25 tr/m²", 500.000 → "500.000 đ/m²"
 */
export function formatPricePerM2(pricePerM2: number | undefined): string {
  if (!pricePerM2) return "Thỏa thuận";

  if (pricePerM2 >= 1_000_000) {
    return `${(pricePerM2 / 1_000_000).toLocaleString("vi-VN", { maximumFractionDigits: 2 })} tr/m²`;
  }
  return `${pricePerM2.toLocaleString("vi-VN")} đ/m²`;
}

/**
 * Format compact hiển thị ở dashboard: "2.5B", "800M", "500K"
 */
export function formatPriceCompact(priceVND: number | undefined): string {
  if (!priceVND) return "Thỏa thuận";

  if (priceVND >= 1_000_000_000) {
    return `${(priceVND / 1_000_000_000).toFixed(1)}B`;
  }
  if (priceVND >= 1_000_000) {
    return `${(priceVND / 1_000_000).toFixed(0)}M`;
  }
  return `${(priceVND / 1000).toFixed(0)}K`;
}

export function formatAcreage(acreage: number | undefined): string {
  if (!acreage) return "—";
  return `${acreage.toLocaleString("vi-VN", { maximumFractionDigits: 2 })} m²`;
}

/**
 * Format tiền theo đơn vị tiền tệ (VND/USD/EUR) với Intl
 */
export function formatCurrency(value: number, unit: string): string {
  const currencyMap: Record<string, string> = {
    vnd: "VND",
    usd: "USD",
    eur: "EUR",
  };
  const currency = currencyMap[unit] ?? "VND";
  return new Intl.NumberFormat("vi-VN", {
    style: "currency",
    currency,
    maximumFractionDigits: 2,
  }).format(value);
}

// ── Format số (dùng cho ô nhập giá / n-input-number) ──────

/**
 * Format số hiển thị theo vi-VN: 3500 → "3.500"
 * Dùng cho prop `format` của n-input-number.
 */
export function formatPriceNumber(value: number | null): string {
  if (value === null) return "";
  return value.toLocaleString("vi-VN");
}

/**
 * Parse từ chuỗi hiển thị về số thực: "3.500" → 3500
 * Dùng cho prop `parse` của n-input-number.
 */
export function parsePriceNumber(input: string): number | null {
  const digits = input.replace(/\D/g, "");
  return digits ? Number(digits) : null;
}

// ── Format thời gian ──────────────────────────────────────

/**
 * Format ngày theo vi-VN: "12/02/2026"
 * VD: "2026-08-08T10:00:00Z" → "08/08/2026"
 */
export function formatDate(dateStr: string): string {
  if (!dateStr) return "—";
  return new Date(dateStr).toLocaleDateString("vi-VN", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

/**
 * Tính "cách đây bao nhiêu" từ date string
 * VD: "Vừa xong", "5 phút trước", "Đăng hôm qua"
 */
export function fromNow(dateStr: string): string {
  if (!dateStr) return "";
  const d = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - d.getTime();
  const diffMin = Math.floor(diffMs / 60000);

  if (diffMin < 1) return "Vừa xong";
  if (diffMin < 60) return `${diffMin} phút trước`;

  const diffHour = Math.floor(diffMin / 60);
  if (diffHour < 24) return `${diffHour} giờ trước`;

  const diffDays = Math.floor(diffHour / 24);
  if (diffDays === 0) return "Đăng hôm nay";
  if (diffDays === 1) return "Đăng hôm qua";
  if (diffDays < 7) return `Đăng ${diffDays} ngày trước`;

  return d.toLocaleDateString("vi-VN");
}

export const formatStatus = (status?: string | boolean): string => {
  switch (status?.toString().toLowerCase()) {
    case "active":
      return "Đang hoạt động";
    case "inactive":
      return "Ngừng hoạt động";
    case "upcoming":
      return "Sắp ra mắt";
    case "completed":
      return "Đã hoàn thành";
    case "paused":
      return "Tạm dừng";
    default:
      return "Không xác định";
  }
};

// Trả về type màu chuẩn Naive UI cho tag (success, warning, error, default)
export const statusTagType = (status?: string) => {
  switch (status?.toLowerCase()) {
    case "active":
      return "success";
    case "inactive":
      return "warning";
    case "upcoming":
      return "info";
    case "completed":
      return "default";
    case "paused":
      return "warning";
    default:
      return "default";
  }
};

export function formatPriceRange(min?: number, max?: number) {
  if (!min && !max) return "Liên hệ";

  const toBillion = (val: number) => {
    return (val / 1000000000).toFixed(1).replace(".0", "") + " tỷ";
  };

  if (min && max) {
    return `${toBillion(min)} - ${toBillion(max)}`;
  }
  if (min) return `Từ ${toBillion(min)}`;
  return `Đến ${toBillion(max!)}`;
}
