import test from "node:test";
import assert from "node:assert/strict";
import { computed, ref } from "vue";

const { UserMenu } = await import("../app/types/window.ts");

/**
 * Test hồi quy cho bug: menu hiện full option với khách hàng và phải load lại vài lần mới đúng.
 *
 * Nguyên nhân: `new UserMenu(auth.user?.roles)` được khởi tạo MỘT LẦN lúc setup, rồi
 * computed chỉ đọc object thường đó → không có dependency reactive → computed cache vĩnh viễn
 * theo state tại thời điểm mount. Khi store nạp xong quyền, menu KHÔNG cập nhật.
 *
 * Cách đúng: bọc việc tạo UserMenu trong computed để nó phụ thuộc vào `user`.
 */
test("Cách SAI: khởi tạo 1 lần rồi đọc trong computed → menu đóng băng, không cập nhật", () => {
  const user = ref({ roles: undefined });

  // Mô phỏng code cũ trong UserProfileMenu.vue
  const frozenMenu = new UserMenu(user.value.roles);
  const frozenOptions = computed(() => frozenMenu.getFilteredOptions().map((o) => o.key));

  assert.deepEqual(frozenOptions.value, ["logout"], "lúc đầu chưa biết role");

  // Store nạp xong quyền khách hàng
  user.value = { roles: ["CUSTOMER"] };

  assert.deepEqual(
    frozenOptions.value,
    ["logout"],
    "computed bị cache nên KHÔNG cập nhật — đây chính là bug làm menu hiện sai",
  );
});

test("Cách ĐÚNG: bọc việc tạo UserMenu trong computed → menu tự cập nhật khi có quyền", () => {
  const user = ref({ roles: undefined });

  // Mô phỏng code mới trong UserProfileMenu.vue
  const userMenu = computed(() => new UserMenu(user.value.roles));
  const options = computed(() => userMenu.value.getFilteredOptions().map((o) => o.key));

  // Giai đoạn chưa nạp được quyền: fail-closed, chỉ thấy mục ai cũng có
  assert.deepEqual(options.value, ["logout"]);

  // Store nạp xong quyền khách hàng → menu phải tự cập nhật, KHÔNG cần load lại trang
  user.value = { roles: ["CUSTOMER"] };
  assert.deepEqual(options.value, ["my-deposits", "manage-favorites", "logout"]);

  // Admin đổi role cho user này (đang mở web) → tải lại quyền là menu đổi theo
  user.value = { roles: ["CUSTOMER", "BROKER"] };
  assert.deepEqual(options.value, [
    "my-deposits",
    "manage-projects",
    "manage-posts",
    "manage-deposits",
    "manage-customers",
    "manage-favorites",
    "logout",
  ]);

  // Mất hết role (bị thu quyền) → chỉ còn mục ai cũng có
  user.value = { roles: [] };
  assert.deepEqual(options.value, ["logout"]);
});

test("Khách hàng không bao giờ thấy mục quản trị/admin dù ở trạng thái nào", () => {
  const user = ref({ roles: undefined });
  const userMenu = computed(() => new UserMenu(user.value.roles));

  const adminKeys = ["admin-escrow", "admin-users"];

  for (const roles of [undefined, [], ["CUSTOMER"], ["BROKER"], ["CUSTOMER", "BROKER"]]) {
    user.value = { roles };
    const keys = userMenu.value.getFilteredOptions().map((o) => o.key);
    for (const adminKey of adminKeys) {
      assert.ok(
        !keys.includes(adminKey),
        `roles=${JSON.stringify(roles)} không được thấy ${adminKey}`,
      );
    }
  }
});
