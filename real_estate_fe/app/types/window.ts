import type { MessageApiInjection } from "naive-ui/es/message/src/MessageProvider";

/**
 * Class Menu chứa user_id hiện tại, gán lên window.menu.
 * Dùng để các nơi khác đọc user_id mà không cần localStorage.
 */

declare global {
  interface Window {
    message?: MessageApiInjection;
    /** Instance Menu chứa user_id hiện tại */
    menu?: Menu;
  }
}

export class Menu {
  user_id?: number;

  constructor(user_id?: number) {
    this.user_id = user_id;
  }
}

const menu = new Menu()

export { menu };
