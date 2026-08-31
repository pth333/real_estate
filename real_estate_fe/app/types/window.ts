import type { MessageApiInjection } from "naive-ui/es/message/src/MessageProvider";
import type { Category } from "~/types/menu";

/**
 * Class Menu chứa user_id + categories hiện tại, gán lên window.menu.
 * Dùng để các nơi khác đọc user_id / category mà không cần localStorage.
 */

declare global {
  interface Window {
    message?: MessageApiInjection;
    menu?: Menu;
    userMenu?: UserMenu;
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
  role?: string; // Role của người dùng hiện tại

  constructor(role?: string) {
    this.role = role;
    this.options = [
      new UserMenuOption("manage-projects", "Quản lý dự án", "/nguoi-ban/quan-ly-du-an"),
      new UserMenuOption("manage-posts", "Quản lý bài viết", "/nguoi-ban/quan-ly-tin-dang"),
      new UserMenuOption("manage-customers", "Quản lý khách hàng", "/nguoi-ban/quan-ly-khach-hang"),
      new UserMenuOption("manage-favorites", "Quản lý yêu thích", "/nguoi-ban/quan-ly-yeu-thich"),
      new UserMenuOption("logout", "Đăng xuất")
    ];
  }

  /**
   * Lấy danh sách tùy chọn menu đã lọc dựa theo role hiện tại của người dùng.
   */
  getFilteredOptions(): UserMenuOption[] {
    if (!this.role) return this.options;
    return this.options.filter(
      (opt) => !opt.roles || opt.roles.includes(this.role!)
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
