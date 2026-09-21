import test from "node:test";
import assert from "node:assert/strict";

const { findActiveItem } = await import("../app/utils/dashboard.ts");

/** Icon chỉ là placeholder — util chỉ quan tâm key/match */
const noIcon = {};

function item(key, match, extra = {}) {
  return { key, label: key, pageTitle: key, path: `/${key}`, icon: noIcon, match, ...extra };
}

/**
 * Cấu hình giống thật của khu vực admin: mục escrow khớp đoạn ngắn "/admin"
 * nằm TRƯỚC các mục con trong danh sách.
 */
const adminItems = [
  item("escrow", ["/admin"]),
  item("disputes", ["disputes"]),
  item("users", ["/admin/users"]),
];

test("Bug đã gặp: mục khớp đoạn ngắn không được 'ăn' route con cụ thể hơn", () => {
  // Cách cũ (.find + includes) trả về "escrow" cho mọi path /admin/* → sai
  const keys = ["/admin", "/admin/disputes", "/admin/users"].map(
    (path) => findActiveItem(adminItems, path)?.key,
  );
  assert.deepEqual(keys, ["escrow", "disputes", "users"]);
});

test("Chọn mục có đoạn khớp dài nhất, không phụ thuộc thứ tự khai báo", () => {
  const reversed = [...adminItems].reverse();
  assert.equal(findActiveItem(reversed, "/admin/users")?.key, "users");
  assert.equal(findActiveItem(reversed, "/admin/disputes")?.key, "disputes");
  assert.equal(findActiveItem(reversed, "/admin")?.key, "escrow");
});

test("Khu vực quản lý môi giới: nhận đúng mục theo cả URL gốc lẫn alias tiếng Việt", () => {
  const managerItems = [
    item("posts", ["quan-ly-tin-dang", "posts"]),
    item("projects", ["quan-ly-du-an", "projects"]),
    item("project-form", ["tao-du-an"], { hidden: true }),
    item("deposits", ["quan-ly-dat-coc", "deposits"]),
    item("customers", ["quan-ly-khach-hang", "customers"]),
    item("favorites", ["quan-ly-yeu-thich", "favorites"]),
  ];

  assert.equal(findActiveItem(managerItems, "/nguoi-ban/quan-ly-tin-dang")?.key, "posts");
  assert.equal(findActiveItem(managerItems, "/manager/posts")?.key, "posts");
  assert.equal(findActiveItem(managerItems, "/nguoi-ban/quan-ly-du-an")?.key, "projects");
  assert.equal(findActiveItem(managerItems, "/manager/projects")?.key, "projects");
  assert.equal(findActiveItem(managerItems, "/nguoi-ban/quan-ly-dat-coc")?.key, "deposits");
  assert.equal(findActiveItem(managerItems, "/nguoi-ban/quan-ly-khach-hang")?.key, "customers");
  assert.equal(findActiveItem(managerItems, "/nguoi-ban/quan-ly-yeu-thich")?.key, "favorites");

  // Trang tạo dự án là mục ẩn — vẫn phải nhận ra để ra đúng tiêu đề và tô sáng mục Dự án
  assert.equal(findActiveItem(managerItems, "/nguoi-ban/tao-du-an")?.key, "project-form");
  assert.equal(findActiveItem(managerItems, "/nguoi-ban/tao-du-an")?.hidden, true);
});

test("Không khớp mục nào thì trả undefined để nơi gọi tự chọn mặc định", () => {
  assert.equal(findActiveItem(adminItems, "/trang-chu"), undefined);
  assert.equal(findActiveItem([], "/admin"), undefined);
});
