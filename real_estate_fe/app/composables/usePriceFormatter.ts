/**
 * Format giá trị tiền tệ VND sang các đơn vị hiển thị:
 * - "2.5 tỷ" / "800 triệu" / "500.000 đ"
 * - "25 tr/m²" / "500.000 đ/m²"
 */
export function usePriceFormatter() {
  /**
   * Format giá BĐS: tỷ → triệu → đ
   */
  function formatPrice(priceVND: number): string {
    if (priceVND >= 1_000_000_000) {
      return `${(priceVND / 1_000_000_000).toLocaleString('vi-VN', { maximumFractionDigits: 2 })} tỷ`
    }
    if (priceVND >= 1_000_000) {
      return `${(priceVND / 1_000_000).toLocaleString('vi-VN', { maximumFractionDigits: 2 })} triệu`
    }
    return `${priceVND.toLocaleString('vi-VN')} đ`
  }

  /**
   * Format giá/m²
   */
  function formatPricePerM2(pricePerM2: number): string {
    if (pricePerM2 >= 1_000_000) {
      return `${(pricePerM2 / 1_000_000).toLocaleString('vi-VN', { maximumFractionDigits: 2 })} tr/m²`
    }
    return `${pricePerM2.toLocaleString('vi-VN')} đ/m²`
  }

  /**
   * Format compact hiển thị ở dashboard: "2.5B"
   */
  function formatPriceCompact(priceVND: number): string {
    if (priceVND >= 1_000_000_000) {
      return `${(priceVND / 1_000_000_000).toFixed(1)}B`
    }
    if (priceVND >= 1_000_000) {
      return `${(priceVND / 1_000_000).toFixed(0)}M`
    }
    return `${(priceVND / 1000).toFixed(0)}K`
  }

  return { formatPrice, formatPricePerM2, formatPriceCompact }
}
