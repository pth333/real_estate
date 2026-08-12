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
