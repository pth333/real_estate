import type { MessageApiInjection } from "naive-ui/es/message/src/MessageProvider";
import type { Category } from "~/types/menu";
import type { UserInfo } from "~/types/auth";

/**
 * Class Menu chứa user_id + categories hiện tại, gán lên window.menu.
 * Dùng để các nơi khác đọc user_id / category mà không cần localStorage.
 */

declare global {
  interface Window {
    message?: MessageApiInjection;
    menu?: Menu;
    userMenu?: UserMenu;
    /** User hiện tại (kèm roles[] + permissions[]) — mirror từ auth store */
    currentUser?: UserInfo;
  }
}

export class Menu {
  settings?: MenuSettings;
  constructor(settings?: MenuSettings) {
    this.settings = settings;
  }
}

export class MenuSettings {
  user_id?: number;
  categories: Category[] = [];
}

/**
 * Class đại diện cho một tùy chọn trong menu người dùng.
 */
export class UserMenuOption {
  key: string;
  label: string;
  path?: string;
  roles?: string[]; // Danh sách các role được phép xem tùy chọn này

  constructor(key: string, label: string, path?: string, roles?: string[]) {
    this.key = key;
    this.label = label;
    this.path = path;
    this.roles = roles;
  }
}

/**
 * Class quản lý danh sách các tùy chọn menu cho người dùng đã đăng nhập.
 * Lưu trữ trực tiếp trên window.userMenu để các nơi khác dễ dàng truy cập và lọc theo role.
 */
export class UserMenu {
  options: UserMenuOption[] = [];
  roles?: string[]; // Danh sách role của người dùng hiện tại

  constructor(roles?: string[]) {
    this.roles = roles;
    this.options = [
      // Mục dành cho khách hàng
      new UserMenuOption("my-deposits", "Đơn đặt cọc của tôi", "/account/deposits", ["CUSTOMER"]),
      // Mục dành cho môi giới
      new UserMenuOption("manage-projects", "Quản lý dự án", "/nguoi-ban/quan-ly-du-an", ["BROKER"]),
      new UserMenuOption("manage-posts", "Quản lý bài viết", "/nguoi-ban/quan-ly-tin-dang", ["BROKER"]),
      new UserMenuOption("manage-deposits", "Đơn đặt cọc xem nhà", "/nguoi-ban/quan-ly-dat-coc",["BROKER"]),
      new UserMenuOption("manage-customers", "Quản lý khách hàng", "/nguoi-ban/quan-ly-khach-hang", ["BROKER"]),
      new UserMenuOption("manage-favorites", "Quản lý yêu thích", "/nguoi-ban/quan-ly-yeu-thich", ["BROKER", "CUSTOMER"]),
      // Mục chỉ admin thấy
      new UserMenuOption("admin-escrow", "Quản trị escrow", "/admin", ["ADMIN"]),
      new UserMenuOption("admin-users", "Người dùng & phân quyền", "/admin/users", ["ADMIN"]),
      new UserMenuOption("logout", "Đăng xuất")
    ];
  }

  /**
   * Lấy danh sách tùy chọn menu đã lọc theo các role hiện tại của người dùng.
   *
   * FAIL-CLOSED: chỉ hiện mục mà user THỰC SỰ có role tương ứng.
   * Trước đây hàm này trả TOÀN BỘ options khi chưa biết role (undefined hoặc mảng rỗng),
   * nên khách hàng vẫn thấy cả mục quản trị/admin nếu quyền chưa kịp nạp.
   * Mục không khai báo `roles` (VD Đăng xuất) thì luôn hiện.
   */
  getFilteredOptions(): UserMenuOption[] {
    const roles = this.roles ?? [];
    return this.options.filter(
      (opt) => !opt.roles || opt.roles.some((role) => roles.includes(role))
    );
  }

  /**
   * Lấy tùy chọn menu bằng key.
   * @param key Key của tùy chọn cần tìm
   */
  getOptionByKey(key: string): UserMenuOption | undefined {
    return this.options.find((opt) => opt.key === key);
  }
}
