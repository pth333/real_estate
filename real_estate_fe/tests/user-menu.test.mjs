import test from "node:test";
import assert from "node:assert/strict";

// Import class UserMenu (sử dụng dynamic import hoặc import bình thường vì đây là ES Module)
const userMenuModule = await import("../app/types/window.ts");
const { UserMenu, UserMenuOption } = userMenuModule;

test("Khởi tạo UserMenu thành công và có đủ 3 options", () => {
  const menu = new UserMenu();
  assert.equal(menu.options.length, 3);

  assert.equal(menu.options[0].key, "manage-posts");
  assert.equal(menu.options[0].label, "Quản lý bài viết");
  assert.equal(menu.options[0].path, "nguoi-ban/quan-ly-tin-dang");

  assert.equal(menu.options[1].key, "manage-customers");
  assert.equal(menu.options[1].label, "Quản lý khách hàng");
  assert.equal(menu.options[1].path, "nguoi-ban/quan-ly-khach-hang");

  assert.equal(menu.options[2].key, "logout");
  assert.equal(menu.options[2].label, "Đăng xuất");
  assert.equal(menu.options[2].path, undefined);
});

test("getOptionByKey trả về đúng option hoặc undefined", () => {
  const menu = new UserMenu();

  const postOpt = menu.getOptionByKey("manage-posts");
  assert.ok(postOpt);
  assert.equal(postOpt.label, "Quản lý bài viết");

  const nonExistentOpt = menu.getOptionByKey("non-existent");
  assert.equal(nonExistentOpt, undefined);
});

test("Lọc tùy chọn dựa trên role thành công", () => {
  const menu = new UserMenu("admin");

  // Gán role chi tiết cho một vài option để test
  menu.options[0].roles = ["seller", "admin"];
  menu.options[1].roles = ["seller"]; // option này admin không xem được

  const filtered = menu.getFilteredOptions();

  // Admin chỉ thấy 2 options: manage-posts và logout
  assert.equal(filtered.length, 2);
  assert.equal(filtered[0].key, "manage-posts");
  assert.equal(filtered[1].key, "logout");
});

