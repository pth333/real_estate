import test from "node:test";
import assert from "node:assert/strict";

const { buildDetailUrl } = await import("../app/utils/slug.ts");

/**
 * Test hồi quy cho bug: bấm "xem" trong trang quản lý tin đăng ra
 * /nguoi-ban/<slug> thay vì /<slug>.
 *
 * Nguyên nhân: `navigateTo(`${slug}`)` truyền chuỗi TƯƠNG ĐỐI (không có "/" đầu).
 * Router resolve tương đối theo route hiện tại, mà trang quản lý đang ở
 * /nguoi-ban/... nên segment cuối bị thay bằng slug → /nguoi-ban/<slug>.
 *
 * Hợp đồng bắt buộc: URL trang chi tiết luôn TUYỆT ĐỐI, bắt đầu bằng "/".
 */
test("buildDetailUrl luôn trả về đường dẫn TUYỆT ĐỐI (bắt đầu bằng /)", () => {
  const slug = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-rs3";
  const url = buildDetailUrl(slug);

  assert.equal(url, `/${slug}`);
  assert.ok(url.startsWith("/"), "thiếu / ở đầu thì router sẽ resolve tương đối");
  assert.ok(!url.startsWith("//"), "// sẽ bị hiểu là URL protocol-relative");
});

test("buildDetailUrl giữ nguyên slug dài có dấu gạch và hậu tố -rs{id}", () => {
  const slug = "nha-pho-2-tang-cau-giay-rs123";
  assert.equal(buildDetailUrl(slug), "/nha-pho-2-tang-cau-giay-rs123");
});

test("Slug dạng trang chi tiết phải khớp regex mà trang [category] dùng để nhận diện", () => {
  // pages/[category]/index.vue dò hậu tố -rs(\d+)$ để biết đây là trang chi tiết BĐS
  const detailPattern = /-rs(\d+)$/;

  const url = buildDetailUrl("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-rs3");
  // Bỏ "/" đầu rồi lấy segment cuối — chính là route.params.category
  const category = url.replace(/^\//, "");
  const match = category.match(detailPattern);

  assert.ok(match, "URL phải kết thúc bằng -rs{id} để trang chi tiết nhận ra");
  assert.equal(Number(match[1]), 3);
});

test("Slug chuẩn hoá bị cắt tiền tố lạp lại vẫn ra đúng 1 dấu /", () => {
  // Nếu nơi gọi lỡ truyền kèm "/" thì không được thành "//..."
  assert.equal(buildDetailUrl("/nha-rieng-rs7"), "/nha-rieng-rs7");
});
