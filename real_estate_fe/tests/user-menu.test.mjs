import test from "node:test";
import assert from "node:assert/strict";

// Import class UserMenu (sử dụng dynamic import hoặc import bình thường vì đây là ES Module)
const userMenuModule = await import("../app/types/window.ts");
const { UserMenu, UserMenuOption } = userMenuModule;

test("Khởi tạo UserMenu thành công và có đủ các mục menu", () => {
  const menu = new UserMenu();
  assert.equal(menu.options.length, 9);

  const byKey = Object.fromEntries(menu.options.map((o) => [o.key, o]));

  // Mục của khách hàng
  assert.equal(byKey["my-deposits"].label, "Đơn đặt cọc của tôi");
  assert.equal(byKey["my-deposits"].path, "/account/deposits");

  // Mục của môi giới
  assert.equal(byKey["manage-posts"].label, "Quản lý bài viết");
  assert.equal(byKey["manage-posts"].path, "/nguoi-ban/quan-ly-tin-dang");
  assert.equal(byKey["manage-deposits"].label, "Đơn đặt cọc xem nhà");
  assert.equal(byKey["manage-deposits"].path, "/nguoi-ban/quan-ly-dat-coc");
  assert.equal(byKey["manage-customers"].label, "Quản lý khách hàng");
  assert.equal(byKey["manage-customers"].path, "/nguoi-ban/quan-ly-khach-hang");

  // Mục chỉ admin thấy
  assert.equal(byKey["admin-escrow"].path, "/admin");
  assert.equal(byKey["admin-users"].path, "/admin/users");

  // Đăng xuất không có path (xử lý riêng)
  assert.equal(byKey["logout"].label, "Đăng xuất");
  assert.equal(byKey["logout"].path, undefined);
});

test("getOptionByKey trả về đúng option hoặc undefined", () => {
  const menu = new UserMenu();

  const postOpt = menu.getOptionByKey("manage-posts");
  assert.ok(postOpt);
  assert.equal(postOpt.label, "Quản lý bài viết");

  const nonExistentOpt = menu.getOptionByKey("non-existent");
  assert.equal(nonExistentOpt, undefined);
});

test("Fail-closed: chưa biết role thì chỉ thấy mục không giới hạn role", () => {
  // undefined = chưa nạp được quyền (VD cookie cũ chưa có roles)
  const unknown = new UserMenu();
  assert.deepEqual(unknown.getFilteredOptions().map((o) => o.key), ["logout"]);

  // [] = đã biết nhưng user không có role nào
  const noRole = new UserMenu([]);
  assert.deepEqual(noRole.getFilteredOptions().map((o) => o.key), ["logout"]);
});

test("UserMenuOption khai báo role dạng mảng để một user giữ nhiều role", () => {
  const option = new UserMenuOption("x", "X", "/x", ["BROKER", "CUSTOMER"]);
  assert.deepEqual(option.roles, ["BROKER", "CUSTOMER"]);
});

test("Admin chỉ thấy mục của admin và mục không giới hạn role", () => {
  const menu = new UserMenu(["ADMIN"]);
  const keys = menu.getFilteredOptions().map((o) => o.key);

  assert.deepEqual(keys, ["admin-escrow", "admin-users", "logout"]);
});

test("Khách hàng thấy mục đơn đặt cọc của mình, không thấy mục môi giới/admin", () => {
  const menu = new UserMenu(["CUSTOMER"]);
  const keys = menu.getFilteredOptions().map((o) => o.key);

  assert.deepEqual(keys, ["my-deposits", "manage-favorites", "logout"]);
});

test("Môi giới thấy mục quản lý của môi giới, không thấy mục admin", () => {
  const menu = new UserMenu(["BROKER"]);
  const keys = menu.getFilteredOptions().map((o) => o.key);

  assert.deepEqual(keys, [
    "manage-projects",
    "manage-posts",
    "manage-deposits",
    "manage-customers",
    "manage-favorites",
    "logout",
  ]);
});

test("User giữ NHIỀU role thấy hợp nhất menu của các role đó", () => {
  const menu = new UserMenu(["CUSTOMER", "BROKER"]);
  const keys = menu.getFilteredOptions().map((o) => o.key);

  assert.deepEqual(keys, [
    "my-deposits",
    "manage-projects",
    "manage-posts",
    "manage-deposits",
    "manage-customers",
    "manage-favorites",
    "logout",
  ]);
});
